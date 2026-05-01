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
	"github.com/energye/energy/v3/ipc"
	engwv "github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	wv "github.com/energye/wv/windows"
	"time"
)

func NewWebviewWindowParent(owner lcl.IWinControl) engwv.IWindowParent {
	windowParent := wv.NewWindowParent(owner)
	windowParent.SetDoubleBuffered(true)
	windowParent.SetAlign(types.AlClient)
	return windowParent
}

func (m *TWebviewEditor) SwitchTabPage(owner lcl.IWinControl, windowParent engwv.IWindowParent) {
	currentWindowParent := m.WVEditor.WindowParent().(wv.IWVWindowParent)
	currentWindowParent.SetBrowser(nil)

	browser := m.WVEditor.Browser().(wv.IWVBrowser)
	newWindowParent := windowParent.(wv.IWVWindowParent)
	if newWindowParent.Parent() == nil {
		newWindowParent.SetParent(owner)
	}
	newWindowParent.SetBrowser(browser)
	browser.SetParentWindow(newWindowParent.Handle())
	browser.SetIsVisible(true)
	browser.SetBounds(owner.ClientRect())
	lcl.RunOnMainThreadAsync(func(id uint32) {
		browser.NotifyParentWindowPositionChanged()
		newWindowParent.UpdateSize()
	})
	ipc.Emit("switch-page", "webview-editor"+time.Now().String())
}
