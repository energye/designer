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
}

func (m *TAppWindow) createMenu() {
	m.mainMenu = lcl.NewMainMenu(m)
	file := lcl.NewMenuItem(m)
	file.SetCaption("文件(&F)")
	fileCreate := lcl.NewMenuItem(m)
	fileCreate.SetCaption("新建(&N)")
	fileCreate.SetShortCut(api.TextToShortCut("Ctrl+N"))
	fileCreate.SetOnClick(func(lcl.IObject) {
		fmt.Println("单击了新建")
	})
	file.Add(fileCreate)

	m.mainMenu.Items().Add(file)

	if tool.IsDarwin() {

	}
}
