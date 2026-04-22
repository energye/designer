package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type ContentLayoutProject struct {
	topBox lcl.IPanel
	title  lcl.ILabel
	box    lcl.IPanel
}

func initContentLayoutProject(owner *ContentLayout) *ContentLayoutProject {
	m := &ContentLayoutProject{}
	m.topBox = lcl.NewPanel(owner.projectPanel)
	m.topBox.SetBorderStyleToBorderStyle(types.BsNone)
	m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	m.topBox.SetParent(owner.projectPanel)

	title := lcl.NewLabel(m.topBox)
	title.SetCaption("项目管理器")
	title.SetLeft(5)
	title.SetTop(5)
	font := title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	title.SetParent(m.topBox)

	m.box = lcl.NewPanel(owner.projectPanel)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlClient)
	m.box.SetParent(owner.projectPanel)

	return m
}
