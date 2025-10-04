package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 顶部工具栏

var toolbar *TopToolbar

type TopToolbar struct {
	page            lcl.IPageControl
	box             lcl.IPanel
	leftTools       lcl.IPanel
	splitter        lcl.ISplitter            // 分割线
	rightTabs       lcl.IPanel               // 组件面板选项卡
	componentTabs   map[string]*ComponentTab // 组件选项卡： 标准，附加，通用等等
	selectComponent *ComponentTabItem        // 选中的控件
}

func (m *TAppWindow) createTopToolbar() {
	bar := &TopToolbar{componentTabs: make(map[string]*ComponentTab)}
	toolbar = bar
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
	//bar.splitter = lcl.NewSplitter(m)
	//bar.splitter.SetParent(bar.box)
	//bar.splitter.SetAlign(types.AlLeft)

	// 工具栏-左 工具按钮
	bar.leftTools = lcl.NewPanel(m)
	bar.leftTools.SetParent(bar.box)
	bar.leftTools.SetBevelOuter(types.BvNone)
	bar.leftTools.SetDoubleBuffered(true)
	bar.leftTools.SetWidth(180)
	bar.leftTools.SetHeight(bar.box.Height())
	bar.leftTools.SetAlign(types.AlLeft)

	// 工具栏-右 组件选项卡
	bar.rightTabs = lcl.NewPanel(m)
	bar.rightTabs.SetParent(bar.box)
	bar.rightTabs.SetBevelOuter(types.BvNone)
	bar.rightTabs.SetDoubleBuffered(true)
	bar.rightTabs.SetHeight(bar.box.Height())
	bar.rightTabs.SetAlign(types.AlClient)

	// 创建工具按钮
	bar.createToolBarBtns()

	// 创建组件选项卡
	bar.createComponentTabs()
}

// 重置Tab组件选项卡按下状态
func (m *TopToolbar) ResetTabComponentDown() {
	for _, comp := range m.componentTabs {
		comp.UnDownComponents()
		comp.DownSelectTool()
	}
}

// 设置当前工具按钮选中
// 之后在设计器里使用
func (m *TopToolbar) SetSelectComponentItem(item *ComponentTabItem) {
	m.selectComponent = item
}

// 工具按钮
func (m *TopToolbar) createToolBarBtns() {
	toolBtnBar := lcl.NewToolBar(m.box)
	toolBtnBar.SetParent(m.leftTools)
	toolBtnBar.SetAlign(types.AlCustom)
	toolBtnBar.SetTop(16)
	toolBtnBar.SetButtonWidth(32)
	toolBtnBar.SetButtonHeight(32)
	toolBtnBar.SetHeight(32)
	toolBtnBar.SetWidth(m.leftTools.Width())
	toolBtnBar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	toolBtnBar.SetEdgeBorders(types.NewSet())
	imageList := LoadImageList(m.leftTools, []string{
		"menu/menu_new_form_150.png",
		"menu/menu_project_open_150.png",
		"actions/laz_save_150.png",
		"menu/menu_save_all_150.png",
		"menu/menu_run_150.png",
	}, 24, 24)
	toolBtnBar.SetImages(imageList)
	//toolBtnBar := lcl.NewFlowPanel(m.box)
	//toolBtnBar.SetParent(m.leftTools)
	//toolBtnBar.SetAlign(types.AlCustom)
	//toolBtnBar.SetWidth(m.leftTools.Width())
	//toolBtnBar.SetHeight(m.leftTools.Height())
	//toolBtnBar.SetBevelOuter(types.BvNone)
	//toolBtnBar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	newSepa := func() {
		seap := lcl.NewToolButton(toolBtnBar)
		seap.SetParent(toolBtnBar)
		seap.SetStyle(types.TbsSeparator)
	}

	newBtn := func(imageIndex int32, hint string, margin int32) lcl.IToolButton {
		btn := lcl.NewToolButton(toolBtnBar)
		btn.SetParent(toolBtnBar)
		btn.SetHint(hint)
		btn.SetImageIndex(imageIndex)
		btn.SetShowHint(true)
		//btn := lcl.NewBitBtn(toolBtnBarf)
		//btn.SetParent(toolBtnBarf)
		//btn.SetWidth(32)
		//btn.SetHeight(32)
		//btn.SetTabStop(true)
		//btn.SetImages(imageList)
		//btn.SetImageIndex(imageIndex)
		//btn.SetMargin(margin)
		return btn
	}

	newFormBtn := newBtn(0, "新建窗体", 0)
	newFormBtn.SetOnClick(func(sender lcl.IObject) {
		newForm := designer.addFormDesignerTab()
		designer.ActiveFormTab(newForm)
		inspector.LoadComponent(newForm.form)
	})
	openFormBtn := newBtn(1, "打开窗体", 1)
	openFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	newSepa()
	saveFormBtn := newBtn(2, "保存窗体", 1)
	saveFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	saveAllFormBtn := newBtn(3, "保存所有窗体", 1)
	saveAllFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	newSepa()
	runFormBtn := newBtn(4, "运行预览窗体", 3)
	runFormBtn.SetOnClick(func(sender lcl.IObject) {

	})

}
