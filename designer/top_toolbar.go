package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 顶部工具栏

type TopToolbar struct {
	page          lcl.IPageControl
	box           lcl.IPanel
	leftTools     lcl.IPanel
	splitter      lcl.ISplitter // 分割线
	rightTabs     lcl.IPanel    // 组件面板选项卡
	componentTabs map[string]*ComponentTab
}

func (m *TAppWindow) createTopToolbar() {
	bar := &TopToolbar{componentTabs: make(map[string]*ComponentTab)}
	m.toolbar = bar
	// 工具栏面板
	bar.box = lcl.NewPanel(m)
	bar.box.SetParent(m)
	bar.box.SetBevelOuter(types.BvNone)
	bar.box.SetDoubleBuffered(true)
	bar.box.SetWidth(m.Width())
	bar.box.SetHeight(toolbarHeight)
	bar.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	bar.box.SetParentColor(true)

	// 工具栏-分隔线
	bar.splitter = lcl.NewSplitter(m)
	bar.splitter.SetParent(bar.box)
	bar.splitter.SetAlign(types.AlLeft)
	bar.splitter.SetWidth(3)

	// 工具栏-左 工具按钮
	bar.leftTools = lcl.NewPanel(m)
	bar.leftTools.SetParent(bar.box)
	bar.leftTools.SetBevelOuter(types.BvNone)
	bar.leftTools.SetDoubleBuffered(true)
	bar.leftTools.SetWidth(200)
	bar.leftTools.SetHeight(bar.box.Height())
	bar.leftTools.SetAlign(types.AlLeft)
	//bar.leftTools.SetColor(colors.ClRed)

	// 工具栏-右 组件选项卡
	bar.rightTabs = lcl.NewPanel(m)
	bar.rightTabs.SetParent(bar.box)
	bar.rightTabs.SetBevelOuter(types.BvNone)
	bar.rightTabs.SetDoubleBuffered(true)
	bar.rightTabs.SetHeight(bar.box.Height())
	bar.rightTabs.SetAlign(types.AlClient)
	//bar.rightTabs.SetColor(colors.ClBlue)

	// 创建工具按钮
	bar.createToolBarBtns()

	// 创建组件选项卡
	bar.createComponentTabs()
}

// 工具按钮
func (m *TopToolbar) createToolBarBtns() {
	toolbar := lcl.NewToolBar(m.box)
	toolbar.SetParent(m.leftTools)
	toolbar.SetAlign(types.AlCustom)
	toolbar.SetTop(15)
	toolbar.SetButtonWidth(32)
	toolbar.SetButtonHeight(32)
	toolbar.SetHeight(32)
	toolbar.SetWidth(m.leftTools.Width())
	toolbar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	toolbar.SetEdgeBorders(types.NewSet())
	toolbar.SetImages(LoadImageList(m.leftTools, []string{
		"menu/menu_new_form_150.png",
		"menu/menu_project_open_150.png",
		"actions/laz_save_150.png",
		"menu/menu_save_all_150.png",
		"menu/menu_run_150.png",
	}, 24, 24))

	newSepa := func() {
		seap := lcl.NewToolButton(toolbar)
		seap.SetParent(toolbar)
		seap.SetStyle(types.TbsSeparator)
	}

	newBtn := func(imageIndex int32, hint string) lcl.IToolButton {
		btn := lcl.NewToolButton(toolbar)
		btn.SetParent(toolbar)
		btn.SetHint(hint)
		btn.SetImageIndex(imageIndex)
		btn.SetShowHint(true)
		return btn
	}

	newFormBtn := newBtn(0, "新建窗体")
	newFormBtn.SetOnClick(func(sender lcl.IObject) {
		designer.ActiveFormTab(designer.newFormDesignerTab())
	})
	openFormBtn := newBtn(1, "打开窗体")
	openFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	newSepa()
	saveFormBtn := newBtn(2, "保存窗体")
	saveFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	saveAllFormBtn := newBtn(3, "保存所有窗体")
	saveAllFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	newSepa()
	runFormBtn := newBtn(4, "运行预览窗体")
	runFormBtn.SetOnClick(func(sender lcl.IObject) {

	})

}
