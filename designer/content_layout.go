package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
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
	layoutWidget   *ContentLayoutWidget

	projectSplitter lcl.ISplitter // 项目面板分隔器
	projectPanel    lcl.IPanel    // 项目面板
	layoutProject   *ContentLayoutProject

	designerPanel lcl.IPanel // 设计器面板

	inspectorSplitter lcl.ISplitter // 属性检查器面板分隔器
	inspectorPanel    lcl.IPanel    // 属性检查器面板
	layoutInspector   *ContentLayoutInspector

	consoleLogSplitter lcl.ISplitter // 日志面板分隔器
	consoleLogPanel    lcl.IPanel    // 控制台输出
	layoutConsoleLog   *ContentLayoutConsoleLog

	contentStatus       lcl.IStatusBar
	contentStatusCenter lcl.IStatusPanel
	contentStatusRight  lcl.IStatusPanel
}

// 底部布局 - 左: 组件库, 左中: 项目查看, 中: 中间画布(自适应), 右: 属性, 下: 日志控制
func initBottomBox(owner lcl.IWinControl) *ContentLayout {
	windowLayout := config.Config.WindowLayout
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
	m.widgetPanel.SetWidth(windowLayout.ContentLayout.WidgetPanelWidth) //动态控制
	m.widgetPanel.SetVisible(windowLayout.MenuView.WidgetsChecked)      //动态控制
	m.widgetPanel.Constraints().SetMinWidth(30)
	m.widgetPanel.Constraints().SetMaxWidth(400)
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
	m.projectPanel.SetBevelOuter(types.BvNone)
	m.projectPanel.SetBorderStyleToBorderStyle(types.BsSingle)
	m.projectPanel.SetWidth(windowLayout.ContentLayout.ProjectPanelWidth) //动态控制
	m.projectPanel.SetVisible(windowLayout.MenuView.ProjectChecked)       //动态控制
	m.projectPanel.SetCaption("项目管理器")
	m.projectPanel.Constraints().SetMinWidth(30)
	m.projectPanel.Constraints().SetMaxWidth(400)
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
	m.designerPanel.Constraints().SetMinWidth(200)
	m.designerPanel.Constraints().SetMinHeight(200)
	m.designerPanel.SetParent(m.rightBox)

	// 查看器面板分隔器
	m.inspectorSplitter = lcl.NewSplitter(owner)
	m.inspectorSplitter.SetAlign(types.AlRight)
	m.inspectorSplitter.SetWidth(defaultSplitterWidth)
	m.inspectorSplitter.SetResizeAnchor(types.AkRight)
	m.inspectorSplitter.SetMinSize(defaultSplitterMinSize)
	m.inspectorSplitter.SetParent(m.rightBox)
	// 属性检查器
	m.inspectorPanel = lcl.NewPanel(owner)
	//m.inspectorPanel.SetColor(wg.LightenColor(colors.ClAqua, 1.2))
	m.inspectorPanel.SetAlign(types.AlRight)
	m.inspectorPanel.SetBevelOuter(types.BvNone)
	m.inspectorPanel.SetBorderStyleToBorderStyle(types.BsSingle)
	m.inspectorPanel.SetCaption("属性检查器")
	m.inspectorPanel.SetWidth(windowLayout.ContentLayout.InspectorPanelWidth) //动态控制
	m.inspectorPanel.SetVisible(windowLayout.MenuView.InspectorChecked)       //动态控制
	m.inspectorPanel.Constraints().SetMinWidth(30)
	m.inspectorPanel.Constraints().SetMaxWidth(400)
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
	m.consoleLogPanel.SetHeight(windowLayout.ContentLayout.ConsoleLogPanelHeight) //动态控制
	m.consoleLogPanel.SetVisible(windowLayout.MenuView.ConsoleChecked)            //动态控制
	m.consoleLogPanel.Constraints().SetMinHeight(30)
	m.consoleLogPanel.SetParent(m.rightBox)
	//m.consoleLogPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("consoleLogPanel.SetOnResize BoundsRect:", m.consoleLogPanel.BoundsRect())
	//})

	m.contentStatus = lcl.NewStatusBar(owner)
	m.contentStatus.SetBorderWidth(0)
	m.contentStatus.SetAutoSize(false)
	m.contentStatus.SetShowHint(true)
	m.contentStatus.SetAutoHint(true)
	m.contentStatus.SetSimplePanel(false)
	m.contentStatus.SetVisible(windowLayout.MenuView.StatusbarChecked)
	m.contentStatus.SetHeight(30)
	m.contentStatus.SetParent(m.box)
	panels := m.contentStatus.Panels()

	m.contentStatusCenter = panels.AddToStatusPanel()
	m.contentStatusCenter.SetAlignment(types.TaCenter)
	m.contentStatusCenter.SetWidth(250)
	m.contentStatusRight = panels.AddToStatusPanel()
	m.contentStatusRight.SetAlignment(types.TaCenter)

	m.layoutWidget = initContentLayoutWidget(m)
	m.layoutProject = initContentLayoutProject(m)
	m.layoutInspector = initContentLayoutInspector(m)
	m.layoutConsoleLog = initContentLayoutConsoleLog(m)

	return m
}

