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
	treeBox    lcl.IPanel                  // 组件树盒子
	treeFilter lcl.ITreeFilterEdit         // 组件树过滤框
	tree       lcl.ITreeView               // 组件树
	images     lcl.IImageList              // 组件树树图标
	root       *DesigningComponent         // 根节点 form 窗体
	nodeData   map[int]*DesigningComponent // 组件树节点数据, key: id
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

// 清除组件树数据
func (m *InspectorComponentTree) Clear() {
	m.tree.Items().Clear()
	m.root = nil
	m.nodeData = make(map[int]*DesigningComponent)
}

// 数据指针转设计组件
func (m *InspectorComponentTree) DataToTreeNodeData(dataPtr uintptr) *DesigningComponent {
	data := (*DesigningComponent)(unsafe.Pointer(dataPtr))
	return data
}

// 添加窗体表单根节点
func (m *InspectorComponentTree) AddFormNode(form *FormTab) {
	if form == nil {
		logs.Error("添加窗体表单节点失败, 窗体表单节点为空")
		return
	} else if m.root != nil {
		logs.Error("添加窗体表单节点失败, 已有窗体表单节点")
		return
	}
	// 窗体 根节点
	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()
	items := m.tree.Items()
	form.designerBox.id = nextTreeDataId()
	m.nodeData[form.designerBox.id] = form.designerBox
	newNode := items.AddChild(nil, form.form.TreeName())
	newNode.SetImageIndex(form.form.IconIndex())    // 显示图标索引
	newNode.SetSelectedIndex(form.form.IconIndex()) // 选中图标索引
	newNode.SetSelected(true)
	newNode.SetData(form.designerBox.instance())
	form.designerBox.node = newNode
	m.root = form.designerBox
}

// 添加组件节点
func (m *InspectorComponentTree) AddComponentNode(parent, child *DesigningComponent) {
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
		//child.parent = parent
		m.nodeData[child.id] = child
		node := items.AddChild(parent.node, child.TreeName())
		child.node = node
		node.SetImageIndex(child.IconIndex())    // 显示图标索引
		node.SetSelectedIndex(child.IconIndex()) // 选中图标索引
		node.SetSelected(true)
		node.SetData(child.instance())
		//parent.child = append(parent.child, child)
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

// 取消选中所有节点
func (m *InspectorComponentTree) UnSelectedAll() {
	for _, node := range m.nodeData {
		node.node.SetSelected(false)
	}
}
