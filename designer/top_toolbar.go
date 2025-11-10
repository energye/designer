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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
	"github.com/energye/lcl/types"
	"widget/wg"
)

// 顶部工具栏

var (
	toolbarHeight int32 = 72
	toolbar       *TopToolbar
)

// 初始化工具栏相关配置
func initConfigToolbar() {
	if tool.IsLinux() {
		toolbarHeight = 72
	} else {
		toolbarHeight = 66
	}
}

type TopToolbar struct {
	//page            lcl.IPageControl
	tab             *wg.TTab
	box             lcl.IPanel
	leftTools       lcl.IPanel
	rightTabs       lcl.IPanel               // 组件面板选项卡
	toolbarBtn      *TToolbarToolBtn         // 工具栏按钮
	componentTabs   map[string]*ComponentTab // 组件选项卡： 标准，附加，通用等等
	selectComponent *ComponentTabItem        // 选中的组件
}

func (m *TAppWindow) createTopToolbar() {
	bar := &TopToolbar{componentTabs: make(map[string]*ComponentTab)}
	toolbar = bar
	// 工具栏面板
	bar.box = lcl.NewPanel(m)
	//bar.box.SetBevelColor(wg.LightenColor(bgDarkColor, 0.3))
	//bar.box.SetBevelOuter(types.BvLowered)
	bar.box.SetBevelOuter(types.BvNone)
	bar.box.SetDoubleBuffered(true)
	bar.box.SetWidth(m.Width())
	bar.box.SetHeight(toolbarHeight)
	bar.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	SetComponentDefaultColor(bar.box)
	bar.box.SetParentColor(true)
	bar.box.SetParent(m)

	// 工具栏-左 工具按钮
	bar.leftTools = lcl.NewPanel(m)
	//bar.leftTools.SetBevelColor(wg.LightenColor(bgDarkColor, 0.3))
	bar.leftTools.SetBevelOuter(types.BvNone)
	bar.leftTools.SetDoubleBuffered(true)
	bar.leftTools.SetWidth(180)
	bar.leftTools.SetHeight(bar.box.Height())
	SetComponentDefaultColor(bar.leftTools)
	bar.leftTools.SetAlign(types.AlLeft)
	bar.leftTools.SetParent(bar.box)

	// 工具栏-右 组件选项卡
	bar.rightTabs = lcl.NewPanel(m)
	bar.rightTabs.SetBevelOuter(types.BvNone)
	bar.rightTabs.SetDoubleBuffered(true)
	bar.rightTabs.SetHeight(bar.box.Height())
	bar.rightTabs.SetWidth(bar.box.Width())
	SetComponentDefaultColor(bar.rightTabs)
	bar.rightTabs.SetAlign(types.AlClient)
	bar.rightTabs.SetParent(bar.box)

	// 创建工具按钮
	bar.createToolBarBtns()

	// 创建组件选项卡
	bar.createComponentTabs()
}

// 重置Tab组件选项卡按下状态
func (m *TopToolbar) ResetTabComponentDown() {
	for _, comp := range m.componentTabs {
		comp.UnDownComponents()
		comp.DownSelectTool()
	}
}

// 设置当前工具按钮选中
// 之后在设计器里使用
func (m *TopToolbar) SetSelectComponentItem(item *ComponentTabItem) {
	m.selectComponent = item
}
