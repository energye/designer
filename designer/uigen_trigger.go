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
	"github.com/energye/designer/pkg/vtedit"
)

// 触发生成代码事件

// 触发UI布局生成事件
//
//	component 更新的组件
//	nodeData  组件的属性数据
//	  nil 时是 xxx.ui.go 更新, 全量更新
//	  非 nil 时指定组件属性, 自引用更新或事件更新
//	type 更新类型: ui 布局, 自引用, 绑定事件
func triggerUIGeneration(component *TDesigningComponent, nodeData *vtedit.TEditNodeData, type_ event.Type) {
	if component == nil {
		return
	}
	payload := event.TPayload{Type: type_, Data: TUIGenerationData{Component: component, NodeData: nodeData}}
	event.Emit(event.TTrigger{Name: event.GenUI, Payload: payload})
}

// TUIGenerationData UI 生成数据载体
type TUIGenerationData struct {
	Component *TDesigningComponent // 更新的组件
	// 组件的属性数据
	// nil 时是 xxx.ui.go 更新, 全量更新
	// 非 nil 时指定组件属性, 自引用更新或事件更新
	NodeData *vtedit.TEditNodeData // 组件的属性数据
}
