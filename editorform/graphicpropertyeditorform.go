package editorform

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 图像加载

type TGraphicPropertyEditorForm struct {
	lcl.TEngForm
	copyButton          lcl.IButton
	pasteButton         lcl.IButton
	loadButton          lcl.IButton
	saveButton          lcl.IButton
	clearButton         lcl.IButton
	okCancelButtonPanel lcl.IButtonPanel
	imagePreview        lcl.IImage
	loadSaveBtnPanel    lcl.IPanel
	openDialog          lcl.IOpenPictureDialog
	saveDialog          lcl.ISavePictureDialog
	groupBox            lcl.IGroupBox
	scrollBox           lcl.IScrollBox
}

func NewGraphicPropertyEditor() {
	designerForm := &TGraphicPropertyEditorForm{}
	lcl.Application.NewForm(designerForm)
}

func (m *TGraphicPropertyEditorForm) FormCreate(sender lcl.IObject) {
	logs.Info("TGraphicPropertyEditorForm FormCreate")
	m.SetCaption("ENERGY Designer 加载图片对话框")
	//m.SetShowInTaskBar(types.StNever)
	m.SetWidth(950)
	m.SetHeight(450)
	m.SetVisible(true)
	m.WorkAreaCenter()
	m.initComponentLayout()
	m.ShowModal()
}

func (m *TGraphicPropertyEditorForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	logs.Info("OnCloseQuery")
}

func (m *TGraphicPropertyEditorForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	logs.Info("OnClose")
}

func (m *TGraphicPropertyEditorForm) initComponentLayout() {
	m.copyButton = lcl.NewButton(m)
	m.pasteButton = lcl.NewButton(m)
	m.loadButton = lcl.NewButton(m)
	m.saveButton = lcl.NewButton(m)
	m.clearButton = lcl.NewButton(m)
	m.okCancelButtonPanel = lcl.NewButtonPanel(m)
	m.imagePreview = lcl.NewImage(m)
	m.loadSaveBtnPanel = lcl.NewPanel(m)
	m.openDialog = lcl.NewOpenPictureDialog(m)
	m.saveDialog = lcl.NewSavePictureDialog(m)
	m.groupBox = lcl.NewGroupBox(m)
	m.scrollBox = lcl.NewScrollBox(m)

	m.groupBox.SetAlign(types.AlClient)
	m.groupBox.SetParent(m)

	m.okCancelButtonPanel.SetAlign(types.AlBottom)
	m.okCancelButtonPanel.SetParent(m)
}
