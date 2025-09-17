package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

var (
	designer *Designer
)

// 窗体设计功能

type Designer struct {
	page    lcl.IPageControl // 设计器 tabs
	tabMenu lcl.IPopupMenu   // tab 菜单
	forms   map[int]*FormTab // 设计器窗体列表
}

type FormTab struct {
	id          int           // 索引, 关联 forms key: index
	name        string        // 窗体名称
	sheet       lcl.ITabSheet // tab sheet
	designerBox lcl.IPanel    // 设计器
}

// 创建主窗口设计器的布局
func (m *BottomBox) createFromDesignerLayout() *Designer {
	des := new(Designer)
	des.forms = make(map[int]*FormTab)
	des.page = lcl.NewPageControl(m.box)
	des.page.SetParent(m.rightBox)
	des.page.SetAlign(types.AlClient)
	des.page.SetTabStop(true)

	des.page.SetOnContextPopup(func(sender lcl.IObject, mousePos types.TPoint, handled *bool) {

	})

	// 创建tab上的右键菜单
	des.createTabMenu()
	return des
}

// 创建tab上的右键菜单
func (m *Designer) createTabMenu() {
	if m.tabMenu != nil {
		return
	}
	m.tabMenu = lcl.NewPopupMenu(m.page)
	m.tabMenu.SetImages(LoadImageList(m.page, []string{"actions/laz_cancel.png"}, 16, 16))
	items := m.tabMenu.Items()
	closeMenuItem := lcl.NewMenuItem(m.page)
	closeMenuItem.SetCaption("关闭窗体")
	closeMenuItem.SetImageIndex(0)
	items.Add(closeMenuItem)

	//m.page.SetPopupMenu(m.tabMenu)
}

// 创建窗体设计器 tab
func (m *Designer) newFormDesignerTab() *FormTab {
	form := new(FormTab)
	id := len(m.forms) + 1
	formName := fmt.Sprintf("Form%d", id) // 默认名
	form.name = formName
	form.id = id
	m.forms[id] = form

	form.sheet = lcl.NewTabSheet(m.page)
	form.sheet.SetParent(m.page)
	form.sheet.SetCaption(formName)
	form.sheet.SetAlign(types.AlClient)

	scroll := lcl.NewScrollBox(form.sheet)
	scroll.SetParent(form.sheet)
	scroll.SetAlign(types.AlClient)
	scroll.SetAutoScroll(true)
	scroll.SetBorderStyleToBorderStyle(types.BsNone)
	scroll.SetColor(colors.ClWhite)
	scroll.SetOnPaint(func(sender lcl.IObject) {

	})

	form.designerBox = lcl.NewPanel(scroll)
	form.designerBox.SetParent(scroll)
	form.designerBox.SetBevelOuter(types.BvNone)
	form.designerBox.SetDoubleBuffered(true)
	form.designerBox.SetParentColor(false)
	form.designerBox.SetColor(colors.ClBtnFace)
	form.designerBox.SetLeft(40)
	form.designerBox.SetTop(40)
	form.designerBox.SetWidth(600)
	form.designerBox.SetHeight(400)
	form.designerBox.SetAlign(types.AlCustom)
	form.designerBox.SetOnPaint(form.OnPaint)
	form.designerBox.SetOnMouseMove(form.OnMouseMove)
	form.designerBox.SetOnMouseDown(form.OnMouseDown)
	form.designerBox.SetOnMouseUp(form.OnMouseUp)
	return form
}

// 激活指定的 tab
func (m *Designer) ActiveFormTab(tab *FormTab) {
	m.page.SetActivePage(tab.sheet)
}

func (m *FormTab) OnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

}

func (m *FormTab) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

}

func (m *FormTab) OnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {

}

func (m *FormTab) OnPaint(sender lcl.IObject) {

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
