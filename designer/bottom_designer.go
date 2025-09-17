package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

var (
	designer *Designer
)

// 窗体设计功能

type Designer struct {
	page    lcl.IPageControl    // 设计器 tabs
	tabMenu lcl.IPopupMenu      // tab 菜单
	forms   map[string]*FormTab // 设计器窗体列表
}

type FormTab struct {
}

func (m *BottomBox) createFromDesignerLayout() *Designer {
	des := new(Designer)
	des.forms = make(map[string]*FormTab)
	des.page = lcl.NewPageControl(m.box)
	des.page.SetParent(m.rightBox)
	des.page.SetAlign(types.AlClient)
	des.page.SetTabStop(true)

	// 创建tab上的右键菜单
	des.createTabMenu()
	return des
}

// 创建tab上的右键菜单
func (m *Designer) createTabMenu() {
	m.tabMenu = lcl.NewPopupMenu(m.page)
	m.tabMenu.SetImages(LoadImageList(m.page, []string{"actions/laz_cancel.png"}, 16, 16))
	items := m.tabMenu.Items()
	closeMenuItem := lcl.NewMenuItem(m.page)
	closeMenuItem.SetCaption("关闭")
	closeMenuItem.SetImageIndex(0)
	items.Add(closeMenuItem)

	m.page.SetPopupMenu(m.tabMenu)
}

// 创建设计器 tab
func (m *Designer) newFormDesignerTab() {

}

// 测试属性
//btn := lcl.NewButton(m)
//btn.SetParent(m)
//testProp := lcl.GetComponentProperties(btn)
//for _, prop := range testProp {
//	fmt.Printf("%+v\n", prop)
//}
//
//font := lcl.NewFont()
//fmt.Println("font.GetNamePath():", font.GetNamePath())
//fs := lcl.Screen.Fonts()
//for i := 0; i < int(fs.Count()); i++ {
//	fmt.Println(fs.Strings(int32(i)))
//}
