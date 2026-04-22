package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/widget/wg"
)

// 固定布局视图面板
type TViewPanel struct {
	lcl.IPanel
	topBox lcl.IPanel
	//bottomBox lcl.IPanel
	title    lcl.ILabel
	closeBtn *wg.TButton
}

func NewViewPanel(owner lcl.IWinControl) *TViewPanel {
	m := &TViewPanel{}
	m.IPanel = lcl.NewPanel(owner)
	m.SetBevelOuter(types.BvNone)
	m.SetDoubleBuffered(true)

	m.topBox = lcl.NewPanel(owner)
	//m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetBorderStyleToBorderStyle(types.BsSingle)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	m.topBox.SetParent(m)
	m.topBox.SetOnResize(func(sender lcl.IObject) {
		br := m.topBox.BoundsRect()
		m.closeBtn.SetLeft(br.Width() - 20)
		m.closeBtn.SetTop(8)
		//fmt.Println("topBox.SetOnResize BoundsRect:", br)
	})

	//m.bottomBox = lcl.NewPanel(owner)
	//m.bottomBox.SetBevelOuter(types.BvNone)
	//m.bottomBox.SetAlign(types.AlClient)
	//m.bottomBox.SetParent(m)

	m.title = lcl.NewLabel(owner)
	font := m.title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	m.title.SetLeft(5)
	m.title.SetTop(5)
	m.title.SetParent(m.topBox)

	m.closeBtn = wg.NewButton(owner)
	m.closeBtn.SetWidth(10)
	m.closeBtn.SetHeight(10)
	m.closeBtn.SetTop(10)
	m.closeBtn.SetParent(m.topBox)
	return m
}

func (m *TViewPanel) SetTitle(title string) {
	m.title.SetCaption(title)
}
