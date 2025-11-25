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
// +build windows

package hook

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/pkgs/win"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
	"syscall"
	"unsafe"
)

func (m *TWindowHook) Hook() {
	if m.Form == nil {
		return
	}
	wndProcCallback := syscall.NewCallback(m.wndProc)
	handle := m.Form.Handle()
	m.oldWndPrc = win.SetWindowLongPtr(handle, win.GWL_WNDPROC, wndProcCallback)
	fmt.Println("hook:", handle)
}

func (m *TWindowHook) wndProc(hwnd types.HWND, message uint32, wParam, lParam uintptr) uintptr {
	fmt.Println("wndProc:", hwnd, message)
	switch message {
	case messages.WM_DPICHANGED:
		if !lcl.Application.Scaled() {
			newWindowSize := (*types.TRect)(unsafe.Pointer(lParam))
			win.SetWindowPos(m.Form.Handle(), uintptr(0),
				newWindowSize.Left, newWindowSize.Top, newWindowSize.Right-newWindowSize.Left, newWindowSize.Bottom-newWindowSize.Top,
				win.SWP_NOZORDER|win.SWP_NOACTIVATE)
		}
	}
	switch message {
	case messages.WM_ACTIVATE:
		win.ExtendFrameIntoClientArea(m.Form.Handle(), win.Margins{CxLeftWidth: 1, CxRightWidth: 1, CyTopHeight: 1, CyBottomHeight: 1})
	case messages.WM_NCCALCSIZE:
		if wParam != 0 {
			isMaximize := uint32(win.GetWindowLong(m.Form.Handle(), win.GWL_STYLE))&win.WS_MAXIMIZE != 0
			if isMaximize {
				rect := (*types.TRect)(unsafe.Pointer(lParam))
				monitor := win.MonitorFromRect(rect, win.MONITOR_DEFAULTTONULL)
				if monitor != 0 {
					var monitorInfo types.TMonitorInfo
					monitorInfo.CbSize = types.DWORD(unsafe.Sizeof(monitorInfo))
					if win.GetMonitorInfo(monitor, &monitorInfo) {
						*rect = monitorInfo.RcWork
					}
				}
			}
			return 0
		}
	}
	return win.CallWindowProc(m.oldWndPrc, uintptr(hwnd), message, wParam, lParam)
}

func (m *TWindowHook) RestoreWndProc() {
	if m.oldWndPrc != 0 {
		win.SetWindowLongPtr(m.Form.Handle(), win.GWL_WNDPROC, m.oldWndPrc)
		m.oldWndPrc = 0
	}
}
