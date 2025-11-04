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
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
)

// 设计 - 组件右键菜单

// 组件菜单
type TComponentMenu struct {
	form             *FormTab
	treePopupMenu    lcl.IPopupMenu // 组件树右键菜单
	zLevel           lcl.IMenuItem  // z 序
	zLevelFront      lcl.IMenuItem  // 移动到最顶层
	zLevelBack       lcl.IMenuItem  // 移动到最底层
	zLevelForwardOne lcl.IMenuItem  // 向前移动一层
	zLevelBackOne    lcl.IMenuItem  // 向后移动一层
	cut              lcl.IMenuItem  // 剪切
	copy             lcl.IMenuItem  // 复制
	paste            lcl.IMenuItem  // 粘贴
	delete           lcl.IMenuItem  // 删除
}

// 返回当前选中组件树节点
func (m *TComponentMenu) ComponentTreeSelectNode() lcl.ITreeNode {
	return m.form.tree.Selected()
}

// 返回当前选中组件
func (m *TComponentMenu) ComponentTreeSelectComponent() *TDesigningComponent {
	node := m.ComponentTreeSelectNode()
	if node != nil && node.IsValid() {
		return m.form.DataToDesigningComponent(node.Data())
	}
	return nil
}

// 移动到最顶层
func (m *TComponentMenu) OnLevelFront(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-移动到最顶层 组件名:", comp.Name())
		comp.WinControl().BringToFront()
		comp.Order(consts.CLevelFront)
	}
}

// 移动到最底层
func (m *TComponentMenu) OnLevelBack(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-移动到最底层 组件名:", comp.Name())
		comp.WinControl().SendToBack()
		comp.Order(consts.CLevelBack)
	}
}

// 向前移动一层
func (m *TComponentMenu) OnLevelForwardOne(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-向前移动一层 组件名:", comp.Name())
		control := comp.WinControl()
		parent := control.Parent()
		if parent == nil || !parent.IsValid() {
			return
		}
		parent.SetControlIndex(control, parent.GetControlIndex(control)+1)
		control.ReAlign()
		parent.ReAlign()
		comp.Order(consts.CLevelForwardOne)
	}
}

// 向后移动一层
func (m *TComponentMenu) OnLevelBackOne(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-向后移动一层 组件名:", comp.Name())
		control := comp.WinControl()
		parent := control.Parent()
		if parent == nil || !parent.IsValid() {
			return
		}
		parent.SetControlIndex(control, parent.GetControlIndex(control)-1)
		control.ReAlign()
		parent.ReAlign()
		comp.Order(consts.CLevelBackOne)
	}
}

// 剪切
func (m *TComponentMenu) OnCut(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-剪切 组件名:", comp.Name())
	}
}

// 复制
func (m *TComponentMenu) OnCopy(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-复制 组件名:", comp.Name())
	}
}

// 粘贴
func (m *TComponentMenu) OnPaste(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-粘贴 组件名:", comp.Name())
	}
}

// 删除
func (m *TComponentMenu) OnDelete(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil && comp.ComponentType != consts.CtForm {
		parent := comp.parent
		lcl.RunOnMainThreadAsync(func(id uint32) {
			// 在删除之前先切换编辑父节组件
			parent.formTab.switchComponentEditing(parent)
			logs.Debug("组件菜单-删除 组件名:", comp.Name())
			comp.Remove() // 删除当前组件
		})
		go triggerUIGeneration(parent) // 更新布局文件
	}
}

// 创建树右键菜单
func (m *FormTab) CreateComponentMenu() {
	if m.componentMenu != nil {
		return
	}
	menu := new(TComponentMenu)
	m.componentMenu = menu
	menu.form = m
	menu.treePopupMenu = lcl.NewPopupMenu(m.tree)
	menu.treePopupMenu.SetImages(imageItem.ImageList100())
	menu.treePopupMenu.SetParent(m.tree)
	menuItems := menu.treePopupMenu.Items()

	// 层级菜单
	zLevel := lcl.NewMenuItem(m.tree)
	zLevel.SetCaption("Z 序")
	menu.zLevel = zLevel
	menuItems.Add(zLevel)

	zLevelFront := lcl.NewMenuItem(m.tree)
	zLevelFront.SetCaption("移动到最顶层")
	zLevelFront.SetImageIndex(imageItem.ImageIndex("order_move_front.png"))
	menu.zLevelFront = zLevelFront
	zLevel.Add(zLevelFront)
	zLevelFront.SetOnClick(menu.OnLevelFront)

	zLevelBack := lcl.NewMenuItem(m.tree)
	zLevelBack.SetCaption("移动到最底层")
	zLevelBack.SetImageIndex(imageItem.ImageIndex("order_move_back.png"))
	menu.zLevelBack = zLevelBack
	zLevel.Add(zLevelBack)
	zLevelBack.SetOnClick(menu.OnLevelBack)

	zLevelForwardOne := lcl.NewMenuItem(m.tree)
	zLevelForwardOne.SetCaption("向前移动一层")
	zLevelForwardOne.SetImageIndex(imageItem.ImageIndex("order_forward_one.png"))
	menu.zLevelForwardOne = zLevelForwardOne
	zLevel.Add(zLevelForwardOne)
	zLevelForwardOne.SetOnClick(menu.OnLevelForwardOne)

	zLevelBackOne := lcl.NewMenuItem(m.tree)
	zLevelBackOne.SetCaption("向后移动一层")
	zLevelBackOne.SetImageIndex(imageItem.ImageIndex("order_back_one.png"))
	menu.zLevelBackOne = zLevelBackOne
	zLevel.Add(zLevelBackOne)
	zLevelBackOne.SetOnClick(menu.OnLevelBackOne)

	line := lcl.NewMenuItem(m.tree)
	line.SetCaption("-")
	menuItems.Add(line)

	cut := lcl.NewMenuItem(m.tree)
	cut.SetCaption("剪切")
	cut.SetImageIndex(imageItem.ImageIndex("item_cut.png"))
	cut.SetOnClick(menu.OnCut)
	cut.SetEnabled(false) // TODO 先禁用未实现
	menu.cut = cut
	menuItems.Add(cut)

	copy := lcl.NewMenuItem(m.tree)
	copy.SetCaption("复制")
	copy.SetImageIndex(imageItem.ImageIndex("item_copy.png"))
	copy.SetOnClick(menu.OnCopy)
	copy.SetEnabled(false) // TODO 先禁用未实现
	menu.copy = copy
	menuItems.Add(copy)

	paste := lcl.NewMenuItem(m.tree)
	paste.SetCaption("粘贴")
	paste.SetImageIndex(imageItem.ImageIndex("item_paste.png"))
	paste.SetOnClick(menu.OnPaste)
	paste.SetEnabled(false) // TODO 先禁用未实现
	menu.paste = paste
	menuItems.Add(paste)

	delete := lcl.NewMenuItem(m.tree)
	delete.SetCaption("删除")
	delete.SetImageIndex(imageItem.ImageIndex("item_delete_selection.png"))
	delete.SetOnClick(menu.OnDelete)
	menu.delete = delete
	menuItems.Add(delete)
}
