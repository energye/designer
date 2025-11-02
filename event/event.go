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

package event

import "github.com/energye/designer/pkg/logs"

// 事件处理

// GenType 事件类型
type GenType int32

const (
	GtUI      GenType = iota // UI布局文件
	GtCode                   // 代码
	GtProject                // 项目配置
	GtOther                  // 其它功能
)

// 定义固定的事件
var (
	GenUI   *TEvent // UI 布局文件实例
	GenCode *TEvent // 代码实例
	Project *TEvent // 项目配置更新实例
	Preview *TEvent // 预览实例
)

// TEventTrigger 事件触发器数据结构
type TEventTrigger struct {
	GenType GenType // 类型
	Payload any     // 数据
}

// TEvent 事件实例
type TEvent struct {
	trigger chan TEventTrigger
	cancel  chan bool
}

// NewEvent 创建事件实例
func NewEvent(trigger chan TEventTrigger, cancel chan bool) *TEvent {
	return &TEvent{
		trigger: trigger,
		cancel:  cancel,
	}
}

// TriggerEvent 触发事件
func (m *TEvent) TriggerEvent(trigger TEventTrigger) {
	if m == nil {
		logs.Error("触发事件失败, 当前实例为空")
		return
	}
	m.trigger <- trigger
}

// Cancel 取消事件
func (m *TEvent) Cancel() {
	if m == nil {
		logs.Error("取消事件失败, 当前实例为空")
		return
	}
	m.cancel <- true
}

// CancelAll 取消所有事件, 在退出时调用
func CancelAll() {
	if GenUI != nil {
		GenUI.Cancel()
	}
	if GenCode != nil {
		GenCode.Cancel()
	}
	if Project != nil {
		Project.Cancel()
	}
	if Preview != nil {
		Preview.Cancel()
	}
}
