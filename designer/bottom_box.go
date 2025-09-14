package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 下面设计器

type BottomBox struct {
	lcl.IPanel
	leftBox  lcl.IPanel    // 左侧-面板
	spliter  lcl.ISplitter // 分割线
	rightBox lcl.IPanel    // 右侧-设计器主体
}

func (m *TAppWindow) createBottomBox() *BottomBox {
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
	// 左侧-面板

	// 右侧-设计器主体

	return box
}
