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
	"sync"
)

type IWebviewEditor interface {
	IEditor
	LoadURL(url string)
	CreateBrowser()
	SwitchTabPage(owner lcl.IWinControl, windowParent engwv.IWindowParent)
	Webview() wv.IWebview
}

type TWebviewEditor struct {
	WVEditor wv.IWebview
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

		initIPCEvent()
	})
}

func initIPCEvent() {
	ipc.On("monaco-inited", func(context ipc.IContext) {
		logs.Info("ipc monaco-inited BrowserId:", context.BrowserId(), context.Data())
	})
	type Params struct {
		File   string `json:"file"`
		Line   int    `json:"line"`
		Column int    `json:"column"`
	}
	ipc.On("gopls-completion", func(context ipc.IContext) {
		logs.Info("ipc gopls-completion BrowserId:", context.BrowserId(), context.Data())

		//json.Unmarshal(msg.Data, &params)
		//result := gLSPClient.Completion(params.File, params.Line, params.Column)
		context.Result("")
	})
}

func NewWebviewEditor(owner lcl.IWinControl) IEditor {
	WebViewInit()
	m := &TWebviewEditor{}
	m.WVEditor = wv.NewWebview(owner)
	m.WVEditor.SetCaption("")
	m.WVEditor.SetAlign(types.AlClient)
	m.WVEditor.SetName("WVEditor")
	m.LoadURL("energy://designer/index.html")
	m.WVEditor.SetOnLoadChange(func(url, title string, load wv.TLoadChange) {
		logs.Info("OnLoadChange:", url, title, load)
	})
	return m
}

func (m *TWebviewEditor) LoadURL(url string) {
	m.WVEditor.SetDefaultURL(url)
}

func (m *TWebviewEditor) CreateBrowser() {
	m.WVEditor.CreateBrowser()
}

func (m *TWebviewEditor) Type() EditType {
	return EtWebview
}

func (m *TWebviewEditor) Webview() wv.IWebview {
	return m.WVEditor
}
