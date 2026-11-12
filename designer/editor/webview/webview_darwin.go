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

//go:build darwin

package webview

import (
	"github.com/energye/energy/v3/core"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	wvDarwin "github.com/energye/wv/darwin"
)

func NewWebviewWindowParent(owner lcl.IWinControl) core.WindowParent {
	windowParent := wvDarwin.NewWebviewParent(owner)
	windowParent.SetDoubleBuffered(true)
	windowParent.SetAlign(types.AlClient)
	return windowParent
}

func (m *TWebviewEditor) PlatformPreProcess() {
	// macOS 无需额外的引用计数处理
}

func (m *TWebviewEditor) PlatformClose() {
	// macOS 无需额外的引用计数处理
}

// SwitchTabPage 把 WKWebView 挂到目标标签页的父容器。
// Cocoa 的 addSubview 会自动把视图从旧 superview 移除，因此不调用
// SetWebview(0)，避免旧容器释放底层 NSView 导致切换后白屏。
func (m *TWebviewEditor) SwitchTabPage(owner lcl.IWinControl, windowParent core.WindowParent) {
	newWindowParent, ok := windowParent.(wvDarwin.IWkWebviewParent)
	if !ok || newWindowParent == nil {
		return
	}
	if part := newWindowParent.Parent(); part == nil || !part.IsVisible() {
		newWindowParent.SetParent(owner)
	}
	m.currentWindowParent = newWindowParent
	if browser, ok := m.WVEditor.Browser().(wvDarwin.IWkWebview); ok {
		newWindowParent.SetWebview(browser.Data())
	}
	lcl.RunOnMainThreadAsync(func(id uint32) {
		newWindowParent.SetBoundsRect(owner.ClientRect())
		newWindowParent.Resize()
		newWindowParent.SetFocus()
	})
}
