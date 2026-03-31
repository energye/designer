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

package options

import (
	"github.com/energye/designer/designer"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
	"strings"
)

func SetWindowCenterByMainWindow(window lcl.IEngForm) {
	SetWindowCenterByRelativeWindow(window, &designer.MainWindow)
}

func SetWindowCenterByRelativeWindow(window, relativeWindow lcl.IEngForm) {
	windowRect := relativeWindow.BoundsRect()
	window.SetLeft(windowRect.Left + (windowRect.Width()-window.Width())/2)
	window.SetTop(windowRect.Top + (windowRect.Height()-window.Height())/2)
}

type TCommonMemoForm struct {
	baseForm       *lcl.TEngForm
	box            lcl.IPanel
	memo           lcl.IMemo
	w, h           int32
	title          string
	relativeWindow lcl.IEngForm
	isDemoBtn      bool
	demoBtn        *wg.TButton
	onOK           func(lines []string)
}

func NewCommonMemoForm(w, h int32, title string, relativeWindow lcl.IEngForm) *TCommonMemoForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TCommonMemoForm{baseForm: newEngForm.(*lcl.TEngForm), w: w, h: h, title: title,
		relativeWindow: relativeWindow}
	newForm.FormCreate(newEngForm)
	return newForm
}

func (m *TCommonMemoForm) SetDemoText(text string) {
	if text != "" && m.demoBtn == nil {
		demoBtnRect := types.TRect{Left: 7, Top: m.memo.Height() + m.memo.Top() + 3}
		demoBtnRect.SetWidth(50)
		demoBtnRect.SetHeight(25)
		m.demoBtn = wg.NewButton(m.baseForm)
		m.demoBtn.SetBoundsRect(demoBtnRect)
		m.demoBtn.SetText("示 例")
		m.demoBtn.SetRadius(3)
		m.demoBtn.SetCursor(types.CrHandPoint)
		m.demoBtn.SetColor(colors.RGBToColor(212, 212, 212))
		m.demoBtn.SetParent(m.box)
		m.demoBtn.SetOnClick(m.demoBtnClick)

		textLines := strings.Split(text, "\n")
		demoTextMemo := lcl.NewMemo(m.baseForm)
		demoTextMemo.Font().SetSize(8)
		lines := demoTextMemo.Lines()
		for _, line := range textLines {
			if line == "" {
				continue
			}
			lines.Add(line)
		}
		memoRect := m.memo.BoundsRect()
		demoTextMemo.SetBorderStyle(types.BsNone)
		demoTextMemo.SetColor(colors.ClInfoBk)
		demoTextMemo.SetReadOnly(true)
		demoTextMemo.SetBounds(0, memoRect.Height()+memoRect.Top+33, m.w, 98)
		demoTextMemo.SetParent(m.box)
	}
}

func (m *TCommonMemoForm) SetOnOK(fn func(lines []string)) {
	m.onOK = fn
}

func (m *TCommonMemoForm) SetDefaultText(text string) {
	textLines := strings.Split(text, "\n")
	lines := m.memo.Lines()
	lines.Clear()
	if m.memo.WantReturns() {
		for _, line := range textLines {
			lines.Add(line)
		}
	} else {
		lines.SetTextToStr(text)
	}
}

func (m *TCommonMemoForm) SetMultipleLine(v bool) {
	m.memo.SetWantReturns(v)
	m.memo.SetWordWrap(v)
}

func (m *TCommonMemoForm) FormCreate(sender lcl.IObject) {
	m.baseForm.SetWidth(m.w)
	m.baseForm.SetHeight(m.h)
	m.baseForm.SetBorderStyleToFormBorderStyle(types.BsNone)
	SetWindowCenterByRelativeWindow(m.baseForm, m.relativeWindow)
	m.box = lcl.NewPanel(m.baseForm)
	m.box.SetAlign(types.AlClient)
	m.box.SetColor(colors.ClWhite)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetBorderStyleToBorderStyle(types.BsSingle)
	m.box.SetParent(m.baseForm)

	memoTop := int32(0)
	if m.title != "" {
		titleLbl := lcl.NewLabel(m.baseForm)
		titleFont := titleLbl.Font()
		titleFont.SetSize(10)
		titleLbl.SetCaption(m.title)
		titleLbl.SetTop(7)
		titleLbl.SetLeft(7)
		titleLbl.SetParent(m.box)
		memoTop = 35
	}

	m.memo = lcl.NewMemo(m.baseForm)
	m.memo.SetBounds(0, memoTop, m.w, m.h-70)
	m.memo.SetBorderStyle(types.BsNone)
	m.memo.SetColor(colors.ClInfoBk)
	m.memo.Font().SetSize(8)
	m.memo.SetParent(m.box)

	cancelBtnRect := types.TRect{Left: m.w - (50 * 2) - 30, Top: m.memo.Height() + m.memo.Top() + 5}
	cancelBtnRect.SetWidth(50)
	cancelBtnRect.SetHeight(25)
	cancelBtn := wg.NewButton(m.baseForm)
	cancelBtn.SetBoundsRect(cancelBtnRect)
	cancelBtn.SetText("关 闭")
	cancelBtn.Font().SetColor(colors.ClWhite)
	cancelBtn.SetRadius(3)
	cancelBtn.SetCursor(types.CrHandPoint)
	cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
	cancelBtn.SetParent(m.box)
	cancelBtn.SetOnClick(func(sender lcl.IObject) {
		m.baseForm.Close()
	})

	okBtnRect := types.TRect{Left: m.w - (50 * 2) + 40, Top: cancelBtnRect.Top}
	okBtnRect.SetWidth(50)
	okBtnRect.SetHeight(25)
	okBtn := wg.NewButton(m.baseForm)
	okBtn.SetBoundsRect(okBtnRect)
	okBtn.SetText("确 定")
	okBtn.Font().SetColor(colors.ClWhite)
	okBtn.SetRadius(3)
	okBtn.SetCursor(types.CrHandPoint)
	okBtn.SetColor(colors.RGBToColor(46, 204, 113))
	okBtn.SetParent(m.box)
	okBtn.SetOnClick(func(sender lcl.IObject) {
		if m.onOK != nil {
			memoLines := m.memo.Lines()
			count := memoLines.Count()
			var lines []string
			for i := int32(0); i < count; i++ {
				lines = append(lines, memoLines.Strings(i))
			}
			m.onOK(lines)
		}
		m.baseForm.Close()
	})
}

