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
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	desTool "github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
)

// 设计器顶部菜单

type TMainMenu struct {
	main    lcl.IMainMenu
	file    lcl.IMenuItem
	edit    lcl.IMenuItem
	view    lcl.IMenuItem
	run     lcl.IMenuItem
	setting lcl.IMenuItem
	helper  lcl.IMenuItem

	// file
	createProject lcl.IMenuItem
	createWindow  lcl.IMenuItem
	open          lcl.IMenuItem
	history       lcl.IMenuItem
	//save          lcl.IMenuItem

	// edit
	editCut       lcl.IMenuItem
	editCopy      lcl.IMenuItem
	editPaste     lcl.IMenuItem
	editSelectAll lcl.IMenuItem
	editUndo      lcl.IMenuItem
	editDel       lcl.IMenuItem

	// view
	viewWidgets   lcl.IMenuItem
	viewProject   lcl.IMenuItem
	viewInspector lcl.IMenuItem
	viewConsole   lcl.IMenuItem
	viewStatusbar lcl.IMenuItem

	// run
	build      lcl.IMenuItem
	buildClean lcl.IMenuItem
	buildAll   lcl.IMenuItem
	runApp     lcl.IMenuItem

	// setting
	buildOption       lcl.IMenuItem
	environmentOption lcl.IMenuItem
	frameworkOption   lcl.IMenuItem
	projectOption     lcl.IMenuItem

	actionList lcl.IActionList
}

// 设计器主菜单
func (m *TAppWindow) initMainMenu() {
	if m.mainMenu != nil {
		return
	}

	mainMenu := new(TMainMenu)
	m.mainMenu = mainMenu

	m.mainMenu.actionList = lcl.NewActionList(m)

	mainMenu.main = lcl.NewMainMenu(m)
	mainMenu.main.SetImages(imageMenu.ImageList100())
	menuItems := mainMenu.main.Items()

	mainMenu.file = lcl.NewMenuItem(m)
	mainMenu.file.SetCaption("文件(&F)")
	menuItems.Add(mainMenu.file)

	mainMenu.edit = lcl.NewMenuItem(m)
	mainMenu.edit.SetCaption("编辑(&E)")
	menuItems.Add(mainMenu.edit)

	mainMenu.view = lcl.NewMenuItem(m)
	mainMenu.view.SetCaption("视图(&V)")
	menuItems.Add(mainMenu.view)

	mainMenu.run = lcl.NewMenuItem(m)
	mainMenu.run.SetCaption("运行(&R)")
	menuItems.Add(mainMenu.run)

	mainMenu.setting = lcl.NewMenuItem(m)
	mainMenu.setting.SetCaption("设置(&S)")
	menuItems.Add(mainMenu.setting)

	mainMenu.helper = lcl.NewMenuItem(m)
	mainMenu.helper.SetCaption("帮助(&H)")
	menuItems.Add(mainMenu.helper)

	mainMenu.fileMenu(m)
	mainMenu.editMenu(m)
	mainMenu.viewMenu(m)
	mainMenu.runMenu(m)
	mainMenu.settingMenu(m)
	mainMenu.helperMenu(m)
	mainMenu.macOS()
}

// SetEnableMenuItems 设置菜单项的启用状态
//
//	enable: 布尔值，true表示启用菜单项，false表示禁用菜单项
//	该函数在主线程中异步执行，确保线程安全性。它会同时设置多个菜单项的启用状态，
//	包括创建窗口、打开、保存、构建、清理构建、运行应用、构建选项、环境选项和项目选项菜单项。
func (m *TMainMenu) SetEnableMenuItems(enable bool) {
	enabled := func() {
		logs.Info("设置菜单项的启用状态 enable:", enable)
		m.createWindow.SetEnabled(enable)
		//m.open.SetEnabled(enable)
		//m.save.SetEnabled(enable)
		m.build.SetEnabled(enable)
		m.buildClean.SetEnabled(enable)
		m.buildAll.SetEnabled(enable)
		m.runApp.SetEnabled(enable)
		m.buildOption.SetEnabled(enable)
		m.environmentOption.SetEnabled(enable)
		m.projectOption.SetEnabled(enable)
	}
	if desTool.IsMainThread() {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			enabled()
		})
	} else {
		lcl.RunOnMainThreadSync(func() {
			enabled()
		})
	}
}

