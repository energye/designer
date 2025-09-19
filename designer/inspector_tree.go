package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
	"unsafe"
)

// 设计 - 组件树

var gTreeId int

func nextTreeDataId() (id int) {
	id = gTreeId
	gTreeId++
	return
}

type InspectorComponentTree struct {
	treeBox    lcl.IPanel            // 组件树盒子
	treeFilter lcl.ITreeFilterEdit   // 组件树过滤框
	tree       lcl.ITreeView         // 组件树
	images     lcl.IImageList        // 树图标
	root       lcl.ITreeNode         // 根 form 窗体
	nodeData   map[int]*TreeNodeData // 组件树节点数据
}

type TreeNodeData struct {
	id        int
	iconIndex int32
}

func (m *TreeNodeData) instance() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *InspectorComponentTree) DataToTreeNodeData(dataPtr uintptr) *TreeNodeData {
	data := (*TreeNodeData)(unsafe.Pointer(dataPtr))
	return data
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

	{
		var images []string
		var eachTabName = func(tab config.Tab) {
			for _, name := range tab.Component {
				images = append(images, fmt.Sprintf("components/%v.png", strings.ToLower(name)))
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
		images = append(images, "components/form.png")
		m.images = LoadImageList(m.treeBox, images, 16, 16)
	}

	m.tree = lcl.NewTreeView(m.treeBox)
	m.tree.SetParent(m.treeBox)
	m.tree.SetTop(35)
	m.tree.SetWidth(leftBoxWidth)
	m.tree.SetHeight(componentTreeHeight - m.tree.Top())
	m.tree.SetReadOnly(true)
	m.tree.SetAlign(types.AlCustom)
	m.tree.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
	m.tree.SetImages(m.images)
	m.tree.SetOnGetSelectedIndex(func(sender lcl.IObject, node lcl.ITreeNode) {
		fmt.Println("SetOnGetSelectedIndex")
		dataPtr := node.Data()
		data := m.DataToTreeNodeData(dataPtr)
		fmt.Println(data)
		node.SetSelectedIndex(data.iconIndex)
	})

	// 测试
	m.AddTreeItem("Form1: TForm")
}

func (m *InspectorComponentTree) AddTreeItem(name string) {
	if m.root == nil {
		data := &TreeNodeData{id: nextTreeDataId(), iconIndex: m.images.Count() - 1}
		m.nodeData[data.id] = data
		m.root = m.tree.Items().Add(nil, name)
		m.root.SetData(data.instance())
		m.root.SetImageIndex(data.iconIndex)
	}
}
