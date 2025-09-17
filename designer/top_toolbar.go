package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
)

// 顶部工具栏

type TopToolbar struct {
	box       lcl.IPanel
	leftTools lcl.IPanel
	splitter  lcl.ISplitter // 分割线
	rightTabs lcl.IPanel    // 组件面板选项卡
}

func (m *TAppWindow) createTopToolbar() {
	bar := &TopToolbar{}
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
	toolbar.SetImages(m.LoadImageList([]string{
		"menu/menu_new_form_150.png",
		"menu/menu_project_open_150.png",
		"menu/menu_saveas_150.png",
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

// 组件选项卡
func (m *TopToolbar) createComponentTabs() {
	page := lcl.NewPageControl(m.box)
	page.SetParent(m.rightTabs)
	page.SetAlign(types.AlClient)
	page.SetTabStop(true)

	// 创建组件选项卡
	newComponentTab := func(tab config.Tab) {
		standard := lcl.NewTabSheet(page)
		standard.SetParent(page)
		standard.SetCaption(tab.Cn)
		standard.SetAlign(types.AlClient)
		// 组件图标
		var imageList []string
		for _, name := range tab.Component {
			imageList = append(imageList, fmt.Sprintf("components/%v_150.png", strings.ToLower(name)))
		}
		// 显示组件工具按钮
		toolbar := lcl.NewToolBar(standard)
		toolbar.SetParent(standard)
		toolbar.SetImages(m.LoadImageList(imageList, 36, 36))
		toolbar.SetButtonWidth(36)
		toolbar.SetButtonHeight(36)
		toolbar.SetHeight(36)
		toolbar.SetEdgeBorders(types.NewSet())
		// 创建组件按钮
		for i, name := range tab.Component {
			toolBtn := lcl.NewToolButton(toolbar)
			toolBtn.SetParent(toolbar)
			toolBtn.SetHint(name)
			toolBtn.SetImageIndex(int32(i))
			toolBtn.SetShowHint(true)
			cIdx, cName := i, name
			toolBtn.SetOnClick(func(sender lcl.IObject) {
				println(cIdx, cName)
			})
		}

	}
	// 创建组件选项卡
	newComponentTab(config.Config.ComponentTabs.Standard)
	newComponentTab(config.Config.ComponentTabs.Additional)
	newComponentTab(config.Config.ComponentTabs.Common)
	newComponentTab(config.Config.ComponentTabs.Dialogs)
	newComponentTab(config.Config.ComponentTabs.Misc)
	newComponentTab(config.Config.ComponentTabs.System)
	newComponentTab(config.Config.ComponentTabs.LazControl)
	newComponentTab(config.Config.ComponentTabs.WebView)

}

func (m *TopToolbar) LoadImageList(imageList []string, width, height int32) lcl.IImageList {
	images := lcl.NewImageList(m.leftTools)
	images.SetWidth(width)
	images.SetHeight(height)
	for _, image := range imageList {
		tool.ImageListAddPng(images, image)
	}
	return images
}
