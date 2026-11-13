// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package designer

import (
	"math"
	"sort"

	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 设计 - 多选对齐/分布/同尺寸

type designerAlignKind int32

const (
	designerAlignLeft designerAlignKind = iota
	designerAlignHCenter
	designerAlignRight
	designerAlignTop
	designerAlignVCenter
	designerAlignBottom
)

type designerSameSizeKind int32

const (
	designerSameWidth designerSameSizeKind = iota
	designerSameHeight
	designerSameSize
)

type designerBoundsChange struct {
	comp *TDesigningComponent
	rect types.TRect
}

// designerAlignRect 计算 target 对齐到 ref 后的目标矩形（窗体客户区坐标）。
func designerAlignRect(kind designerAlignKind, ref, target types.TRect) types.TRect {
	width, height := target.Width(), target.Height()
	switch kind {
	case designerAlignLeft:
		target.Left = ref.Left
	case designerAlignHCenter:
		target.Left = ref.Left + (ref.Width()-width)/2
	case designerAlignRight:
		target.Left = ref.Right - width
	case designerAlignTop:
		target.Top = ref.Top
	case designerAlignVCenter:
		target.Top = ref.Top + (ref.Height()-height)/2
	case designerAlignBottom:
		target.Top = ref.Bottom - height
	}
	target.Right = target.Left + width
	target.Bottom = target.Top + height
	return target
}

// designerSameSizeRect 计算 target 使用 ref 尺寸后的目标矩形，位置保持不变。
func designerSameSizeRect(kind designerSameSizeKind, ref, target types.TRect) types.TRect {
	if kind == designerSameWidth || kind == designerSameSize {
		target.Right = target.Left + ref.Width()
	}
	if kind == designerSameHeight || kind == designerSameSize {
		target.Bottom = target.Top + ref.Height()
	}
	return target
}

// designerDistributePositions 计算等间距分布后各矩形的前导坐标。
// rects 必须已按左/上坐标升序排列；首尾矩形保持不动。
func designerDistributePositions(rects []types.TRect, horizontal bool) []int32 {
	n := len(rects)
	positions := make([]int32, n)
	if n == 0 {
		return positions
	}
	lead := func(r types.TRect) int32 {
		if horizontal {
			return r.Left
		}
		return r.Top
	}
	size := func(r types.TRect) int32 {
		if horizontal {
			return r.Width()
		}
		return r.Height()
	}
	for i, r := range rects {
		positions[i] = lead(r)
	}
	if n < 3 {
		return positions
	}
	start := lead(rects[0])
	end := lead(rects[n-1]) + size(rects[n-1])
	total := int32(0)
	for _, r := range rects {
		total += size(r)
	}
	gap := float64(end-start-total) / float64(n-1)
	prefix := float64(0)
	for i, r := range rects {
		switch i {
		case 0:
			positions[i] = start
		case n - 1:
			positions[i] = lead(rects[n-1])
		default:
			positions[i] = int32(math.Round(float64(start) + prefix + gap*float64(i)))
		}
		prefix += float64(size(r))
	}
	return positions
}

func designerComponentDepth(comp *TDesigningComponent) int {
	depth := 0
	for parent := comp.parent; parent != nil; parent = parent.parent {
		depth++
	}
	return depth
}

// alignableSelectedComponents 返回可参与对齐的选中组件，忽略窗体根节点和非可视组件。
func (m *FormTab) alignableSelectedComponents() []*TDesigningComponent {
	comps := m.SelectedComponents()
	result := make([]*TDesigningComponent, 0, len(comps))
	for _, comp := range comps {
		if comp == nil || comp == m.FormRoot || comp.ComponentType == consts.CtNonVisual {
			continue
		}
		control := comp.Control()
		if control == nil || !control.IsValid() {
			continue
		}
		result = append(result, comp)
	}
	return result
}

func (m *FormTab) alignReference(comps []*TDesigningComponent) *TDesigningComponent {
	if len(comps) == 0 {
		return nil
	}
	if m.primary != nil {
		for _, comp := range comps {
			if comp == m.primary {
				return comp
			}
		}
	}
	return comps[0]
}

// AlignSelected 将选中组件按主选组件对齐。
func (m *FormTab) AlignSelected(kind designerAlignKind) {
	if m.formDesigner == nil || m.formDesigner.Form == nil {
		return
	}
	comps := m.alignableSelectedComponents()
	if len(comps) < 2 {
		return
	}
	ref := m.alignReference(comps)
	refRect := m.formDesigner.GetRelativeBounds(ref)
	changes := make([]designerBoundsChange, 0, len(comps))
	for _, comp := range comps {
		rect := m.formDesigner.GetRelativeBounds(comp)
		changes = append(changes, designerBoundsChange{comp: comp, rect: designerAlignRect(kind, refRect, rect)})
	}
	m.applyDesignerBounds(changes)
}

// SameSizeSelected 将选中组件的宽/高/尺寸统一为主选组件的尺寸。
func (m *FormTab) SameSizeSelected(kind designerSameSizeKind) {
	if m.formDesigner == nil || m.formDesigner.Form == nil {
		return
	}
	comps := m.alignableSelectedComponents()
	if len(comps) < 2 {
		return
	}
	ref := m.alignReference(comps)
	refRect := m.formDesigner.GetRelativeBounds(ref)
	changes := make([]designerBoundsChange, 0, len(comps))
	for _, comp := range comps {
		rect := m.formDesigner.GetRelativeBounds(comp)
		changes = append(changes, designerBoundsChange{comp: comp, rect: designerSameSizeRect(kind, refRect, rect)})
	}
	m.applyDesignerBounds(changes)
}

// DistributeSelected 在首尾组件之间等间距分布选中组件。
func (m *FormTab) DistributeSelected(horizontal bool) {
	if m.formDesigner == nil || m.formDesigner.Form == nil {
		return
	}
	comps := m.alignableSelectedComponents()
	if len(comps) < 3 {
		return
	}
	items := make([]designerBoundsChange, 0, len(comps))
	for _, comp := range comps {
		items = append(items, designerBoundsChange{comp: comp, rect: m.formDesigner.GetRelativeBounds(comp)})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if horizontal {
			return items[i].rect.Left < items[j].rect.Left
		}
		return items[i].rect.Top < items[j].rect.Top
	})
	rects := make([]types.TRect, len(items))
	for i := range items {
		rects[i] = items[i].rect
	}
	positions := designerDistributePositions(rects, horizontal)
	changes := make([]designerBoundsChange, 0, len(items))
	for i := range items {
		rect := items[i].rect
		width, height := rect.Width(), rect.Height()
		if horizontal {
			rect.Left = positions[i]
			rect.Right = rect.Left + width
		} else {
			rect.Top = positions[i]
			rect.Bottom = rect.Top + height
		}
		changes = append(changes, designerBoundsChange{comp: items[i].comp, rect: rect})
	}
	m.applyDesignerBounds(changes)
}

