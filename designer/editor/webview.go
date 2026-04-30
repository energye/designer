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
	"fmt"
	"github.com/energye/designer/resources/editor"
	"github.com/energye/energy/v3/application"
	"github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"sync"
)

type IWebviewEditor interface {
	IEditor
	LoadURL(url string)
	CreateBrowser()
}

type TWebviewEditor struct {
	WVEditor wv.IWebview
}

var (
	wvInitOnce sync.Once
	gWVApp     *wv.Application
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
	})
}

func NewWebviewEditor(owner lcl.IWinControl) IEditor {
	WebViewInit()
	m := &TWebviewEditor{}
	m.WVEditor = wv.NewWebview(owner)
	m.WVEditor.SetCaption("")
	m.WVEditor.SetName(owner.Name() + "_WVEditor")
	m.WVEditor.SetHeight(owner.Height())
	m.WVEditor.SetWidth(owner.Width())
	m.WVEditor.SetAlign(types.AlClient)
	m.WVEditor.SetParent(owner)
	//m.WVEditor.SetDefaultURL("auto:blank")
	m.WVEditor.SetOnLoadChange(func(url, title string, load wv.TLoadChange) {
		fmt.Println(url, title, load)
	})
	m.WVEditor.SetOnBrowserAfterCreated(func(sender lcl.IObject) {
		fmt.Println("SetOnBrowserAfterCreated")
	})
	return m
}

func (m *TWebviewEditor) LoadURL(url string) {
	m.WVEditor.SetDefaultURL(url)
}

func (m *TWebviewEditor) CreateBrowser() {
	m.WVEditor.CreateBrowser()
}
