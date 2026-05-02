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
	"github.com/energye/designer/resources/editor"
	"github.com/energye/energy/v3/application"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/energy/v3/wv"
	engwv "github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type IWebviewEditor interface {
	IEditor
	LoadURL(url string)
	CreateBrowser()
	SwitchTabPage(owner lcl.IWinControl, windowParent engwv.IWindowParent)
	Webview() wv.IWebview
	Initialized() bool
	SetCanLoadChan(canLoad chan error)
}

type TWebviewEditor struct {
	WVEditor    wv.IWebview
	fileManager *FileManager
	checkTimer  *time.Ticker
	stopChan    chan struct{}
	canLoadChan chan error
	initialized bool
}

var (
	wvInitOnce sync.Once
	gWVApp     *wv.Application
	gLSPClient *gopls.LSPClient
)

func WebViewInit() {
	wvInitOnce.Do(func() {
		gWVApp = wv.Init(nil, nil)
		gWVApp.SetLocalLoad(application.LocalLoad{
			Scheme:     "energy",
			Domain:     "designer",
			ResRootDir: "assets",
			FS:         editor.Assets,
		})
		gWVApp.Start()
		var err error
		gLSPClient, err = gopls.NewLSPClient(bean.GPath)
		if err != nil {
			logs.Error("NewLSPClient:", err)
		} else {
			gLSPClient.Initialize(bean.GPath)
		}
	})
}

func (m *TWebviewEditor) initIPCEvent() {
	ipc.On("monaco-inited", func(context ipc.IContext) {
		logs.Info("ipc monaco-inited BrowserId:", context.BrowserId(), context.Data())
	})

	type CompletionParams struct {
		File   string `json:"file"`
		Line   int    `json:"line"`
		Column int    `json:"column"`
	}

	ipc.On("gopls-completion", func(context ipc.IContext) {
		logs.Info("ipc gopls-completion BrowserId:", context.BrowserId(), context.Data())

		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("")
			return
		}

		var params CompletionParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			logs.Error("解析完成请求参数失败:", err)
			context.Result("")
			return
		}

		if gLSPClient != nil {
			result, err := gLSPClient.Completion(params.File, params.Line, params.Column)
			if err != nil {
				logs.Error("gopls-completion:", err)
			}
			jsonData, _ := json.Marshal(result)
			context.Result(string(jsonData))
		} else {
			context.Result("")
		}
	})

	type FileData struct {
		File     string `json:"file"`
		Content  string `json:"content"`
		Language string `json:"language"`
		ModTime  int64  `json:"modTime"`
	}

	ipc.On("open-file-request", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			logs.Error("无效的数据格式")
			context.Result("")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			logs.Error("文件路径格式错误")
			context.Result("")
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			logs.Error("读取文件失败:", err)
			context.Result("")
			return
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logs.Error("获取文件信息失败:", err)
			context.Result("")
			return
		}

		result := FileData{
			File:     filePath,
			Content:  string(content),
			Language: detectLanguage(filePath),
			ModTime:  fileInfo.ModTime().UnixMilli(),
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			logs.Error("序列化文件数据失败:", err)
			context.Result("")
			return
		}

		context.Result(string(jsonData))
	})

	ipc.On("save-file", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			logs.Error("无效的数据格式")
			context.Result("error: invalid data")
			return
		}

		var fileData FileData
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &fileData); err != nil {
			logs.Error("解析保存数据失败:", err)
			context.Result("error: " + err.Error())
			return
		}

		if err := os.WriteFile(fileData.File, []byte(fileData.Content), 0644); err != nil {
			logs.Error("保存文件失败:", err)
			context.Result("error: " + err.Error())
			return
		}

		fileInfo, _ := os.Stat(fileData.File)
		if fileInfo != nil {
			m.fileManager.UpdateModTime(fileData.File, fileInfo.ModTime())
			m.fileManager.SetDirty(fileData.File, false)
		}

		logs.Info("文件已保存:", fileData.File)
		context.Result("ok")
	})

	ipc.On("reload-file-request", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			logs.Error("无效的数据格式")
			context.Result("")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			logs.Error("文件路径格式错误")
			context.Result("")
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			logs.Error("读取文件失败:", err)
			context.Result("")
			return
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logs.Error("获取文件信息失败:", err)
			context.Result("")
			return
		}

		result := FileData{
			File:     filePath,
			Content:  string(content),
			Language: detectLanguage(filePath),
			ModTime:  fileInfo.ModTime().UnixMilli(),
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			logs.Error("序列化文件数据失败:", err)
			context.Result("")
			return
		}

		context.Result(string(jsonData))
	})

	ipc.On("register-opened-file", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			logs.Error("无效的数据格式")
			context.Result("error")
			return
		}

		type RegData struct {
			File    string `json:"file"`
			ModTime int64  `json:"modTime"`
		}
		var regData RegData
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &regData); err != nil {
			logs.Error("解析注册文件数据失败:", err)
			context.Result("error")
			return
		}

		m.fileManager.RegisterFile(regData.File)
		context.Result("ok")
	})

	ipc.On("unregister-opened-file", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			logs.Error("无效的数据格式")
			context.Result("error")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			logs.Error("文件路径格式错误")
			context.Result("error")
			return
		}

		m.fileManager.UnregisterFile(filePath)
		context.Result("ok")
	})

	ipc.On("set-file-dirty", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			logs.Error("无效的数据格式")
			context.Result("error")
			return
		}

		type DirtyData struct {
			File    string `json:"file"`
			IsDirty bool   `json:"isDirty"`
		}
		var dirtyData DirtyData
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &dirtyData); err != nil {
			logs.Error("解析脏状态数据失败:", err)
			context.Result("error")
			return
		}

		m.fileManager.SetDirty(dirtyData.File, dirtyData.IsDirty)
		context.Result("ok")
	})
}