func (m *TMainMenu) macOS() {
	if tool.IsDarwin() {
		// macOS
	}
}

func (m *TMainMenu) fileMenu(owner lcl.IComponent) {
	create := lcl.NewMenuItem(owner)
	create.SetCaption("新建(&N)")
	create.SetImageIndex(imageMenu.ImageIndex("menu_project_create.png"))
	m.file.Add(create)

	m.createProject = lcl.NewMenuItem(owner)
	m.createProject.SetCaption("新建项目")
	m.createProject.SetShortCut(api.TextToShortCut("Ctrl+P"))
	m.createProject.SetImageIndex(imageMenu.ImageIndex("menu_project_add.png"))
	m.createProject.SetOnClick(func(lcl.IObject) {
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectCreate}})
	})
	create.Add(m.createProject)

	sep := lcl.NewMenuItem(owner)
	sep.SetCaption("-")
	create.Add(sep)

	m.createWindow = lcl.NewMenuItem(owner)
	m.createWindow.SetCaption("新建窗体")
	m.createWindow.SetShortCut(api.TextToShortCut("Ctrl+N"))
	m.createWindow.SetImageIndex(imageMenu.ImageIndex("menu_new_form.png"))
	m.createWindow.SetOnClick(func(sender lcl.IObject) {
		logs.Debug("新建窗体")
		MainWindow.toolLayout.toolbarBtn.onNewForm(sender)
	})
	create.Add(m.createWindow)

	m.open = lcl.NewMenuItem(owner)
	m.open.SetCaption("打开(&O)")
	m.open.SetShortCut(api.TextToShortCut("Ctrl+O"))
	m.open.SetImageIndex(imageMenu.ImageIndex("menu_project_open.png"))
	m.open.SetOnClick(func(sender lcl.IObject) {
		MainWindow.toolLayout.toolbarBtn.onOpenForm(sender)
	})
	m.file.Add(m.open)

	//m.save = lcl.NewMenuItem(owner)
	//m.save.SetCaption("保存(&S)")
	//m.save.SetShortCut(api.TextToShortCut("Ctrl+S"))
	//m.save.SetImageIndex(imageMenu.ImageIndex("menu_save.png"))
	//m.save.SetOnClick(func(sender lcl.IObject) {
	//	logs.Debug("保存窗体")
	//})
	//m.file.Add(m.save)

	//saveAllWindow := lcl.NewMenuItem(owner)
	//saveAllWindow.SetCaption("保存所有窗体(&L)")
	//saveAllWindow.SetShortCut(api.TextToShortCut("Shift+Ctrl+L"))
	//saveAllWindow.SetImageIndex(imageMenu.ImageIndex("menu_save_all.png"))
	//saveAllWindow.SetOnClick(func(sender lcl.IObject) {
	//	logs.Debug("保存所有窗体")
	//})
	//m.file.Add(saveAllWindow)

	m.history = lcl.NewMenuItem(owner)
	m.history.SetCaption("历史项目")
	m.history.SetImageIndex(imageMenu.ImageIndex("menu_project_history.png"))
	m.file.Add(m.history)
	m.fileHistoryProjectMenu()

	exitWindow := lcl.NewMenuItem(owner)
	exitWindow.SetCaption("退出(&Q)")
	exitWindow.SetShortCut(api.TextToShortCut("Ctrl+Q"))
	exitWindow.SetImageIndex(imageMenu.ImageIndex("menu_exit.png"))
	exitWindow.SetOnClick(func(sender lcl.IObject) {
		logs.Debug("退出")
		MainWindow.Close()
	})
	m.file.Add(exitWindow)
}

func (m *TMainMenu) fileHistoryProjectMenu() {
	m.history.Clear()
	for _, project := range config.Config.HistoryProject {
		egpFilePath := project
		item := lcl.NewMenuItem(m.history)
		item.SetCaption(egpFilePath)
		item.SetOnClick(func(sender lcl.IObject) {
			logs.Debug("打开项目:", egpFilePath)
			event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: egpFilePath}})
		})
		m.history.Add(item)
	}
}

