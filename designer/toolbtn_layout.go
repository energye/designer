package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type TToolBtnLayout struct {
	box        lcl.IPanel
	toolbarBtn *TToolbarToolBtn // 工具栏按钮
}

func initToolBtnLayout(window lcl.IWinControl) *TToolBtnLayout {
	windowBr := window.ClientRect()
	m := &TToolBtnLayout{}
	m.box = lcl.NewPanel(window)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	if switchAutoContentLayoutAlign {
		m.box.SetAlign(types.AlTop)
	} else {
		m.box.SetAlign(types.AlCustom)
		m.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
		m.box.SetBounds(0, 0, windowBr.Width(), toolBarHeight)
		m.box.SetHeight(toolBarHeight)
	}
	m.box.SetParent(window)

	m.initToolBarBtns()
	return m
}

// 工具按钮
func (m *TToolBtnLayout) initToolBarBtns() {
	m.toolbarBtn = new(TToolbarToolBtn)

	toolBtnBar := lcl.NewToolBar(m.box)
	toolBtnBar.SetAlign(types.AlClient)
	toolBtnBar.SetButtonWidth(24)
	toolBtnBar.SetButtonHeight(24)
	toolBtnBar.BorderSpacing().SetLeft(3)
	toolBtnBar.BorderSpacing().SetTop(3)
	toolBtnBar.BorderSpacing().SetBottom(3)
	toolBtnBar.SetHeight(24)
	toolBtnBar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	toolBtnBar.SetEdgeBorders(types.NewSet())
	toolBtnBar.SetImages(imageMenu.ImageList150())
	toolBtnBar.SetShowCaptions(true)
	toolBtnBar.SetList(true)
	toolBtnBar.SetParent(m.box)
	m.toolbarBtn.toolBtnBar = toolBtnBar

	newSep := func() {
		sep := lcl.NewToolButton(toolBtnBar)
		sep.SetParent(toolBtnBar)
		sep.SetStyle(types.TbsSeparator)
	}

	newBtn := func(imageIndex int32, hint string, margin int32, text string) lcl.IToolButton {
		btn := lcl.NewToolButton(toolBtnBar)
		btn.SetParent(toolBtnBar)
		btn.SetHint(hint)
		btn.SetImageIndex(imageIndex)
		btn.SetShowHint(true)
		btn.SetCaption(text)
		return btn
	}

	m.toolbarBtn.newWindowBtn = newBtn(imageMenu.ImageIndex("menu_new_form.png"), "新建窗体(Ctrl+N)", 0, "新建")
	m.toolbarBtn.newWindowBtn.SetOnClick(m.toolbarBtn.onNewForm)

	m.toolbarBtn.openBtn = newBtn(imageMenu.ImageIndex("menu_project_open.png"), "打开(Ctrl+O)", 1, "打开")
	m.toolbarBtn.openBtn.SetOnClick(m.toolbarBtn.onOpenForm)
	newSep()

	//toolbarBtn.saveBtn = newBtn(imageMenu.ImageIndex("menu_save_150.png"), "保存(Ctrl+S)", 1)
	//toolbarBtn.saveBtn.SetOnClick(toolbarBtn.onSaveForm)

	//toolbarBtn.saveAllFormBtn = newBtn(imageMenu.ImageIndex("menu_save_all_150.png"), "保存所有窗体", 1)
	//toolbarBtn.saveAllFormBtn.SetOnClick(toolbarBtn.onSaveAllForm)
	//newSep()

	m.toolbarBtn.runPreviewBtn = newBtn(imageMenu.ImageIndex("menu_run.png"), "运行(F9)", 3, "运行(F9)")
	m.toolbarBtn.runPreviewBtn.SetOnClick(m.toolbarBtn.onRunPreviewForm)
}
