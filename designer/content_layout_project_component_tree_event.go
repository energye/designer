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

// 树选择与画布选择互相同步时的重入保护
var gTreeSelectionSyncing bool

// 组件树鼠标按下事件
func (m *ContentLayoutProject) TreeOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if button == types.MbRight {
		selectNode := ProjectTreeGetNodeAt(X, Y)
		if selectNode != nil && selectNode.IsValid() && !treeNodeInSelection(m.tree, selectNode) {
			// 右键点击未选中节点时切换为单选，点击已选中节点保留多选
			gTreeSelectionSyncing = true
			m.tree.ClearSelection(false)
			m.tree.SetSelected(selectNode)
			gTreeSelectionSyncing = false
			m.syncCanvasFromTree()
		}
	}
}

// 组件树鼠标抬起事件：原生多选处理完成后同步到画布
func (m *ContentLayoutProject) TreeOnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if button != types.MbLeft || gTreeSelectionSyncing {
		return
	}
	m.syncCanvasFromTree()
}

func treeNodeInSelection(tree lcl.ITreeView, node lcl.ITreeNode) bool {
	if tree == nil || node == nil || !node.IsValid() {
		return false
	}
	for i := int32(0); i < int32(tree.SelectionCount()); i++ {
		item := tree.Selections(i)
		if item != nil && item.IsValid() && item.Instance() == node.Instance() {
			return true
		}
	}
	return false
}

// syncCanvasFromTree 将组件树的多选状态同步到画布
func (m *ContentLayoutProject) syncCanvasFromTree() {
	primaryNode := ProjectTreeSelectNode()
	if primaryNode == nil || !primaryNode.IsValid() {
		return
	}
	primaryComp := TreeNodeDataToDesigningComponent(primaryNode.Data())
	if primaryComp == nil || primaryComp.FormTab == nil {
		// 源码文件节点：保持原有打开文件行为
		m.TreeOnChange(nil, primaryNode)
		return
	}
	tree := m.tree
	formTab := primaryComp.FormTab
	comps := make([]*TDesigningComponent, 0, tree.SelectionCount())
	for i := int32(0); i < int32(tree.SelectionCount()); i++ {
		node := tree.Selections(i)
		if node == nil || !node.IsValid() {
			continue
		}
		comp := TreeNodeDataToDesigningComponent(node.Data())
		if comp == nil || comp.FormTab != formTab {
			continue
		}
		comps = append(comps, comp)
	}
	if len(comps) == 0 {
		comps = append(comps, primaryComp)
	}
	formTab.SelectComponents(comps, primaryComp)
}

// ProjectTreeSyncSelection 将画布选中集合同步到项目树（调用方必须位于 UI 线程）
func ProjectTreeSyncSelection(formTab *FormTab) {
	if formTab == nil || gTreeSelectionSyncing || MainWindow.contentLayout == nil {
		return
	}
	tree := MainWindow.contentLayout.layoutProject.tree
	if tree == nil || !tree.IsValid() {
		return
	}
	gTreeSelectionSyncing = true
	defer func() {
		gTreeSelectionSyncing = false
	}()
	tree.BeginUpdate()
	defer tree.EndUpdate()
	tree.ClearSelection(false)
	for _, comp := range formTab.SelectedComponents() {
		if comp.node != nil && comp.node.IsValid() {
			comp.node.SetMultiSelected(true)
		}
	}
	if primary := formTab.PrimarySelected(); primary != nil && primary.node != nil && primary.node.IsValid() {
		tree.SetSelected(primary.node)
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
