// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

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
func (m *FormTab) NewFormDesigner() *TDesigningComponent {
	dc := new(TDesigningComponent)
	dc.ComponentType = CtForm
	dc.createComponentPropertyPage()
	m.FormRoot = dc

	//designerForm := lcl.NewEngForm(nil)
	//designerForm := lcl.NewForm(nil)
	designerForm := &TDesignerForm{}
	lcl.Application.NewForm(designerForm)
	designerForm.SetName(m.Name)
	designerForm.SetCaption(m.Name)
	// 创建窗体设计器处理器
	formDesigner := NewEngFormDesigner(m)
	m.formDesigner = formDesigner
	designerForm.SetDesigner(formDesigner.Designer())
	//SetDesignMode(designerForm)
	designerForm.SetParent(m.scroll)
	designerForm.SetVisible(true)
	//designerForm.SetOnMouseMove(m.designerOnMouseMove)
	//designerForm.SetOnMouseDown(m.designerOnMouseDown)
	//designerForm.SetOnMouseUp(m.designerOnMouseUp)

	formRoot := lcl.NewPanel(designerForm)
	formRoot.SetBevelOuter(types.BvNone)
	formRoot.SetBorderStyleToBorderStyle(types.BsSingle)
	formRoot.SetDoubleBuffered(true)
	formRoot.SetParentColor(false)
	formRoot.SetColor(colors.ClBtnFace)
	formRoot.SetName(m.Name)
	formRoot.SetCaption("")
	formRoot.SetAlign(types.AlClient)
	formRoot.SetShowHint(true)
	//m.designerOnPaint(FormRoot)
	formRoot.SetOnMouseMove(m.designerOnMouseMove)
	formRoot.SetOnMouseDown(m.designerOnMouseDown)
	formRoot.SetOnMouseUp(m.designerOnMouseUp)
	formRoot.SetParent(designerForm)
	//SetDesignMode(FormRoot)

	// 设计面板
	dc.originObject = designerForm
	dc.object = designerForm
	dc.formTab = m

	// 窗体拖拽大小
	dc.drag = newDrag(m.scroll, DsRightBottom)
	//m.FormRoot.drag.SetParent(m.scroll)
	dc.drag.SetRelation(dc)
	dc.drag.Show()
	dc.drag.Follow()
	return dc
}
