package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
	"unsafe"
)

// 设计 - 组件树

var (
	gTreeId        int              // 维护组件树全局数据id
	gTreeImageList map[string]int32 // 组件树树节点图标索引 key: 组件类名 value: 索引
)

func init() {
	gTreeImageList = make(map[string]int32)
}

// 返回查看器组件树节点使用的图标
func CompTreeIcon(name string) int32 {
	return gTreeImageList[name]
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
	images       lcl.IImageList      // 组件树树图标
	images150    lcl.IImageList      // 组件树树图标
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

	// 组件树图标 20x20
	{
		var (
			width, height     int32 = 20, 20
			images, images150 []string
		)
		var eachTabName = func(tab config.Tab) {
			for _, name := range tab.Component {
				images = append(images, fmt.Sprintf("components/%v.png", strings.ToLower(name)))
				images150 = append(images150, fmt.Sprintf("components/%v_150.png", strings.ToLower(name)))
				gTreeImageList[name] = int32(len(images) - 1)
			}
		}
		eachTabName(config.Config.ComponentTabs.Standard)
		eachTabName(config.Config.ComponentTabs.Additional)
		eachTabName(config.Config.ComponentTabs.Common)
		eachTabName(config.Config.ComponentTabs.Dialogs)
		eachTabName(config.Config.ComponentTabs.Misc)
		eachTabName(config.Config.ComponentTabs.System)
		eachTabName(config.Config.ComponentTabs.LazControl)
		eachTabName(config.Config.ComponentTabs.WebView)
		// 最后一个图标是 TForm 窗体图标
		images = append(images, "components/form.png")
		gTreeImageList["TForm"] = int32(len(images) - 1)
		gTreeImageList["TEngForm"] = int32(len(images) - 1)
		// 加载所有图标
		m.images = tool.LoadImageList(m.treeBox, images, width, height)
		m.images150 = tool.LoadImageList(m.treeBox, images150, 36, 36)
	}

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
	cut := lcl.NewMenuItem(m.tree)
	cut.SetCaption("剪切")
	cut.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(cut)

	copy := lcl.NewMenuItem(m.tree)
	copy.SetCaption("复制")
	copy.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(copy)

	paste := lcl.NewMenuItem(m.tree)
	paste.SetCaption("粘贴")
	paste.SetOnClick(func(lcl.IObject) {
	})
	m.treePopupMenu.Items().Add(paste)

	delete := lcl.NewMenuItem(m.tree)
	delete.SetCaption("删除")
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