func (m *TMainMenu) editMenu(owner lcl.IComponent) {
	cutAction := lcl.NewEditCut(m.actionList)
	cutAction.SetShortCut(api.TextToShortCut(platformShortcut("X")))
	cutAction.SetCaption("剪切")
	cutAction.SetOnExecute(func(sender lcl.IObject) {
		activeControl := lcl.Screen.ActiveControl()
		cutAction.ExecuteTarget(activeControl) // 示例: 执行默认行为
	})

	copyAction := lcl.NewEditCopy(m.actionList)
	copyAction.SetShortCut(api.TextToShortCut(platformShortcut("C")))
	copyAction.SetCaption("复制")

	pasteAction := lcl.NewEditPaste(m.actionList)
	pasteAction.SetShortCut(api.TextToShortCut(platformShortcut("V")))
	pasteAction.SetCaption("粘贴")

	selectAllAction := lcl.NewEditSelectAll(m.actionList)
	selectAllAction.SetShortCut(api.TextToShortCut(platformShortcut("A")))
	selectAllAction.SetCaption("全选")

	undoAction := lcl.NewEditUndo(m.actionList)
	undoAction.SetShortCut(api.TextToShortCut(platformShortcut("Z")))
	undoAction.SetCaption("撤销")

	undoDelete := lcl.NewEditUndo(m.actionList)
	undoDelete.SetShortCut(api.TextToShortCut(platformShortcut("Del")))
	undoDelete.SetCaption("删除")

	m.editCut = lcl.NewMenuItem(owner)
	m.editCut.SetCaption(cutAction.Caption())
	m.editCut.SetAction(cutAction)
	m.edit.Add(m.editCut)

	m.editCopy = lcl.NewMenuItem(owner)
	m.editCopy.SetCaption(copyAction.Caption())
	m.editCopy.SetAction(copyAction)
	m.edit.Add(m.editCopy)

	m.editPaste = lcl.NewMenuItem(owner)
	m.editPaste.SetCaption(pasteAction.Caption())
	m.editPaste.SetAction(pasteAction)
	m.edit.Add(m.editPaste)

	separator1 := lcl.NewMenuItem(owner)
	separator1.SetCaption("-")
	m.edit.Add(separator1)

	m.editSelectAll = lcl.NewMenuItem(owner)
	m.editSelectAll.SetCaption(selectAllAction.Caption())
	m.editSelectAll.SetAction(selectAllAction)
	m.edit.Add(m.editSelectAll)

	m.editUndo = lcl.NewMenuItem(owner)
	m.editUndo.SetCaption(selectAllAction.Caption())
	m.editUndo.SetAction(undoAction)
	m.edit.Add(m.editUndo)

	m.editDel = lcl.NewMenuItem(owner)
	m.editDel.SetCaption(undoDelete.Caption())
	m.editDel.SetAction(undoDelete)
	m.edit.Add(m.editDel)
}

func platformShortcut(key string) string {
	if desTool.IsDarwin {
		return "Meta+" + key
	}
	return "Ctrl+" + key
}

