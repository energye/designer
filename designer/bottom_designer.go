package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strconv"
)

var (
	designer *Designer
	margin   int32 = 40
)

// 窗体设计功能

type Designer struct {
	page    lcl.IPageControl // 设计器 tabs
	tabMenu lcl.IPopupMenu   // tab 菜单
	forms   map[int]*FormTab // 设计器窗体列表
}

type FormTab struct {
	id          int            // 索引, 关联 forms key: index
	name        string         // 窗体名称
	scroll      lcl.IScrollBox // 外 滚动条
	bg          lcl.IPanel
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

	form.scroll = lcl.NewScrollBox(form.sheet)
	form.scroll.SetParent(form.sheet)
	form.scroll.SetAlign(types.AlClient)
	form.scroll.SetAutoScroll(true)
	form.scroll.SetBorderStyleToBorderStyle(types.BsNone)
	//form.scroll.SetColor(colors.ClWhite)
	form.scroll.SetDoubleBuffered(true)
	form.scroll.HorzScrollBar().SetIncrement(1)
	form.scroll.VertScrollBar().SetIncrement(1)

	form.bg = lcl.NewPanel(form.scroll)
	form.bg.SetParent(form.scroll)
	form.bg.SetAlign(types.AlClient)

	form.designerBox = lcl.NewPanel(form.bg)
	form.designerBox.SetParent(form.bg)
	form.designerBox.SetBevelOuter(types.BvNone)
	form.designerBox.SetDoubleBuffered(true)
	form.designerBox.SetParentColor(false)
	form.designerBox.SetColor(colors.ClBtnFace)
	form.designerBox.SetLeft(margin)
	form.designerBox.SetTop(margin)
	form.designerBox.SetWidth(600)
	form.designerBox.SetHeight(400)
	form.designerBox.SetAlign(types.AlCustom)
	form.designerBox.SetOnPaint(form.designerOnPaint)
	form.designerBox.SetOnMouseMove(form.designerOnMouseMove)
	form.designerBox.SetOnMouseDown(form.designerOnMouseDown)
	form.designerBox.SetOnMouseUp(form.designerOnMouseUp)

	// 测试控件
	testBtn := lcl.NewButton(form.designerBox)
	testBtn.SetParent(form.designerBox)
	testBtn.SetLeft(50)
	testBtn.SetTop(50)
	testBtn.SetCaption("测试按钮")

	return form
}

// 激活指定的 tab
func (m *Designer) ActiveFormTab(tab *FormTab) {
	m.page.SetActivePage(tab.sheet)
}

// 绘制刻度尺, 在外层 scroll 上
func (m *FormTab) scrollDrawRuler() {
	gridSize := 5 // 10px 一个小刻度
	canvas := m.bg.Canvas()
	canvas.PenToPen().SetColor(colors.ClBlack)
	width, height := m.designerBox.Width(), m.designerBox.Height()
	println("width, height:", width, height)
	// X
	for i := 0; i <= int(width)/gridSize; i++ {
		x := int32(i * gridSize)
		x = x + margin
		if i%10 == 0 { // 长
			canvas.LineWithIntX4(x, margin-35, x, margin-10)
			text := strconv.Itoa(i * gridSize)
			textWidth := canvas.TextWidthWithUnicodestring(text)
			canvas.TextOutWithIntX2Unicodestring(x-(textWidth/2), 0, text)
		} else if i%5 == 0 { // 中
			canvas.LineWithIntX4(x, margin-25, x, margin-10)
		} else { // 小
			canvas.LineWithIntX4(x, margin-15, x, margin-10)
		}
	}
	// Y
	for i := 0; i <= int(height)/gridSize; i++ {
		y := int32(i * gridSize)
		y = y + margin
		if i%10 == 0 { // 长
			canvas.LineWithIntX4(margin-35, y, margin-10, y)
			text := strconv.Itoa(i * gridSize)
			textWidth := canvas.TextWidthWithUnicodestring(text)
			canvas.TextOutWithIntX2Unicodestring(0, y-(textWidth/2), text)
		} else if i%5 == 0 { // 中
			canvas.LineWithIntX4(margin-25, y, margin-10, y)
		} else { // 小
			canvas.LineWithIntX4(margin-15, y, margin-10, y)
		}
	}
}

func (m *FormTab) designerOnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

}

func (m *FormTab) designerOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

}

func (m *FormTab) designerOnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {

}

func (m *FormTab) designerOnPaint(sender lcl.IObject) {
	//println("designerOnPaint")
	// 绘制刻度
	m.scrollDrawRuler()
	// 绘制网格
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
