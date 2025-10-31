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
	fileCreateWindow := lcl.NewMenuItem(owner)
	fileCreateWindow.SetCaption("新建窗体(&N)")
	fileCreateWindow.SetShortCut(api.TextToShortCut("Ctrl+N"))
	fileCreateWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("单击了新建窗体")
	})
	m.file.Add(fileCreateWindow)
}

func (m *TMainMenu) editMenu(owner lcl.IComponent) {

}

func (m *TMainMenu) settingMenu(owner lcl.IComponent) {

}

func (m *TMainMenu) projectMenu(owner lcl.IComponent) {

}

func (m *TMainMenu) helperMenu(owner lcl.IComponent) {
	_, _, _, _, _, v := api.LCLVersion()
	helperAbout := lcl.NewMenuItem(owner)
	helperAbout.SetCaption("关于")
	helperAbout.SetOnClick(func(sender lcl.IObject) {
		versionInfo := api.PasStr("ENERGY Designer " + config.Config.Version + "\nLCL " + v)
		lcl.Application.MessageBox(versionInfo, versionInfo, 0)
	})
	m.helper.Add(helperAbout)
}
