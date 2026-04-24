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
)

// 写入控制台
func WriteConsole(text string) {
	MainWindow.contentLayout.layoutConsoleLog.WriteDesignerLog(text)
}

// 清空控制台
func ClearConsole() {
	MainWindow.contentLayout.layoutConsoleLog.ClearConsole()
}

// 初始化消息控制台事件
func initConsoleEvent() {
	logs.Println("初始化消息控制台事件")
	event.On(event.Console, func(trigger event.TTrigger) {
		payload, ok := trigger.Payload.(event.TPayload)
		if ok {
			call := func() {
				if payload.Type == event.ConsoleInfo {
					WriteConsole(payload.Data.(string))
				} else {
					//ClearConsole()
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
		logs.Println("停止控制台消息处理器")
	})
}
