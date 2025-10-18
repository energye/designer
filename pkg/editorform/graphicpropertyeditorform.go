package editorform

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 图片加载

type DialogCallback func(filePath string, ok bool)

type TGraphicPropertyEditorForm struct {
	lcl.TEngForm
	loadButton          lcl.IButton
	saveButton          lcl.IButton
	clearButton         lcl.IButton
	copyButton          lcl.IButton
	pasteButton         lcl.IButton
	okCancelButtonPanel lcl.IButtonPanel
	imagePreview        lcl.IImage
	loadSaveBtnPanel    lcl.IPanel
	openDialog          lcl.IOpenPictureDialog
	saveDialog          lcl.ISavePictureDialog
	groupBox            lcl.IGroupBox
	scrollBox           lcl.IScrollBox
	dialogCallback      DialogCallback
	imageFilePath       string
}

func NewGraphicPropertyEditor(dialogCallback DialogCallback) *TGraphicPropertyEditorForm {
	designerForm := &TGraphicPropertyEditorForm{dialogCallback: dialogCallback}
	lcl.Application.NewForm(designerForm)
	//designerForm.IEngForm = lcl.NewEngForm(nil)
	//designerForm.FormCreate(nil)
	//designerForm.SetOnCloseQuery(designerForm.OnCloseQuery)
	//designerForm.SetOnClose(designerForm.OnClose)
	return designerForm
}

func (m *TGraphicPropertyEditorForm) FormCreate(sender lcl.IObject) {
	logs.Info("TGraphicPropertyEditorForm FormCreate")
	m.SetCaption("ENERGY Designer 加载图片对话框")
	//m.SetShowInTaskBar(types.StNever)
	m.SetWidth(950)
	m.SetHeight(450)
	m.WorkAreaCenter()
	//m.SetPopupParent(m.form)
	//m.SetPopupMode(types.PmExplicit)
	m.initComponentLayout()
}

func (m *TGraphicPropertyEditorForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	logs.Info("TGraphicPropertyEditorForm OnCloseQuery")
}

func (m *TGraphicPropertyEditorForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	logs.Info("TGraphicPropertyEditorForm OnClose")
}

func (m *TGraphicPropertyEditorForm) initComponentLayout() {
	// 组件创建
	m.loadButton = lcl.NewButton(m)
	m.saveButton = lcl.NewButton(m)
	m.clearButton = lcl.NewButton(m)
	m.copyButton = lcl.NewButton(m)
	m.pasteButton = lcl.NewButton(m)
	m.okCancelButtonPanel = lcl.NewButtonPanel(m)
	m.imagePreview = lcl.NewImage(m)
	m.loadSaveBtnPanel = lcl.NewPanel(m)
	m.openDialog = lcl.NewOpenPictureDialog(m)
	m.saveDialog = lcl.NewSavePictureDialog(m)
	m.groupBox = lcl.NewGroupBox(m)
	m.scrollBox = lcl.NewScrollBox(m)

	{
		// dialog
		m.openDialog.SetTitle("打开图片文件")
		m.openDialog.SetFilter(config.DialogFilter.ImageFilter())
		m.openDialog.SetFilterIndex(1)

		m.saveDialog.SetTitle("保存图片文件")
		m.saveDialog.SetFilter(config.DialogFilter.ImageFilter())
		m.saveDialog.SetFilterIndex(1)
	}

	// 图片显示主体
	m.groupBox.SetAlign(types.AlClient)
	m.groupBox.BorderSpacing().SetAround(6)
	m.groupBox.SetCaption("图片预览")
	m.groupBox.SetParent(m)

	// 确定取消按钮
	m.okCancelButtonPanel.SetAlign(types.AlBottom)
	m.okCancelButtonPanel.SetShowButtons(types.NewSet(types.PbOK, types.PbCancel))
	m.okCancelButtonPanel.OKButton().SetCaption("确定")
	m.okCancelButtonPanel.CancelButton().SetCaption("取消")
	m.okCancelButtonPanel.OKButton().SetOnClick(func(sender lcl.IObject) {
		logs.Debug("OKButton().SetOnClick")
		m.SetModalResult(types.MrOk)
		if m.dialogCallback != nil {
			m.dialogCallback(m.imageFilePath, true)
		}
		m.Close()
	})
	m.okCancelButtonPanel.CancelButton().SetOnClick(func(sender lcl.IObject) {
		logs.Debug("CancelButton().SetOnClick")
		m.SetModalResult(types.MrCancel)
		if m.dialogCallback != nil {
			m.dialogCallback(m.imageFilePath, false)
		}
		m.Close()
	})
	m.okCancelButtonPanel.SetParent(m)

	// 图片滚动条
	m.scrollBox.SetAlign(types.AlClient)
	m.scrollBox.SetAutoScroll(true)
	m.scrollBox.BorderSpacing().SetAround(6)
	m.scrollBox.HorzScrollBar().SetTracking(true)
	m.scrollBox.HorzScrollBar().SetVisible(true)
	m.scrollBox.VertScrollBar().SetTracking(true)
	m.scrollBox.VertScrollBar().SetVisible(true)
	m.scrollBox.SetParent(m.groupBox)
	m.scrollBox.SetOnResize(m.imagePreviewOnPictureChanged)

	// 图片显示
	m.imagePreview.SetAutoSize(true)
	m.imagePreview.SetCenter(true)
	m.imagePreview.SetParent(m.scrollBox)
	m.imagePreview.SetOnPictureChanged(m.imagePreviewOnPictureChanged)
	m.imagePreview.SetOnPaintBackground(m.imagePreviewOnPaintBackground)

	// 图片功能按钮
	m.loadSaveBtnPanel.SetWidth(90)
	m.loadSaveBtnPanel.SetBevelOuter(types.BvNone)
	m.loadSaveBtnPanel.SetAlign(types.AlRight)
	m.loadSaveBtnPanel.SetParent(m.groupBox)
	{
		var left int32 = 10
		m.loadButton.SetCaption("加 载")
		m.loadButton.SetTop(20)
		m.loadButton.SetLeft(left)
		m.loadButton.SetParent(m.loadSaveBtnPanel)
		m.loadButton.SetOnClick(m.loadImageBtnOnClick)

		m.saveButton.SetCaption("保 存")
		m.saveButton.SetLeft(left)
		m.saveButton.SetTop(m.loadButton.Height() + m.loadButton.Top() + 20)
		m.saveButton.SetParent(m.loadSaveBtnPanel)
		m.saveButton.SetOnClick(m.saveImageBtnOnClick)

		m.clearButton.SetCaption("清 除")
		m.clearButton.SetLeft(left)
		m.clearButton.SetTop(m.saveButton.Height() + m.saveButton.Top() + 20)
		m.clearButton.SetParent(m.loadSaveBtnPanel)
		m.clearButton.SetOnClick(m.clearImageBtnOnClick)

		m.copyButton.SetCaption("复 制")
		m.copyButton.SetLeft(left)
		m.copyButton.SetTop(m.clearButton.Height() + m.clearButton.Top() + 20)
		m.copyButton.SetParent(m.loadSaveBtnPanel)
		m.copyButton.SetOnClick(m.copyImageBtnOnClick)

		m.pasteButton.SetCaption("粘 贴")
		m.pasteButton.SetLeft(left)
		m.pasteButton.SetTop(m.copyButton.Height() + m.copyButton.Top() + 20)
		m.pasteButton.SetParent(m.loadSaveBtnPanel)
		m.pasteButton.SetOnClick(m.pasteImageBtnOnClick)
	}
}

