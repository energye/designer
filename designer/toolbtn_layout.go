package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

type TToolBtnLayout struct {
	box lcl.IPanel
}

func initToolBtnLayout(owner lcl.IWinControl) *TToolBtnLayout {
	m := &TToolBtnLayout{}
	m.box = lcl.NewPanel(owner)
	m.box.SetColor(colors.ClLegacySkyBlue)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlTop)
	m.box.SetHeight(40)
	m.box.SetParent(owner)
	return m
}
