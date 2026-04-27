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
	topBox lcl.IPanel
	title  lcl.ILabel
	box    lcl.IPanel
	// 项目栏 tree, 1个项目节点 projectRoot 2个项目管理根节点 componentRoot assetsRoot
	tree lcl.ITreeView
	// 项目根节点
	projectRoot lcl.ITreeNode
	// 组件根节点(所有窗体和组件)
	componentRoot lcl.ITreeNode
	// 资源目录和文件根节点（所有代码和文件）
	assetsRoot lcl.ITreeNode
	// 组件菜单
	componentMenu *TComponentMenu
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

	m.tree = lcl.NewTreeView(owner.projectPanel)
	m.tree.SetAutoExpand(true)
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
	m.tree.SetOnGetSelectedIndex(m.TreeOnGetSelectedIndex)
	m.tree.SetOnMouseDown(m.TreeOnMouseDown)
	m.tree.SetOnContextPopup(m.TreeOnContextPopup)
	//m.CreateComponentMenu()
	m.tree.SetPopupMenu(m.componentMenu.treePopupMenu)
	m.tree.SetParent(m.box)

	items := m.tree.Items()

	m.projectRoot = items.AddChild(nil, "当前项目名")
	m.projectRoot.SetImageIndex(imageComponents.ImageIndex("folder.png"))
	m.projectRoot.SetSelectedIndex(imageComponents.ImageIndex("folder.png"))
	m.projectRoot.SetExpanded(true)

	m.componentRoot = items.AddChild(m.projectRoot, "窗体")
	m.componentRoot.SetImageIndex(imageComponents.ImageIndex("tform.png"))
	m.componentRoot.SetSelectedIndex(imageComponents.ImageIndex("tform.png"))
	m.componentRoot.SetExpanded(true)

	m.assetsRoot = items.AddChild(m.projectRoot, "资源")
	m.assetsRoot.SetImageIndex(imageComponents.ImageIndex("folder.png"))
	m.assetsRoot.SetSelectedIndex(imageComponents.ImageIndex("folder.png"))
	m.assetsRoot.SetExpanded(false)

	//m.assetsRoot.DeleteChildren()

	return m
}

func (m *ContentLayoutProject) AddComponentNode(parent, component *TDesigningComponent) lcl.ITreeNode {
	if m.tree == nil {
		return nil
	}
	parentNode := m.componentRoot
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
	node.SetExpanded(true)
	return node
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