func (m *TGraphicPropertyEditorForm) imagePreviewOnPictureChanged(sender lcl.IObject) {
	imgBr := m.imagePreview.BoundsRect()
	if imgBr.Width() < m.scrollBox.ClientWidth() {
		m.imagePreview.SetLeft((m.scrollBox.ClientWidth() - imgBr.Width()) / 2)
	} else {
		m.imagePreview.SetLeft(0)
	}
	if imgBr.Height() < m.scrollBox.ClientHeight() {
		m.imagePreview.SetTop((m.scrollBox.ClientHeight() - imgBr.Height()) / 2)
	} else {
		m.imagePreview.SetTop(0)
	}
}

func (m *TGraphicPropertyEditorForm) imagePreviewOnPaintBackground(sender lcl.IObject, canvas lcl.ICanvas, rect types.TRect) {
	cell := int32(8)
	bmp := lcl.NewBitmap()
	bmp.SetPixelFormat(types.Pf24bit)
	bmp.SetSize(rect.Width(), rect.Height())
	bmpCanvas := bmp.Canvas()
	bmpCanvas.BrushToBrush().SetColor(colors.ClWhite)
	bmpCanvas.FillRectWithIntX4(0, 0, bmp.Width(), bmp.Height())
	bmpCanvas.BrushToBrush().SetColor(colors.ClLtGray)
	for i := 0; i < int(bmp.Width()/cell); i++ {
		for j := 0; j < int(bmp.Height()/cell); j++ {
			if (i%2 != 0) == (j%2 != 0) {
				bmpCanvas.FillRectWithIntX4(int32(i)*cell, int32(j)*cell, int32(i+1)*cell, int32(j+1)*cell)
			}
		}
	}
	sourceRect := types.TRect{Left: 0, Top: 0}
	sourceRect.SetWidth(bmp.Width())
	sourceRect.SetHeight(bmp.Height())
	canvas.CopyRectWithRectX2Canvas(rect, bmpCanvas, sourceRect)
}

func (m *TGraphicPropertyEditorForm) loadImageBtnOnClick(sender lcl.IObject) {
	logs.Debug("TGraphicPropertyEditorForm loadImageBtnOnClick")
	if m.openDialog.Execute() {
		fileName := m.openDialog.FileName()
		m.imagePreview.Picture().LoadFromFile(fileName)
		m.imageFilePath = fileName
	}
}

func (m *TGraphicPropertyEditorForm) saveImageBtnOnClick(sender lcl.IObject) {
	logs.Debug("TGraphicPropertyEditorForm saveImageBtnOnClick")
	if m.saveDialog.Execute() {
		fileName := m.saveDialog.FileName()
		m.imagePreview.Picture().SaveToFile(fileName, m.saveDialog.GetFilterExt())
		m.imageFilePath = fileName
	}
}

func (m *TGraphicPropertyEditorForm) clearImageBtnOnClick(sender lcl.IObject) {
	m.imagePreview.Picture().Clear()
	m.imageFilePath = ""
}

func (m *TGraphicPropertyEditorForm) copyImageBtnOnClick(sender lcl.IObject) {
	lcl.Clipboard.Assign(m.imagePreview.Picture().Graphic())
}

func (m *TGraphicPropertyEditorForm) pasteImageBtnOnClick(sender lcl.IObject) {
	m.imagePreview.Picture().Assign(lcl.Clipboard)
}
