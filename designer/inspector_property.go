package designer

import (
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
)

// 设计 - 组件属性

type InspectorComponentProperty struct {
	box           lcl.IPanel                // 组件属性盒子
	filter        lcl.ITreeFilterEdit       // 组件属性过滤框
	page          lcl.IPageControl          // 属性和事件页
	propertySheet lcl.ITabSheet             // 属性页
	eventSheet    lcl.ITabSheet             // 事件页
	propertyTree  lcl.ILazVirtualStringTree // 组件属性
	eventTree     lcl.ILazVirtualStringTree // 组件事件
}

func (m *InspectorComponentProperty) init(leftBoxWidth int32) {
	componentPropertyTitle := lcl.NewLabel(m.box)
	componentPropertyTitle.SetParent(m.box)
	componentPropertyTitle.SetCaption("属性")
	componentPropertyTitle.Font().SetStyle(types.NewSet(types.FsBold))
	componentPropertyTitle.SetTop(5)
	componentPropertyTitle.SetLeft(5)

	m.filter = lcl.NewTreeFilterEdit(m.box)
	m.filter.SetParent(m.box)
	m.filter.SetTop(2)
	m.filter.SetLeft(30)
	m.filter.SetWidth(leftBoxWidth - m.filter.Left())
	m.filter.SetAlign(types.AlCustom)
	m.filter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	// 选项卡
	{
		m.page = lcl.NewPageControl(m.box)
		m.page.SetParent(m.box)
		m.page.SetTabStop(true)
		m.page.SetTop(32)
		m.page.SetWidth(leftBoxWidth)
		m.page.SetHeight(m.box.Height() - m.page.Top())
		m.page.SetAlign(types.AlCustom)
		m.page.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))

		m.propertySheet = lcl.NewTabSheet(m.page)
		m.propertySheet.SetParent(m.page)
		m.propertySheet.SetCaption("  属性  ")
		m.propertySheet.SetAlign(types.AlClient)

		m.eventSheet = lcl.NewTabSheet(m.page)
		m.eventSheet.SetParent(m.page)
		m.eventSheet.SetCaption("  事件  ")
		m.eventSheet.SetAlign(types.AlClient)
	}
	// 组件的属性列表和事件列表"树"
	{
		// 组件属性树列表
		m.propertyTree = lcl.NewLazVirtualStringTree(m.propertySheet)
		m.propertyTree.SetParent(m.propertySheet)
		m.propertyTree.SetBorderStyleToBorderStyle(types.BsNone)
		m.propertyTree.SetAlign(types.AlClient)
		// 组件事件树列表
		m.eventTree = lcl.NewLazVirtualStringTree(m.eventSheet)
		m.eventTree.SetParent(m.eventSheet)
		m.eventTree.SetBorderStyleToBorderStyle(types.BsNone)
		m.eventTree.SetAlign(types.AlClient)
	}
	// 初始化组件树
	m.initComponentTree()

	// 测试
	{

	}
}

// 初始化组件属性树
func (m *InspectorComponentProperty) initComponentTree() {
	m.propertyTree.SetOnCreateEditor(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, outEditLink *lcl.IVTEditLink) {
		log.Println("OnCreateEditor")
		*outEditLink = vtedit.NewStringEditLink()
	})
}
