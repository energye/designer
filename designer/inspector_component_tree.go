package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"unsafe"
)

// 设计 - 组件树

var (
	gTreeId int // 维护组件树全局数据id
)

func init() {
}

// 获取下一个树数据ID
func nextTreeDataId() (id int) {
	id = gTreeId
	gTreeId++
	return
}

// 查看器组件树
type InspectorComponentTree struct {
	treeBox      lcl.IPanel          // 组件树盒子
	treeFilter   lcl.ITreeFilterEdit // 组件树过滤框
	componentBox lcl.IPanel          // 组件盒子
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

	m.componentBox = lcl.NewPanel(m.treeBox)
	m.componentBox.SetParent(m.treeBox)
	m.componentBox.SetTop(35)
	m.componentBox.SetWidth(leftBoxWidth)
	m.componentBox.SetHeight(componentTreeHeight - m.componentBox.Top())
	m.componentBox.SetBevelOuter(types.BvNone)
	m.componentBox.SetDoubleBuffered(true)
	m.componentBox.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
}

// 清除组件树数据
//func (m *InspectorComponentTree) Clear() {
//m.tree.Items().Clear()
//m.root = nil
//m.nodeData = make(map[int]*DesigningComponent)
//}

// FormTab

// 创建树右键菜单
func (m *FormTab) TreePopupMenu() lcl.IPopupMenu {
	m.treePopupMenu = lcl.NewPopupMenu(m.tree)
	m.treePopupMenu.SetImages(imageItem.ImageList100())
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

func (m *FormTab) TreeOnContextPopup(sender lcl.IObject, mousePos types.TPoint, handled *bool) {
	logs.Debug("TreeOnContextPopup pos:", mousePos)
}

func (m *FormTab) TreeOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("TreeOnMouseDown x,y:", X, Y)
	if button == types.MbRight {
		selectNode := m.tree.GetNodeAt(X, Y)
		if selectNode.IsValid() {
			m.tree.SetSelected(selectNode)
		}
	}
}

// 数据指针转设计组件
func (m *FormTab) DataToDesigningComponent(data uintptr) *DesigningComponent {
	dc := (*DesigningComponent)(unsafe.Pointer(data))
	return dc
}

// 组件树选择事件
func (m *FormTab) TreeOnGetSelectedIndex(sender lcl.IObject, node lcl.ITreeNode) {
	data := node.Data()
	component := m.DataToDesigningComponent(data)
	if component != nil {
		component.ownerFormTab.hideAllDrag() // 隐藏所有 drag
		component.drag.Show()                // 显示当前设计组件 drag
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			component.LoadPropertyToInspector()
		})
	}
	logs.Info("Inspector-component-tree OnGetSelectedIndex name:", node.Text(), "id:", component.id)
}

// 取消选中所有节点
//func (m *InspectorComponentTree) UnSelectedAll() {
//	for _, node := range m.nodeData {
//		node.node.SetSelected(false)
//	}
//}