// applyDesignerBounds 按窗体客户区坐标应用一组组件边界。
func (m *FormTab) applyDesignerBounds(changes []designerBoundsChange) {
	if m.formDesigner == nil || m.formDesigner.Form == nil {
		return
	}
	// 祖先先应用，子组件再按祖先移动后的父原点计算局部坐标。
	sort.SliceStable(changes, func(i, j int) bool {
		return designerComponentDepth(changes[i].comp) < designerComponentDepth(changes[j].comp)
	})
	changed := false
	var first *TDesigningComponent
	for _, change := range changes {
		comp := change.comp
		control := comp.Control()
		if control == nil || !control.IsValid() {
			continue
		}
		br := control.BoundsRect()
		if br.Left == change.rect.Left && br.Top == change.rect.Top &&
			br.Width() == change.rect.Width() && br.Height() == change.rect.Height() {
			continue
		}
		origin := m.formDesigner.GetRelativeTopLeft(comp)
		parentLeft := origin.X - br.Left
		parentTop := origin.Y - br.Top
		left := change.rect.Left - parentLeft
		top := change.rect.Top - parentTop
		control.SetLeft(left)
		control.SetTop(top)
		// 改变 Left/Top 时 LCL 可能因锚点保持右/下边缘，必须再显式写回宽高。
		control.SetWidth(change.rect.Width())
		control.SetHeight(change.rect.Height())
		actual := control.BoundsRect()
		if comp.SyncBoundsToProperties(actual) {
			changed = true
			if first == nil {
				first = comp
			}
		}
	}
	if !changed {
		return
	}
	if m.primary != nil && m.primary.drag != nil {
		m.primary.drag.Follow()
		m.primary.drag.BringToFront()
	}
	m.syncSelectionFrames()
	triggerUIGeneration(first, nil, event.CodeGenUI)
	if m.formDesigner.Form.IsValid() {
		m.formDesigner.Form.SetFocus()
	}
}

