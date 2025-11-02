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

package codegen

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

// Go代码生成回调函数

// 生成器实例
var gen = &TGenCode{trigger: make(chan event.TGeneratorTrigger, 1), cancel: make(chan bool, 1)}

// TGenCode	生成器
type TGenCode struct {
	trigger chan event.TGeneratorTrigger // 触发代码生成事件
	cancel  chan bool                    // 取消代码生成事件
}

// Start 启动代码生成器
func (m *TGenCode) Start() {
	for {
		select {
		case trigger := <-m.trigger:
			// 触发代码生成事件
			if trigger.GenType == event.GtCode { //增强判断, 确保是代码生成事件
				if uiFilePath, ok := trigger.Payload.(string); ok {
					err := GenerateCode(uiFilePath)
					if err != nil {
						logs.Error("代码生成错误:", err.Error())
					}
				}
			}
		case <-m.cancel:
			// 停止代码生成器
			logs.Info("停止代码生成器")
			return
		}
	}
}

func init() {
	// 注册代码生成事件
	event.GenCode = event.NewGenerator(gen.trigger, gen.cancel)
	// 启动代码生成器
	go gen.Start()
}
