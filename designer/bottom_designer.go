package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

var (
	designer                    *Designer
	margin                      int32 = 5
	borderWidth                 int32 = 8
	defaultWidth, defaultHeight int32 = 600, 400
)

// 窗体设计功能

type Designer struct {
	page          lcl.IPageControl // 设计器 tabs
	tabMenu       lcl.IPopupMenu   // tab 菜单
	designerForms map[int]*FormTab // 设计器窗体列表
}

type FormTab struct {
	id                   int                   // 索引, 关联 forms key: index
	name                 string                // 窗体名称
	scroll               lcl.IScrollBox        // 外 滚动条
	isDesigner           bool                  // 是否下在设计
	sheet                lcl.ITabSheet         // tab sheet
	designerBox          *DesigningComponent   // 设计器
	isDown, isUp, isMove bool                  // 鼠标事件
	dragForm             *drag                 // 拖拽窗体控制器
	componentName        map[string]int        // 组件分类名
	form                 *DesigningComponent   // 窗体
	componentList        []*DesigningComponent // 组件列表
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

// 添加一个窗体设计器 tab
func (m *Designer) addFormDesignerTab() *FormTab {
	form := new(FormTab)
	id := len(m.designerForms) + 1
	formName := fmt.Sprintf("Form%d", id) // 默认名
	form.name = formName
	form.id = id
	form.componentName = make(map[string]int)
	m.designerForms[id] = form

	form.sheet = lcl.NewTabSheet(m.page)
	form.sheet.SetParent(m.page)
	form.sheet.SetCaption(formName)
	//form.sheet.SetAlign(types.AlClient)

	form.scroll = lcl.NewScrollBox(form.sheet)
	form.scroll.SetParent(form.sheet)
	form.scroll.SetAlign(types.AlClient)
	form.scroll.SetAutoScroll(true)
	form.scroll.SetBorderStyleToBorderStyle(types.BsNone)
	form.scroll.SetDoubleBuffered(true)
	//form.scroll.HorzScrollBar().SetIncrement(1)
	//form.scroll.VertScrollBar().SetIncrement(1)

	//newStatusBar(form.scroll)

	//form.bg = lcl.NewPanel(form.scroll)
	//form.bg.SetParent(form.scroll)
	//form.bg.SetAlign(types.AlClient)

	form.designerBox = new(DesigningComponent)
	designerBox := lcl.NewPanel(form.scroll)
	designerBox.SetParent(form.scroll)
	designerBox.SetBevelOuter(types.BvNone)
	designerBox.SetBorderStyleToBorderStyle(types.BsSingle)
	designerBox.SetDoubleBuffered(true)
	designerBox.SetParentColor(false)
	designerBox.SetColor(colors.ClBtnFace)
	designerBox.SetLeft(margin)
	designerBox.SetTop(margin)
	designerBox.SetWidth(defaultWidth)
	designerBox.SetHeight(defaultHeight)
	designerBox.SetAlign(types.AlCustom)
	designerBox.SetOnPaint(form.designerOnPaint)
	designerBox.SetOnMouseMove(form.designerOnMouseMove)
	designerBox.SetOnMouseDown(form.designerOnMouseDown)
	designerBox.SetOnMouseUp(form.designerOnMouseUp)
	form.designerBox.object = designerBox

	// 创建一个隐藏的窗体用于获取属性
	form.form = NewFormDesigner(form)

	// 窗体拖拽大小
	form.dragForm = newDrag(form.scroll, DsRightBottom)
	form.dragForm.SetRelation(form.designerBox)
	form.dragForm.Show()
	form.dragForm.Follow()

	// 测试控件
	//NewButtonDesigner(form, 50, 50)
	//NewEditDesigner(form, 150, 150)
	//NewCheckBoxDesigner(form, 150, 200)

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

func (m *FormTab) addDesignerComponent(component *DesigningComponent) {
	m.componentList = append(m.componentList, component)
}

// 隐藏所有控件的 drag
func (m *FormTab) hideAllDrag() {
	for _, component := range m.componentList {
		component.drag.Hide()
	}
}

// 窗体设计界面 鼠标按下
func (m *FormTab) designerOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	m.hideAllDrag()
	// 创建组件
	logs.Debug("鼠标点击设计器")
	if toolbar.selectComponent != nil {
		componentName := toolbar.selectComponent.name
		logs.Debug("当前选中控件:", toolbar.selectComponent.index, toolbar.selectComponent.name)
		// 创建
		if create := GetRegisterComponent(componentName); create != nil {
			create(m, x, y).LoadPropertyToInspector()
		} else {
			logs.Warn("当前选中设计组件", toolbar.selectComponent.name, "未实现或未注册")
		}
		// 重置
		toolbar.ResetTabComponentDown()
	} else {
		logs.Debug("加载窗体")
		inspector.LoadComponent(m.form)
	}
}

// 窗体设计界面 鼠标抬起
func (m *FormTab) designerOnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {

}

// 窗体设计界面 鼠标移动
func (m *FormTab) designerOnMouseMove(sender lcl.IObject, shift types.TShiftState, x, y int32) {
	//width, height := m.designerBox.Width(), m.designerBox.Height()
}

// 获取组件名 Caption
func (m *FormTab) GetComponentCaptionName(component string) string {
	if c, ok := m.componentName[component]; ok {
		m.componentName[component] = c + 1
	} else {
		m.componentName[component] = 1
	}
	component = fmt.Sprintf("%v%d", component, m.componentName[component])
	return component
}

func (m *FormTab) designerOnPaint(sender lcl.IObject) {
	// 绘制网格
	m.drawGrid()
}

// 绘制风格线
func (m *FormTab) drawGrid() {
	gridSize := 9 // 小刻度
	designerBox := m.designerBox.object.(lcl.IPanel)
	canvas := designerBox.Canvas()
	canvas.PenToPen().SetColor(colors.ClBlack)
	width, height := designerBox.Width(), designerBox.Height()
	for i := 1; i < int(width)/gridSize; i++ {
		x := int32(i * gridSize)
		for j := 1; j < int(height)/gridSize; j++ {
			y := int32(j * gridSize)
			canvas.SetPixels(x, y, colors.ClBlack)
		}
	}
}

// 绘制刻度尺, 在外层 scroll 上
//
//	func (m *FormTab) scrollDrawRuler() {
//		gridSize := 5 // 小刻度
//		//canvas := m.bg.Canvas()
//		canvas := m.scroll.Canvas()
//		canvas.PenToPen().SetColor(colors.ClBlack)
//		width, height := m.designerBox.Width(), m.designerBox.Height()
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

func SetDesignMode(component lcl.IControl) {
	lcl.DesigningComponent().SetComponentDesignMode(component, true)
	lcl.DesigningComponent().SetComponentDesignInstanceMode(component, true)
	lcl.DesigningComponent().SetComponentInlineMode(component, true)
	lcl.DesigningComponent().SetWidgetSetDesigning(component)
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