// ==== 设计页对齐/分布工具栏 ====

type designerAlignBar struct {
	toolbar   lcl.IToolBar
	formTab   *FormTab
	alignBtns []lcl.IToolButton
	sizeBtns  []lcl.IToolButton
	distBtns  []lcl.IToolButton
}

func newDesignerAlignBar(page *TFormDesignPage) *designerAlignBar {
	if page == nil || page.formDesignPage == nil || !page.formDesignPage.IsValid() {
		return nil
	}
	bar := &designerAlignBar{formTab: page.formTab}
	toolbar := lcl.NewToolBar(page.formDesignPage)
	toolbar.SetAlign(types.AlTop)
	toolbar.SetHeight(30)
	toolbar.SetButtonWidth(24)
	toolbar.SetButtonHeight(24)
	toolbar.SetList(true)
	toolbar.SetShowCaptions(true)
	toolbar.SetEdgeBorders(types.NewSet())
	toolbar.BorderSpacing().SetLeft(3)
	toolbar.BorderSpacing().SetTop(3)
	toolbar.BorderSpacing().SetBottom(3)
	toolbar.SetParent(page.formDesignPage)
	bar.toolbar = toolbar

	newSep := func() {
		sep := lcl.NewToolButton(toolbar)
		sep.SetParent(toolbar)
		sep.SetStyle(types.TbsSeparator)
	}
	newBtn := func(caption, hint string, onClick lcl.TNotifyEvent) lcl.IToolButton {
		btn := lcl.NewToolButton(toolbar)
		btn.SetCaption(caption)
		btn.SetHint(hint)
		btn.SetShowHint(true)
		btn.SetImageIndex(-1)
		btn.SetCursor(types.CrHandPoint)
		btn.SetOnClick(onClick)
		btn.SetParent(toolbar)
		return btn
	}

	bar.alignBtns = []lcl.IToolButton{
		newBtn("左对齐", "左对齐", func(lcl.IObject) { bar.formTab.AlignSelected(designerAlignLeft) }),
		newBtn("水平居中", "水平居中对齐", func(lcl.IObject) { bar.formTab.AlignSelected(designerAlignHCenter) }),
		newBtn("右对齐", "右对齐", func(lcl.IObject) { bar.formTab.AlignSelected(designerAlignRight) }),
		newBtn("上对齐", "上对齐", func(lcl.IObject) { bar.formTab.AlignSelected(designerAlignTop) }),
		newBtn("垂直居中", "垂直居中对齐", func(lcl.IObject) { bar.formTab.AlignSelected(designerAlignVCenter) }),
		newBtn("下对齐", "下对齐", func(lcl.IObject) { bar.formTab.AlignSelected(designerAlignBottom) }),
	}
	newSep()
	bar.sizeBtns = []lcl.IToolButton{
		newBtn("同宽", "统一宽度", func(lcl.IObject) { bar.formTab.SameSizeSelected(designerSameWidth) }),
		newBtn("同高", "统一高度", func(lcl.IObject) { bar.formTab.SameSizeSelected(designerSameHeight) }),
		newBtn("同尺寸", "统一宽高", func(lcl.IObject) { bar.formTab.SameSizeSelected(designerSameSize) }),
	}
	newSep()
	bar.distBtns = []lcl.IToolButton{
		newBtn("水平分布", "水平等间距分布", func(lcl.IObject) { bar.formTab.DistributeSelected(true) }),
		newBtn("垂直分布", "垂直等间距分布", func(lcl.IObject) { bar.formTab.DistributeSelected(false) }),
	}
	bar.updateState()
	return bar
}

// updateState 根据当前选中组件数量更新按钮可用状态。
func (b *designerAlignBar) updateState() {
	if b == nil || b.toolbar == nil || !b.toolbar.IsValid() {
		return
	}
	count := 0
	if b.formTab != nil {
		count = len(b.formTab.alignableSelectedComponents())
	}
	alignEnabled := count >= 2
	distEnabled := count >= 3
	for _, btn := range b.alignBtns {
		btn.SetEnabled(alignEnabled)
	}
	for _, btn := range b.sizeBtns {
		btn.SetEnabled(alignEnabled)
	}
	for _, btn := range b.distBtns {
		btn.SetEnabled(distEnabled)
	}
}
