package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

type TDesignerForm struct {
	lcl.TEngForm
}

func (m *TDesignerForm) FormCreate(sender lcl.IObject) {
	logs.Info("TDesignerForm FormCreate")
	m.SetLeft(margin)
	m.SetTop(margin)
	m.SetWidth(defaultWidth)
	m.SetHeight(defaultHeight)
	m.SetAlign(types.AlCustom)
	m.SetShowInTaskBar(types.StNever)
	m.SetControlStyle(m.ControlStyle().Include(types.CsNoDesignVisible))
	m.SetBorderStyleToFormBorderStyle(types.BsNone)
	m.SetFormStyle(types.FsNormal)
}

func (m *TDesignerForm) CreateParams(params *types.TCreateParams) {
	logs.Info("TDesignerForm CreateParams ", *params)
}

// 创建设计窗体-隐藏
func (m *FormTab) NewFormDesigner() *DesigningComponent {
	dc := new(DesigningComponent)
	dc.componentType = CtForm
	m.designerBox = dc

	//designerForm := lcl.NewEngForm(nil)
	//designerForm := lcl.NewForm(nil)
	designerForm := &TDesignerForm{}
	lcl.Application.NewForm(designerForm)
	designerForm.SetName(m.name)
	designerForm.SetCaption("")
	// 创建窗体设计器处理器
	formDesigner := NewEngFormDesigner(m)
	m.formDesigner = formDesigner
	designerForm.SetDesigner(formDesigner.Designer())
	//SetDesignMode(designerForm)
	designerForm.SetParent(m.scroll)
	designerForm.SetVisible(true)

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
	//m.designerOnPaint(designerBox)
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
	//m.designerBox.drag.SetParent(m.scroll)
	m.designerBox.drag.SetRelation(m.designerBox)
	m.designerBox.drag.Show()
	m.designerBox.drag.Follow()
	return dc
}
