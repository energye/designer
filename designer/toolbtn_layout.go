package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type TToolBtnLayout struct {
	box        lcl.IPanel
	toolbarBtn *TToolbarToolBtn // 工具栏按钮
}

func initToolBtnLayout(owner lcl.IWinControl) *TToolBtnLayout {
	m := &TToolBtnLayout{}
	m.box = lcl.NewPanel(owner)
	//m.box.SetColor(colors.ClLegacySkyBlue)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlTop)
	m.box.SetHeight(30)
	m.box.SetParent(owner)

	m.initToolBarBtns()
	return m
}

// 工具按钮
func (m *TToolBtnLayout) initToolBarBtns() {
	m.toolbarBtn = new(TToolbarToolBtn)

	toolBtnBar := lcl.NewToolBar(m.box)
	toolBtnBar.SetAlign(types.AlClient)
	toolBtnBar.SetButtonWidth(32)
	toolBtnBar.SetButtonHeight(32)
	toolBtnBar.SetHeight(32)
	toolBtnBar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	toolBtnBar.SetEdgeBorders(types.NewSet())
	toolBtnBar.SetImages(imageMenu.ImageList150())
	toolBtnBar.SetParent(m.box)
	m.toolbarBtn.toolBtnBar = toolBtnBar

	newSep := func() {
		sep := lcl.NewToolButton(toolBtnBar)
		sep.SetParent(toolBtnBar)
		sep.SetStyle(types.TbsSeparator)
	}

	newBtn := func(imageIndex int32, hint string, margin int32) lcl.IToolButton {
		btn := lcl.NewToolButton(toolBtnBar)
		btn.SetParent(toolBtnBar)
		btn.SetHint(hint)
		btn.SetImageIndex(imageIndex)
		btn.SetShowHint(true)
		return btn
	}

	m.toolbarBtn.newWindowBtn = newBtn(imageMenu.ImageIndex("menu_new_form_150.png"), "新建窗体(Ctrl+N)", 0)
	m.toolbarBtn.newWindowBtn.SetOnClick(m.toolbarBtn.onNewForm)

	m.toolbarBtn.openBtn = newBtn(imageMenu.ImageIndex("menu_project_open_150.png"), "打开(Ctrl+O)", 1)
	m.toolbarBtn.openBtn.SetOnClick(m.toolbarBtn.onOpenForm)
	newSep()

	//toolbarBtn.saveBtn = newBtn(imageMenu.ImageIndex("menu_save_150.png"), "保存(Ctrl+S)", 1)
	//toolbarBtn.saveBtn.SetOnClick(toolbarBtn.onSaveForm)

	//toolbarBtn.saveAllFormBtn = newBtn(imageMenu.ImageIndex("menu_save_all_150.png"), "保存所有窗体", 1)
	//toolbarBtn.saveAllFormBtn.SetOnClick(toolbarBtn.onSaveAllForm)
	//newSep()

	m.toolbarBtn.runPreviewBtn = newBtn(imageMenu.ImageIndex("menu_run_150.png"), "运行(F9)", 3)
	m.toolbarBtn.runPreviewBtn.SetOnClick(m.toolbarBtn.onRunPreviewForm)
}
