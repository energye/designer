package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
)

// 顶部菜单
func (m *TAppWindow) createMenu() {
	m.mainMenu = lcl.NewMainMenu(m)
	file := lcl.NewMenuItem(m)
	file.SetCaption("文件(&F)")
	m.mainMenu.Items().Add(file)
	fileCreateWindow := lcl.NewMenuItem(m)
	fileCreateWindow.SetCaption("新建窗体(&N)")
	fileCreateWindow.SetShortCut(api.TextToShortCut("Ctrl+N"))
	fileCreateWindow.SetOnClick(func(lcl.IObject) {
		logs.Debug("单击了新建窗体")
	})
	file.Add(fileCreateWindow)

	edit := lcl.NewMenuItem(m)
	edit.SetCaption("编辑(&E)")
	m.mainMenu.Items().Add(edit)

	setting := lcl.NewMenuItem(m)
	setting.SetCaption("设置(&S)")
	m.mainMenu.Items().Add(setting)

	helper := lcl.NewMenuItem(m)
	helper.SetCaption("帮助(&H)")
	m.mainMenu.Items().Add(helper)
	helperAbout := lcl.NewMenuItem(m)
	helperAbout.SetCaption("关于")
	helperAbout.SetOnClick(func(sender lcl.IObject) {
		versionInfo := api.PasStr("ENERGY Designer " + config.Config.Version)
		lcl.Application.MessageBox(versionInfo, versionInfo, 0)
	})
	helper.Add(helperAbout)

	if tool.IsDarwin() {

	}
}