func (m *TMainMenu) viewMenu(owner lcl.IComponent) {
	windowLayout := config.Config.WindowLayout
	m.viewWidgets = lcl.NewMenuItem(owner)
	m.viewWidgets.SetCaption("组件库")
	m.viewWidgets.SetChecked(windowLayout.MenuView.WidgetsChecked) // 动态控制
	m.viewWidgets.SetOnClick(func(sender lcl.IObject) {
		if MainWindow.contentLayout != nil {
			checked := !m.viewWidgets.Checked()
			m.viewWidgets.SetChecked(checked)
			MainWindow.contentLayout.widgetPanel.SetVisible(checked)
			MainWindow.contentLayout.widgetSplitter.SetVisible(checked)
		}
	})
	m.view.Add(m.viewWidgets)

	m.viewProject = lcl.NewMenuItem(owner)
	m.viewProject.SetCaption("项目管理器")
	m.viewProject.SetChecked(windowLayout.MenuView.ProjectChecked) // 动态控制
	m.viewProject.SetOnClick(func(sender lcl.IObject) {
		if MainWindow.contentLayout != nil {
			checked := !m.viewProject.Checked()
			m.viewProject.SetChecked(checked)
			MainWindow.contentLayout.projectPanel.SetVisible(checked)
			MainWindow.contentLayout.projectSplitter.SetVisible(checked)
		}
	})
	m.view.Add(m.viewProject)

	m.viewInspector = lcl.NewMenuItem(owner)
	m.viewInspector.SetCaption("属性检查器")
	m.viewInspector.SetChecked(windowLayout.MenuView.InspectorChecked) // 动态控制
	m.viewInspector.SetOnClick(func(sender lcl.IObject) {
		if MainWindow.contentLayout != nil {
			checked := !m.viewInspector.Checked()
			m.viewInspector.SetChecked(checked)
			MainWindow.contentLayout.inspectorPanel.SetVisible(checked)
			MainWindow.contentLayout.inspectorSplitter.SetVisible(checked)
		}
	})
	m.view.Add(m.viewInspector)

	m.viewConsole = lcl.NewMenuItem(owner)
	m.viewConsole.SetCaption("日志")
	m.viewConsole.SetChecked(windowLayout.MenuView.ConsoleChecked) // 动态控制
	m.viewConsole.SetOnClick(func(sender lcl.IObject) {
		if MainWindow.contentLayout != nil {
			checked := !m.viewConsole.Checked()
			m.viewConsole.SetChecked(checked)
			MainWindow.contentLayout.consoleLogPanel.SetVisible(checked)
			MainWindow.contentLayout.consoleLogSplitter.SetVisible(checked)
		}
	})
	m.view.Add(m.viewConsole)

	m.viewStatusbar = lcl.NewMenuItem(owner)
	m.viewStatusbar.SetCaption("状态栏")
	m.viewStatusbar.SetChecked(windowLayout.MenuView.StatusbarChecked) // 动态控制
	m.viewStatusbar.SetOnClick(func(sender lcl.IObject) {
		if MainWindow.contentLayout != nil {
			checked := !m.viewStatusbar.Checked()
			m.viewStatusbar.SetChecked(checked)
			MainWindow.contentLayout.contentStatus.SetVisible(checked)
		}
	})
	m.view.Add(m.viewStatusbar)

}

func (m *TMainMenu) runMenu(owner lcl.IComponent) {
	m.build = lcl.NewMenuItem(owner)
	m.build.SetCaption("构建")
	m.build.SetImageIndex(imageMenu.ImageIndex("menu_build.png"))
	m.build.SetShortCut(api.TextToShortCut("Ctrl+F8"))
	m.build.SetOnClick(func(lcl.IObject) {
		event.ConsoleWriteInfo("Build Start")
		SetEnableFuncComponent(false)
		go func() {
			build.Run()
			SetEnableFuncComponent(true)
			event.ConsoleWriteInfo("Build End")
		}()
	})
	m.run.Add(m.build)

	m.buildClean = lcl.NewMenuItem(owner)
	m.buildClean.SetCaption("清理构建")
	m.buildClean.SetImageIndex(imageMenu.ImageIndex("menu_build_clean.png"))
	m.buildClean.SetShortCut(api.TextToShortCut("Ctrl+Shift+F8"))
	m.buildClean.SetOnClick(func(lcl.IObject) {
		event.ConsoleWriteInfo("Build Clean Start")
		SetEnableFuncComponent(false)
		go func() {
			build.RunClean()
			SetEnableFuncComponent(true)
			event.ConsoleWriteInfo("Build Clean End")
		}()
	})
	m.run.Add(m.buildClean)

	m.buildAll = lcl.NewMenuItem(owner)
	m.buildAll.SetCaption("构建所有")
	m.buildAll.SetImageIndex(imageMenu.ImageIndex("menu_build.png"))
	m.buildAll.SetShortCut(api.TextToShortCut("Ctrl+Shift+F9"))
	m.buildAll.SetOnClick(func(lcl.IObject) {
		event.ConsoleWriteInfo("Build ALL Start")
		SetEnableFuncComponent(false)
		go func() {
			build.RunAll()
			SetEnableFuncComponent(true)
			event.ConsoleWriteInfo("Build ALL End")
		}()
	})
	m.run.Add(m.buildAll)

	sep := lcl.NewMenuItem(owner)
	sep.SetCaption("-")
	m.run.Add(sep)

	m.runApp = lcl.NewMenuItem(owner)
	m.runApp.SetCaption("运行应用")
	m.runApp.SetImageIndex(imageMenu.ImageIndex("menu_run.png"))
	m.runApp.SetShortCut(api.TextToShortCut("F9"))
	m.runApp.SetOnClick(func(lcl.IObject) {
		logs.Debug("运行")
		MainWindow.toolLayout.toolbarBtn.onRunPreviewForm(m.runApp)
	})
	m.run.Add(m.runApp)
}