func NewWebviewEditor(owner lcl.IWinControl) IEditor {
	WebViewInit()
	m := &TWebviewEditor{
		fileManager: NewFileManager(),
		stopChan:    make(chan struct{}),
	}
	m.WVEditor = wv.NewWebview(owner)
	m.WVEditor.SetCaption("")
	m.WVEditor.SetAlign(types.AlClient)
	m.WVEditor.SetName("WVEditor")
	m.LoadURL("energy://designer/index.html")
	m.WVEditor.SetOnLoadChange(func(url, title string, load wv.TLoadChange) {
		logs.Info("OnLoadChange:", url, title, load)
		if load == wv.LcFinish {
			m.initialized = true
			m.sendCanLoadSignal()
		}
	})

	m.initIPCEvent()
	m.startFileChangeChecker()

	SetCurrentEditor(m)

	return m
}

func (m *TWebviewEditor) SetCanLoadChan(canLoad chan error) {
	m.canLoadChan = canLoad
}

func (m *TWebviewEditor) LoadURL(url string) {
	m.WVEditor.SetDefaultURL(url)
}

func (m *TWebviewEditor) CreateBrowser() {
	m.WVEditor.CreateBrowser()
	if m.initialized {
		m.sendCanLoadSignal()
	}
}

func (m *TWebviewEditor) sendCanLoadSignal() {
	if m.canLoadChan == nil {
		return
	}
	select {
	case m.canLoadChan <- nil:
		// 发送成功
	default:
		// channel 已满，跳过
	}
}

func (m *TWebviewEditor) Type() EditType {
	return EtWebview
}

func (m *TWebviewEditor) Webview() wv.IWebview {
	return m.WVEditor
}

func (m *TWebviewEditor) Initialized() bool {
	return m.initialized
}

func (m *TWebviewEditor) startFileChangeChecker() {
	m.checkTimer = time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-m.stopChan:
				return
			case <-m.checkTimer.C:
				m.checkFileChanges()
			}
		}
	}()
}

func (m *TWebviewEditor) checkFileChanges() {
	changedFiles := m.fileManager.CheckExternalChanges()
	if len(changedFiles) == 0 {
		return
	}

	for _, filePath := range changedFiles {
		state, _ := m.fileManager.GetFileState(filePath)
		if state != nil && state.IsDirty {
			ipc.Emit("file-conflict-detected", filePath)
			logs.Warn("文件冲突:", filePath)
		} else {
			ipc.Emit("file-changed-externally", filePath)
			logs.Info("文件被外部修改:", filePath)
		}
	}
}

func (m *TWebviewEditor) Stop() {
	if m.checkTimer != nil {
		m.checkTimer.Stop()
		close(m.stopChan)
	}
}

func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".jsx":
		return "javascript"
	case ".tsx":
		return "typescript"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".scss":
		return "scss"
	case ".less":
		return "less"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".md":
		return "markdown"
	case ".py":
		return "python"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".hpp", ".cc":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".php":
		return "php"
	case ".rb":
		return "ruby"
	case ".sql":
		return "sql"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	case ".sh":
		return "shell"
	case ".bat":
		return "batch"
	default:
		return "plaintext"
	}
}
