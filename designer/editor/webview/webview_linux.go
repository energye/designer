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

//go:build linux

package webview

import (
	. "github.com/energye/energy/v3/platform/linux/types"
	wvEng "github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
	wvLinux "github.com/energye/wv/linux"
)

func NewWebviewWindowParent(owner lcl.IWinControl) wvEng.IWindowParent {
	windowParent := wvLinux.NewWebviewParent(owner)
	windowParent.SetDoubleBuffered(true)
	windowParent.SetAlign(types.AlClient)
	return windowParent
}

func (m *TWebviewEditor) PlatformPreProcess() {
	if m.refd.CompareAndSwap(false, true) {
		webkit2wv := m.WVEditor.(*wvEng.TWebview).GtkWebview()
		webkit2wv.Ref()
		//webkit2wv.OpenDevTools() // debug
		currentWindowParentHwnd := func() types.HWND {
			return m.currentWindowParent.(wvLinux.IWkWebviewParent).Handle()
		}
		// 为了让 webview 获取焦点，并让 menu editing 菜单生效
		// 可以让webview控件为活动控件
		webkit2wv.SetOnFocusIn(func(sender PGtkWidget, event PGdkEventFocus, userData GPointer) bool {
			hwnd := currentWindowParentHwnd()
			//println("[debug] webkit2wv.SetOnFocusIn hwnd:", hwnd)
			if hwnd != 0 {
				msg := &types.TLMessage{Msg: messages.LM_SETFOCUS}
				api.Gtk3Widget_DeliverMessage(currentWindowParentHwnd(), msg, false)
			}
			return false
		})
		webkit2wv.SetOnFocusOut(func(sender PGtkWidget, event PGdkEventFocus, userData GPointer) bool {
			hwnd := currentWindowParentHwnd()
			//println("[debug] webkit2wv.SetOnFocusOut hwnd:", hwnd)
			if hwnd != 0 {
				msg := &types.TLMessage{Msg: messages.LM_KILLFOCUS}
				api.Gtk3Widget_DeliverMessage(hwnd, msg, false)
			}
			return false
		})
	}
}

func (m *TWebviewEditor) PlatformClose() {
	if m.refd.Load() {
		webkit2wv := m.WVEditor.(*wvEng.TWebview).GtkWebview()
		webkit2wv.Unref()
	}
}

func (m *TWebviewEditor) SwitchTabPage(owner lcl.IWinControl, windowParent wvEng.IWindowParent) {
	newWindowParent := windowParent.(wvLinux.IWkWebviewParent)
	if part := newWindowParent.Parent(); part == nil || !part.IsVisible() {
		newWindowParent.SetParent(owner)
	}
	if m.currentWindowParent != nil {
		m.currentWindowParent.(wvLinux.IWkWebviewParent).SetWebview(nil)
	}
	m.currentWindowParent = newWindowParent
	browser := m.WVEditor.Browser().(wvLinux.IWkWebview)
	newWindowParent.SetWebview(browser)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		cr := owner.ClientRect()
		newWindowParent.SetBoundsRect(cr)
		m.WVEditor.(*wvEng.TWebview).GtkWebview().GrabFocus()
	})
}
