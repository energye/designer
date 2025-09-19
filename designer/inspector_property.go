package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 设计 - 组件属性

type InspectorComponentProperty struct {
	propertyBox    lcl.IPanel              // 组件属性盒子
	propertyFilter lcl.ITreeFilterEdit     // 组件属性过滤框
	property       lcl.ILazVirtualDrawTree // 组件属性
}

func (m *InspectorComponentProperty) init(leftBoxWidth int32) {
	componentPropertyTitle := lcl.NewLabel(m.propertyBox)
	componentPropertyTitle.SetParent(m.propertyBox)
	componentPropertyTitle.SetCaption("属性")
	componentPropertyTitle.Font().SetStyle(types.NewSet(types.FsBold))
	componentPropertyTitle.SetTop(5)
	componentPropertyTitle.SetLeft(5)

	m.propertyFilter = lcl.NewTreeFilterEdit(m.propertyBox)
	m.propertyFilter.SetParent(m.propertyBox)
	m.propertyFilter.SetTop(2)
	m.propertyFilter.SetLeft(30)
	m.propertyFilter.SetWidth(leftBoxWidth - m.propertyFilter.Left())
	m.propertyFilter.SetAlign(types.AlCustom)
	m.propertyFilter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))

	// 属性树列表
	m.property = lcl.NewLazVirtualDrawTree(m.propertyBox)
	m.property.SetParent(m.propertyBox)
	m.property.SetTop(32)
	m.property.SetWidth(leftBoxWidth)
	m.property.SetHeight(m.propertyBox.Height() - m.property.Top())
	m.property.SetAlign(types.AlCustom)
	m.property.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))

	// 测试
}
