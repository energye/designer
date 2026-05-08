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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

type ContentLayoutProject struct {
	topBox            lcl.IPanel
	title             lcl.ILabel
	box               lcl.IPanel
	tree              lcl.ITreeView
	projectTreeNode   lcl.ITreeNode
	componentTreeNode lcl.ITreeNode
	srcTreeNode       lcl.ITreeNode
	componentMenu     *TComponentMenu
	srcFileMenu       *TSrcFileMenu
}

func initContentLayoutProject(owner *ContentLayout) *ContentLayoutProject {
	m := &ContentLayoutProject{}
	m.topBox = lcl.NewPanel(owner.projectPanel)
	m.topBox.SetBorderStyleToBorderStyle(types.BsNone)
	m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	m.topBox.SetParent(owner.projectPanel)

	title := lcl.NewLabel(m.topBox)
	title.SetCaption("项目管理器")
	title.SetLeft(5)
	title.SetTop(5)
	font := title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	title.SetParent(m.topBox)

	m.box = lcl.NewPanel(owner.projectPanel)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlClient)
	m.box.SetParent(owner.projectPanel)

	m.componentMenu = initComponentMenu(m.box)
	m.srcFileMenu = initSrcFileMenu(m.box)

	m.tree = lcl.NewTreeView(owner.projectPanel)
	m.tree.SetAutoExpand(false)
	m.tree.SetReadOnly(true)
	m.tree.SetDoubleBuffered(true)
	m.tree.SetTreeLineColor(colors.RGBToColor(128, 128, 128))
	m.tree.SetAlign(types.AlClient)
	m.tree.SetRowSelect(true)
	m.tree.SetBorderStyleToBorderStyle(types.BsNone)
	m.tree.SetTreeLinePenStyle(types.PsClear)
	m.tree.Font().SetHeight(-11)
	m.tree.SetImages(imageComponents.ImageList50())
	//m.tree.SetImages(imageComponents.ImageList100())
	//m.tree.SetMultiSelect(true) // 多选控制
	//m.tree.SetOnChange(m.TreeOnChange)
	m.tree.SetOnExpanding(m.TreeOnExpanding)
	m.tree.SetOnMouseDown(m.TreeOnMouseDown)
	m.tree.SetOnContextPopup(m.TreeOnContextPopup)
	//m.tree.SetOnAdvancedCustomDrawItem(m.TreeOnAdvancedCustomDrawItem)
	//m.CreateComponentMenu()
	m.tree.SetPopupMenu(m.componentMenu.treePopupMenu)
	m.tree.SetParent(m.box)

	items := m.tree.Items()

	m.projectTreeNode = items.AddChild(nil, "-")
	m.projectTreeNode.SetImageIndex(imageComponents.ImageIndex("folder.png"))
	m.projectTreeNode.SetSelectedIndex(imageComponents.ImageIndex("folder.png"))
	m.projectTreeNode.SetExpanded(true)

	m.componentTreeNode = items.AddChild(m.projectTreeNode, "Forms")
	m.componentTreeNode.SetImageIndex(imageComponents.ImageIndex("design.png"))
	m.componentTreeNode.SetSelectedIndex(imageComponents.ImageIndex("design.png"))
	m.componentTreeNode.SetExpanded(true)

	m.srcTreeNode = items.AddChild(m.projectTreeNode, "src")
	m.srcTreeNode.SetImageIndex(imageComponents.ImageIndex("folder.png"))
	m.srcTreeNode.SetSelectedIndex(imageComponents.ImageIndex("folder.png"))
	m.srcTreeNode.SetExpanded(true)

	return m
}

func (m *ContentLayoutProject) AddComponentNode(parent, component *TDesigningComponent) lcl.ITreeNode {
	if m.tree == nil {
		return nil
	}
	parentNode := m.componentTreeNode
	if parent != nil {
		parentNode = parent.node
	}
	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()
	items := m.tree.Items()
	node := items.AddChild(parentNode, component.TreeName())
	node.SetImageIndex(component.IconIndex())
	node.SetSelectedIndex(component.IconIndex())
	node.SetData(component.instance())
	node.SetExpanded(false)
	return node
}

func ProjectTreeSetProjectName(name string) {
	if MainWindow.contentLayout == nil {
		return
	}
	MainWindow.contentLayout.layoutProject.projectTreeNode.SetText(name)
}

func ProjectTreeClearComponentTreeNode() {
	if MainWindow.contentLayout == nil {
		return
	}
	MainWindow.contentLayout.layoutProject.componentTreeNode.DeleteChildren()
}

func ProjectTreeClearSrcTreeNode() {
	if MainWindow.contentLayout == nil {
		return
	}
	gProjectSrcTree.stop()
	MainWindow.contentLayout.layoutProject.srcTreeNode.DeleteChildren()
	gProjectSrcTree.reset()
}

func ProjectTreeSrcTreeNode() lcl.ITreeNode {
	if MainWindow.contentLayout == nil {
		return nil
	}
	return MainWindow.contentLayout.layoutProject.srcTreeNode
}

func ProjectTree() lcl.ITreeView {
	if MainWindow.contentLayout == nil {
		return nil
	}
	return MainWindow.contentLayout.layoutProject.tree
}

func ProjectTreeBeginUpdate() {
	if MainWindow.contentLayout == nil {
		return
	}
	MainWindow.contentLayout.layoutProject.tree.BeginUpdate()
}

func ProjectTreeEndUpdate() {
	if MainWindow.contentLayout == nil {
		return
	}
	MainWindow.contentLayout.layoutProject.tree.EndUpdate()
}

func ProjectTreeAddComponentNode(parent, component *TDesigningComponent) lcl.ITreeNode {
	if MainWindow.contentLayout == nil {
		return nil
	}
	return MainWindow.contentLayout.layoutProject.AddComponentNode(parent, component)
}

// 返回当前选中组件树节点
func ProjectTreeSelectNode() lcl.ITreeNode {
	if MainWindow.contentLayout == nil {
		return nil
	}
	return MainWindow.contentLayout.layoutProject.tree.Selected()
}

func ProjectTreeGetNodeAt(X, Y int32) lcl.ITreeNode {
	if MainWindow.contentLayout == nil {
		return nil
	}
	return MainWindow.contentLayout.layoutProject.tree.GetNodeAt(X, Y)
}

func ProjectTreeSetSelected(node lcl.ITreeNode) {
	if MainWindow.contentLayout == nil {
		return
	}
	MainWindow.contentLayout.layoutProject.tree.SetSelected(node)
}

func ProjectTreeComponentMenuPopUp(X, Y int32) {
	if MainWindow.contentLayout == nil {
		return
	}
	MainWindow.contentLayout.layoutProject.componentMenu.treePopupMenu.PopUpWithIntX2(X, Y)
}
