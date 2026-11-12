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

//go:build windows

package webview

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/core"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/pkgs/win"
	"github.com/energye/lcl/types"
	wv "github.com/energye/wv/windows"
)

func NewWebviewWindowParent(owner lcl.IWinControl) core.WindowParent {
	windowParent := wv.NewWindowParent(owner)
	windowParent.SetDoubleBuffered(true)
	windowParent.SetAlign(types.AlClient)
	return windowParent
}

func (m *TWebviewEditor) PlatformPreProcess() {
	browser, ok := m.WVEditor.Browser().(wv.IWVBrowser)
	if !ok {
		return
	}
	// 覆盖 energy 默认的导航/新窗口处理器，由编辑器统一决定
	// 应用内导航还是交给系统默认浏览器。
	browser.SetOnNavigationStarting(func(sender lcl.IObject, webView wv.ICoreWebView2, args wv.ICoreWebView2NavigationStartingEventArgs) {
		args = wv.NewCoreWebView2NavigationStartingEventArgs(args)
		defer args.Free()
		uri := args.URI()
		cancel, openExternal := decideNavigation(uri, m.assetsOrigin)
		if openExternal {
			openExternalAsync(uri)
		}
		if cancel {
			args.SetCancel(true)
		}
	})
	browser.SetOnNewWindowRequested(func(sender lcl.IObject, webView wv.ICoreWebView2, args wv.ICoreWebView2NewWindowRequestedEventArgs) {
		args = wv.NewCoreWebView2NewWindowRequestedEventArgs(args)
		defer args.Free()
		uri := args.URI()
		if decideNewWindow(uri, m.assetsOrigin) {
			openExternalAsync(uri)
		}
		args.SetHandled(true)
	})
}

func (m *TWebviewEditor) PlatformClose() {
	// no impl
}

func (m *TWebviewEditor) SwitchTabPage(owner lcl.IWinControl, windowParent core.WindowParent) {
	newWindowParent := windowParent.(wv.IWVWindowParent)
	currentWindowParent := m.WVEditor.WindowParent().(wv.IWVWindowParent)
	currentWindowParent.SetBrowser(nil)

	browser := m.WVEditor.Browser().(wv.IWVBrowser)
	newWindowParent.SetBrowser(browser)
	if newWindowParent.Parent() == nil {
		newWindowParent.SetParent(owner)
	}
	browser.SetParentWindow(newWindowParent.Handle())
	browser.SetIsVisible(true)
	browser.SetBounds(owner.ClientRect())
	lcl.RunOnMainThreadAsync(func(id uint32) {
		browser.NotifyParentWindowPositionChanged()
		newWindowParent.UpdateSize()
	})
}

func openExternalAsync(rawURL string) {
	logs.Debug("webview open external url:", rawURL)
	go func() {
		if ret := win.ShellExecute(0, "open", rawURL, "", "", win.SW_SHOWNORMAL); ret <= 32 {
			logs.Error("webview open external url failed:", rawURL, "code:", ret)
		}
	}()
}