func (m *TMainMenu) switchRunMenuItem(status consts.PreviewState) {
	m.runApp.SetEnabled(true)
	if status == consts.PsStarted {
		m.runApp.SetCaption("停止")
		m.runApp.SetImageIndex(imageMenu.ImageIndex("menu_stop.png"))
	} else if status == consts.PsStarting {
		m.runApp.SetEnabled(false)
		m.runApp.SetCaption("停止")
		m.runApp.SetImageIndex(imageMenu.ImageIndex("menu_stop.png"))
	} else {
		m.runApp.SetCaption("运行")
		m.runApp.SetImageIndex(imageMenu.ImageIndex("menu_run.png"))
	}
}

func (m *TMainMenu) settingMenu(owner lcl.IComponent) {
	m.buildOption = lcl.NewMenuItem(owner)
	m.buildOption.SetCaption("构建选项")
	m.buildOption.SetImageIndex(imageMenu.ImageIndex("menu_compile.png"))
	m.buildOption.SetShortCut(api.TextToShortCut("Ctrl+F9"))
	m.buildOption.SetOnClick(func(lcl.IObject) {
		logs.Debug("构建选项")
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.BuildConfig}})
	})
	m.setting.Add(m.buildOption)

	m.environmentOption = lcl.NewMenuItem(owner)
	m.environmentOption.SetCaption("环境配置")
	m.environmentOption.SetImageIndex(imageMenu.ImageIndex("menu_environment_options.png"))
	m.environmentOption.SetShortCut(api.TextToShortCut("Ctrl+F10"))
	m.environmentOption.SetOnClick(func(lcl.IObject) {
		logs.Debug("环境配置")
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.EnvConfig}})
	})
	m.setting.Add(m.environmentOption)

	//m.frameworkOption = lcl.NewMenuItem(owner)
	//m.frameworkOption.SetCaption("框架安装目录")
	//m.frameworkOption.SetImageIndex(imageMenu.ImageIndex("menu_environment_options.png"))
	//m.frameworkOption.SetOnClick(func(lcl.IObject) {
	//	logs.Debug("框架安装目录")
	//	form := NewInstallFrameworkForm()
	//	form.ShowModal()
	//})
	//m.setting.Add(m.frameworkOption)

	m.projectOption = lcl.NewMenuItem(owner)
	m.projectOption.SetCaption("应用配置")
	m.projectOption.SetImageIndex(imageMenu.ImageIndex("menu_app_config.png"))
	m.projectOption.SetShortCut(api.TextToShortCut("Ctrl+F11"))
	m.projectOption.SetOnClick(func(lcl.IObject) {
		logs.Debug("应用配置")
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectConfig}})
	})
	m.setting.Add(m.projectOption)
}

func (m *TMainMenu) helperMenu(owner lcl.IComponent) {
	_, _, _, _, _, v := api.LCLVersion()
	about := lcl.NewMenuItem(owner)
	about.SetCaption("关于")
	about.SetImageIndex(imageMenu.ImageIndex("menu_project_about.png"))
	about.SetOnClick(func(sender lcl.IObject) {
		versionInfo := api.PasStr("ENERGY Designer " + config.DesignerConfig.Version + "\nLCL " + v)
		lcl.Application.MessageBox(versionInfo, versionInfo, 0)
	})
	m.helper.Add(about)
}
