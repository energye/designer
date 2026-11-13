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
	treePopupMenu    lcl.IPopupMenu // 组件树右键菜单
	zLevelFront      lcl.IMenuItem  // 移动到最顶层
	zLevelBack       lcl.IMenuItem  // 移动到最底层
	zLevelForwardOne lcl.IMenuItem  // 向前移动一层
	zLevelBackOne    lcl.IMenuItem  // 向后移动一层
	cut              lcl.IMenuItem  // 剪切
	copy             lcl.IMenuItem  // 复制
	paste            lcl.IMenuItem  // 粘贴
	delete           lcl.IMenuItem  // 删除

	alignMenu       lcl.IMenuItem   // 对齐子菜单
	alignItems      []lcl.IMenuItem // 对齐菜单项
	sameSizeMenu    lcl.IMenuItem   // 同尺寸子菜单
	sameSizeItems   []lcl.IMenuItem // 同尺寸菜单项
	distributeMenu  lcl.IMenuItem   // 分布子菜单
	distributeItems []lcl.IMenuItem // 分布菜单项
}

func (m *TComponentMenu) Free() {
	m.zLevelFront.SetOnClick(nil)
	m.zLevelBack.SetOnClick(nil)
	m.zLevelForwardOne.SetOnClick(nil)
	m.zLevelBackOne.SetOnClick(nil)
	m.cut.SetOnClick(nil)
	m.copy.SetOnClick(nil)
	m.paste.SetOnClick(nil)
	m.delete.SetOnClick(nil)
	if m.treePopupMenu != nil {
		m.treePopupMenu.SetOnPopup(nil)
	}
	for _, item := range m.alignItems {
		if item != nil {
			item.SetOnClick(nil)
		}
	}
	for _, item := range m.sameSizeItems {
		if item != nil {
			item.SetOnClick(nil)
		}
	}
	for _, item := range m.distributeItems {
		if item != nil {
			item.SetOnClick(nil)
		}
	}
}

// 返回当前选中组件
func (m *TComponentMenu) ComponentTreeSelectComponent() *TDesigningComponent {
	node := ProjectTreeSelectNode()
	if node != nil && node.IsValid() {
		return TreeNodeDataToDesigningComponent(node.Data())
	}
	return nil
}

// selectedFormTab 返回右键选中组件所属的设计页。
func (m *TComponentMenu) selectedFormTab() *FormTab {
	comp := m.ComponentTreeSelectComponent()
	if comp == nil || comp.FormTab == nil {
		return nil
	}
	return comp.FormTab
}

func setMenuItemsEnabled(items []lcl.IMenuItem, enabled bool) {
	for _, item := range items {
		if item != nil {
			item.SetEnabled(enabled)
		}
	}
}

// OnPopup 弹出前按当前选中数量刷新对齐菜单可用状态。
func (m *TComponentMenu) OnPopup(sender lcl.IObject) {
	count := 0
	if formTab := m.selectedFormTab(); formTab != nil {
		count = len(formTab.alignableSelectedComponents())
	}
	alignEnabled := count >= 2
	distEnabled := count >= 3
	if m.alignMenu != nil {
		m.alignMenu.SetEnabled(alignEnabled)
	}
	if m.sameSizeMenu != nil {
		m.sameSizeMenu.SetEnabled(alignEnabled)
	}
	if m.distributeMenu != nil {
		m.distributeMenu.SetEnabled(distEnabled)
	}
	setMenuItemsEnabled(m.alignItems, alignEnabled)
	setMenuItemsEnabled(m.sameSizeItems, alignEnabled)
	setMenuItemsEnabled(m.distributeItems, distEnabled)
}

func (m *TComponentMenu) alignSelected(kind designerAlignKind) {
	if formTab := m.selectedFormTab(); formTab != nil {
		formTab.AlignSelected(kind)
	}
}

func (m *TComponentMenu) sameSizeSelected(kind designerSameSizeKind) {
	if formTab := m.selectedFormTab(); formTab != nil {
		formTab.SameSizeSelected(kind)
	}
}

