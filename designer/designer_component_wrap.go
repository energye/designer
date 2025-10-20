package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 组件设计 非呈现组件的包裹

var (
	nonWrapW, nonWrapH int32 = 38, 38
)

type NonVisualComponentWrap struct {
	wrap lcl.IPanel
	icon lcl.IImage
	text lcl.ILabel
	comp *DesigningComponent
}

func NewNonVisualComponentWrap(owner lcl.IWinControl, comp *DesigningComponent) *NonVisualComponentWrap {
	m := new(NonVisualComponentWrap)
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

func (m *NonVisualComponentWrap) TextFollowHide() {
	m.text.SetVisible(false)
}

func (m *NonVisualComponentWrap) TextFollowShow() {
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

func (m *NonVisualComponentWrap) SetHint(hint string) {
	m.wrap.SetHint(hint)
}

func (m *NonVisualComponentWrap) SetParent(parent lcl.IWinControl) {
	m.wrap.SetParent(parent)
	m.text.SetParent(parent)
}

func (m *NonVisualComponentWrap) ClientToParent(point types.TPoint, parent lcl.IWinControl) types.TPoint {
	return m.wrap.ClientToParent(point, parent)
}

func (m *NonVisualComponentWrap) SetLeftTop(x, y int32) {
	m.wrap.SetBounds(x, y, nonWrapW, nonWrapH)
}

func (m *NonVisualComponentWrap) BoundsRect() types.TRect {
	return m.wrap.BoundsRect()
}

func (m *NonVisualComponentWrap) Instance() uintptr {
	return m.icon.Instance()
}
