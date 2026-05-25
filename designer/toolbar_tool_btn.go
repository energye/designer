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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
)

// 工具按钮功能

type TToolbarToolBtn struct {
	toolBtnBar   lcl.IToolBar
	newWindowBtn lcl.IToolButton
	openBtn      lcl.IToolButton
	saveBtn      lcl.IToolButton

	undoBtn lcl.IToolButton
	redoBtn lcl.IToolButton

	runPreviewBtn lcl.IToolButton

	previewState consts.PreviewState // 预览状态
}

// SetEnableToolButtons 设置工具栏按钮的启用状态
//
//	enable: 布尔值，true表示启用按钮，false表示禁用按钮
func (m *TToolbarToolBtn) SetEnableToolButtons(enable bool) {
	enabled := func() {
		m.newWindowBtn.SetEnabled(enable)
		//m.openBtn.SetEnabled(enable)
		//m.saveBtn.SetEnabled(enable)
		m.runPreviewBtn.SetEnabled(enable)
	}
	if tool.IsMainThread() {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			enabled()
		})
	} else {
		lcl.RunOnMainThreadSync(func() {
			enabled()
		})
	}
}

// 新建窗体
func (m *TToolbarToolBtn) onNewForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 新建窗体")
	go lcl.RunOnMainThreadAsync(func(id uint32) {
		// 创建窗体后执行一次 go mod tidy 禁用功能按钮, TODO 先这样
		//SetEnableFuncComponent(false)

		// 创建新的 form tab
		newForm := designer.addDesignFormTab(nil)
		if newForm == nil {
			event.ConsoleWriteError("new design form error")
			return
		}
		// 激活显示 新的 form tab
		designer.ActiveFormTab(newForm)
		// 1. 加载属性到设计器
		// 此步骤会初始化并填充设计组件实例
		newForm.FormRoot.LoadPropertyToInspector()
		// 2. 添加到组件树
		newNode := newForm.AddFormNode()
		newNode.SetSelected(true)
		// 触发 ui 生成事件
		triggerUIGeneration(newForm.FormRoot, nil, event.CodeGenUI)

		designer.tab.HideAllActivated()
		// 显示 tab page
		newForm.mainPage.SetActive(true)
		designer.tab.RecalculatePosition()

		// 创建窗体后执行一次 go mod tidy, TODO 先这样
		//go func() {
		//	cmd := command.NewCMD()
		//	cmd.HideWindow = true
		//	cmd.Dir = projBean.GPath
		//	cmd.Console = func(data string, level command.Level) {
		//		event.ConsoleWriteInfo(data)
		//	}
		//	cmd.Command("go", "mod", "tidy")
		//	// 新建窗体后执行一次 go mod tidy 恢复功能按钮
		//	SetEnableFuncComponent(true)
		//}()
	})
}

func (m *TToolbarToolBtn) onOpenForm(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 打开项目/打开UI布局")
	// 停止自动关联项目加载
	stopAutoAssociateProjectLoad()
	MainWindow.openDialog.SetTitle("打开项目/打开UI布局")
	MainWindow.openDialog.SetFilter(config.DialogFilter.UIFilter())
	MainWindow.openDialog.SetFilterIndex(1)
	if MainWindow.openDialog.Execute() {
		ProjectTreeClearComponentTreeNode()
		ProjectTreeClearSrcTreeNode()
		filePath := MainWindow.openDialog.FileName()
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: filePath}})
	}
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
	changeStatus := func() {
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
		MainWindow.mainMenu.switchRunMenuItem(status)
	}
	lcl.RunOnMainThreadSync(func() {
		changeStatus()
	})
}

func (m *TToolbarToolBtn) onSave(sender lcl.IObject) {
	logs.Debug("工具栏按钮, 保存窗体")
}

func (m *TToolbarToolBtn) onUndo(sender lcl.IObject) {

}

func (m *TToolbarToolBtn) onRedo(sender lcl.IObject) {

}