func (m *TComponentMenu) distributeSelected(horizontal bool) {
	if formTab := m.selectedFormTab(); formTab != nil {
		formTab.DistributeSelected(horizontal)
	}
}

// 移动到最顶层
func (m *TComponentMenu) OnLevelFront(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-移动到最顶层 组件名:", comp.Name())
		control := comp.Control()
		if control == nil {
			return
		}
		control.BringToFront()
		comp.Order(consts.CLevelFront)
	}
}

// 移动到最底层
func (m *TComponentMenu) OnLevelBack(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-移动到最底层 组件名:", comp.Name())
		control := comp.Control()
		if control == nil {
			return
		}
		control.SendToBack()
		comp.Order(consts.CLevelBack)
	}
}

// 向前移动一层
func (m *TComponentMenu) OnLevelForwardOne(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-向前移动一层 组件名:", comp.Name())
		control := comp.Control()
		if control == nil {
			return
		}
		parent := control.Parent()
		if parent == nil || !parent.IsValid() {
			return
		}
		parent.SetControlIndex(control, parent.GetControlIndex(control)+1)
		comp.Order(consts.CLevelForwardOne)
	}
}

// 向后移动一层
func (m *TComponentMenu) OnLevelBackOne(sender lcl.IObject) {
	comp := m.ComponentTreeSelectComponent()
	if comp != nil {
		logs.Debug("组件菜单-向后移动一层 组件名:", comp.Name())
		control := comp.Control()
		if control == nil {
			return
		}
		parent := control.Parent()
		if parent == nil || !parent.IsValid() {
			return
		}
		parent.SetControlIndex(control, parent.GetControlIndex(control)-1)
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
		logs.Debug("组件菜单-删除 组件名:", comp.Name())
		DeleteDesignComponent(comp)
	}
}

