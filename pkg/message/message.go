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

package message

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"time"
)

// 消息弹框

type TMessage struct {
	alphaStep   byte
	showTimer   lcl.ITimer
	afterTime   *time.Timer
	displayTime int32
	hintWindow  lcl.IHintWindow
}

var (
	message *TMessage
)

func mustMessage() {
	if message == nil {
		message = new(TMessage)
		message.alphaStep = 5
		// 当鼠标指针悬停在某个对象上时，弹出的包含辅助信息的对话框。
		// THintWindow 是 TCustomForm 的子类，用于显示文本提示。它不适合与子控件搭配使用。
		//
		// HintWindow := THintWindow.Create(nil);
		// Rect := HintWindow.CalcHintRect(0, 'This is the hint',nil);
		// HintWindow.ActivateHint(Rect, 'This is the hint');

		message.hintWindow = lcl.NewHintWindow(nil)
		message.hintWindow.SetAlphaBlend(true)
		message.hintWindow.SetAlphaBlendValue(0)
		message.hintWindow.Canvas().SetAntialiasingMode(types.AmOn)
		message.hintWindow.SetControlStyle(message.hintWindow.ControlStyle().Include(types.CsParentBackground))
		showTimer := lcl.NewTimer(message.hintWindow)
		showTimer.SetEnabled(false)
		showTimer.SetInterval(15)
		showTimer.SetOnTimer(message.OnShowTimer)
		message.showTimer = showTimer
	}
}

func Follow(content string) {
	mustMessage()
	cursorPos := lcl.Mouse.CursorPos()
	hintRect := message.hintWindow.CalcHintRect(0, content, 0)
	w, h := hintRect.Width(), hintRect.Height()
	hintRect.Left = cursorPos.X + 20
	hintRect.Top = cursorPos.Y + 20
	hintRect.SetWidth(w)
	hintRect.SetHeight(h)
	message.hintWindow.SetAlphaBlendValue(255)

	message.hintWindow.SetBoundsRect(hintRect)
	message.hintWindow.SetCaption(content)
	message.hintWindow.Show()
	message.hintWindow.Invalidate()
	//message.hintWindow.ActivateHintWithRectString(hintRect, content)
}

func FollowHide() {
	mustMessage()
	if message != nil {
		message.hintWindow.SetAlphaBlendValue(0)
		message.hintWindow.Hide()
	}
}

func Info(title, content string, w, h int32) {
	mustMessage()
	msg := title + "\n  " + content
	message.displayTime = 3 // 秒
	message.showTimer.SetEnabled(true)
	windowCenterRect := rect(w, h)
	message.hintWindow.ActivateHintWithRectString(windowCenterRect, msg)
}

func (m *TMessage) OnShowTimer(sender lcl.IObject) {
	abv := m.hintWindow.AlphaBlendValue()
	if abv >= 255 {
		m.showTimer.SetEnabled(false)
		// displayTime 秒后关闭
		m.afterTime = time.AfterFunc(time.Second*time.Duration(m.displayTime), func() {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				FollowHide()
			})
		})
		return
	}
	m.hintWindow.SetAlphaBlendValue(abv + m.alphaStep)
}

func rect(width, height int32) types.TRect {
	rect := lcl.Application.MainForm().BoundsRect()
	rect.Left = rect.Left + rect.Width()/2 - width/2
	rect.Top = rect.Top + rect.Height()/2 - height/2
	rect.SetWidth(width)
	rect.SetHeight(height)
	return rect
}
