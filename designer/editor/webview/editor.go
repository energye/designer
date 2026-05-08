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

package webview

import (
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/pkg/logs"
	reseditor "github.com/energye/designer/resources/editor"
	"github.com/energye/energy/v3/application"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/energy/v3/wv"
	engwv "github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"os"
	"sync"
	"time"
)

type IWebviewEditor interface {
	editor.IEditor
	LoadURL(url string)
	CreateBrowser()
	SwitchTabPage(owner lcl.IWinControl, windowParent engwv.IWindowParent)
	Webview() wv.IWebview
	Initialized() bool
	SetCanLoadChan(canLoad chan error)
}

type TWebviewEditor struct {
	WVEditor    wv.IWebview
	fileManager *editor.FileManager
	checkTimer  *time.Ticker
	stopChan    chan struct{}
	canLoadChan chan error
	initialized bool
}

var (
	wvInitOnce sync.Once
	gWVApp     *wv.Application
)

func Init() {
	wvInitOnce.Do(func() {
		gWVApp = wv.Init(nil, nil)
		gWVApp.SetLocalLoad(application.LocalLoad{
			Scheme:     "energy",
			Domain:     "designer",
			ResRootDir: "monaco",
			FS:         reseditor.Assets,
		})
		gWVApp.Start()
	})
}

func NewWebviewEditor(owner lcl.IWinControl) editor.IEditor {
	Init()
	editor.InitPLS() // TODO 功能未完善
	m := &TWebviewEditor{
		fileManager: editor.NewFileManager(),
		stopChan:    make(chan struct{}),
	}
	m.WVEditor = wv.NewWebview(owner)
	m.WVEditor.SetCaption("")
	m.WVEditor.SetAlign(types.AlClient)
	m.WVEditor.SetName("WVEditor")
	m.LoadURL("energy://designer/index.html")
	m.WVEditor.SetOnLoadChange(func(url, title string, load wv.TLoadChange) {
		if load == wv.LcFinish {
			m.initialized = true
			m.sendCanLoadSignal()
		}
	})

	m.initIPCEvent()
	m.initDefinitionIPC()

	m.startFileChangeChecker()

	editor.StopFormFileWatcher()
	editor.StartFormFileWatcher()

	editor.SetCurrentEditor(m)

	return m
}

func init() {
	editor.RegisterEditorFactory(editor.EtWebview, NewWebviewEditor)
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
	default:
	}
}

func (m *TWebviewEditor) Type() editor.EditType {
	return editor.EtWebview
}

func (m *TWebviewEditor) OpenFile(filePath string, readOnly ...bool) {
	isReadOnly := len(readOnly) > 0 && readOnly[0]
	ipc.Emit("open-file", filePath, isReadOnly)
}

func (m *TWebviewEditor) CloseFile(filePath string) {
	ipc.Emit("close-file", filePath)
}

func (m *TWebviewEditor) SaveCurrentFile() {
	ipc.Emit("save-current-file", "")
}

func (m *TWebviewEditor) FileManager() *editor.FileManager {
	return m.fileManager
}

func (m *TWebviewEditor) PLSClient() *gopls.PLSClient {
	return editor.PLSClient()
}

func (m *TWebviewEditor) Webview() wv.IWebview {
	return m.WVEditor
}

func (m *TWebviewEditor) Initialized() bool {
	return m.initialized
}

func (m *TWebviewEditor) startFileChangeChecker() {
	m.checkTimer = time.NewTicker(time.Second)
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

		fi, err := os.Stat(filePath)
		if err == nil {
			m.fileManager.UpdateModTime(filePath, fi.ModTime())
			editor.UpdateSavedModTime(filePath)
		}

		if state != nil && state.IsDirty {
			logs.Info("文件有未保存修改且被外部变更:", filePath)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("file-conflict-detected", filePath)
			})
		} else {
			logs.Info("文件被外部修改，通知前端重新加载:", filePath)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("file-changed-externally", filePath)
			})
		}
	}
}

func (m *TWebviewEditor) Stop() {
	if m.checkTimer != nil {
		m.checkTimer.Stop()
		close(m.stopChan)
	}
	editor.StopFormFileWatcher()
}
