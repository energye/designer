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
// 主选中组件显示拖拽控制点和属性面板，其余选中组件的 1px 蓝色边框
// 由设计画布 OnPaint 直接绘制。不再为每个次选组件创建 4 个覆盖面板:
// 选中数量很大时, 每次刷新对数百个面板做 SetBounds/BringToFront/SetVisible
// 会拖垮 UI 线程 (98 个次选组件实测 4-10 秒), 画布绘制只需一次 Invalidate。

var designerSelectionColor = colors.RGBToColor(0, 125, 255)

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
	if m.formDesignPage != nil && m.formDesignPage.alignBar != nil {
		m.formDesignPage.alignBar.updateState()
	}
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

// beginSelectionDrag 拖动开始时暂停绘制次选边框。
// 拖动中组件逐帧移动，边框跟随绘制会在画布残留旧位置色带；拖动结束统一恢复。
func (m *FormTab) beginSelectionDrag() {
	if m.selectionDragActive {
		return
	}
	m.selectionDragActive = true
	// 拖动期间隐藏覆盖层(边框/控制点), 结束时统一恢复
	if m.overlayPane != nil && m.overlayPane.IsValid() && m.overlayPane.Visible() {
		m.overlayPane.SetVisible(false)
	}
}

// endSelectionDrag 拖动结束后统一恢复一次边框绘制。
func (m *FormTab) endSelectionDrag() {
	if !m.selectionDragActive {
		return
	}
	m.selectionDragActive = false
	m.syncSelectionFrames()
}

// syncSelectionFrames 按当前选择刷新次选边框与主选控制点。
// 绘制与命中都在页签级覆盖层(玻璃层)上, 这里重建覆盖层窗口区域并重绘;
// 无窗口区域平台回退容器画布绘制(paintContainerOverlays)。
func (m *FormTab) syncSelectionFrames() {
	if m.selectionDragActive {
		return
	}
	m.updateOverlayPane()
}

// hideSelectionFrames 清除画布上已绘制的选中边框
func (m *FormTab) hideSelectionFrames() {
	if box := m.designerBox; box != nil && box.IsValid() {
		box.Invalidate()
	}
}
