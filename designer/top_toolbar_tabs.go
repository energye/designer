// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"widget/wg"
)

// 组件选项卡
type ComponentTab struct {
	//sheet         lcl.ITabSheet
	sheet         *wg.TPage
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
	//page := lcl.NewPageControl(m.box)
	tab := wg.NewTab(m.box)
	m.tab = tab
	tab.SetBounds(0, 0, m.rightTabs.Width(), m.rightTabs.Height())
	tab.SetAlign(types.AlClient)
	tab.EnableScrollButton(false)
	tab.SetColor(colors.ClGray)
	tab.SetParent(m.rightTabs)
	tab.SetOnChange(func(sender lcl.IObject) {
		logs.Debug("Toolbar Tabs Change")
		m.ResetTabComponentDown()
	})

	inspectorTreeImageIndex := 0 // 查看器组件树图片索引
	// 创建组件选项卡
	newComponentTab := func(tab config.Tab) {
		compTab := &ComponentTab{components: make(map[string]*ComponentTabItem)}
		m.componentTabs[tab.En] = compTab
		//sheet := lcl.NewTabSheet(m.tab)
		sheet := m.tab.NewPage()
		sheet.Button().SetText(tab.Cn)                                // 设置标签按钮显示文本
		sheet.Button().SetColorGradient(colors.ClGray, colors.ClGray) // 设置标签按钮过度颜色
		sheet.SetDefaultColor(colors.ClGray)                          // 设置默认颜色
		sheet.SetActiveColor(wg.LightenColor(colors.ClGray, 0.3))     // 设置激活颜色
		sheet.SetColor(colors.ClGray)                                 // 设置背景色
		//sheet.SetAlign(types.AlClient)
		//sheet.SetParent(m.tab)
		compTab.sheet = sheet

		// 显示组件工具按钮
		componentToolbar := lcl.NewToolBar(sheet)
		componentToolbar.SetImages(imageComponents.ImageList150())
		componentToolbar.SetButtonWidth(36)
		componentToolbar.SetButtonHeight(36)
		componentToolbar.SetHeight(36)
		componentToolbar.SetTop(3)
		componentToolbar.SetWidth(sheet.Width())
		componentToolbar.SetEdgeBorders(types.NewSet())
		componentToolbar.SetAlign(types.AlCustom)
		componentToolbar.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
		componentToolbar.SetBorderStyleToBorderStyle(types.BsNone)
		componentToolbar.SetParent(sheet)
		compTab.toolbar = componentToolbar

		// 选择工具 鼠标
		selectToolBtn := lcl.NewToolButton(componentToolbar)
		selectToolBtn.SetHint("选择工具")
		selectToolBtn.SetImageIndex(imageComponents.ImageIndex("cursortool_150.png"))
		selectToolBtn.SetShowHint(true)
		selectToolBtn.SetDown(true)
		selectToolBtn.SetParent(componentToolbar)
		comp := &ComponentTabItem{owner: compTab, index: 0, name: "SelectTool", btn: selectToolBtn}
		compTab.components[comp.name] = comp
		compTab.selectToolBtn = selectToolBtn

		sep := lcl.NewToolButton(componentToolbar)
		sep.SetStyle(types.TbsSeparator)
		sep.SetParent(componentToolbar)

		// 创建组件按钮
		for i, name := range tab.Component {
			btn := lcl.NewToolButton(componentToolbar)
			btn.SetHint(name)
			btn.SetImageIndex(imageComponents.ImageIndex(name + "_150.png")) // 36px
			btn.SetShowHint(true)
			btn.SetParent(componentToolbar)
			comp = &ComponentTabItem{owner: compTab, index: i, inspectorTreeImageIndex: inspectorTreeImageIndex, name: name, btn: btn}
			compTab.components[name] = comp
			inspectorTreeImageIndex++
		}
		go compTab.BindToolBtnEvent()
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
	lcl.RunOnMainThreadAsync(func(id uint32) {
		tab.RecalculatePosition()
		m.componentTabs[config.Config.ComponentTabs.Standard.En].sheet.Active(true)
	})
}

// 绑定事件
func (m *ComponentTab) BindToolBtnEvent() {
	m.selectToolBtn.SetOnClick(m.SelectToolBtnOnClick)
	for _, comp := range m.components {
		comp.btn.SetOnClick(comp.ComponentBtnOnClick)
	}
}

// 工具栏上的按钮 取消按下
func (m *ComponentTab) UnDownComponents() {
	for _, com := range m.components {
		com.btn.SetDown(false)
	}
	toolbar.SetSelectComponentItem(nil)
}

// 取消选择工具按下
func (m *ComponentTab) UnDownSelectTool() {
	m.selectToolBtn.SetDown(false)
	toolbar.SetSelectComponentItem(nil)
}

// 设置选择工具按下
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
	// 设置当前工具按钮按下
	m.btn.SetDown(true)
	// 设置当前工具按钮选中
	toolbar.SetSelectComponentItem(m)
}
