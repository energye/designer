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

import (
	"github.com/energye/designer/pkg/logs"
)

// 事件中转处理

// 定义一些默认事件名
const (
	GenUI   = "GenUI"   // 生成UI事件
	GenCode = "GenCode" // 生成代码事件
	Project = "Project" // 项目管理事件
	Preview = "Preview" // 预览事件
	Console = "Console" // 控制台事件
)

// 事件列表
var (
	event *TEvent // 事件对象
	cache = 5     // 通道缓冲容量, 和缓存事件数量一致
)

// TTrigger 事件触发器数据结构
type TTrigger struct {
	Name    string     // 事件名
	Payload any        // 数据
	Result  chan<- any // 返回值
}

// TEvent 事件实例
type TEvent struct {
	list    map[string]*TCallback //事件列表
	trigger chan TTrigger         // 触发事件通道
	cancel  chan bool             // 取消事件通道
}

// TCallback 事件回调
type TCallback struct {
	trigger func(trigger TTrigger) // 触发事件回调函数
	cancel  func()                 // 取消事件回调函数
}

// On 注册一个事件监听器
// name: 事件名称，用于标识要监听的事件类型
// trigger: 触发回调函数，当事件被触发时执行的具体逻辑
// cancel: 取消回调函数，当需要取消监听时执行的清理逻辑
func On(name string, trigger func(trigger TTrigger), cancel func()) {
	if name == "" || trigger == nil || cancel == nil {
		logs.Error("监听事件失败, 必要参数为空")
		return
	}
	callback := &TCallback{trigger: trigger, cancel: cancel}
	event.list[name] = callback
}

// Emit 触发事件
func Emit(trigger TTrigger) {
	event.Trigger(trigger)
}

// 运行事件监听
func (m *TEvent) run() {
	logs.Info("运行事件监听服务")
	for {
		select {
		case trigger := <-m.trigger:
			if callback, ok := m.list[trigger.Name]; ok {
				go callback.trigger(trigger)
			}
		case <-m.cancel:
			logs.Info("停止所有事件监听服务")
			for _, callback := range m.list {
				go callback.cancel()
			}
			return
		}
	}
}

// Trigger 触发事件
func (m *TEvent) Trigger(trigger TTrigger) {
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
	close(m.trigger)
	close(m.cancel)
}

// CancelAll 取消所有事件, 在退出时调用
func CancelAll() {
	event.Cancel()
}

func init() {
	event = new(TEvent)
	event.list = make(map[string]*TCallback)
	event.trigger = make(chan TTrigger, cache)
	event.cancel = make(chan bool, cache)
	go event.run()
}
