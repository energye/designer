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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 工具按钮功能

type TToolbarToolBtn struct {
	toolBtnBar     lcl.IToolBar
	newFormBtn     lcl.IToolButton
	openFormBtn    lcl.IToolButton
	saveFormBtn    lcl.IToolButton
	saveAllFormBtn lcl.IToolButton
	runPreviewBtn  lcl.IToolButton
	previewState   consts.PreviewState // 预览状态
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

	toolbarBtn.newFormBtn = newBtn(imageMenu.ImageIndex("menu_new_form_150.png"), "新建窗体", 0)
	toolbarBtn.newFormBtn.SetOnClick(toolbarBtn.onNewForm)

	toolbarBtn.openFormBtn = newBtn(imageMenu.ImageIndex("menu_project_open_150.png"), "打开", 1)
	toolbarBtn.openFormBtn.SetOnClick(toolbarBtn.onOpenForm)
	newSep()

	toolbarBtn.saveFormBtn = newBtn(imageMenu.ImageIndex("menu_save_150.png"), "保存窗体", 1)
	toolbarBtn.saveFormBtn.SetOnClick(toolbarBtn.onSaveForm)

	toolbarBtn.saveAllFormBtn = newBtn(imageMenu.ImageIndex("menu_save_all_150.png"), "保存所有窗体", 1)
	toolbarBtn.saveAllFormBtn.SetOnClick(toolbarBtn.onSaveAllForm)
	newSep()

	toolbarBtn.runPreviewBtn = newBtn(imageMenu.ImageIndex("menu_run_150.png"), "运行(F9)", 3)
	toolbarBtn.runPreviewBtn.SetOnClick(toolbarBtn.onRunPreviewForm)
}

func (m *TToolbarToolBtn) onNewForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 新建窗体")
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
}

func (m *TToolbarToolBtn) onOpenForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 打开")
	mainWindow.openDialog.SetTitle("打开项目/UI")
	mainWindow.openDialog.SetFilter(config.DialogFilter.UIFilter())
	mainWindow.openDialog.SetFilterIndex(1)
	if mainWindow.openDialog.Execute() {
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			filePath := mainWindow.openDialog.FileName()
			event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: filePath}})
		})
	}
}

func (m *TToolbarToolBtn) onSaveForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 保存窗体")
}

func (m *TToolbarToolBtn) onSaveAllForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 保存所有窗体")
}

func (m *TToolbarToolBtn) onRunPreviewForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 预览窗体")
	if m.previewState == consts.PsStarted {
		logs.Debug("工具栏按钮, 停止预览窗体")
		event.Emit(event.TTrigger{Name: event.Preview, Payload: consts.PsStop})
	} else {
		logs.Debug("工具栏按钮, 运行预览窗体")
		result := make(chan any)
		go func() {
			logs.Debug("状态监听开始")
			for res := range result {
				logs.Debug("预览窗口结果:", res)
				if status, ok := res.(consts.PreviewState); ok {
					m.switchPreviewBtn(status)
					if status == consts.PsStop {
						// 运行结束
						break
					}
				} else {
					logs.Error("预览窗口结果错误:", res)
					// 运行结束
					m.switchPreviewBtn(consts.PsStop)
					break
				}
			}
			logs.Debug("状态监听结束")
		}()
		// 启动运行预览
		event.Emit(event.TTrigger{Name: event.Preview, Payload: consts.PsStarted, Result: result})
	}
}

// 切换预览按钮状态, 在运行和结束运行之间切换
func (m *TToolbarToolBtn) switchPreviewBtn(status consts.PreviewState) {
	logs.Debug("切换预览按钮状态 status:", status)
	m.previewState = status
	m.runPreviewBtn.SetEnabled(true)
	if m.previewState == consts.PsStarted {
		m.runPreviewBtn.SetHint("停止(F9)")
		m.runPreviewBtn.SetImageIndex(imageMenu.ImageIndex("menu_stop_150.png"))
	} else if m.previewState == consts.PsStarting {
		m.runPreviewBtn.SetEnabled(false)
		m.runPreviewBtn.SetHint("停止(F9)")
		m.runPreviewBtn.SetImageIndex(imageMenu.ImageIndex("menu_stop_150.png"))
	} else {
		m.runPreviewBtn.SetHint("运行(F9)")
		m.runPreviewBtn.SetImageIndex(imageMenu.ImageIndex("menu_run_150.png"))
	}
	mainWindow.mainMenu.switchRunMenuItem(status)
}
