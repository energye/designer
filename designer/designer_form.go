package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 创建设计窗体-隐藏
func (m *FormTab) NewFormDesigner() *DesigningComponent {
	dc := new(DesigningComponent)
	dc.componentType = CtForm

	m.designerBox = dc

	designerForm := lcl.NewForm(lcl.Application)
	designerForm.SetLeft(margin)
	designerForm.SetTop(margin)
	designerForm.SetWidth(defaultWidth)
	designerForm.SetHeight(defaultHeight)
	designerForm.SetBorderStyleToFormBorderStyle(types.BsNone)
	designerForm.SetName(m.name)
	designerForm.SetCaption(m.name)
	designerForm.SetAlign(types.AlCustom)
	designerForm.SetVisible(true)
	// 创建窗体设计器处理器
	formDesigner := NewEngFormDesigner(m)
	m.formDesigner = formDesigner
	designerForm.SetDesigner(formDesigner.Designer())
	designerForm.SetParent(m.scroll)

	designerBox := lcl.NewPanel(designerForm)
	designerBox.SetBevelOuter(types.BvNone)
	designerBox.SetBorderStyleToBorderStyle(types.BsSingle)
	designerBox.SetDoubleBuffered(true)
	designerBox.SetParentColor(false)
	designerBox.SetColor(colors.ClBtnFace)
	designerBox.SetName(m.name)
	designerBox.SetCaption("")
	designerBox.SetAlign(types.AlClient)
	designerBox.SetShowHint(true)
	m.designerOnPaint(designerBox)
	designerBox.SetOnMouseMove(m.designerOnMouseMove)
	designerBox.SetOnMouseDown(m.designerOnMouseDown)
	designerBox.SetOnMouseUp(m.designerOnMouseUp)
	designerBox.SetParent(designerForm)
	//SetDesignMode(designerBox)
	dc.designerBox = designerBox

	// 设计面板
	m.designerBox.originObject = designerForm
	m.designerBox.object = designerForm
	m.designerBox.ownerFormTab = m
	m.designerBox.GetProps()

	// 窗体拖拽大小
	m.designerBox.drag = newDrag(m.scroll, DsRightBottom)
	m.designerBox.drag.SetParent(m.scroll)
	m.designerBox.drag.SetRelation(m.designerBox)
	m.designerBox.drag.Show()
	m.designerBox.drag.Follow()
	return dc
}
