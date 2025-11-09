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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 组件设计 非呈现组件的包裹

var (
	nonWrapW, nonWrapH int32 = 38, 38
)

// 非可视化组件
type TNonVisualComponentWrap struct {
	wrap lcl.IPanel
	icon lcl.IImage
	text lcl.ILabel
	comp *TDesigningComponent
}

func NewNonVisualComponentWrap(owner lcl.IWinControl, comp *TDesigningComponent) *TNonVisualComponentWrap {
	m := new(TNonVisualComponentWrap)
	wrap := lcl.NewPanel(owner)
	wrap.SetWidth(nonWrapW)
	wrap.SetHeight(nonWrapH)
	wrap.SetCursor(types.CrSize)
	wrap.SetShowHint(true)
	icon := lcl.NewImage(owner)
	icon.SetAlign(types.AlClient)
	icon.SetImages(imageComponents.ImageList150())
	icon.SetCursor(types.CrSize)
	icon.SetParent(wrap)
	text := lcl.NewLabel(owner)
	m.wrap = wrap
	m.icon = icon
	m.text = text
	m.comp = comp
	return m
}
func (m *TNonVisualComponentWrap) Free() {
	m.wrap.Free()
	m.text.Free()
	m.icon.Free()
	m.comp = nil
}

func (m *TNonVisualComponentWrap) TextFollowHide() {
	m.text.SetVisible(false)
}

func (m *TNonVisualComponentWrap) TextFollowShow() {
	m.icon.SetImageIndex(m.comp.IconIndex())
	caption := m.comp.Name()
	m.text.SetCaption(caption)
	br := m.wrap.BoundsRect()
	textWidth := m.text.Canvas().TextWidthWithUnicodestring(caption)
	x := br.Left + br.Width()/2
	y := br.Top + br.Height()
	m.text.SetLeft(x - textWidth/2)
	m.text.SetTop(y)
	m.text.SetVisible(true)
}

func (m *TNonVisualComponentWrap) SetHint(hint string) {
	m.wrap.SetHint(hint)
}

func (m *TNonVisualComponentWrap) SetParent(parent lcl.IWinControl) {
	m.wrap.SetParent(parent)
	m.text.SetParent(parent)
}

func (m *TNonVisualComponentWrap) ClientToParent(point types.TPoint, parent lcl.IWinControl) types.TPoint {
	return m.wrap.ClientToParent(point, parent)
}

func (m *TNonVisualComponentWrap) SetLeftTop(x, y int32) {
	m.wrap.SetBounds(x, y, nonWrapW, nonWrapH)
}

func (m *TNonVisualComponentWrap) BoundsRect() types.TRect {
	return m.wrap.BoundsRect()
}

func (m *TNonVisualComponentWrap) Instance() uintptr {
	return m.icon.Instance()
}
