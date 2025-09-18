package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type StatusBar struct {
	statusBar lcl.IStatusBar
	left      lcl.IStatusPanel
	right     lcl.IStatusPanel
}

func (m *TAppWindow) createStatusBar() {
	bar := new(StatusBar)
	m.bar = bar
	statusBar := lcl.NewStatusBar(m)
	statusBar.SetParent(m)
	statusBar.SetSimplePanel(false)
	statusBar.SetAutoHint(true)
	bar.statusBar = statusBar

	bar.left = statusBar.Panels().AddToStatusPanel()
	bar.left.SetWidth(200)
	bar.left.SetAlignment(types.TaCenter)

	bar.right = statusBar.Panels().AddToStatusPanel()
	bar.right.SetWidth(200)
	bar.right.SetAlignment(types.TaCenter)
}