// 初始化组件树右键菜单
func initComponentMenu(owner lcl.IComponent) *TComponentMenu {
	menu := new(TComponentMenu)
	menu.treePopupMenu = lcl.NewPopupMenu(owner)
	menu.treePopupMenu.SetImages(imageItem.ImageList100())
	menu.treePopupMenu.SetParent(owner)
	menuItems := menu.treePopupMenu.Items()

	zLevelFront := lcl.NewMenuItem(owner)
	zLevelFront.SetCaption("移动到最顶层")
	zLevelFront.SetImageIndex(imageItem.ImageIndex("order_move_front.png"))
	menu.zLevelFront = zLevelFront
	menuItems.Add(zLevelFront)
	zLevelFront.SetOnClick(menu.OnLevelFront)

	zLevelBack := lcl.NewMenuItem(owner)
	zLevelBack.SetCaption("移动到最底层")
	zLevelBack.SetImageIndex(imageItem.ImageIndex("order_move_back.png"))
	menu.zLevelBack = zLevelBack
	menuItems.Add(zLevelBack)
	zLevelBack.SetOnClick(menu.OnLevelBack)

	zLevelForwardOne := lcl.NewMenuItem(owner)
	zLevelForwardOne.SetCaption("向前移动一层")
	zLevelForwardOne.SetImageIndex(imageItem.ImageIndex("order_forward_one.png"))
	menu.zLevelForwardOne = zLevelForwardOne
	menuItems.Add(zLevelForwardOne)
	zLevelForwardOne.SetOnClick(menu.OnLevelForwardOne)

	zLevelBackOne := lcl.NewMenuItem(owner)
	zLevelBackOne.SetCaption("向后移动一层")
	zLevelBackOne.SetImageIndex(imageItem.ImageIndex("order_back_one.png"))
	menu.zLevelBackOne = zLevelBackOne
	menuItems.Add(zLevelBackOne)
	zLevelBackOne.SetOnClick(menu.OnLevelBackOne)

	line := lcl.NewMenuItem(owner)
	line.SetCaption("-")
	menuItems.Add(line)

	newSubMenu := func(caption string) lcl.IMenuItem {
		item := lcl.NewMenuItem(owner)
		item.SetCaption(caption)
		menuItems.Add(item)
		return item
	}
	newSubItem := func(parent lcl.IMenuItem, caption string, onClick lcl.TNotifyEvent) lcl.IMenuItem {
		item := lcl.NewMenuItem(owner)
		item.SetCaption(caption)
		item.SetOnClick(onClick)
		parent.Add(item)
		return item
	}

	menu.alignMenu = newSubMenu("对齐")
	menu.alignItems = []lcl.IMenuItem{
		newSubItem(menu.alignMenu, "左对齐", func(lcl.IObject) { menu.alignSelected(designerAlignLeft) }),
		newSubItem(menu.alignMenu, "水平居中", func(lcl.IObject) { menu.alignSelected(designerAlignHCenter) }),
		newSubItem(menu.alignMenu, "右对齐", func(lcl.IObject) { menu.alignSelected(designerAlignRight) }),
		newSubItem(menu.alignMenu, "上对齐", func(lcl.IObject) { menu.alignSelected(designerAlignTop) }),
		newSubItem(menu.alignMenu, "垂直居中", func(lcl.IObject) { menu.alignSelected(designerAlignVCenter) }),
		newSubItem(menu.alignMenu, "下对齐", func(lcl.IObject) { menu.alignSelected(designerAlignBottom) }),
	}

	menu.sameSizeMenu = newSubMenu("同尺寸")
	menu.sameSizeItems = []lcl.IMenuItem{
		newSubItem(menu.sameSizeMenu, "同宽", func(lcl.IObject) { menu.sameSizeSelected(designerSameWidth) }),
		newSubItem(menu.sameSizeMenu, "同高", func(lcl.IObject) { menu.sameSizeSelected(designerSameHeight) }),
		newSubItem(menu.sameSizeMenu, "同尺寸", func(lcl.IObject) { menu.sameSizeSelected(designerSameSize) }),
	}

	menu.distributeMenu = newSubMenu("分布")
	menu.distributeItems = []lcl.IMenuItem{
		newSubItem(menu.distributeMenu, "水平分布", func(lcl.IObject) { menu.distributeSelected(true) }),
		newSubItem(menu.distributeMenu, "垂直分布", func(lcl.IObject) { menu.distributeSelected(false) }),
	}

	lineAlign := lcl.NewMenuItem(owner)
	lineAlign.SetCaption("-")
	menuItems.Add(lineAlign)

	cut := lcl.NewMenuItem(owner)
	cut.SetCaption("剪切")
	cut.SetImageIndex(imageItem.ImageIndex("item_cut.png"))
	cut.SetOnClick(menu.OnCut)
	cut.SetEnabled(false) // TODO 先禁用未实现
	menu.cut = cut
	menuItems.Add(cut)

	copy := lcl.NewMenuItem(owner)
	copy.SetCaption("复制")
	copy.SetImageIndex(imageItem.ImageIndex("item_copy.png"))
	copy.SetOnClick(menu.OnCopy)
	copy.SetEnabled(false) // TODO 先禁用未实现
	menu.copy = copy
	menuItems.Add(copy)

	paste := lcl.NewMenuItem(owner)
	paste.SetCaption("粘贴")
	paste.SetImageIndex(imageItem.ImageIndex("item_paste.png"))
	paste.SetOnClick(menu.OnPaste)
	paste.SetEnabled(false) // TODO 先禁用未实现
	menu.paste = paste
	menuItems.Add(paste)

	delete := lcl.NewMenuItem(owner)
	delete.SetCaption("删除")
	delete.SetImageIndex(imageItem.ImageIndex("item_delete_selection.png"))
	delete.SetOnClick(menu.OnDelete)
	menu.delete = delete
	menuItems.Add(delete)

	menu.treePopupMenu.SetOnPopup(menu.OnPopup)
	menu.OnPopup(nil)
	return menu
}
