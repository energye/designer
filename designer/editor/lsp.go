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
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
	"os"
	"sync"
	"time"
)

var (
	plsInitOnce  sync.Once
	plsStartOnce sync.Once
	gPLSClient   *gopls.PLSClient
	plsReady     bool
	plsMu        sync.RWMutex
)

// InitPLS 启动异步 PLS 初始化，不阻塞调用方
func InitPLS() {
	plsStartOnce.Do(func() {
		go initPLSAsync()
	})
}

func initPLSAsync() {
	plsInitOnce.Do(func() {
		var err error
		gPLSClient, err = gopls.NewPLSClient(bean.GPath)
		if err != nil {
			logs.Error("NewPLSClient:", err)
			return
		}
		rootURI := filePathToURI(bean.GPath)
		logs.Info("gopls 初始化, rootURI:", rootURI, "GPath:", bean.GPath)
		if err := gPLSClient.Initialize(rootURI); err != nil {
			logs.Error("gopls Initialize 失败:", err)
			return
		}
		logs.Info("gopls 等待初始索引完成...")
		time.Sleep(2 * time.Second)
		logs.Info("gopls 初始化就绪")

		gPLSClient.SetDiagnosticsHandler(func(uri string, diagnostics []gopls.Diagnostic) {
			filePath := uriToFilePath(uri)
			if filePath == "" {
				return
			}
			logs.Info("gopls 诊断: file=", filePath, "count=", len(diagnostics))
			diagData, _ := json.Marshal(diagnostics)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-diagnostics", filePath, string(diagData))
			})
		})

		plsMu.Lock()
		plsReady = true
		plsMu.Unlock()
	})
}

// PLSClient 返回全局 PLS 客户端实例，未就绪时返回 nil
func PLSClient() *gopls.PLSClient {
	return gPLSClient
}

// IsPLSReady 返回 gopls 是否已初始化完成
func IsPLSReady() bool {
	plsMu.RLock()
	defer plsMu.RUnlock()
	return plsReady
}

// SetDiagnosticsHandler 设置诊断处理器，允许原生编辑器覆盖默认行为
func SetDiagnosticsHandler(handler func(uri string, diagnostics []gopls.Diagnostic)) {
	if gPLSClient != nil {
		gPLSClient.SetDiagnosticsHandler(handler)
	}
}

// notifyFileChanged 通知 gopls 文件已被外部修改
// 使用 DidClose+DidOpen 强制 gopls 重新索引，确保补全和诊断信息更新
func notifyFileChanged(filePath string) {
	if gPLSClient == nil {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	fileURI := filePathToURI(filePath)
	go func() {
		gPLSClient.DidClose(fileURI)
		gPLSClient.DidOpen(fileURI, detectLanguage(filePath), string(content), 1)
	}()
}
