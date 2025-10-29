package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 设计 - 组件树

var (
	gTreeId int // 维护组件树全局数据id
)

// 获取下一个树数据ID
func nextTreeDataId() (id int) {
	id = gTreeId
	gTreeId++
	return
}

// 查看器组件树
type InspectorComponentTree struct {
	treeBox    lcl.IPanel          // 组件树盒子
	treeFilter lcl.ITreeFilterEdit // 组件树过滤框
}

func (m *InspectorComponentTree) init(leftBoxWidth int32) {
	componentTreeTitle := lcl.NewLabel(m.treeBox)
	componentTreeTitle.SetParent(m.treeBox)
	componentTreeTitle.SetCaption("组件")
	componentTreeTitle.Font().SetStyle(types.NewSet(types.FsBold))
	componentTreeTitle.SetTop(8)
	componentTreeTitle.SetLeft(5)

	m.treeFilter = lcl.NewTreeFilterEdit(m.treeBox)
	m.treeFilter.SetParent(m.treeBox)
	m.treeFilter.SetTop(5)
	m.treeFilter.SetLeft(30)
	m.treeFilter.SetWidth(leftBoxWidth - m.treeFilter.Left())
	m.treeFilter.SetAlign(types.AlCustom)
	m.treeFilter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
}

// 组件菜单
type ComponentMenu struct {
	form          *FormTab
	treePopupMenu lcl.IPopupMenu // 组件树右键菜单
}

func (m *ComponentMenu) ComponentTreeSelectNode() lcl.ITreeNode {
	return m.form.tree.Selected()
}
func (m *ComponentMenu) OnLevelFront(sender lcl.IObject) {
	node := m.ComponentTreeSelectNode()
	if node != nil {

	}
}
func (m *ComponentMenu) OnLevelBack(sender lcl.IObject) {

}
func (m *ComponentMenu) OnLevelForwardOne(sender lcl.IObject) {

}
func (m *ComponentMenu) OnLevelBackOne(sender lcl.IObject) {

}
func (m *ComponentMenu) OnCut(sender lcl.IObject) {

}
func (m *ComponentMenu) OnCopy(sender lcl.IObject) {

}
func (m *ComponentMenu) OnPaste(sender lcl.IObject) {

}
func (m *ComponentMenu) OnDelete(sender lcl.IObject) {

}

// 创建树右键菜单
func (m *FormTab) CreateComponentMenu() {
	if m.componentMenu != nil {
		return
	}
	menu := new(ComponentMenu)
	menu.form = m
	menu.treePopupMenu = lcl.NewPopupMenu(m.tree)
	menu.treePopupMenu.SetImages(imageItem.ImageList100())
	menu.treePopupMenu.SetParent(m.tree)
	m.componentMenu = menu
	menuItems := menu.treePopupMenu.Items()

	// 层级菜单
	zLevel := lcl.NewMenuItem(m.tree)
	zLevel.SetCaption("Z 序")
	menuItems.Add(zLevel)

	zLevelFront := lcl.NewMenuItem(m.tree)
	zLevelFront.SetCaption("移动到最顶层")
	zLevelFront.SetImageIndex(imageItem.ImageIndex("order_move_front.png"))
	zLevel.Add(zLevelFront)
	zLevelFront.SetOnClick(menu.OnLevelFront)

	zLevelBack := lcl.NewMenuItem(m.tree)
	zLevelBack.SetCaption("移动到最底层")
	zLevelBack.SetImageIndex(imageItem.ImageIndex("order_move_back.png"))
	zLevel.Add(zLevelBack)
	zLevelBack.SetOnClick(menu.OnLevelBack)

	zLevelForwardOne := lcl.NewMenuItem(m.tree)
	zLevelForwardOne.SetCaption("向前移动一层")
	zLevelForwardOne.SetImageIndex(imageItem.ImageIndex("order_forward_one.png"))
	zLevel.Add(zLevelForwardOne)
	zLevelForwardOne.SetOnClick(menu.OnLevelForwardOne)

	zLevelBackOne := lcl.NewMenuItem(m.tree)
	zLevelBackOne.SetCaption("向后移动一层")
	zLevelBackOne.SetImageIndex(imageItem.ImageIndex("order_back_one.png"))
	zLevel.Add(zLevelBackOne)
	zLevelBackOne.SetOnClick(menu.OnLevelBackOne)

	line := lcl.NewMenuItem(m.tree)
	line.SetCaption("-")
	menuItems.Add(line)

	cut := lcl.NewMenuItem(m.tree)
	cut.SetCaption("剪切")
	cut.SetImageIndex(imageItem.ImageIndex("item_cut.png"))
	cut.SetOnClick(menu.OnCut)
	menuItems.Add(cut)

	copy := lcl.NewMenuItem(m.tree)
	copy.SetCaption("复制")
	copy.SetImageIndex(imageItem.ImageIndex("item_copy.png"))
	copy.SetOnClick(menu.OnCopy)
	menuItems.Add(copy)

	paste := lcl.NewMenuItem(m.tree)
	paste.SetCaption("粘贴")
	paste.SetImageIndex(imageItem.ImageIndex("item_paste.png"))
	paste.SetOnClick(menu.OnPaste)
	menuItems.Add(paste)

	delete := lcl.NewMenuItem(m.tree)
	delete.SetCaption("删除")
	delete.SetImageIndex(imageItem.ImageIndex("item_delete_selection.png"))
	delete.SetOnClick(menu.OnDelete)
	menuItems.Add(delete)

}
