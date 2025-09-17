package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 对象查看器

var (
	componentTreeHeight int32 = 150
)

type Inspector struct {
	componentTreeFilter lcl.IListFilterEdit
	componentTree       lcl.ITreeView           // 组件树
	splitter            lcl.ISplitter           // 分割线
	componentProperty   lcl.ILazVirtualDrawTree // 组件属性
}

func (m *BottomBox) createInspector() *Inspector {
	ins := new(Inspector)

	// 分隔线
	ins.splitter = lcl.NewSplitter(m.leftBox)
	ins.splitter.SetParent(m.leftBox)
	ins.splitter.SetAlign(types.AlTop)
	ins.splitter.SetWidth(3)

	// 组件树
	ins.componentTree = lcl.NewTreeView(m.leftBox)
	ins.componentTree.SetParent(m.leftBox)
	ins.componentTree.SetWidth(m.leftBox.Width())
	ins.componentTree.SetHeight(componentTreeHeight)
	ins.componentTree.SetAlign(types.AlTop)

	// 组件树搜索过滤
	ins.componentTreeFilter = lcl.NewListFilterEdit(m.leftBox)
	ins.componentTreeFilter.SetParent(m.leftBox)
	ins.componentTreeFilter.SetAlign(types.AlTop)

	// 组件属性
	return ins
}
