package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

var (
	defaultSplitterWidth   = int32(2)
	defaultSplitterMinSize = int32(1)
)

type ContentLayout struct {
	box      lcl.IPanel // 主内容布局盒子
	rightBox lcl.IPanel // 右侧盒子

	widgetSplitter lcl.ISplitter // 组件面板分隔器
	widgetPanel    lcl.IPanel    // 组件面板

	projectSplitter lcl.ISplitter // 项目面板分隔器
	projectPanel    lcl.IPanel    // 项目面板

	designerPanel lcl.IPanel // 设计器面板

	inspectorSplitter lcl.ISplitter // 属性检查器面板分隔器
	inspectorPanel    lcl.IPanel    // 属性检查器面板

	consoleLogSplitter lcl.ISplitter // 日志面板分隔器
	consoleLogPanel    lcl.IPanel    // 控制台输出

	contentStatus lcl.IStatusBar
}

// 底部布局 - 左: 组件库, 左中: 项目查看, 中: 中间画布(自适应), 右: 属性, 下: 日志控制
func initBottomBox(owner lcl.IWinControl) *ContentLayout {
	m := &ContentLayout{}

	// 主内容布局盒子
	m.box = lcl.NewPanel(owner)
	//m.box.SetColor(colors.Cl3DDkShadow)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlClient)
	m.box.SetCaption("主内容布局")
	m.box.SetParent(owner)

	// 组件面板分隔器
	m.widgetSplitter = lcl.NewSplitter(owner)
	m.widgetSplitter.SetAlign(types.AlLeft)
	m.widgetSplitter.SetWidth(defaultSplitterWidth)
	m.widgetSplitter.SetMinSize(defaultSplitterMinSize)
	m.widgetSplitter.SetParent(m.box)
	// 组件面板
	m.widgetPanel = lcl.NewPanel(owner)
	//m.widgetPanel.SetColor(wg.LightenColor(colors.ClAqua, 0.3))
	m.widgetPanel.SetBevelOuter(types.BvNone)
	m.widgetPanel.SetBorderStyleToBorderStyle(types.BsSingle)
	m.widgetPanel.SetDoubleBuffered(true)
	m.widgetPanel.SetAlign(types.AlLeft)
	m.widgetPanel.SetCaption("组件库")
	m.widgetPanel.SetWidth(170)    //动态控制
	m.widgetPanel.SetVisible(true) //动态控制
	m.widgetPanel.Constraints().SetMinWidth(30)
	m.widgetPanel.SetParent(m.box)
	//m.widgetPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("widgetPanel.SetOnResize BoundsRect:", m.widgetPanel.BoundsRect())
	//})

	// 右侧盒子
	m.rightBox = lcl.NewPanel(owner)
	//m.rightBox.SetColor(wg.LightenColor(colors.ClAqua, 1.5))
	m.rightBox.SetBevelOuter(types.BvNone)
	m.rightBox.SetDoubleBuffered(true)
	m.rightBox.SetAlign(types.AlClient)
	m.rightBox.SetCaption("右内容")
	m.rightBox.SetHeight(150)
	m.rightBox.SetParent(m.box)

	// 项目面板分隔器
	m.projectSplitter = lcl.NewSplitter(owner)
	m.projectSplitter.SetAlign(types.AlLeft)
	m.projectSplitter.SetWidth(defaultSplitterWidth)
	m.projectSplitter.SetMinSize(defaultSplitterMinSize)
	m.projectSplitter.SetParent(m.rightBox)
	// 项目面板
	m.projectPanel = lcl.NewPanel(owner)
	//m.projectPanel.SetColor(wg.LightenColor(colors.ClAqua, 0.6))
	m.projectPanel.SetAlign(types.AlLeft)
	m.projectPanel.SetWidth(150)    //动态控制
	m.projectPanel.SetVisible(true) //动态控制
	m.projectPanel.SetCaption("项目管理器")
	m.projectPanel.Constraints().SetMinWidth(30)
	m.projectPanel.SetParent(m.rightBox)
	//m.projectPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("projectPanel.SetOnResize BoundsRect:", m.projectPanel.BoundsRect())
	//})

	// 设计器
	m.designerPanel = lcl.NewPanel(owner)
	m.designerPanel.SetBevelOuter(types.BvNone)
	m.designerPanel.SetBorderStyleToBorderStyle(types.BsSingle)
	m.designerPanel.SetDoubleBuffered(true)
	m.designerPanel.SetCaption("设计器画布")
	//m.designerPanel.SetColor(wg.LightenColor(colors.ClAqua, 0.9))
	m.designerPanel.SetAlign(types.AlClient)
	m.designerPanel.SetParent(m.rightBox)

	// 查看器面板分隔器
	m.inspectorSplitter = lcl.NewSplitter(owner)
	m.inspectorSplitter.SetAlign(types.AlRight)
	m.inspectorSplitter.SetWidth(defaultSplitterWidth)
	m.inspectorSplitter.SetMinSize(defaultSplitterMinSize)
	m.inspectorSplitter.SetParent(m.rightBox)
	// 属性检查器
	m.inspectorPanel = lcl.NewPanel(owner)
	//m.inspectorPanel.SetColor(wg.LightenColor(colors.ClAqua, 1.2))
	m.inspectorPanel.SetAlign(types.AlRight)
	m.inspectorPanel.SetBorderStyleToBorderStyle(types.BsSingle)
	m.inspectorPanel.SetCaption("属性检查器")
	m.inspectorPanel.SetWidth(225)    //动态控制
	m.inspectorPanel.SetVisible(true) //动态控制
	m.inspectorPanel.Constraints().SetMinWidth(30)
	m.inspectorPanel.SetParent(m.rightBox)
	//m.inspectorPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("inspectorPanel.SetOnResize BoundsRect:", m.inspectorPanel.BoundsRect())
	//})

	// 日志面板分隔器
	m.consoleLogSplitter = lcl.NewSplitter(owner)
	m.consoleLogSplitter.SetAlign(types.AlBottom)
	m.consoleLogSplitter.SetHeight(defaultSplitterWidth)
	m.consoleLogSplitter.SetMinSize(defaultSplitterMinSize)
	//m.consoleLogSplitter.SetBorderStyleToBorderStyle(types.BsSingle)
	m.consoleLogSplitter.SetParent(m.rightBox)
	// 日志
	m.consoleLogPanel = lcl.NewPanel(owner)
	//m.consoleLogPanel.SetColor(wg.LightenColor(colors.ClAqua, 1.5))
	m.consoleLogPanel.SetBevelOuter(types.BvNone)
	m.consoleLogPanel.SetBorderStyleToBorderStyle(types.BsSingle)
	m.consoleLogPanel.SetDoubleBuffered(true)
	m.consoleLogPanel.SetAlign(types.AlBottom)
	m.consoleLogPanel.SetCaption("控制台输出")
	m.consoleLogPanel.SetHeight(150)   //动态控制
	m.consoleLogPanel.SetVisible(true) //动态控制
	m.consoleLogPanel.Constraints().SetMinHeight(30)
	m.consoleLogPanel.SetParent(m.rightBox)
	//m.consoleLogPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("consoleLogPanel.SetOnResize BoundsRect:", m.consoleLogPanel.BoundsRect())
	//})

	m.contentStatus = lcl.NewStatusBar(owner)
	m.contentStatus.SetBorderWidth(0)
	m.contentStatus.SetShowHint(true)
	m.contentStatus.SetParent(m.box)

	initContentLayoutWidget(m)

	initContentInspector(m)

	return m
}
