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

	// 设计模式边框底板: 先于设计窗体创建(同为滚动容器子控件, 创建序靠前 Z 序在后),
	// 窗体覆盖底板中部, 四周露出环带+常驻抓手, 标识设计中的窗体边界
	frame := lcl.NewPanel(m.formDesignPage.formDesignScroll)
	frame.SetParent(m.formDesignPage.formDesignScroll)
	frame.SetBevelOuter(types.BvNone)
	frame.SetColor(designerFrameBgColor)
	frame.SetDoubleBuffered(true)
	frame.SetOnPaint(func(sender lcl.IObject) {
		m.paintDesignFrame(frame)
	})
	frame.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X, Y int32) {
		// 点击设计边框/抓手: 选中窗体根节点(与 Delphi 点击窗体边框选中窗体一致)
		if button == types.MbLeft {
			m.SwitchComponentEditing(m.FormRoot)
		}
	})
	frame.SetBounds(0, 0, 1, 1)
	frame.SendToBack() // 显式压底, 不依赖创建序的隐式 Z 序规则(跨平台)
	m.designFrame = frame

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
	m.FormRoot.drag = newDrag(consts.DsRightBottom)
	m.FormRoot.drag.SetRelation(m.FormRoot)

	// 设计模式边框跟随设计窗体初始边界
	m.syncDesignFrame()
}

// syncDesignFrame 设计模式边框跟随设计窗体边界:
// 边框底板外扩 designFrameMargin 像素置于窗体后(先于窗体创建),
// 四周露出环带+常驻抓手, 常驻标识设计中的窗体
// (设计窗体为 BsNone 无边框, 否则与工作区背景融为一体)。
// 边框与窗体同为滚动容器子控件, 滚动时一起平移, 任意滚动位置下边界始终贴合。
// 全平台生效(覆盖层平台与 Cocoa 容器画布回退路径); 几何变化后经
// updateOverlayPane 重新调用。
func (m *FormTab) syncDesignFrame() {
	if m.designFrame == nil || !m.designFrame.IsValid() ||
		m.FormRoot == nil || m.formDesignPage == nil {
		return
	}
	root := m.FormRoot.Control()
	if root == nil || !root.IsValid() {
		return
	}
	br := root.BoundsRect() // 滚动容器客户区坐标(随滚动一起平移)
	m.designFrame.SetBounds(br.Left-designFrameMargin, br.Top-designFrameMargin,
		br.Width()+designFrameMargin*2, br.Height()+designFrameMargin*2)
}

// paintDesignFrame 绘制设计模式边框环带(Delphi 风格设计态外框):
// 底板中部被窗体盖住, 可见部分 = 浅蓝环带底色 + 紧贴窗体边缘的设计蓝边界线
// + 四角/四边中点的 8 个常驻抓手方块(全部位于窗体外侧环带内)。
// "蓝色 = 设计器元素"与运行时窗体外观明确区分, 避免误认为窗体的真实样子。
// 抓手边长与窗体根节点缩放锚点(dragBorder)一致且位置重合 ——
// 窗体选中时覆盖层黑色锚点精确叠加在右下三个抓手之上(与 Delphi 一致);
// 未选中时抓手为纯装饰, 点击底板选中窗体(见底板 OnMouseDown)。
func (m *FormTab) paintDesignFrame(p lcl.IPanel) {
	canvas := p.Canvas()
	if canvas == nil || !canvas.IsValid() || m.FormRoot == nil {
		return
	}
	root := m.FormRoot.Control()
	if root == nil || !root.IsValid() {
		return
	}
	brush := canvas.BrushToBrush()
	// 环带底色
	brush.SetColor(designerFrameBgColor)
	canvas.FillRectWithRect(types.Rect(0, 0, p.ClientWidth(), p.ClientHeight()))
	br := root.BoundsRect()
	mg := int32(designFrameMargin)
	// 窗体矩形(底板客户区坐标): 底板 = 窗体 ± designFrameMargin
	f := types.Rect(mg, mg, mg+br.Width(), mg+br.Height())
	// 窗体边界线(紧贴窗体边缘, 全部位于窗体外侧不被遮挡)
	brush.SetColor(designerFrameColor)
	canvas.FillRectWithRect(types.Rect(f.Left, f.Top-1, f.Right, f.Top))
	canvas.FillRectWithRect(types.Rect(f.Left, f.Bottom, f.Right, f.Bottom+1))
	canvas.FillRectWithRect(types.Rect(f.Left-1, f.Top-1, f.Left, f.Bottom+1))
	canvas.FillRectWithRect(types.Rect(f.Right, f.Top-1, f.Right+1, f.Bottom+1))
	// 常驻抓手: 四角 + 四边中点, 白底设计蓝边框
	cx, cy := f.Left+f.Width()/2, f.Top+f.Height()/2
	g, h := int32(designGripSize), int32(designGripSize/2)
	for _, r := range []types.TRect{
		types.Rect(f.Left-g, f.Top-g, f.Left, f.Top),         // 左上角
		types.Rect(f.Right, f.Top-g, f.Right+g, f.Top),       // 右上角
		types.Rect(f.Left-g, f.Bottom, f.Left, f.Bottom+g),   // 左下角
		types.Rect(f.Right, f.Bottom, f.Right+g, f.Bottom+g), // 右下角
		types.Rect(cx-h, f.Top-g, cx+h, f.Top),               // 上中
		types.Rect(cx-h, f.Bottom, cx+h, f.Bottom+g),         // 下中
		types.Rect(f.Left-g, cy-h, f.Left, cy+h),             // 左中
		types.Rect(f.Right, cy-h, f.Right+g, cy+h),           // 右中
	} {
		brush.SetColor(designerFrameColor)
		canvas.FillRectWithRect(r)
		brush.SetColor(designerGripFillColor)
		canvas.FillRectWithRect(types.Rect(r.Left+1, r.Top+1, r.Right-1, r.Bottom-1))
	}
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
		// 方向键移动撤销点: 连按合并(换键或超时才捕获一次)
		if m.undoPushDue(keyCode) {
			m.saveUndoPoint()
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

// deletableSelectedComponents 返回可删除的选中组件:
// 排除窗体根节点与被选中祖先的子组件, 包含非可视组件
// (OpenDialog 等非显示控件与可视组件一样可被删除/剪切)。
func (m *FormTab) deletableSelectedComponents() []*TDesigningComponent {
	var result []*TDesigningComponent
	for _, comp := range m.SelectedComponents() {
		if comp == m.FormRoot || comp.ComponentType == consts.CtForm {
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
	targets := m.deletableSelectedComponents()
	if len(targets) == 0 {
		return false
	}
	// 删除撤销点: 删除前捕获
	m.saveUndoPoint()
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
		// 删除撤销点: 删除前捕获
		formTab.saveUndoPoint()
		formTab.RemoveComponentFromSelection(comp)
	}
	parent.FormTab.SwitchComponentEditing(parent)
	comp.Remove()
	triggerUIGeneration(parent, nil, event.CodeGenUI)
}
