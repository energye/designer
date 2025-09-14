package designer

import (
	"fmt"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
)

// 创建设计器布局
func (m *TAppWindow) createDesignerLayout() {
	// 顶部菜单
	m.createMenu()
	// 工具栏
	m.createToolbar()
}

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
		fmt.Println("单击了新建窗体")
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

	if tool.IsDarwin() {

	}
}

// 创建工具栏
func (m *TAppWindow) createToolbar() {

}
