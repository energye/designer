// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package editor

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
)

var (
	plsMu      sync.RWMutex
	gPLSClient *gopls.PLSClient
	plsGen     uint64
)

// InitPLS 关闭旧实例并异步启动新的 gopls 客户端，不阻塞调用方。
func InitPLS() {
	plsMu.Lock()
	old := gPLSClient
	gPLSClient = nil
	plsGen++
	gen := plsGen
	plsMu.Unlock()
	if old != nil {
		old.Close()
	}
	// 旧客户端的文档状态随它一起作废，新页面会重新 didOpen。
	ResetPLSQueue()
	StartPLSDispatcher()
	setPLSStatus(false, "gopls 正在启动")
	go initPLSAsync(gen)
}

// StopPLS 关闭 gopls 客户端，并让尚未完成的异步初始化丢弃其结果。
func StopPLS() {
	plsMu.Lock()
	plsGen++
	client := gPLSClient
	gPLSClient = nil
	plsMu.Unlock()
	if client != nil {
		client.Close()
	}
	ResetPLSQueue()
	setPLSIndexing(false)
}

// plsGenerationCurrent 异步初始化失败时，不覆盖新一轮启动的状态。
func plsGenerationCurrent(gen uint64) bool {
	plsMu.RLock()
	defer plsMu.RUnlock()
	return gen == plsGen
}

func initPLSAsync(gen uint64) {
	client, err := gopls.NewPLSClient(bean.GPath)
	if err != nil {
		logs.Error("NewPLSClient:", err)
		if plsGenerationCurrent(gen) {
			setPLSStatus(false, err.Error())
		}
		return
	}
	// 首次索引期间请求会排队，gopls 可能需要几十秒才响应；
	// 状态要同步给前端，前端据此延长 pending 超时。
	// 客户端被替换后，旧客户端的兜底定时器不能再改全局状态。
	client.SetIndexingHandler(func(indexing bool) {
		if PLSClient() != client {
			return
		}
		setPLSIndexing(indexing)
	})

	rootURI := FilePathToURI(bean.GPath)
	logs.Debug("gopls 初始化, rootURI:", rootURI, "GPath:", bean.GPath)
	if err := client.Initialize(rootURI); err != nil {
		logs.Error("gopls Initialize 失败:", err)
		client.Close()
		if plsGenerationCurrent(gen) {
			setPLSStatus(false, "gopls 初始化失败: "+err.Error())
		}
		return
	}

	// 只有 initialize 成功之后才发布客户端：在此之前所有请求快速返回空结果，
	// 文档通知也停在派发队列里，不会发给尚未初始化的 gopls。
	client.SetDiagnosticsHandler(func(uri string, diagnostics []gopls.Diagnostic) {
		// 客户端已被替换或停止后不再向已销毁的浏览器发消息。
		if PLSClient() != client {
			return
		}
		filePath := URIToFilePath(uri)
		if filePath == "" {
			return
		}
		logs.Debug("gopls 诊断: file=", filePath, "count=", len(diagnostics))
		diagData, _ := json.Marshal(diagnostics)
		lcl.RunOnMainThreadAsync(func(id uint32) {
			ipc.Emit("gopls-diagnostics", filePath, string(diagData))
		})
	})

	plsMu.Lock()
	if gen != plsGen {
		plsMu.Unlock()
		client.Close()
		return
	}
	gPLSClient = client
	plsMu.Unlock()
	if !plsGenerationCurrent(gen) {
		client.Close()
		return
	}
	setPLSStatus(true, "")
	client.BeginIndexing()
	logs.Debug("gopls 初始化就绪")
}

// PLSClient 返回全局 PLS 客户端实例，未就绪时返回 nil
func PLSClient() *gopls.PLSClient {
	plsMu.RLock()
	defer plsMu.RUnlock()
	return gPLSClient
}

// SetDiagnosticsHandler 设置诊断处理器，允许原生编辑器覆盖默认行为
func SetDiagnosticsHandler(handler func(uri string, diagnostics []gopls.Diagnostic)) {
	if client := PLSClient(); client != nil {
		client.SetDiagnosticsHandler(handler)
	}
}

// notifyFileChanged 通知 gopls 文件已被外部修改
// 使用 DidClose+DidOpen 强制 gopls 重新索引，确保补全和诊断信息更新
func notifyFileChanged(filePath string) {
	// 编辑器里有未保存内容时，不能用磁盘内容重置 gopls，
	// 否则用户正在编辑的文本会被静默替换（代码生成/外部改文件都会走这里）。
	if IsFileDirty(filePath) {
		logs.Debug("文件在编辑器中有未保存内容，跳过 gopls 重载:", filePath)
		return
	}
	client := PLSClient()
	if client == nil {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	fileURI := FilePathToURI(filePath)
	// 与前端 didClose/didOpen 走同一条 FIFO 队列，且不再把 version 重置成 1：
	// 版本由派发器统一单调分配，避免与前端仍在递增的版本号冲突。
	EnqueuePLSNotification(PLSNotification{Kind: PLSNotifyClose, URI: fileURI})
	EnqueuePLSNotification(PLSNotification{
		Kind:       PLSNotifyOpen,
		URI:        fileURI,
		LanguageID: DetectLanguage(filePath),
		Content:    string(content),
	})
}
