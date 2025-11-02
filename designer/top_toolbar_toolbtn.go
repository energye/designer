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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 工具按钮功能

type TToolbarToolBtn struct {
	toolBtnBar     lcl.IToolBar
	openFormBtn    lcl.IToolButton
	saveFormBtn    lcl.IToolButton
	saveAllFormBtn lcl.IToolButton
	runPreviewBtn  lcl.IToolButton
}

// 工具按钮
func (m *TopToolbar) createToolBarBtns() {
	if m.toolbarBtn != nil {
		return
	}
	toolbarBtn := new(TToolbarToolBtn)
	m.toolbarBtn = toolbarBtn

	toolBtnBar := lcl.NewToolBar(m.box)
	toolBtnBar.SetAlign(types.AlCustom)
	toolBtnBar.SetTop(16)
	toolBtnBar.SetButtonWidth(32)
	toolBtnBar.SetButtonHeight(32)
	toolBtnBar.SetHeight(32)
	toolBtnBar.SetWidth(m.leftTools.Width())
	toolBtnBar.SetAnchors(types.NewSet(types.AkLeft, types.AkRight))
	toolBtnBar.SetEdgeBorders(types.NewSet())
	toolBtnBar.SetImages(imageMenu.ImageList150())
	toolBtnBar.SetParent(m.leftTools)
	toolbarBtn.toolBtnBar = toolBtnBar

	newSep := func() {
		sep := lcl.NewToolButton(toolBtnBar)
		sep.SetParent(toolBtnBar)
		sep.SetStyle(types.TbsSeparator)
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

	toolbarBtn.openFormBtn = newBtn(imageMenu.ImageIndex("menu_project_open_150.png"), "打开窗体", 1)
	toolbarBtn.openFormBtn.SetOnClick(toolbarBtn.onOpenForm)
	newSep()

	toolbarBtn.saveFormBtn = newBtn(imageMenu.ImageIndex("menu_save_150.png"), "保存窗体", 1)
	toolbarBtn.saveFormBtn.SetOnClick(toolbarBtn.onSaveForm)

	toolbarBtn.saveAllFormBtn = newBtn(imageMenu.ImageIndex("menu_save_all_150.png"), "保存所有窗体", 1)
	toolbarBtn.saveAllFormBtn.SetOnClick(toolbarBtn.onSaveAllForm)
	newSep()

	toolbarBtn.runPreviewBtn = newBtn(imageMenu.ImageIndex("menu_run_150.png"), "运行预览窗体", 3)
	toolbarBtn.runPreviewBtn.SetOnClick(toolbarBtn.onRunPreviewForm)
}

func (m *TToolbarToolBtn) onOpenForm(sender lcl.IObject) {

}

func (m *TToolbarToolBtn) onSaveForm(sender lcl.IObject) {

}

func (m *TToolbarToolBtn) onSaveAllForm(sender lcl.IObject) {

}

func (m *TToolbarToolBtn) onRunPreviewForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 预览")
	event.Preview.TriggerEvent(event.TEventTrigger{})
}
