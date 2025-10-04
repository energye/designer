package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
	"unsafe"
)

// 设计 - 组件树

var (
	gTreeId int // 维护全局树数据id
)

// 获取下一个树数据ID
func nextTreeDataId() (id int) {
	id = gTreeId
	gTreeId++
	return
}

// 查看器组件树
type InspectorComponentTree struct {
	treeBox    lcl.IPanel                  // 组件树盒子
	treeFilter lcl.ITreeFilterEdit         // 组件树过滤框
	tree       lcl.ITreeView               // 组件树
	images     lcl.IImageList              // 树图标
	root       *DesigningComponent         // 根节点 form 窗体
	nodeData   map[int]*DesigningComponent // 组件树节点数据, key: id
}

// 树节点数据
//type TreeNodeData struct {
//	owner     *InspectorComponentTree // 所属组件树
//	parent    *TreeNodeData           // 所属父节点
//	child     []*TreeNodeData         // 拥有的子节点列表
//	id        int                     // id 标识
//	iconIndex int32                   // 图标
//	node      lcl.ITreeNode           // 节点对象
//	component *DesigningComponent     // 设计的组件
//}

func (m *DesigningComponent) instance() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *InspectorComponentTree) DataToTreeNodeData(dataPtr uintptr) *DesigningComponent {
	data := (*DesigningComponent)(unsafe.Pointer(dataPtr))
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
		var (
			width, height int32 = 20, 20
			images        []string
		)
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
		m.images = LoadImageList(m.treeBox, images, width, height)
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
	//root := m.AddTreeNodeItem(nil, "Form1: TForm", -1)
	//s1 := root.AddChild("Test", 1)
	//s1.AddChild("Tes1t", 2)
	//s1.AddChild("Test2", 3)
}

// 删除当前节点
func (m *DesigningComponent) Remove() {
	//owner:=m.owner
	//m.owner=nil
}

// 向当前组件节点添加子组件节点
func (m *DesigningComponent) AddChild(child *DesigningComponent, name string, iconIndex int32) {
	inspector.componentTree.AddComponentNode(m, child, name, iconIndex)
}

// 添加窗体表单根节点
func (m *InspectorComponentTree) AddFormNode(node *DesigningComponent, name string, iconIndex int32) {
	if node == nil {
		logs.Error("添加窗体表单节点失败, 窗体表单节点为空")
		return
	} else if m.root != nil {
		logs.Error("添加窗体表单节点失败, 已有窗体表单节点")
		return
	}
	// 窗体 根节点
	if node.componentType == CtForm {
		m.tree.BeginUpdate()
		defer m.tree.EndUpdate()
		items := m.tree.Items()
		node.id = nextTreeDataId()
		node.iconIndex = m.images.Count() - 1
		m.nodeData[node.id] = node
		newNode := items.AddChild(nil, name)
		newNode.SetImageIndex(node.iconIndex)    // 显示图标索引
		newNode.SetSelectedIndex(node.iconIndex) // 选中图标索引
		newNode.SetSelected(true)
		newNode.SetData(node.instance())
		node.node = newNode
		m.root = node
	} else {
		logs.Error("添加窗体表单节点失败, 当前节点非窗体表单节点")
		return
	}
}

// 添加组件节点
func (m *InspectorComponentTree) AddComponentNode(parent, child *DesigningComponent, name string, iconIndex int32) {
	if parent == nil {
		logs.Error("添加组件节点失败, 父节点为空")
		return
	} else if child == nil {
		logs.Error("添加组件节点失败, 子节点为空")
		return
	}
	if child.componentType == CtOther {
		m.tree.BeginUpdate()
		defer m.tree.EndUpdate()
		items := m.tree.Items()
		// 控件 子节点
		child.id = nextTreeDataId()
		child.iconIndex = iconIndex
		child.parent = parent
		m.nodeData[child.id] = child
		node := items.AddChild(parent.node, name)
		child.node = node
		node.SetImageIndex(child.iconIndex)    // 显示图标索引
		node.SetSelectedIndex(child.iconIndex) // 选中图标索引
		node.SetSelected(true)
		node.SetData(child.instance())
		// 添加到子节点
		parent.child = append(parent.child, child)
	} else {
		logs.Error("添加组件节点失败, 子节点非组件节点")
	}
}

// 组件树选择事件
func (m *InspectorComponentTree) TreeOnGetSelectedIndex(sender lcl.IObject, node lcl.ITreeNode) {
	dataPtr := node.Data()
	data := m.DataToTreeNodeData(dataPtr)
	//node.SetSelectedIndex(node.ImageIndex())

	logs.Info("Inspector-component-tree OnGetSelectedIndex name:", node.Text(), "id:", data.id)
}
