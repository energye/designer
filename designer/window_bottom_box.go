package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 下面设计器

type BottomBox struct {
	lcl.IPanel
	left  lcl.IPanel
	right lcl.IPanel
}

func (m *TAppWindow) CreateBottomBox() *BottomBox {
	box := &BottomBox{}
	box.IPanel = lcl.NewPanel(m)
	box.IPanel.SetParent(m)
	box.IPanel.SetBevelOuter(types.BvNone)
	box.IPanel.SetDoubleBuffered(true)
	box.IPanel.SetTop(toolbarHeight)
	box.IPanel.SetWidth(m.Width())
	box.IPanel.SetHeight(m.Height() - toolbarHeight)
	box.IPanel.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	box.IPanel.SetColor(colors.RGBToColor(100, 120, 140))
	m.box = box
	return box
}
