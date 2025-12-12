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

// 组件设计 组件的包裹

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
	if m.comp != nil {
		m.comp = nil
		m.icon.Free()
		m.text.Free()
	}
}

func (m *TNonVisualComponentWrap) TextFollowHide() {
	//m.text.SetVisible(false) // 注释: 因为需要一个中间组件才起作用
	m.text.SetCaption("") // 设计模式使用空字串来隐藏
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

func (m *TNonVisualComponentWrap) SetParent(parent lcl.IWinControl) {
	m.icon.SetParent(parent)
	m.text.SetParent(parent)
}

func (m *TNonVisualComponentWrap) SetLeftTop(x, y int32) {
	m.icon.SetBounds(x, y, nonWrapW, nonWrapH)
}

func (m *TNonVisualComponentWrap) Instance() uintptr {
	return m.icon.Instance()
}
