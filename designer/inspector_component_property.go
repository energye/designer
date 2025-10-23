package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 设计 - 组件属性

type InspectorComponentProperty struct {
	box    lcl.IPanel          // 组件属性盒子
	filter lcl.ITreeFilterEdit // 组件属性过滤框
	//page         lcl.IPageControl    // 属性和事件页
	//pageProperty lcl.ITabSheet       // 属性页
	//pageEvent    lcl.ITabSheet       // 事件页
	//propertyTree     lcl.ILazVirtualStringTree // 组件属性
	//eventTree        lcl.ILazVirtualStringTree // 组件事件
	currentComponent *DesigningComponent // 当前正在设计的组件
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
		//m.page = lcl.NewPageControl(m.box)
		//m.page.SetParent(m.box)
		//m.page.SetTabStop(true)
		//m.page.SetTop(32)
		//m.page.SetWidth(leftBoxWidth)
		//m.page.SetHeight(m.box.Height() - m.page.Top())
		//m.page.SetAlign(types.AlCustom)
		//m.page.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
		//
		//m.pageProperty = lcl.NewTabSheet(m.page)
		//m.pageProperty.SetParent(m.page)
		//m.pageProperty.SetCaption("  属性  ")
		//m.pageProperty.SetAlign(types.AlClient)
		//m.pageProperty.SetBorderWidth(0)
		//
		//m.pageEvent = lcl.NewTabSheet(m.page)
		//m.pageEvent.SetParent(m.page)
		//m.pageEvent.SetCaption("  事件  ")
		//m.pageEvent.SetAlign(types.AlClient)
	}
	// 组件的属性列表和事件列表"树"
	{
		// 组件属性树列表
		//m.propertyTree = lcl.NewLazVirtualStringTree(m.pageProperty)
		//m.propertyTree.SetParent(m.pageProperty)
		//vstConfig(m.propertyTree)

		// 组件事件树列表
		//m.eventTree = lcl.NewLazVirtualStringTree(m.pageEvent)
		//m.eventTree.SetParent(m.pageEvent)
		//vstConfig(m.eventTree)

	}
	// 初始化组件属性事件
	//m.initComponentPropertyTreeEvent()

}

func vstConfig(tree lcl.ILazVirtualStringTree) {
	tree.SetBorderStyleToBorderStyle(types.BsNone)
	tree.SetAlign(types.AlClient)
	tree.SetLineStyle(types.LsSolid)
	tree.SetDefaultNodeHeight(28)
	tree.SetIndent(8)

	// options
	propTreeOptions := tree.TreeOptions()
	propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Exclude(types.ToShowTreeLines))
	propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Include(types.ToShowVertGridLines, types.ToShowHorzGridLines))
	propTreeOptions.SetSelectionOptions(propTreeOptions.SelectionOptions().Include(types.ToFullRowSelect, types.ToLevelSelectConstraint))
	propTreeOptions.SetMiscOptions(propTreeOptions.MiscOptions().Include(types.ToEditable, types.ToEditOnClick, types.ToEditOnDblClick))

	// 颜色
	propColors := tree.Colors()
	propColors.SetFocusedSelectionColor(colors.RGBToColor(43, 169, 241))
	propColors.SetUnfocusedSelectionColor(colors.RGBToColor(43, 169, 241))

	// header
	header := tree.Header()
	header.SetOptions(header.Options().Include(types.HoVisible, types.HoAutoSpring)) //types.HoAutoResize
	header.Font().SetStyle(header.Font().Style().Include(types.FsBold))
	header.Font().SetColor(colors.ClGray)
	columns := header.Columns()
	columns.Clear()
	propNameCol := columns.AddToVirtualTreeColumn()
	propNameCol.SetText("名")
	propNameCol.SetWidth(125)
	propNameCol.SetAlignment(types.TaLeftJustify)
	//propNameCol.SetOptions(propNameCol.Options().Include(types.CoDisableAnimatedResize))

	propValueCol := columns.AddToVirtualTreeColumn()
	propValueCol.SetText("值")
	//propValueCol.SetWidth(leftBoxWidth - 125)
	propValueCol.SetAlignment(types.TaLeftJustify)
	propValueCol.SetOptions(propValueCol.Options().Include(types.CoAutoSpring))
}
