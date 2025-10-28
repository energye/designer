package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

var (
	designer                    *Designer
	margin                      int32 = 0
	borderWidth                 int32 = 8
	defaultWidth, defaultHeight int32 = 600, 400
)

// 窗体设计功能

type Designer struct {
	page          lcl.IPageControl // 设计器 tabs
	tabMenu       lcl.IPopupMenu   // tab 菜单
	designerForms map[int]*FormTab // 设计器窗体列表
}

// 获取所有设计窗体
func GetDesignerForms() map[int]*FormTab {
	return designer.designerForms
}

// 创建主窗口设计器的布局
func (m *BottomBox) createFromDesignerLayout() *Designer {
	des := new(Designer)
	des.designerForms = make(map[int]*FormTab)
	des.page = lcl.NewPageControl(m.box)
	des.page.SetParent(m.rightBox)
	des.page.SetAlign(types.AlClient)
	des.page.SetTabStop(true)

	des.page.SetOnContextPopup(func(sender lcl.IObject, mousePos types.TPoint, handled *bool) {

	})
	des.page.SetOnClick(func(sender lcl.IObject) {
		logs.Debug("Designer PageControl click")
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
	m.tabMenu.SetImages(imageActions.ImageList100())
	items := m.tabMenu.Items()
	closeMenuItem := lcl.NewMenuItem(m.page)
	closeMenuItem.SetCaption("关闭窗体")
	closeMenuItem.SetImageIndex(imageActions.ImageIndex("laz_cancel.png"))
	items.Add(closeMenuItem)

	//m.page.SetPopupMenu(m.tabMenu)
}

func (m *Designer) hideFormTabs() {
	for _, formTab := range m.designerForms {
		formTab.tree.SetVisible(false)
	}
}

// 添加一个窗体设计器 tab
func (m *Designer) addDesignerFormTab() *FormTab {
	m.hideFormTabs()
	form := new(FormTab)
	form.componentName = make(map[string]int)
	// 组件树
	form.tree = lcl.NewTreeView(inspector.componentTree.treeBox)
	form.tree.SetAutoExpand(true)
	form.tree.SetReadOnly(true)
	form.tree.SetDoubleBuffered(true)
	//m.tree.SetMultiSelect(true) // 多选控制
	form.tree.SetTop(35)
	//form.tree.SetWidth(leftBoxWidth)
	form.tree.SetWidth(inspector.componentTree.treeBox.Width())
	//form.tree.SetHeight(componentTreeHeight - form.tree.Top())
	form.tree.SetHeight(inspector.componentTree.treeBox.Height() - form.tree.Top())
	form.tree.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
	form.tree.SetAlign(types.AlCustom)
	//form.tree.SetAlign(types.AlClient)
	form.tree.SetVisible(true)
	form.tree.SetImages(imageComponents.ImageList100())
	//form.tree.SetOnGetSelectedIndex(form.TreeOnGetSelectedIndex)
	form.tree.SetOnGetSelectedIndex(form.TreeOnGetSelectedIndex)
	form.tree.SetOnMouseDown(form.TreeOnMouseDown)
	form.tree.SetOnContextPopup(form.TreeOnContextPopup)
	// 树菜单
	form.tree.SetPopupMenu(form.TreePopupMenu())
	form.tree.SetParent(inspector.componentTree.treeBox)

	// 默认名
	form.id = len(m.designerForms) + 1
	form.name = fmt.Sprintf("Form%v", form.id)
	// 窗体ID
	m.designerForms[form.id] = form

	form.sheet = lcl.NewTabSheet(m.page)
	form.sheet.SetCaption(form.name)
	form.sheet.SetOnHide(form.tabSheetOnHide)
	form.sheet.SetOnShow(form.tabSheetOnShow)
	//form.sheet.SetAlign(types.AlClient)
	form.sheet.SetParent(m.page)

	form.scroll = lcl.NewScrollBox(form.sheet)
	form.scroll.SetAlign(types.AlClient)
	form.scroll.SetAutoScroll(true)
	form.scroll.SetBorderStyleToBorderStyle(types.BsNone)
	form.scroll.SetDoubleBuffered(true)
	form.scroll.SetParent(form.sheet)

	//newStatusBar(form.scroll)

	// 创建设计窗体
	form.NewFormDesigner()

	return form
}

// 激活指定的 tab
func (m *Designer) ActiveFormTab(tab *FormTab) {
	m.page.SetActivePage(tab.sheet)
	for _, form := range m.designerForms {
		form.isDesigner = false
	}
	tab.isDesigner = true
}

// 绘制刻度尺, 在外层 scroll 上
//
//	func (m *FormTab) scrollDrawRuler() {
//		gridSize := 5 // 小刻度
//		//canvas := m.bg.Canvas()
//		canvas := m.scroll.Canvas()
//		canvas.PenToPen().SetColor(colors.ClBlack)
//		width, height := m.formRoot.Width(), m.formRoot.Height()
//		println("width, height:", width, height)
//		// X
//		for i := 0; i <= int(width)/gridSize; i++ {
//			x := int32(i * gridSize)
//			x = x + margin
//			if i%10 == 0 { // 长
//				canvas.LineWithIntX4(x, margin-35, x, margin-10)
//				text := strconv.Itoa(i * gridSize)
//				textWidth := canvas.TextWidthWithUnicodestring(text)
//				canvas.TextOutWithIntX2Unicodestring(x-(textWidth/2), 0, text)
//			} else if i%5 == 0 { // 中
//				canvas.LineWithIntX4(x, margin-25, x, margin-10)
//			} else { // 小
//				canvas.LineWithIntX4(x, margin-15, x, margin-10)
//			}
//		}
//		// Y
//		for i := 0; i <= int(height)/gridSize; i++ {
//			y := int32(i * gridSize)
//			y = y + margin
//			if i%10 == 0 { // 长
//				canvas.LineWithIntX4(margin-35, y, margin-10, y)
//				text := strconv.Itoa(i * gridSize)
//				textWidth := canvas.TextWidthWithUnicodestring(text)
//				canvas.TextOutWithIntX2Unicodestring(0, y-(textWidth/2), text)
//			} else if i%5 == 0 { // 中
//				canvas.LineWithIntX4(margin-25, y, margin-10, y)
//			} else { // 小
//				canvas.LineWithIntX4(margin-15, y, margin-10, y)
//			}
//		}
//	}

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
