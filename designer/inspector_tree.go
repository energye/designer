package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 设计 - 组件树

type InspectorComponentTree struct {
	treeBox    lcl.IPanel          // 组件树盒子
	treeFilter lcl.ITreeFilterEdit // 组件树过滤框
	tree       lcl.ITreeView       // 组件树
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

	m.tree = lcl.NewTreeView(m.treeBox)
	m.tree.SetParent(m.treeBox)
	m.tree.SetTop(35)
	m.tree.SetWidth(leftBoxWidth)
	m.tree.SetHeight(componentTreeHeight - m.tree.Top())
	m.tree.SetAlign(types.AlCustom)
	m.tree.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
}
