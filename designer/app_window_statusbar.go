package designer

import (
	"github.com/energye/lcl/lcl"
)

type StatusBar struct {
	statusBar lcl.IStatusBar
	left      lcl.IStatusPanel
	right     lcl.IStatusPanel
}

func newStatusBar(owner lcl.IWinControl) *StatusBar {
	bar := new(StatusBar)
	statusBar := lcl.NewStatusBar(owner)
	statusBar.SetParent(owner)
	//statusBar.SetSimplePanel(false)
	statusBar.SetAutoHint(true)
	bar.statusBar = statusBar
	return bar
}
