package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
)

// 组件选项卡
type ComponentTab struct {
	sheet         lcl.ITabSheet
	toolbar       lcl.IToolBar
	selectToolBtn lcl.IToolButton
	components    map[string]*ComponentTabItem
}

// 组件选项卡 组件项
type ComponentTabItem struct {
	owner                   *ComponentTab
	index                   int // tab 面板组件图标索引
	inspectorTreeImageIndex int // 查看器组件树图标索引
	name                    string
	btn                     lcl.IToolButton
}

// 组件选项卡
func (m *TopToolbar) createComponentTabs() {
	page := lcl.NewPageControl(m.box)
	page.SetParent(m.rightTabs)
	page.SetAlign(types.AlClient)
	page.SetTabStop(true)
	m.page = page
	m.page.SetOnChange(func(sender lcl.IObject) {
		logs.Debug("Toolbar Tabs Change")
		m.ResetTabComponentDown()
	})

	inspectorTreeImageIndex := 0 // 查看器组件树图片索引
	// 创建组件选项卡
	newComponentTab := func(tab config.Tab) {
		compTab := &ComponentTab{components: make(map[string]*ComponentTabItem)}
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
		componentToolbar := lcl.NewToolBar(sheet)
		componentToolbar.SetParent(sheet)
		componentToolbar.SetImages(LoadImageList(m.rightTabs, imageList, 36, 36))
		componentToolbar.SetButtonWidth(36)
		componentToolbar.SetButtonHeight(36)
		componentToolbar.SetHeight(36)
		componentToolbar.SetEdgeBorders(types.NewSet())
		compTab.toolbar = componentToolbar

		// 选择工具 鼠标
		selectToolBtn := lcl.NewToolButton(componentToolbar)
		selectToolBtn.SetParent(componentToolbar)
		selectToolBtn.SetHint("选择工具")
		selectToolBtn.SetImageIndex(int32(0))
		selectToolBtn.SetShowHint(true)
		selectToolBtn.SetDown(true)
		comp := &ComponentTabItem{owner: compTab, index: 0, name: "SelectTool", btn: selectToolBtn}
		compTab.components[comp.name] = comp
		compTab.selectToolBtn = selectToolBtn

		seap := lcl.NewToolButton(componentToolbar)
		seap.SetParent(componentToolbar)
		seap.SetStyle(types.TbsSeparator)

		// 创建组件按钮
		for i, name := range tab.Component {
			imageIndex := i + 1
			btn := lcl.NewToolButton(componentToolbar)
			btn.SetParent(componentToolbar)
			btn.SetHint(name)
			btn.SetImageIndex(int32(imageIndex))
			btn.SetShowHint(true)
			comp = &ComponentTabItem{owner: compTab, index: i, inspectorTreeImageIndex: inspectorTreeImageIndex, name: name, btn: btn}
			compTab.components[name] = comp
			inspectorTreeImageIndex++
		}
		compTab.BindToolBtnEvent()
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

// 绑定事件
func (m *ComponentTab) BindToolBtnEvent() {
	m.selectToolBtn.SetOnClick(m.SelectToolBtnOnClick)
	for _, comp := range m.components {
		comp.btn.SetOnClick(comp.ComponentBtnOnClick)
	}
}

func (m *ComponentTab) UnDownComponents() {
	for _, com := range m.components {
		com.btn.SetDown(false)
	}
	toolbar.SetSelectComponentItem(nil)
}

func (m *ComponentTab) UnDownSelectTool() {
	m.selectToolBtn.SetDown(false)
	toolbar.SetSelectComponentItem(nil)
}

func (m *ComponentTab) DownSelectTool() {
	m.selectToolBtn.SetDown(true)
	toolbar.SetSelectComponentItem(nil)
}

// 选择工具按钮事件
func (m *ComponentTab) SelectToolBtnOnClick(sender lcl.IObject) {
	logs.Debug("SelectToolBtnOnClick")
	m.UnDownComponents()
	m.DownSelectTool()
	toolbar.SetSelectComponentItem(nil)
}

// 组件按钮事件
func (m *ComponentTabItem) ComponentBtnOnClick(sender lcl.IObject) {
	logs.Debug("ToolBtnOnClick", m.index, m.name)
	m.owner.UnDownComponents()
	m.owner.UnDownSelectTool()
	m.btn.SetDown(true)
	toolbar.SetSelectComponentItem(m)
}
