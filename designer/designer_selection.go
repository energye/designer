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
// See the License for the specific language governing permissions and
// limitations under the License.

package designer

import (
	"github.com/energye/designer/consts"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 设计 - 画布选择模型
// 主选中组件显示拖拽控制点和属性面板，其余选中组件显示 1px 蓝色边框。

var designerSelectionColor = colors.RGBToColor(0, 125, 255)

type designerSelectionFrame struct {
	top    lcl.IPanel
	bottom lcl.IPanel
	left   lcl.IPanel
	right  lcl.IPanel
}

func newOverlayPanel(box lcl.IWinControl, color types.TColor) lcl.IPanel {
	p := lcl.NewPanel(box)
	p.SetBounds(0, 0, 0, 0)
	p.SetVisible(false)
	p.SetEnabled(false)
	p.SetBevelOuter(types.BvNone)
	p.SetColor(color)
	p.SetParent(box)
	return p
}

func newDesignerSelectionFrame(box lcl.IWinControl) *designerSelectionFrame {
	return &designerSelectionFrame{
		top:    newOverlayPanel(box, designerSelectionColor),
		bottom: newOverlayPanel(box, designerSelectionColor),
		left:   newOverlayPanel(box, designerSelectionColor),
		right:  newOverlayPanel(box, designerSelectionColor),
	}
}

func (f *designerSelectionFrame) valid() bool {
	return f != nil && f.top.IsValid() && f.bottom.IsValid() && f.left.IsValid() && f.right.IsValid()
}

// show 绘制在组件外侧，避免盖住控件自身边框。
func (f *designerSelectionFrame) show(r types.TRect) {
	l, t, w, h := r.Left, r.Top, r.Width(), r.Height()
	panels := []struct {
		p          lcl.IPanel
		x, y, w, h int32
	}{
		{f.top, l - 1, t - 1, w + 2, 1},
		{f.bottom, l - 1, t + h, w + 2, 1},
		{f.left, l - 1, t, 1, h},
		{f.right, l + w, t, 1, h},
	}
	for _, item := range panels {
		item.p.SetBounds(item.x, item.y, item.w, item.h)
		item.p.BringToFront()
		item.p.SetVisible(true)
	}
}

func (f *designerSelectionFrame) hide() {
	if !f.valid() {
		return
	}
	for _, p := range []lcl.IPanel{f.top, f.bottom, f.left, f.right} {
		if p.Visible() {
			p.SetVisible(false)
		}
		p.SetBounds(0, 0, 0, 0)
	}
}

// IsComponentSelected 组件是否在选中集合内
func (m *FormTab) IsComponentSelected(comp *TDesigningComponent) bool {
	if comp == nil || m.selected == nil {
		return false
	}
	_, ok := m.selected[comp]
	return ok
}

// PrimarySelected 返回主选中组件
func (m *FormTab) PrimarySelected() *TDesigningComponent {
	return m.primary
}

// SelectedComponents 按组件树顺序返回所有选中组件
func (m *FormTab) SelectedComponents() []*TDesigningComponent {
	if len(m.selected) == 0 || m.FormRoot == nil {
		return nil
	}
	result := make([]*TDesigningComponent, 0, len(m.selected))
	var walk func(comp *TDesigningComponent)
	walk = func(comp *TDesigningComponent) {
		if _, ok := m.selected[comp]; ok {
			result = append(result, comp)
		}
		for _, child := range comp.Child {
			walk(child)
		}
	}
	walk(m.FormRoot)
	return result
}

func (m *FormTab) firstSelectedInTree(set map[*TDesigningComponent]struct{}) *TDesigningComponent {
	if len(set) == 0 || m.FormRoot == nil {
		return nil
	}
	var walk func(comp *TDesigningComponent) *TDesigningComponent
	walk = func(comp *TDesigningComponent) *TDesigningComponent {
		if _, ok := set[comp]; ok {
			return comp
		}
		for _, child := range comp.Child {
			if found := walk(child); found != nil {
				return found
			}
		}
		return nil
	}
	return walk(m.FormRoot)
}

// selectComponents 更新选中集合。
// syncTree 为 true 时同步项目树多选状态。
func (m *FormTab) selectComponents(comps []*TDesigningComponent, primary *TDesigningComponent, syncTree bool) {
	next := make(map[*TDesigningComponent]struct{}, len(comps))
	for _, comp := range comps {
		if comp == nil || comp.FormTab != m {
			continue
		}
		next[comp] = struct{}{}
	}
	if len(next) == 0 {
		primary = nil
	} else if primary == nil || !m.IsComponentSelectedIn(next, primary) {
		primary = m.firstSelectedInTree(next)
	}
	for comp := range m.selected {
		if _, ok := next[comp]; !ok {
			comp.HideDesignHelpers()
		}
	}
	m.selected = next
	m.primary = primary
	if primary != nil {
		primary.mustComponentPropertyPage()
		if primary.drag != nil {
			primary.drag.mustDS()
		}
		if m.State == FtsHide {
			designer.tab.HideAllActivated()
			m.ShowTabPage()
		}
		if !m.IsDesigner {
			designer.tab.HideAllActivated()
			designer.ActiveFormTab(m)
		}
		primary.ShowDesignHelpers()
		primary.LoadPropertyToInspector()
	}
	for comp := range next {
		if comp != primary {
			comp.HideDesignHelpers()
		}
	}
	m.syncSelectionFrames()
	if syncTree {
		ProjectTreeSyncSelection(m)
	}
	if m.formDesigner != nil && m.formDesigner.Form != nil && m.formDesigner.Form.IsValid() {
		m.formDesigner.Form.SetFocus()
	}
}

func (m *FormTab) IsComponentSelectedIn(set map[*TDesigningComponent]struct{}, comp *TDesigningComponent) bool {
	_, ok := set[comp]
	return ok
}

// SelectComponents 选中一组组件，primary 为主选中组件
func (m *FormTab) SelectComponents(comps []*TDesigningComponent, primary *TDesigningComponent) {
	m.selectComponents(comps, primary, true)
}

// ClearSelection 清空选择
func (m *FormTab) ClearSelection() {
	m.selectComponents(nil, nil, true)
}

// AddComponentSelection 追加一个组件到选择集（Shift 点击）
func (m *FormTab) AddComponentSelection(comp *TDesigningComponent) {
	if comp == nil || comp.FormTab != m || m.IsComponentSelected(comp) {
		return
	}
	comps := m.SelectedComponents()
	comps = append(comps, comp)
	m.selectComponents(comps, comp, true)
}

// ToggleComponentSelection 切换选中状态（Ctrl 点击），返回切换后是否选中
func (m *FormTab) ToggleComponentSelection(comp *TDesigningComponent) bool {
	if comp == nil || comp.FormTab != m {
		return false
	}
	comps := m.SelectedComponents()
	if m.IsComponentSelected(comp) {
		filtered := make([]*TDesigningComponent, 0, len(comps))
		primary := m.primary
		if primary == comp {
			primary = nil
		}
		for _, item := range comps {
			if item != comp {
				filtered = append(filtered, item)
			}
		}
		m.selectComponents(filtered, primary, true)
		return false
	}
	comps = append(comps, comp)
	m.selectComponents(comps, comp, true)
	return true
}

// RemoveComponentFromSelection 从选择集中移除一个组件
func (m *FormTab) RemoveComponentFromSelection(comp *TDesigningComponent) {
	if !m.IsComponentSelected(comp) {
		return
	}
	comps := m.SelectedComponents()
	primary := m.primary
	if primary == comp {
		primary = nil
	}
	filtered := make([]*TDesigningComponent, 0, len(comps))
	for _, item := range comps {
		if item != comp {
			filtered = append(filtered, item)
		}
	}
	m.selectComponents(filtered, primary, true)
}

// SelectAllComponents 全选设计组件（不含窗体根节点和非可视组件）
func (m *FormTab) SelectAllComponents() {
	var comps []*TDesigningComponent
	var walk func(comp *TDesigningComponent)
	walk = func(comp *TDesigningComponent) {
		for _, child := range comp.Child {
			if child.ComponentType != consts.CtNonVisual {
				control := child.Control()
				if control != nil && control.IsValid() && control.Visible() {
					comps = append(comps, child)
				}
			}
			walk(child)
		}
	}
	if m.FormRoot != nil {
		walk(m.FormRoot)
	}
	if len(comps) == 0 {
		return
	}
	m.selectComponents(comps, comps[0], true)
}

// HasSelectedAncestor 组件是否存在被选中的祖先（多选拖动时避免父子重复位移）
func (m *FormTab) HasSelectedAncestor(comp *TDesigningComponent) bool {
	if comp == nil {
		return false
	}
	for parent := comp.parent; parent != nil; parent = parent.parent {
		if m.IsComponentSelected(parent) {
			return true
		}
	}
	return false
}

// IsComponentAncestor 判断 ancestor 是否为 comp 的祖先
func (m *FormTab) IsComponentAncestor(ancestor, comp *TDesigningComponent) bool {
	if ancestor == nil || comp == nil {
		return false
	}
	for parent := comp.parent; parent != nil; parent = parent.parent {
		if parent == ancestor {
			return true
		}
	}
	return false
}

// syncSelectionFrames 按当前选择刷新非主选组件的选中边框
func (m *FormTab) syncSelectionFrames() {
	box := m.designerBox
	if box == nil || !box.IsValid() || m.formDesigner == nil {
		return
	}
	var secondary []*TDesigningComponent
	for _, comp := range m.SelectedComponents() {
		if comp == m.primary || comp == m.FormRoot || comp.ComponentType == consts.CtNonVisual {
			continue
		}
		secondary = append(secondary, comp)
	}
	if m.selFrameBox != box.Instance() {
		m.selFrames = nil
		m.selFrameBox = box.Instance()
	}
	for len(m.selFrames) < len(secondary) {
		m.selFrames = append(m.selFrames, newDesignerSelectionFrame(box))
	}
	for i, frame := range m.selFrames {
		if i >= len(secondary) || !frame.valid() {
			frame.hide()
			continue
		}
		control := secondary[i].Control()
		if control == nil || !control.IsValid() {
			frame.hide()
			continue
		}
		frame.show(m.formDesigner.GetRelativeBounds(secondary[i]))
	}
}

// hideSelectionFrames 隐藏所有选中边框
func (m *FormTab) hideSelectionFrames() {
	for _, frame := range m.selFrames {
		frame.hide()
	}
}
