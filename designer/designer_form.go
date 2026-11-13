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
	designerForm.SetDesigner(formDesigner.Designer())
	// 拖拽过程中失去焦点(如 Alt+Tab)时收不到鼠标抬起事件，避免对齐线/框选残留。
	designerForm.SetOnDeactivate(func(sender lcl.IObject) {
		m.hideGuides()
		if m.formDesigner != nil {
			m.formDesigner.bandCancel()
		}
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
	designerKeyA      = 65
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
	switch keyCode {
	case designerKeyDelete:
		return m.DeleteSelectedComponents()
	case designerKeyLeft, designerKeyUp, designerKeyRight, designerKeyDown:
		comps := m.movableSelectedComponents()
		if len(comps) == 0 {
			return false
		}
		step := int32(1)
		if shift.In(types.SsShift) {
			step = int32(designerGridSize)
		}
		dx, dy := int32(0), int32(0)
		switch keyCode {
		case designerKeyLeft:
			dx = -step
		case designerKeyRight:
			dx = step
		case designerKeyUp:
			dy = -step
		case designerKeyDown:
			dy = step
		}
		for _, comp := range comps {
			br := comp.BoundsRect()
			comp.SetBounds(br.Left+dx, br.Top+dy, br.Width(), br.Height())
			actual := comp.BoundsRect()
			comp.UpdateNodeDataPoint(actual.Left, actual.Top)
		}
		m.syncSelectionFrames()
		return true
	case designerKeyEscape:
		m.SwitchComponentEditing(m.FormRoot)
		return true
	case designerKeyA:
		if shift.In(types.SsCtrl) {
			m.SelectAllComponents()
			return true
		}
	}
	return false
}

// movableSelectedComponents 返回可移动的选中组件（排除窗体根、非可视和被选中祖先的子组件）
func (m *FormTab) movableSelectedComponents() []*TDesigningComponent {
	var result []*TDesigningComponent
	for _, comp := range m.SelectedComponents() {
		if comp == m.FormRoot || comp.ComponentType == consts.CtForm || comp.ComponentType == consts.CtNonVisual {
			continue
		}
		if m.HasSelectedAncestor(comp) {
			continue
		}
		result = append(result, comp)
	}
	return result
}

// DeleteSelectedComponents 删除所有选中组件并更新布局
func (m *FormTab) DeleteSelectedComponents() bool {
	targets := m.movableSelectedComponents()
	if len(targets) == 0 {
		return false
	}
	targetSet := make(map[*TDesigningComponent]struct{}, len(targets))
	for _, comp := range targets {
		targetSet[comp] = struct{}{}
	}
	// 只删除未被其它待删组件包含的组件，避免父组件释放后访问子组件
	var top []*TDesigningComponent
	for _, comp := range targets {
		contained := false
		for parent := comp.parent; parent != nil; parent = parent.parent {
			if _, ok := targetSet[parent]; ok {
				contained = true
				break
			}
		}
		if !contained {
			top = append(top, comp)
		}
	}
	m.ClearSelection()
	parents := make(map[*TDesigningComponent]struct{})
	for _, comp := range top {
		if comp.parent == nil {
			continue
		}
		parents[comp.parent] = struct{}{}
		comp.Remove()
	}
	for parent := range parents {
		triggerUIGeneration(parent, nil, event.CodeGenUI)
	}
	// 单选删除后保持原有行为：选中父容器
	if len(parents) == 1 {
		for parent := range parents {
			m.SwitchComponentEditing(parent)
		}
	}
	return true
}

// DeleteDesignComponent 删除设计组件并更新布局。
func DeleteDesignComponent(comp *TDesigningComponent) {
	if comp == nil || comp.parent == nil || comp.ComponentType == consts.CtForm {
		return
	}
	parent := comp.parent
	formTab := comp.FormTab
	if formTab != nil {
		formTab.RemoveComponentFromSelection(comp)
	}
	parent.FormTab.SwitchComponentEditing(parent)
	comp.Remove()
	triggerUIGeneration(parent, nil, event.CodeGenUI)
}
