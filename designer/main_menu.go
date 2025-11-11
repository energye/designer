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
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
	"github.com/energye/lcl/tool/exec"
)

// 设计器顶部菜单

type TMainMenu struct {
	main    lcl.IMainMenu
	file    lcl.IMenuItem
	edit    lcl.IMenuItem
	run     lcl.IMenuItem
	setting lcl.IMenuItem
	helper  lcl.IMenuItem
	runItem lcl.IMenuItem

	createProject lcl.IMenuItem
	createWindow  lcl.IMenuItem
	open          lcl.IMenuItem
	save          lcl.IMenuItem
}

// 设计器主菜单
func (m *TAppWindow) createMainMenu() {
	if m.mainMenu != nil {
		return
	}
	mainMenu := new(TMainMenu)
	m.mainMenu = mainMenu
	mainMenu.main = lcl.NewMainMenu(m)
	mainMenu.main.SetImages(imageMenu.ImageList100())
	menuItems := mainMenu.main.Items()

	mainMenu.file = lcl.NewMenuItem(m)
	mainMenu.file.SetCaption("文件(&F)")
	menuItems.Add(mainMenu.file)

	mainMenu.edit = lcl.NewMenuItem(m)
	mainMenu.edit.SetCaption("编辑(&E)")
	menuItems.Add(mainMenu.edit)

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
	mainMenu.runMenu(m)
	mainMenu.settingMenu(m)
	mainMenu.helperMenu(m)
	mainMenu.macOS()
}

func (m *TMainMenu) macOS() {
	if tool.IsDarwin() {
		// macOS
	}
}

func (m *TMainMenu) fileMenu(owner lcl.IComponent) {
	create := lcl.NewMenuItem(owner)
	create.SetCaption("新建(&N)")
	m.file.Add(create)

	m.createProject = lcl.NewMenuItem(owner)
	m.createProject.SetCaption("新建项目")
	m.createProject.SetShortCut(api.TextToShortCut("Ctrl+P"))
	m.createProject.SetImageIndex(imageMenu.ImageIndex("menu_project_add.png"))
	m.createProject.SetOnClick(func(lcl.IObject) {
		mainWindow.selectDirectoryDialog.SetTitle("新建项目")
		history := mainWindow.selectDirectoryDialog.HistoryList()
		//for i := int32(0); i < history.Count(); i++ {
		//}
		if history.Count() == 0 {
			mainWindow.selectDirectoryDialog.SetInitialDir(exec.Dir)
		}
		if mainWindow.selectDirectoryDialog.Execute() {
			dir := mainWindow.selectDirectoryDialog.FileName()
			event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectCreate, Data: dir}})
		}
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
	})
	create.Add(m.createWindow)

	m.open = lcl.NewMenuItem(owner)
	m.open.SetCaption("打开(&O)")
	m.open.SetShortCut(api.TextToShortCut("Ctrl+O"))
	m.open.SetImageIndex(imageMenu.ImageIndex("menu_project_open.png"))
	m.open.SetOnClick(func(sender lcl.IObject) {
		toolbar.toolbarBtn.onOpenForm(sender)
	})
	m.file.Add(m.open)

	m.save = lcl.NewMenuItem(owner)
	m.save.SetCaption("保存(&S)")
	m.save.SetShortCut(api.TextToShortCut("Ctrl+S"))
	m.save.SetImageIndex(imageMenu.ImageIndex("menu_save.png"))
	m.save.SetOnClick(func(sender lcl.IObject) {
		logs.Debug("保存窗体")
	})
	m.file.Add(m.save)

	//saveAllWindow := lcl.NewMenuItem(owner)
	//saveAllWindow.SetCaption("保存所有窗体(&L)")
	//saveAllWindow.SetShortCut(api.TextToShortCut("Shift+Ctrl+L"))
	//saveAllWindow.SetImageIndex(imageMenu.ImageIndex("menu_save_all.png"))
	//saveAllWindow.SetOnClick(func(sender lcl.IObject) {
	//	logs.Debug("保存所有窗体")
	//})
	//m.file.Add(saveAllWindow)

	exitWindow := lcl.NewMenuItem(owner)
	exitWindow.SetCaption("退出(&Q)")
	exitWindow.SetShortCut(api.TextToShortCut("Ctrl+Q"))
	exitWindow.SetImageIndex(imageMenu.ImageIndex("menu_exit.png"))
	exitWindow.SetOnClick(func(sender lcl.IObject) {
		logs.Debug("退出")
	})
	m.file.Add(exitWindow)
}

