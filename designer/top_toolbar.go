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
	"github.com/energye/lcl/types"
)

// 顶部工具栏

var toolbar *TopToolbar

type TopToolbar struct {
	page            lcl.IPageControl
	box             lcl.IPanel
	leftTools       lcl.IPanel
	rightTabs       lcl.IPanel               // 组件面板选项卡
	componentTabs   map[string]*ComponentTab // 组件选项卡： 标准，附加，通用等等
	selectComponent *ComponentTabItem        // 选中的组件
}

func (m *TAppWindow) createTopToolbar() {
	bar := &TopToolbar{componentTabs: make(map[string]*ComponentTab)}
	toolbar = bar
	// 工具栏面板
	bar.box = lcl.NewPanel(m)
	bar.box.SetParent(m)
	bar.box.SetBevelOuter(types.BvNone)
	bar.box.SetDoubleBuffered(true)
	bar.box.SetWidth(m.Width())
	bar.box.SetHeight(toolbarHeight)
	bar.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	bar.box.SetParentColor(true)

	// 工具栏-左 工具按钮
	bar.leftTools = lcl.NewPanel(m)
	bar.leftTools.SetParent(bar.box)
	bar.leftTools.SetBevelOuter(types.BvNone)
	bar.leftTools.SetDoubleBuffered(true)
	bar.leftTools.SetWidth(180)
	bar.leftTools.SetHeight(bar.box.Height())
	bar.leftTools.SetAlign(types.AlLeft)

	// 工具栏-右 组件选项卡
	bar.rightTabs = lcl.NewPanel(m)
	bar.rightTabs.SetParent(bar.box)
	bar.rightTabs.SetBevelOuter(types.BvNone)
	bar.rightTabs.SetDoubleBuffered(true)
	bar.rightTabs.SetHeight(bar.box.Height())
	bar.rightTabs.SetAlign(types.AlClient)

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

// 工具按钮
func (m *TopToolbar) createToolBarBtns() {
	toolBtnBar := lcl.NewToolBar(m.box)
	toolBtnBar.SetParent(m.leftTools)
	toolBtnBar.SetAlign(types.AlCustom)
	toolBtnBar.SetTop(16)
	toolBtnBar.SetButtonWidth(32)
	toolBtnBar.SetButtonHeight(32)
	toolBtnBar.SetHeight(32)
	toolBtnBar.SetWidth(m.leftTools.Width())
	toolBtnBar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	toolBtnBar.SetEdgeBorders(types.NewSet())
	toolBtnBar.SetImages(imageMenu.ImageList150())
	newSepa := func() {
		seap := lcl.NewToolButton(toolBtnBar)
		seap.SetParent(toolBtnBar)
		seap.SetStyle(types.TbsSeparator)
	}

	newBtn := func(imageIndex int32, hint string, margin int32) lcl.IToolButton {
		btn := lcl.NewToolButton(toolBtnBar)
		btn.SetParent(toolBtnBar)
		btn.SetHint(hint)
		btn.SetImageIndex(imageIndex)
		btn.SetShowHint(true)
		//btn := lcl.NewBitBtn(toolBtnBarf)
		//btn.SetParent(toolBtnBarf)
		//btn.SetWidth(32)
		//btn.SetHeight(32)
		//btn.SetTabStop(true)
		//btn.SetImages(imageList)
		//btn.SetImageIndex(imageIndex)
		//btn.SetMargin(margin)
		return btn
	}

	newFormBtn := newBtn(imageMenu.ImageIndex("menu_new_form_150.png"), "新建窗体", 0)
	newFormBtn.SetOnClick(func(sender lcl.IObject) {
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			newForm := designer.addDesignerFormTab()
			designer.ActiveFormTab(newForm)
			// 1. 加载属性到设计器
			// 此步骤会初始化并填充设计组件实例
			newForm.FormRoot.LoadPropertyToInspector()
			// 2. 添加到组件树
			newForm.AddFormNode()
			go triggerUIGeneration(newForm.FormRoot)
		})
	})
	openFormBtn := newBtn(imageMenu.ImageIndex("menu_project_open_150.png"), "打开窗体", 1)
	openFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	newSepa()
	saveFormBtn := newBtn(imageMenu.ImageIndex("menu_save_150.png"), "保存窗体", 1)
	saveFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	saveAllFormBtn := newBtn(imageMenu.ImageIndex("menu_save_all_150.png"), "保存所有窗体", 1)
	saveAllFormBtn.SetOnClick(func(sender lcl.IObject) {

	})
	newSepa()
	runFormBtn := newBtn(imageMenu.ImageIndex("menu_run_150.png"), "运行预览窗体", 3)
	runFormBtn.SetOnClick(func(sender lcl.IObject) {

	})

}
