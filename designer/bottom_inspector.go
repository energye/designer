package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 对象查看器

var (
	componentTreeHeight int32 = 150
	inspector           *Inspector
)

// 组件树和对象查看器
type Inspector struct {
	boxSplitter       lcl.ISplitter               // 分割线
	componentTree     *InspectorComponentTree     // 组件树
	componentProperty *InspectorComponentProperty // 组件属性
}

func (m *BottomBox) createInspectorLayout() *Inspector {
	ins := new(Inspector)
	// 面板 对象查看器分隔
	{
		ins.boxSplitter = lcl.NewSplitter(m.leftBox)
		ins.boxSplitter.SetParent(m.leftBox)
		ins.boxSplitter.SetAlign(types.AlTop)

		tree := new(InspectorComponentTree)
		tree.nodeData = make(map[int]*TreeNodeData)
		tree.treeBox = lcl.NewPanel(m.leftBox)
		tree.treeBox.SetParent(m.leftBox)
		tree.treeBox.SetBevelOuter(types.BvNone)
		tree.treeBox.SetDoubleBuffered(true)
		tree.treeBox.SetWidth(m.leftBox.Width())
		tree.treeBox.SetHeight(componentTreeHeight)
		tree.treeBox.Constraints().SetMinWidth(50)
		tree.treeBox.Constraints().SetMinHeight(50)
		tree.treeBox.SetAlign(types.AlTop)
		ins.componentTree = tree

		property := new(InspectorComponentProperty)
		property.propertyBox = lcl.NewPanel(m.leftBox)
		property.propertyBox.SetParent(m.leftBox)
		property.propertyBox.SetBevelOuter(types.BvNone)
		property.propertyBox.SetDoubleBuffered(true)
		property.propertyBox.SetWidth(m.leftBox.Width())
		property.propertyBox.SetHeight(componentTreeHeight)
		property.propertyBox.Constraints().SetMinWidth(50)
		property.propertyBox.Constraints().SetMinHeight(50)
		property.propertyBox.SetAlign(types.AlClient)
		ins.componentProperty = property
		//ins.componentPropertyBox.SetColor(colors.Cl3DDkShadow)
	}
	// 组件树
	{
		ins.componentTree.init(m.leftBox.Width())
	}
	// 组件属性
	{
		ins.componentProperty.init(m.leftBox.Width())
	}
	return ins
}