func (m *TCommonMemoForm) ShowModal() {
	m.baseForm.ShowModal()
}

func (m *TCommonMemoForm) demoBtnClick(sender lcl.IObject) {
	m.isDemoBtn = !m.isDemoBtn
	rect := m.baseForm.BoundsRect()
	if m.isDemoBtn {
		rect.SetHeight(rect.Height() + 100)
	} else {
		rect.SetHeight(rect.Height() - 100)
	}
	m.baseForm.SetBoundsRect(rect)
}

type TCommonMemoBox struct {
	parent lcl.IWinControl
	box    lcl.IPanel
	memo   lcl.IMemo
	rect   types.TRect
	title  string
	change func(lines []string)
}

func NewCommonMemoBox(rect types.TRect, title string, parent lcl.IWinControl) *TCommonMemoBox {
	newMemo := &TCommonMemoBox{rect: rect, title: title, parent: parent}
	newMemo.Create()
	return newMemo
}

func (m *TCommonMemoBox) SetOnChange(fn func(lines []string)) {
	m.change = fn
}

func (m *TCommonMemoBox) Create() {
	m.box = lcl.NewPanel(m.parent)
	m.box.SetBoundsRect(m.rect)
	m.box.SetColor(m.parent.Color())
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetBorderStyleToBorderStyle(types.BsNone)
	m.box.SetParent(m.parent)

	memoTop := int32(0)
	if m.title != "" {
		titleLbl := lcl.NewLabel(m.parent)
		titleFont := titleLbl.Font()
		titleFont.SetSize(10)
		titleLbl.SetCaption(m.title)
		titleLbl.SetTop(7)
		titleLbl.SetLeft(7)
		titleLbl.SetParent(m.box)
		memoTop += 35
	}
	memoRect := types.TRect{Left: 0, Top: memoTop}
	memoRect.SetWidth(m.rect.Width())
	memoRect.SetHeight(m.rect.Height() - (memoRect.Top + m.rect.Top))
	m.memo = lcl.NewMemo(m.parent)
	m.memo.SetBoundsRect(memoRect)
	m.memo.SetBorderStyle(types.BsNone)
	m.memo.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	m.memo.SetColor(colors.ClInfoBk)
	m.memo.Font().SetSize(8)
	m.memo.SetOnChange(func(sender lcl.IObject) {
		if m.change != nil {
			m.change(m.Lines())
		}
	})
}

func (m *TCommonMemoBox) SetAnchors(v types.TAnchors) {
	m.box.SetAnchors(v)
}

func (m *TCommonMemoBox) SetMultipleLine(v bool) {
	m.memo.SetWantReturns(v)
	m.memo.SetWordWrap(v)
}

func (m *TCommonMemoBox) SetDefaultText(text string) {
	textLines := strings.Split(text, "\n")
	lines := m.memo.Lines()
	lines.Clear()
	if m.memo.WantReturns() {
		for _, line := range textLines {
			lines.Add(line)
		}
	} else {
		lines.SetTextToStr(text)
	}
}

func (m *TCommonMemoBox) SetDemoText(text string) {
	if text != "" {
		textLines := strings.Split(text, "\n")
		demoTextMemo := lcl.NewMemo(m.parent)
		demoTextMemo.Font().SetSize(8)
		lines := demoTextMemo.Lines()
		for _, line := range textLines {
			if line == "" {
				continue
			}
			lines.Add(line)
		}
		height := int32(100)
		memoRect := m.memo.BoundsRect()
		demoTextMemoRect := types.TRect{Left: 0, Top: (memoRect.Height() + memoRect.Top) - height}
		demoTextMemoRect.SetWidth(memoRect.Width())
		demoTextMemoRect.SetHeight(100)
		demoTextMemo.SetBorderStyle(types.BsNone)
		demoTextMemo.SetColor(colors.ClInfoBk)
		demoTextMemo.SetReadOnly(true)
		demoTextMemo.SetAnchors(types.NewSet(types.AkLeft, types.AkRight, types.AkBottom))
		demoTextMemo.SetBoundsRect(demoTextMemoRect)
		demoTextMemo.SetParent(m.box)

		m.memo.SetHeight(memoRect.Height() - demoTextMemoRect.Height())
	}
}

func (m *TCommonMemoBox) Show() {
	m.memo.SetParent(m.box)
}

func (m *TCommonMemoBox) Lines() []string {
	memoLines := m.memo.Lines()
	count := memoLines.Count()
	var lines []string
	for i := int32(0); i < count; i++ {
		lines = append(lines, memoLines.Strings(i))
	}
	return lines
}