func (m *ContentLayout) SetStatusCenterText(s string) {
	m.contentStatusCenter.SetText(s)
}

func (m *ContentLayout) SetStatusRightText(s string) {
	m.contentStatusRight.SetText(s)
}

func SetStatusCenterText(s string) {
	if MainWindow.contentLayout != nil {
		MainWindow.contentLayout.SetStatusCenterText(s)
	}
}

func SetStatusRightText(s string) {
	if MainWindow.contentLayout != nil {
		MainWindow.contentLayout.SetStatusRightText(s)
	}
}

// 初始化设计器布局
func (m *ContentLayout) initFromDesignerLayout() *Designer {
	des := new(Designer)
	des.designerForms = make(map[int]*FormTab)
	des.tab = wg.NewTab(m.designerPanel)
	des.tab.SetBounds(0, 0, m.rightBox.Width(), m.rightBox.Height())
	des.tab.SetAlign(types.AlClient)
	des.tab.ScrollLeft().SetTop(3)
	des.tab.ScrollLeft().SetHeight(20)
	des.tab.ScrollLeft().SetColor(wg.DarkenColor(bgLightColor, 0.1))
	des.tab.ScrollRight().SetTop(3)
	des.tab.ScrollRight().SetHeight(20)
	des.tab.ScrollRight().SetColor(wg.DarkenColor(bgLightColor, 0.1))
	des.tab.EnableScrollButton(false)
	des.tab.SetParent(m.designerPanel)

	des.defaultTip = wg.NewButton(des.tab)
	des.defaultTip.SetDisabledColor(colors.RGBToColor(204, 232, 255), colors.RGBToColor(204, 232, 255))
	defaultTipBr := des.defaultTip.BoundsRect()
	dtw, dth := int32(140), int32(140)
	defaultTipBr.Left = m.rightBox.Width()/2 - dtw/2
	defaultTipBr.Top = m.rightBox.Height()/2 - dth/2
	defaultTipBr.SetWidth(dtw)
	defaultTipBr.SetHeight(dth)
	des.defaultTip.SetBoundsRect(defaultTipBr)
	des.defaultTip.SetAlpha(80)
	des.defaultTip.SetRadius(15)
	des.defaultTip.TextAlign = wg.TextAlignLeft
	des.defaultTip.TextLineSpacing = 8
	des.defaultTip.TextOffSetX = 10
	des.defaultTip.Font().SetSize(10)
	des.defaultTip.Font().SetColor(wg.DarkenColor(colors.ClGray, 0.2))
	des.defaultTip.SetDisable(true)
	des.defaultTip.SetAnchors(types.NewSet())
	des.defaultTip.SetText("新建项目 (Ctrl+P)\n打开项目 (Ctrl+O)\n新建窗体 (Ctrl+N)\n应用配置 (Ctrl+F11)\n　　运行 (F9)")
	des.defaultTip.SetParent(des.tab)

	des.createTabMenu()
	return des
}
