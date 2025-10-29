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

// 清除组件树数据
//func (m *InspectorComponentTree) Clear() {
//m.tree.Items().Clear()
//m.root = nil
//m.nodeData = make(map[int]*TDesigningComponent)
//}

// FormTab

// 创建树右键菜单
func (m *FormTab) TreePopupMenu() lcl.IPopupMenu {
	m.treePopupMenu = lcl.NewPopupMenu(m.tree)
	m.treePopupMenu.SetImages(imageItem.ImageList100())

	// 层级菜单
	zLevel := lcl.NewMenuItem(m.tree)
	zLevel.SetCaption("Z 序")
	m.treePopupMenu.Items().Add(zLevel)

	zLevelFront := lcl.NewMenuItem(m.tree)
	zLevelFront.SetCaption("移动到最顶层")
	zLevelFront.SetImageIndex(imageItem.ImageIndex("order_move_front.png"))
	zLevel.Add(zLevelFront)
	zLevelFront.SetOnClick(func(sender lcl.IObject) {

	})

	zLevelBack := lcl.NewMenuItem(m.tree)
	zLevelBack.SetCaption("移动到最底层")
	zLevelBack.SetImageIndex(imageItem.ImageIndex("order_move_back.png"))
	zLevel.Add(zLevelBack)
	zLevelBack.SetOnClick(func(sender lcl.IObject) {

	})

	zLevelForwardOne := lcl.NewMenuItem(m.tree)
	zLevelForwardOne.SetCaption("向前移动一层")
	zLevelForwardOne.SetImageIndex(imageItem.ImageIndex("order_forward_one.png"))
	zLevel.Add(zLevelForwardOne)
	zLevelForwardOne.SetOnClick(func(sender lcl.IObject) {

	})

	zLevelBackOne := lcl.NewMenuItem(m.tree)
	zLevelBackOne.SetCaption("向后移动一层")
	zLevelBackOne.SetImageIndex(imageItem.ImageIndex("order_back_one.png"))
	zLevel.Add(zLevelBackOne)
	zLevelBackOne.SetOnClick(func(sender lcl.IObject) {

	})

	line := lcl.NewMenuItem(m.tree)
	line.SetCaption("-")
	m.treePopupMenu.Items().Add(line)

	cut := lcl.NewMenuItem(m.tree)
	cut.SetCaption("剪切")
	cut.SetImageIndex(imageItem.ImageIndex("item_cut.png"))
	cut.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(cut)

	copy := lcl.NewMenuItem(m.tree)
	copy.SetCaption("复制")
	copy.SetImageIndex(imageItem.ImageIndex("item_copy.png"))
	copy.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(copy)

	paste := lcl.NewMenuItem(m.tree)
	paste.SetCaption("粘贴")
	paste.SetImageIndex(imageItem.ImageIndex("item_paste.png"))
	paste.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(paste)

	delete := lcl.NewMenuItem(m.tree)
	delete.SetCaption("删除")
	delete.SetImageIndex(imageItem.ImageIndex("item_delete_selection.png"))
	delete.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(delete)

	m.treePopupMenu.SetParent(m.tree)
	return m.treePopupMenu
}
