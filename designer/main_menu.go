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
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
)

// 设计器顶部菜单

type TMainMenu struct {
	main    lcl.IMainMenu
	file    lcl.IMenuItem
	edit    lcl.IMenuItem
	setting lcl.IMenuItem
	project lcl.IMenuItem
	helper  lcl.IMenuItem
}

// 设计器主菜单
func (m *TAppWindow) createMainMenu() {
	if m.mainMenu != nil {
		return
	}
	mainMenu := new(TMainMenu)
	m.mainMenu = mainMenu
	mainMenu.main = lcl.NewMainMenu(m)
	menuItems := mainMenu.main.Items()

	mainMenu.file = lcl.NewMenuItem(m)
	mainMenu.file.SetCaption("文件(&F)")
	menuItems.Add(mainMenu.file)

	mainMenu.edit = lcl.NewMenuItem(m)
	mainMenu.edit.SetCaption("编辑(&E)")
	menuItems.Add(mainMenu.edit)

	mainMenu.setting = lcl.NewMenuItem(m)
	mainMenu.setting.SetCaption("设置(&S)")
	menuItems.Add(mainMenu.setting)

	mainMenu.project = lcl.NewMenuItem(m)
	mainMenu.project.SetCaption("项目(&P)")
	menuItems.Add(mainMenu.project)

	mainMenu.helper = lcl.NewMenuItem(m)
	mainMenu.helper.SetCaption("帮助(&H)")
	menuItems.Add(mainMenu.helper)

	mainMenu.fileMenu(m)
	mainMenu.editMenu(m)
	mainMenu.settingMenu(m)
	mainMenu.projectMenu(m)
	mainMenu.helperMenu(m)
	mainMenu.macOS()
}

func (m *TMainMenu) macOS() {
	if tool.IsDarwin() {
		// macOS
	}
}

func (m *TMainMenu) fileMenu(owner lcl.IComponent) {
	createWindow := lcl.NewMenuItem(owner)
	createWindow.SetCaption("新建窗体(&N)")
	createWindow.SetShortCut(api.TextToShortCut("Ctrl+N"))
	createWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("新建窗体")
	})
	m.file.Add(createWindow)
	openWindow := lcl.NewMenuItem(owner)
	openWindow.SetCaption("打开窗体(&O)")
	openWindow.SetShortCut(api.TextToShortCut("Ctrl+O"))
	openWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("打开窗体")
	})
	m.file.Add(openWindow)
	saveWindow := lcl.NewMenuItem(owner)
	saveWindow.SetCaption("保存窗体(&S)")
	saveWindow.SetShortCut(api.TextToShortCut("Ctrl+S"))
	saveWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("保存窗体")
	})
	m.file.Add(saveWindow)
	saveAllWindow := lcl.NewMenuItem(owner)
	saveAllWindow.SetCaption("保存所有窗体(&L)")
	saveAllWindow.SetShortCut(api.TextToShortCut("Shift+Ctrl+L"))
	saveAllWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("保存所有窗体")
	})
	m.file.Add(saveAllWindow)
	exitWindow := lcl.NewMenuItem(owner)
	exitWindow.SetCaption("退出(&Q)")
	exitWindow.SetShortCut(api.TextToShortCut("Ctrl+Q"))
	exitWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("退出")
	})
	m.file.Add(exitWindow)
}

func (m *TMainMenu) editMenu(owner lcl.IComponent) {

}

func (m *TMainMenu) settingMenu(owner lcl.IComponent) {

}

func (m *TMainMenu) projectMenu(owner lcl.IComponent) {
	createProject := lcl.NewMenuItem(owner)
	createProject.SetCaption("新建项目(&P)")
	createProject.SetShortCut(api.TextToShortCut("Ctrl+P"))
	createProject.SetOnClick(func(lcl.IObject) {
		logs.Debug("新建项目")
	})
	m.project.Add(createProject)
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
