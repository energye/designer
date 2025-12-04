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
	icon lcl.IImage
	text lcl.ILabel
	comp *TDesigningComponent
}

func NewNonVisualComponentWrap(owner lcl.IWinControl, comp *TDesigningComponent) *TNonVisualComponentWrap {
	m := new(TNonVisualComponentWrap)
	icon := lcl.NewImage(owner)
	icon.SetAlign(types.AlCustom)
	icon.SetWidth(nonWrapW)
	icon.SetHeight(nonWrapH)
	icon.SetImages(imageComponents.ImageList150())
	text := lcl.NewLabel(owner)
	m.icon = icon
	m.text = text
	m.comp = comp
	return m
}
func (m *TNonVisualComponentWrap) Free() {
	m.comp = nil
}

func (m *TNonVisualComponentWrap) TextFollowHide() {
	m.text.SetVisible(false)
}

func (m *TNonVisualComponentWrap) SetImage() {
	m.icon.SetImageIndex(m.comp.IconIndex())
}

func (m *TNonVisualComponentWrap) TextFollowShow() {
	caption := m.comp.Name()
	m.text.SetCaption(caption)
	br := m.icon.BoundsRect()
	textWidth := m.text.Canvas().TextWidthWithUnicodestring(caption)
	x := br.Left + br.Width()/2
	y := br.Top + br.Height()
	m.text.SetLeft(x - textWidth/2)
	m.text.SetTop(y)
	m.text.SetVisible(true)
}

func (m *TNonVisualComponentWrap) SetHint(hint string) {
	m.icon.SetHint(hint)
}

func (m *TNonVisualComponentWrap) SetParent(parent lcl.IWinControl) {
	m.icon.SetParent(parent)
	m.text.SetParent(parent)
}

func (m *TNonVisualComponentWrap) ClientToParent(point types.TPoint, parent lcl.IWinControl) types.TPoint {
	return m.icon.ClientToParent(point, parent)
}

func (m *TNonVisualComponentWrap) SetLeftTop(x, y int32) {
	m.icon.SetBounds(x, y, nonWrapW, nonWrapH)
}

func (m *TNonVisualComponentWrap) BoundsRect() types.TRect {
	return m.icon.BoundsRect()
}

func (m *TNonVisualComponentWrap) Instance() uintptr {
	return m.icon.Instance()
}
