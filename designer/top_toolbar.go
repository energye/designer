package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 顶部工具栏

type TopToolbar struct {
	toolbarPnl     lcl.IPanel
	spliter        lcl.ISplitter // 分割线
	componentsTabs lcl.IPanel    // 组件面板选项卡
}

func (m *TAppWindow) createTopToolbar() *TopToolbar {
	bar := &TopToolbar{}
	bar.toolbarPnl = lcl.NewPanel(m)
	bar.toolbarPnl.SetParent(m)
	bar.toolbarPnl.SetBevelOuter(types.BvNone)
	bar.toolbarPnl.SetDoubleBuffered(true)
	bar.toolbarPnl.SetWidth(m.Width())
	bar.toolbarPnl.SetHeight(toolbarHeight)
	bar.toolbarPnl.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	bar.toolbarPnl.SetColor(colors.RGBToColor(56, 57, 60))
	m.toolbar = bar
	// 创建工具按钮
	toolbar := lcl.NewToolBar(m)
	toolbar.SetParent(bar.toolbarPnl)

	toolBtn := lcl.NewToolButton(toolbar)
	toolBtn.SetParent(toolbar)
	toolBtn.SetCaption("asdf")
	toolBtn.SetHint("提示")
	return bar
}
