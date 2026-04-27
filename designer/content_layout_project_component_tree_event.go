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
	"github.com/energye/designer/consts"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"unsafe"
)

// 设计 - 组件树 - 事件

// 组件树右键菜单
func (m *ContentLayoutProject) TreeOnContextPopup(sender lcl.IObject, mousePos types.TPoint, handled *bool) {
	node := ProjectTreeSelectNode()
	if node != nil && node.IsValid() {
		data := node.Data()
		component := TreeNodeDataToDesigningComponent(data)
		if component != nil && component.ComponentType == consts.CtForm {
			*handled = true
		}
	}
}

// 组件树鼠标按下事件
func (m *ContentLayoutProject) TreeOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if button == types.MbRight {
		selectNode := ProjectTreeGetNodeAt(X, Y)
		if selectNode != nil && selectNode.IsValid() {
			ProjectTreeSetSelected(selectNode)
		}
	}
}

// 组件树选择事件
func (m *ContentLayoutProject) TreeOnGetSelectedIndex(sender lcl.IObject, node lcl.ITreeNode) {
	data := node.Data()
	component := TreeNodeDataToDesigningComponent(data)
	if component != nil && component.FormTab != nil {
		component.FormTab.SwitchComponentEditing(component)
	}
}

// 数据指针转设计组件
func TreeNodeDataToDesigningComponent(data uintptr) *TDesigningComponent {
	dc := (*TDesigningComponent)(unsafe.Pointer(data))
	return dc
}

// 取消选中所有节点
//func (m *InspectorComponentTree) UnSelectedAll() {
//	for _, node := range m.nodeData {
//		node.node.SetSelected(false)
//	}
//}
