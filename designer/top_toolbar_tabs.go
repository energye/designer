package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
)

type ComponentTab struct {
	sheet      lcl.ITabSheet
	toolbar    lcl.IToolBar
	toolBtn    lcl.IToolButton
	components map[string]*component
}

type component struct {
	index int
	name  string
	btn   lcl.IToolButton
}

// 组件选项卡
func (m *TopToolbar) createComponentTabs() {
	page := lcl.NewPageControl(m.box)
	page.SetParent(m.rightTabs)
	page.SetAlign(types.AlClient)
	page.SetTabStop(true)
	m.page = page

	// 创建组件选项卡
	newComponentTab := func(tab config.Tab) {
		compTab := &ComponentTab{components: make(map[string]*component)}
		m.componentTabs[tab.En] = compTab

		sheet := lcl.NewTabSheet(page)
		sheet.SetParent(page)
		sheet.SetCaption(tab.Cn)
		sheet.SetAlign(types.AlClient)
		compTab.sheet = sheet
		// 组件图标
		var imageList []string
		imageList = append(imageList, "components/cursortool_150.png")
		for _, name := range tab.Component {
			imageList = append(imageList, fmt.Sprintf("components/%v_150.png", strings.ToLower(name)))
		}
		// 显示组件工具按钮
		toolbar := lcl.NewToolBar(sheet)
		toolbar.SetParent(sheet)
		toolbar.SetImages(LoadImageList(m.rightTabs, imageList, 36, 36))
		toolbar.SetButtonWidth(36)
		toolbar.SetButtonHeight(36)
		toolbar.SetHeight(36)
		toolbar.SetEdgeBorders(types.NewSet())
		compTab.toolbar = toolbar

		// 选择工具 鼠标
		toolBtn := lcl.NewToolButton(toolbar)
		toolBtn.SetParent(toolbar)
		toolBtn.SetHint("选择工具")
		toolBtn.SetImageIndex(int32(0))
		toolBtn.SetShowHint(true)
		compTab.toolBtn = toolBtn

		seap := lcl.NewToolButton(toolbar)
		seap.SetParent(toolbar)
		seap.SetStyle(types.TbsSeparator)

		// 创建组件按钮
		for i, name := range tab.Component {
			imageIndex := i + 1
			btn := lcl.NewToolButton(toolbar)
			btn.SetParent(toolbar)
			btn.SetHint(name)
			btn.SetImageIndex(int32(imageIndex))
			btn.SetShowHint(true)
			compTab.components[name] = &component{index: i, name: name, btn: btn}
		}
	}
	// 创建组件选项卡
	newComponentTab(config.Config.ComponentTabs.Standard)
	newComponentTab(config.Config.ComponentTabs.Additional)
	newComponentTab(config.Config.ComponentTabs.Common)
	newComponentTab(config.Config.ComponentTabs.Dialogs)
	newComponentTab(config.Config.ComponentTabs.Misc)
	newComponentTab(config.Config.ComponentTabs.System)
	newComponentTab(config.Config.ComponentTabs.LazControl)
	newComponentTab(config.Config.ComponentTabs.WebView)
}
