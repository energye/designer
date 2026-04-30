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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

type TDesignerForm struct {
	//lcl.TEngForm
	*lcl.TEngForm
}

func (m *TDesignerForm) FormCreate(sender lcl.IObject) {
	logs.Info("TDesignerForm FormCreate")
	m.SetWidth(designerDefaultWidth)
	m.SetHeight(designerDefaultHeight)
	m.SetAlign(types.AlCustom)
	m.SetBorderStyleToFormBorderStyle(types.BsNone)
	//m.SetShowInTaskBar(types.StNever)
	//m.SetControlStyle(m.ControlStyle().Include(types.CsNoDesignVisible))
	//m.SetFormStyle(types.FsNormal)
}

// 创建设计窗体
func (m *FormTab) NewFormDesigner() {
	// 创建组件属性树
	m.FormRoot.mustComponentPropertyPage()

	// 创建设计窗体实例
	formInstance := api.NewInstanceByComponentClass(lcl.TEngFormClass())
	//SetComponentDesignMode(formInstance)
	api.CreateObjectByComponent(formInstance, 0)
	newDesForm := lcl.AsEngForm(formInstance)
	designerForm := &TDesignerForm{TEngForm: newDesForm.(*lcl.TEngForm)}
	//newDesForm.SetControlStyle(newDesForm.ControlStyle().Include(types.CsOwnedChildrenNotSelectable))
	designerForm.FormCreate(designerForm)
	designerForm.SetName(m.FormRoot.Name())
	designerForm.SetCaption(m.FormRoot.Name())
	// 创建窗体设计器处理器
	formDesigner := NewEngFormDesigner(m)
	m.formDesigner = formDesigner
	formDesigner.LookupRoot = m.FormRoot // formRoot
	formDesigner.Form = designerForm
	designerForm.SetDesigner(formDesigner.Designer())
	designerForm.SetParent(m.formDesignPage.formDesignScroll)

	designerForm.SetKeyPreview(true)
	designerForm.SetOnKeyDown(m.DesignFormOnKeyDown)

	designerForm.Show()

	// 创建窗体设计面板, 放置实际设计的组件
	designerFormBoxInstance := api.NewInstanceByComponentClass(lcl.TCustomPanelClass())
	SetComponentDesignMode(designerFormBoxInstance) // 设置为设计模式
	api.CreateObjectByComponent(designerFormBoxInstance, formInstance)
	designerFormBox := lcl.AsCustomPanel(designerFormBoxInstance)
	//designerFormBox.SetControlStyle(designerFormBox.ControlStyle().Include(types.CsOwnedChildrenNotSelectable))
	designerFormBox.SetAlign(types.AlClient)
	designerFormBox.SetBevelOuter(types.BvNone)
	designerFormBox.SetColor(colors.ClForm)
	designerFormBox.SetParent(designerForm)

	// 设计面板
	m.FormRoot.originObject = designerForm
	m.FormRoot.object = designerFormBox
	m.FormRoot.FormTab = m
	formDesigner.AddComponentToList(m.FormRoot)

	// 窗体拖拽
	m.FormRoot.drag = newDrag(m.formDesignPage.formDesignScroll, consts.DsRightBottom)
	m.FormRoot.drag.mustDS()
	m.FormRoot.drag.SetRelation(m.FormRoot)
}

func (m *FormTab) DesignFormOnKeyDown(sender lcl.IObject, key *types.Word, shift types.TShiftState) {
	a := m.FindDesignComponent(m.FormRoot)
	if a == nil {
		a = m.FormRoot
	}
	//println("[Designer KeyDown] key:", *key, "shift:", shift, shift.In(types.SsCtrl), "Form:", m.FormRoot.Name(), "Component:", a.Name())
}
