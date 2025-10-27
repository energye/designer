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

// 创建设计窗体
func (m *FormTab) NewFormDesigner() *DesigningComponent {
	dc := new(DesigningComponent)
	dc.componentType = CtForm
	dc.createComponentPropertyPage()
	m.formRoot = dc

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

	formRoot := lcl.NewPanel(designerForm)
	formRoot.SetBevelOuter(types.BvNone)
	formRoot.SetBorderStyleToBorderStyle(types.BsSingle)
	formRoot.SetDoubleBuffered(true)
	formRoot.SetParentColor(false)
	formRoot.SetColor(colors.ClBtnFace)
	formRoot.SetName(m.name)
	formRoot.SetCaption("")
	formRoot.SetAlign(types.AlClient)
	formRoot.SetShowHint(true)
	//m.designerOnPaint(formRoot)
	formRoot.SetOnMouseMove(m.designerOnMouseMove)
	formRoot.SetOnMouseDown(m.designerOnMouseDown)
	formRoot.SetOnMouseUp(m.designerOnMouseUp)
	formRoot.SetParent(designerForm)
	//SetDesignMode(formRoot)

	// 设计面板
	m.formRoot.originObject = designerForm
	m.formRoot.object = designerForm
	m.formRoot.formTab = m
	m.formRoot.GetProps()

	// 窗体拖拽大小
	m.formRoot.drag = newDrag(m.scroll, DsRightBottom)
	//m.formRoot.drag.SetParent(m.scroll)
	m.formRoot.drag.SetRelation(m.formRoot)
	m.formRoot.drag.Show()
	m.formRoot.drag.Follow()
	return dc
}
