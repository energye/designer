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
)

// 设计画布框选

const designerBandThreshold = 4

type designerBand struct {
	down        bool
	active      bool
	shown       bool // 面板已置顶显示, 拖动期间无需重复设置 Z 序与可见性
	start       types.TPoint
	panels      [4]lcl.IPanel
	boxInstance uintptr
}

func (m *TEngFormDesigner) bandBegin(start types.TPoint) {
	m.band.down = true
	m.band.active = false
	m.band.start = start
	m.bandHide()
}

func (m *TEngFormDesigner) ensureBandPanels() bool {
	box := m.LookupRoot
	if box == nil || box.object == nil || !box.object.IsValid() {
		return false
	}
	if m.band.boxInstance == box.object.Instance() && m.band.panels[0] != nil && m.band.panels[0].IsValid() {
		return true
	}
	m.band.boxInstance = box.object.Instance()
	// 复用选择框颜色，画在控件之上
	for i := range m.band.panels {
		m.band.panels[i] = newOverlayPanel(box.object, designerSelectionColor)
	}
	return true
}

func (m *TEngFormDesigner) bandHide() {
	for _, p := range m.band.panels {
		if p == nil || !p.IsValid() {
			continue
		}
		if p.Visible() {
			p.SetVisible(false)
		}
		p.SetBounds(0, 0, 0, 0)
	}
	m.band.shown = false
}

// bandUpdate 更新框选矩形；位移小于阈值时不显示
func (m *TEngFormDesigner) bandUpdate(current types.TPoint) {
	if !m.ensureBandPanels() {
		return
	}
	box := m.LookupRoot.object
	w, h := box.ClientWidth(), box.ClientHeight()
	rect := normalizeRect(m.band.start, current)
	if !m.band.active {
		if rect.Width() < designerBandThreshold && rect.Height() < designerBandThreshold {
			return
		}
		m.band.active = true
	}
	if rect.Left < 0 {
		rect.Left = 0
	}
	if rect.Top < 0 {
		rect.Top = 0
	}
	if rect.Right > w {
		rect.Right = w
	}
	if rect.Bottom > h {
		rect.Bottom = h
	}
	rw, rh := rect.Width(), rect.Height()
	if rw < 1 || rh < 1 {
		m.bandHide()
		return
	}
	bounds := [4][4]int32{
		{rect.Left, rect.Top, rw, 1},        // top
		{rect.Left, rect.Bottom - 1, rw, 1}, // bottom
		{rect.Left, rect.Top, 1, rh},        // left
		{rect.Right - 1, rect.Top, 1, rh},   // right
	}
	for i, p := range m.band.panels {
		b := bounds[i]
		p.SetBounds(b[0], b[1], b[2], b[3])
		if !m.band.shown {
			// Z 序与可见性只在首次显示时设置: 拖动期间 Z 序不变,
			// 每次 mousemove 重复 BringToFront 会触发 z 序变更与
			// 被遮挡区域的重排重绘, 是框选卡顿的主要原因
			p.BringToFront()
			p.SetVisible(true)
		}
	}
	m.band.shown = true
}

// bandEnd 结束框选并应用选择；未产生拖动时视为点击空白画布
func (m *TEngFormDesigner) bandEnd(current types.TPoint) {
	active := m.band.active
	rect := normalizeRect(m.band.start, current)
	m.band.down = false
	m.band.active = false
	m.bandHide()

	root := m.LookupRoot
	if root == nil || root.FormTab == nil {
		return
	}
	formTab := root.FormTab
	if !active {
		// 点击空白画布：选中窗体
		formTab.SelectComponents([]*TDesigningComponent{root}, root)
		return
	}
	var comps []*TDesigningComponent
	selectedSet := make(map[*TDesigningComponent]struct{})
	var walk func(comp *TDesigningComponent)
	walk = func(comp *TDesigningComponent) {
		if comp != root && comp.ComponentType != consts.CtNonVisual {
			control := comp.Control()
			if control != nil && control.IsValid() && control.Visible() {
				if rectsIntersect(rect, m.GetRelativeBounds(comp)) && !hasSelectedAncestorInSet(comp, selectedSet) {
					comps = append(comps, comp)
					selectedSet[comp] = struct{}{}
				}
			}
		}
		for _, child := range comp.Child {
			walk(child)
		}
	}
	walk(root)
	if len(comps) == 0 {
		comps = append(comps, root)
	}
	formTab.SelectComponents(comps, comps[0])
}

// bandCancel 放弃框选
func (m *TEngFormDesigner) bandCancel() {
	m.band.down = false
	m.band.active = false
	m.bandHide()
}

func normalizeRect(a, b types.TPoint) types.TRect {
	left, right := a.X, b.X
	if left > right {
		left, right = right, left
	}
	top, bottom := a.Y, b.Y
	if top > bottom {
		top, bottom = bottom, top
	}
	return types.TRect{Left: left, Top: top, Right: right, Bottom: bottom}
}

func rectsIntersect(a, b types.TRect) bool {
	return a.Left < b.Right && b.Left < a.Right && a.Top < b.Bottom && b.Top < a.Bottom
}

func hasSelectedAncestorInSet(comp *TDesigningComponent, set map[*TDesigningComponent]struct{}) bool {
	for parent := comp.parent; parent != nil; parent = parent.parent {
		if _, ok := set[parent]; ok {
			return true
		}
	}
	return false
}
