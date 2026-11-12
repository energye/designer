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
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/pkg/assetserve"
	"github.com/energye/designer/pkg/logs"
	reseditor "github.com/energye/designer/resources/editor"
	"github.com/energye/energy/v3/application"
	"github.com/energye/energy/v3/core"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type IWebviewEditor interface {
	editor.IEditor
	LoadURL(url string)
	CreateBrowser()
	SwitchTabPage(owner lcl.IWinControl, windowParent core.WindowParent)
	Webview() wv.IWebview
	Initialized() bool
	SetCanLoadChan(canLoad chan error)
}

type TWebviewEditor struct {
	WVEditor            wv.IWebview
	currentWindowParent core.WindowParent
	fileManager         *editor.FileManager
	checkTimer          *time.Ticker
	stopChan            chan struct{}
	stopOnce            sync.Once
	canLoadChan         chan error
	initialized         bool
	refd                atomic.Bool
	assetsOrigin        string
}

var (
	// WebView 运行时是进程级全局状态，只启动一次。
	wvAppOnce sync.Once
	gWVApp    *wv.Application

	// Monaco 资源服务单独管理，Stop 时关闭，下次 Init 可重新启动。
	assetsMu      sync.Mutex
	gAssetsServer *http.Server
	gAssetsPort   int
)

// Init 启动 WebView 运行时与 Monaco 资源服务，返回资源服务端口。
func Init() int {
	wvAppOnce.Do(func() {
		gWVApp = wv.Init()
		gWVApp.SetOptions(application.Options{
			Linux: application.Linux{HardwareGPU: application.HGPUEnable},
		})
		gWVApp.Start()
	})

	assetsMu.Lock()
	defer assetsMu.Unlock()
	if gAssetsServer != nil {
		return gAssetsPort
	}

	asset := assetserve.New()
	asset.AssetsFSName = "monaco"
	asset.Assets = reseditor.Assets
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		logs.Error("monaco assets server listen failed:", err)
		return 0
	}
	server := &http.Server{Handler: asset}
	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Error("monaco assets server stopped:", err)
		}
	}()
	gAssetsServer = server
	gAssetsPort = ln.Addr().(*net.TCPAddr).Port
	return gAssetsPort
}

// StopAssetsServer 关闭 Monaco 资源服务；可在下次 Init 时重新拉起。
func StopAssetsServer() {
	assetsMu.Lock()
	server := gAssetsServer
	gAssetsServer = nil
	gAssetsPort = 0
	assetsMu.Unlock()
	if server == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		_ = server.Close()
	}
}

func monacoURL(port int) string {
	return monacoAssetsOrigin(port) + "/index.html"
}

func NewWebviewEditor(owner lcl.IWinControl) editor.IEditor {
	port := Init()
	editor.InitPLS()
	m := &TWebviewEditor{
		fileManager: editor.NewFileManager(),
		stopChan:    make(chan struct{}),
	}
	m.WVEditor = wv.NewWebview(owner)
	m.WVEditor.SetCaption("")
	m.WVEditor.SetAlign(types.AlClient)
	m.WVEditor.SetName("WVEditor")
	m.assetsOrigin = monacoAssetsOrigin(port)
	m.LoadURL(monacoURL(port))
	m.WVEditor.SetOnLoadChange(func(url, title string, load core.TLoadChange) {
		if load == core.LcFinish {
			m.initialized = true
			m.sendCanLoadSignal()
			// 页面加载完成后补推一次 gopls 状态，覆盖「状态上报早于页面就绪」的情况。
			ok, reason := editor.PLSStatus()
			ipc.Emit("gopls-status", ok, reason)
			ipc.Emit("gopls-indexing", editor.PLSIndexing())
		}
	})
	m.PlatformPreProcess()

	ipc.SetMainDefaultBrowserId(0)
	ipc.RegisterProcessMessage(m.WVEditor)

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
		isDirty := m.fileManager.IsDirty(filePath)

		fi, err := os.Stat(filePath)
		if err == nil {
			// 同步基准，避免同一变更每秒重复通知；冲突的取舍交给前端决定。
			m.fileManager.UpdateModTime(filePath, fi.ModTime())
			editor.UpdateSavedModTime(filePath)
		}

		if isDirty {
			logs.Debug("文件有未保存修改且被外部变更:", filePath)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("file-conflict-detected", filePath)
			})
		} else {
			logs.Debug("文件被外部修改，通知前端重新加载:", filePath)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("file-changed-externally", filePath)
			})
		}
	}
}

// Stop 幂等地释放本实例占用的资源。WebView 运行时本身是进程级的，
// 不在这里销毁；资源服务与 gopls 会在下次创建编辑器时重新拉起。
func (m *TWebviewEditor) Stop() {
	m.stopOnce.Do(func() {
		m.PlatformClose()
		ipc.UnRegisterProcessMessage(m.WVEditor)
		if m.checkTimer != nil {
			m.checkTimer.Stop()
		}
		close(m.stopChan)
		editor.StopFormFileWatcher()
		editor.StopPLS()
		StopAssetsServer()
		if editor.GetCurrentEditor() == m {
			editor.SetCurrentEditor(nil)
		}
	})
}
