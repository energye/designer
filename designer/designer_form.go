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
	"github.com/energye/designer/event"
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
	designerForm.SetOnKeyUp(m.DesignFormOnKeyUp)
	// 拖拽过程中失去焦点(如 Alt+Tab)时收不到鼠标抬起事件，避免对齐线残留。
	designerForm.SetOnDeactivate(func(sender lcl.IObject) {
		m.hideGuides()
	})

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
	designerFormBox.SetDoubleBuffered(true)
	designerFormBox.SetParent(designerForm)
	// 网格只画在设计画布背景一层，用固定瓦片平铺，拖动/缩放时无需重建。
	m.designerBox = designerFormBox
	m.gridVisible = true
	m.designerOnPaint(designerFormBox)

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

const (
	designerKeyMenu   = 18
	designerKeyEscape = 27
	designerKeyLeft   = 37
	designerKeyUp     = 38
	designerKeyRight  = 39
	designerKeyDown   = 40
	designerKeyDelete = 46
)

func (m *FormTab) DesignFormOnKeyDown(sender lcl.IObject, key *types.Word, shift types.TShiftState) {
	if key == nil {
		return
	}
	if int32(*key) == designerKeyMenu || shift.In(types.SsAlt) {
		designerAltDown.Store(true)
	}
	m.handleDesignKey(int32(*key), shift)
}

func (m *FormTab) DesignFormOnKeyUp(sender lcl.IObject, key *types.Word, shift types.TShiftState) {
	if key == nil {
		return
	}
	if int32(*key) == designerKeyMenu || !shift.In(types.SsAlt) {
		designerAltDown.Store(false)
	}
}

// handleDesignKey 处理设计态常用快捷键；返回 true 表示已消费。
func (m *FormTab) handleDesignKey(keyCode int32, shift types.TShiftState) bool {
	comp := m.FindDesignComponent(m.FormRoot)
	switch keyCode {
	case designerKeyDelete:
		if comp == nil || comp == m.FormRoot || comp.ComponentType == consts.CtForm {
			return false
		}
		DeleteDesignComponent(comp)
		return true
	case designerKeyLeft, designerKeyUp, designerKeyRight, designerKeyDown:
		if comp == nil || comp == m.FormRoot {
			return false
		}
		step := int32(1)
		if shift.In(types.SsShift) {
			step = int32(designerGridSize)
		}
		br := comp.BoundsRect()
		left, top := br.Left, br.Top
		switch keyCode {
		case designerKeyLeft:
			left -= step
		case designerKeyRight:
			left += step
		case designerKeyUp:
			top -= step
		case designerKeyDown:
			top += step
		}
		comp.SetBounds(left, top, br.Width(), br.Height())
		actual := comp.BoundsRect()
		comp.UpdateNodeDataPoint(actual.Left, actual.Top)
		return true
	case designerKeyEscape:
		m.SwitchComponentEditing(m.FormRoot)
		return true
	}
	return false
}

// DeleteDesignComponent 删除设计组件并更新布局。
func DeleteDesignComponent(comp *TDesigningComponent) {
	if comp == nil || comp.parent == nil || comp.ComponentType == consts.CtForm {
		return
	}
	parent := comp.parent
	parent.FormTab.SwitchComponentEditing(parent)
	comp.Remove()
	triggerUIGeneration(parent, nil, event.CodeGenUI)
}
