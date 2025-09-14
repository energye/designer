package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 顶部工具栏

type TopToolbar struct {
	lcl.IPanel
	toolbar        lcl.IPanel    // 工具按钮
	spliter        lcl.ISplitter // 分割线
	componentsTabs lcl.IPanel    // 组件面板选项卡
}

func (m *TAppWindow) createTopToolbar() *TopToolbar {
	bar := &TopToolbar{}
	bar.IPanel = lcl.NewPanel(m)
	bar.IPanel.SetParent(m)
	bar.IPanel.SetBevelOuter(types.BvNone)
	bar.IPanel.SetDoubleBuffered(true)
	bar.IPanel.SetWidth(m.Width())
	bar.IPanel.SetHeight(toolbarHeight)
	bar.IPanel.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	bar.IPanel.SetColor(colors.RGBToColor(56, 57, 60))
	m.toolbar = bar
	// 创建工具按钮

	return bar
}
