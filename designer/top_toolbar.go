package designer

import (
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
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
	toolbar.SetBorderStyleToBorderStyle(types.BsNone)
	toolbar.SetImages(m.LoadImageList())
	toolbar.SetAlign(types.AlTop)

	toolBtn := lcl.NewToolButton(toolbar)
	toolBtn.SetParent(toolbar)
	toolBtn.SetHint("提示")
	toolBtn.SetImageIndex(0)
	toolBtn.SetShowHint(true)

	sepc := lcl.NewToolButton(toolbar)
	sepc.SetParent(toolbar)
	sepc.SetStyle(types.TbsSeparator)

	toolBtn = lcl.NewToolButton(toolbar)
	toolBtn.SetParent(toolbar)
	toolBtn.SetHint("提示")
	toolBtn.SetImageIndex(0)
	toolBtn.SetShowHint(true)
}

// 组件选项卡
func (m *TopToolbar) createComponentTabs() {
	page := lcl.NewPageControl(m.box)
	page.SetParent(m.rightTabs)
	page.SetAlign(types.AlClient)
	page.SetTabStop(true)

	standard := lcl.NewTabSheet(page)
	standard.SetParent(page)
	standard.SetCaption("标准控件")
	standard.SetAlign(types.AlClient)

	additional := lcl.NewTabSheet(page)
	additional.SetParent(page)
	additional.SetCaption("额外控件")
	additional.SetAlign(types.AlClient)
}

func (m *TopToolbar) LoadImageList() lcl.IImageList {
	images := lcl.NewImageList(m.leftTools)
	tool.ImageListAddPng(images, "components/default.png")
	return images
}