func (m *TMainMenu) editMenu(owner lcl.IComponent) {
}

func (m *TMainMenu) runMenu(owner lcl.IComponent) {
	build := lcl.NewMenuItem(owner)
	build.SetCaption("构建")
	build.SetImageIndex(imageMenu.ImageIndex("menu_build.png"))
	build.SetOnClick(func(lcl.IObject) {
		logs.Debug("构建")
	})
	m.run.Add(build)

	cleanBuild := lcl.NewMenuItem(owner)
	cleanBuild.SetCaption("清理构建")
	cleanBuild.SetImageIndex(imageMenu.ImageIndex("menu_build_clean.png"))
	cleanBuild.SetOnClick(func(lcl.IObject) {
		logs.Debug("清理构建")
	})
	m.run.Add(cleanBuild)

	sep := lcl.NewMenuItem(owner)
	sep.SetCaption("-")
	m.run.Add(sep)

	m.runItem = lcl.NewMenuItem(owner)
	m.runItem.SetCaption("运行")
	m.runItem.SetImageIndex(imageMenu.ImageIndex("menu_run.png"))
	m.runItem.SetOnClick(func(lcl.IObject) {
		logs.Debug("运行")
		toolbar.toolbarBtn.onRunPreviewForm(m.runItem)
	})
	m.run.Add(m.runItem)
}

func (m *TMainMenu) switchRunMenuItem(status consts.PreviewState) {
	m.runItem.SetEnabled(true)
	if status == consts.PsStarted {
		m.runItem.SetCaption("停止")
		m.runItem.SetImageIndex(imageMenu.ImageIndex("menu_stop.png"))
	} else if status == consts.PsStarting {
		m.runItem.SetEnabled(false)
		m.runItem.SetCaption("停止")
		m.runItem.SetImageIndex(imageMenu.ImageIndex("menu_stop.png"))
	} else {
		m.runItem.SetCaption("运行")
		m.runItem.SetImageIndex(imageMenu.ImageIndex("menu_run.png"))
	}
}

func (m *TMainMenu) settingMenu(owner lcl.IComponent) {
	buildOption := lcl.NewMenuItem(owner)
	buildOption.SetCaption("构建选项")
	buildOption.SetImageIndex(imageMenu.ImageIndex("menu_compile.png"))
	buildOption.SetOnClick(func(lcl.IObject) {
		logs.Debug("构建选项")
	})
	m.setting.Add(buildOption)

	environmentOption := lcl.NewMenuItem(owner)
	environmentOption.SetCaption("环境配置")
	environmentOption.SetImageIndex(imageMenu.ImageIndex("menu_environment_options_200.png"))
	environmentOption.SetOnClick(func(lcl.IObject) {
		logs.Debug("环境配置")
	})
	m.setting.Add(environmentOption)

	projectOption := lcl.NewMenuItem(owner)
	projectOption.SetCaption("项目配置")
	projectOption.SetImageIndex(imageMenu.ImageIndex("menu_environment_options_200.png"))
	projectOption.SetOnClick(func(lcl.IObject) {
		logs.Debug("项目配置")
	})
	m.setting.Add(projectOption)
}

func (m *TMainMenu) helperMenu(owner lcl.IComponent) {
	_, _, _, _, _, v := api.LCLVersion()
	about := lcl.NewMenuItem(owner)
	about.SetCaption("关于")
	about.SetOnClick(func(sender lcl.IObject) {
		versionInfo := api.PasStr("ENERGY Designer " + config.Config.Version + "\nLCL " + v)
		lcl.Application.MessageBox(versionInfo, versionInfo, 0)
	})
	m.helper.Add(about)
}
