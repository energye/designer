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
)

var (
	gPLSClient *gopls.PLSClient
)

// InitPLS 启动异步 PLS 初始化，不阻塞调用方
func InitPLS() {
	if gPLSClient != nil {
		gPLSClient.Close()
		gPLSClient = nil
	}
	go initPLSAsync()
}

func initPLSAsync() {
	var err error
	gPLSClient, err = gopls.NewPLSClient(bean.GPath)
	if err != nil {
		logs.Error("NewPLSClient:", err)
		return
	}
	rootURI := FilePathToURI(bean.GPath)
	logs.Info("gopls 初始化, rootURI:", rootURI, "GPath:", bean.GPath)
	if err := gPLSClient.Initialize(rootURI); err != nil {
		logs.Error("gopls Initialize 失败:", err)
		return
	}

	gPLSClient.SetDiagnosticsHandler(func(uri string, diagnostics []gopls.Diagnostic) {
		filePath := URIToFilePath(uri)
		if filePath == "" {
			return
		}
		logs.Info("gopls 诊断: file=", filePath, "count=", len(diagnostics))
		diagData, _ := json.Marshal(diagnostics)
		lcl.RunOnMainThreadAsync(func(id uint32) {
			ipc.Emit("gopls-diagnostics", filePath, string(diagData))
		})
	})
	logs.Info("gopls 初始化就绪")
}

// PLSClient 返回全局 PLS 客户端实例，未就绪时返回 nil
func PLSClient() *gopls.PLSClient {
	return gPLSClient
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
	fileURI := FilePathToURI(filePath)
	go func() {
		gPLSClient.DidClose(fileURI)
		gPLSClient.DidOpen(fileURI, DetectLanguage(filePath), string(content), 1)
	}()
}
