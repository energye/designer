package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
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
	root       *TreeNodeData         // 根 form 窗体
	nodeData   map[int]*TreeNodeData // 组件树节点数据
}

type TreeNodeData struct {
	owner     *InspectorComponentTree
	id        int
	iconIndex int32
	node      lcl.ITreeNode
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

	// 组件树图标 20x20
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
		// 最后一个图标是 form 窗体图标
		images = append(images, "components/form.png")
		// 加载所有图标
		m.images = LoadImageList(m.treeBox, images, 20, 20)
	}

	// 创建组件树
	m.tree = lcl.NewTreeView(m.treeBox)
	m.tree.SetParent(m.treeBox)
	m.tree.SetAutoExpand(true)
	m.tree.SetTop(35)
	m.tree.SetWidth(leftBoxWidth)
	m.tree.SetHeight(componentTreeHeight - m.tree.Top())
	m.tree.SetReadOnly(true)
	//m.tree.SetMultiSelect(true) // 多选控制
	m.tree.SetAlign(types.AlCustom)
	m.tree.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
	m.tree.SetImages(m.images)
	m.tree.SetOnGetSelectedIndex(m.TreeOnGetSelectedIndex)

	// 测试
	root := m.AddTreeNodeItem(nil, "Form1: TForm", -1)
	s1 := root.AddChild("Test", 1)
	s1.AddChild("Tes1t", 2)
	s1.AddChild("Test2", 3)
}

// 向当前组件节点添加子组件节点
func (m *TreeNodeData) AddChild(name string, iconIndex int32) *TreeNodeData {
	return m.owner.AddTreeNodeItem(m, name, iconIndex)
}

// 删除当前节点
func (m *TreeNodeData) Remove() {
	//owner:=m.owner
	//m.owner=nil
}

// 添加一个组件到节点
func (m *InspectorComponentTree) AddTreeNodeItem(parent *TreeNodeData, name string, iconIndex int32) *TreeNodeData {
	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()
	items := m.tree.Items()
	if m.root == nil && parent == nil { // 窗体 根菜单
		data := &TreeNodeData{owner: m, id: nextTreeDataId(), iconIndex: m.images.Count() - 1}
		m.nodeData[data.id] = data
		node := items.AddChild(nil, name)
		data.node = node
		node.SetImageIndex(data.iconIndex)    // 显示图标索引
		node.SetSelectedIndex(data.iconIndex) // 选中图标索引
		node.SetSelected(true)
		node.SetData(data.instance())
		m.root = data
		return data
	} else { // 控件 子节点
		data := &TreeNodeData{owner: m, id: nextTreeDataId(), iconIndex: iconIndex}
		m.nodeData[data.id] = data
		node := items.AddChild(parent.node, name)
		data.node = node
		node.SetImageIndex(data.iconIndex)    // 显示图标索引
		node.SetSelectedIndex(data.iconIndex) // 选中图标索引
		node.SetSelected(true)
		node.SetData(data.instance())
		return data
	}
}

// 组件树选择事件
func (m *InspectorComponentTree) TreeOnGetSelectedIndex(sender lcl.IObject, node lcl.ITreeNode) {
	dataPtr := node.Data()
	data := m.DataToTreeNodeData(dataPtr)
	//node.SetSelectedIndex(node.ImageIndex())

	log.Println("Inspector-component-tree OnGetSelectedIndex name:", node.Text(), "id:", data.id)
}
