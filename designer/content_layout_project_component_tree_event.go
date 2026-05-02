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
	"os"

	"github.com/energye/designer/consts"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"unsafe"
)

// 设计 - 组件树 - 事件

// 组件树右键菜单
func (m *ContentLayoutProject) TreeOnContextPopup(sender lcl.IObject, mousePos types.TPoint, handled *bool) {
	node := ProjectTreeSelectNode()
	if node == nil || !node.IsValid() {
		*handled = true
		return
	}
	data := node.Data()
	component := TreeNodeDataToDesigningComponent(data)
	if component != nil {
		// 组件节点：切换回组件菜单
		m.tree.SetPopupMenu(m.componentMenu.treePopupMenu)
		if component.ComponentType == consts.CtForm {
			// Form 根节点不显示菜单
			*handled = true
		}
		return
	}
	// 非组件节点，检查是否为源码节点
	srcPath := ProjectSrcTreeNodePath(node)
	if srcPath == "" {
		*handled = true
		return
	}
	// 源码节点：根据是文件还是文件夹设置菜单项
	info, err := os.Stat(srcPath)
	if err != nil {
		*handled = true
		return
	}
	if info.IsDir() {
		m.srcFileMenu.SetupForFolder()
	} else {
		m.srcFileMenu.SetupForFile()
	}
	// 动态切换弹出菜单
	m.tree.SetPopupMenu(m.srcFileMenu.treePopupMenu)
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
func (m *ContentLayoutProject) TreeOnChange(sender lcl.IObject, node lcl.ITreeNode) {
	//println("[DEBUG] TreeOnChange")
	if node == nil {
		return
	}

	// 先检查是否为组件节点
	data := node.Data()
	component := TreeNodeDataToDesigningComponent(data)
	if component != nil && component.FormTab != nil {
		component.FormTab.SwitchComponentEditing(component)
		return
	}

	// 非组件节点, 检查是否为源码文件节点
	srcPath := ProjectSrcTreeNodePath(node)
	if srcPath == "" {
		return
	}

	// 只处理文件节点 (非目录)
	info, err := os.Stat(srcPath)
	if err != nil || info.IsDir() {
		return
	}

	openFileInAppropriateTab(srcPath)
}

func (m *ContentLayoutProject) TreeOnAdvancedCustomDrawItem(sender lcl.ICustomTreeView, node lcl.ITreeNode, state types.TCustomDrawState, stage types.TCustomDrawStage, paintImages *bool, defaultDraw *bool) {

}

// 数据指针转设计组件
func TreeNodeDataToDesigningComponent(data uintptr) *TDesigningComponent {
	if data == 0 {
		return nil
	}
	dc := (*TDesigningComponent)(unsafe.Pointer(data))
	return dc
}

// 取消选中所有节点
//func (m *InspectorComponentTree) UnSelectedAll() {
//	for _, node := range m.nodeData {
//		node.node.SetSelected(false)
//	}
//}
