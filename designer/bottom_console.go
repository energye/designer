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

package designer

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"widget/wg"
)

type TConsole struct {
	split      lcl.ISplitter
	consoleBox lcl.IPanel
	console    lcl.IMemo
	closeBtn   lcl.IBitBtn
}

// createConsole 创建控制台界面组件
// 该函数初始化控制台相关的UI元素，包括分割器、控制台面板和文本显示区域
func (m *BottomBox) createConsole() {
	console := new(TConsole)
	m.console = console

	console.split = lcl.NewSplitter(m.box)
	console.split.SetAlign(types.AlBottom)
	//console.split.SetResizeStyle(types.RsNone)
	console.split.SetHeight(splitterWidth)
	console.split.SetParent(m.rightBox)

	console.consoleBox = lcl.NewPanel(m.box)
	console.consoleBox.SetBevelColor(wg.LightenColor(bgDarkColor, 0.3))
	console.consoleBox.SetBevelOuter(types.BvLowered)
	console.consoleBox.SetDoubleBuffered(true)
	console.consoleBox.SetAlign(types.AlBottom)
	//console.consoleBox.SetColor(colors.ClBlue)
	console.consoleBox.SetHeight(100)
	console.consoleBox.SetParent(m.rightBox)

	console.console = lcl.NewMemo(m.box)
	console.console.SetAlign(types.AlClient)
	console.console.SetBorderStyle(types.BsNone)
	console.console.SetScrollBars(types.SsAutoBoth)
	console.console.SetDoubleBuffered(true)
	console.console.SetReadOnly(true)
	SetComponentDefaultColor(console.console)
	console.console.Font().SetColor(bgDarkColor)
	console.console.SetParent(console.consoleBox)
}

// 写入控制台
func WriteConsole(text string) {
	mainWindow.box.WriteConsole(text)
}

// 清空控制台
func ClearConsole() {
	mainWindow.box.ClearConsole()
}

// 写入控制台
func (m *BottomBox) WriteConsole(text string) {
	m.console.console.Lines().Add(text)
}

// 清空控制台
func (m *BottomBox) ClearConsole() {
	m.console.console.Lines().Clear()
}

// 初始化消息控制台事件
func initConsoleEvent() {
	event.On(event.Console, func(trigger event.TTrigger) {
		payload, ok := trigger.Payload.(event.TPayload)
		if ok {
			call := func() {
				if payload.Type == event.ConsoleInfo {
					WriteConsole(payload.Data.(string))
				} else {
					ClearConsole()
				}
			}
			if tool.IsMainThread() {
				call()
			} else {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					call()
				})
			}
		}
	}, func() {
		logs.Info("停止控制台消息处理器")
	})
}
