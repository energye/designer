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

package uigen

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

// 生成器实例
var gen = &TGenUI{trigger: make(chan event.TEventTrigger, 1), cancel: make(chan bool, 1)}

// TGenUI UI生成器
type TGenUI struct {
	trigger chan event.TEventTrigger // 触发UI生成事件
	cancel  chan bool                // 取消UI生成事件
}

// Start 启动UI生成器
func (m *TGenUI) Start() {
	for {
		select {
		case trigger := <-m.trigger:
			// 处理UI生成事件
			if formTab, ok := trigger.Payload.(*designer.FormTab); ok {
				runDebouncedGenerate(formTab)
			}
		case <-m.cancel:
			// 停止UI生成器
			logs.Info("停止UI生成器")
			return
		}
	}
}

func init() {
	// 注册UI生成事件
	event.GenUI = event.NewEvent(gen.trigger, gen.cancel)
	// 启动UI生成器
	go gen.Start()
}
