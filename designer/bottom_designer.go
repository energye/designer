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
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"widget/wg"
)

// 窗体设计功能

var (
	designer                    *Designer
	margin                      int32 = 0
	borderWidth                 int32 = 8
	defaultWidth, defaultHeight int32 = 600, 400
)

// 设计器
type Designer struct {
	//page          lcl.IPageControl // 设计器 tabs
	tab           *wg.TTab         // 设计器 tabs
	tabMenu       lcl.IPopupMenu   // tab 菜单
	designerForms map[int]*FormTab // 设计器窗体列表
}

// 创建设计器的布局
func (m *BottomBox) createFromDesignerLayout() *Designer {
	des := new(Designer)
	des.designerForms = make(map[int]*FormTab)
	//des.page = lcl.NewPageControl(m.box)
	des.tab = wg.NewTab(m.box)
	des.tab.SetBounds(0, 0, m.rightBox.Width(), m.rightBox.Height())
	des.tab.SetAlign(types.AlClient)
	des.tab.SetBevelColor(wg.LightenColor(bgDarkColor, 0.3))
	des.tab.SetBevelOuter(types.BvLowered)
	des.tab.ScrollLeft().SetTop(3)
	des.tab.ScrollLeft().SetHeight(20)
	des.tab.ScrollLeft().SetColor(wg.DarkenColor(bgLightColor, 0.1))
	des.tab.ScrollRight().SetTop(3)
	des.tab.ScrollRight().SetHeight(20)
	des.tab.ScrollRight().SetColor(wg.DarkenColor(bgLightColor, 0.1))
	des.tab.EnableScrollButton(false)
	//des.page.SetTabStop(true)
	des.tab.SetParent(m.rightBox)
	// 右键菜单
	//des.page.SetOnContextPopup(func(sender lcl.IObject, mousePos types.TPoint, handled *bool) {
	//
	//})
	des.tab.SetOnClick(func(sender lcl.IObject) {
		logs.Debug("Designer PageControl click")
	})
	lcl.RunOnMainThreadAsync(func(id uint32) {
		des.tab.RecalculatePosition()
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
	m.tabMenu = lcl.NewPopupMenu(m.tab)
	m.tabMenu.SetImages(imageActions.ImageList100())
	items := m.tabMenu.Items()
	closeMenuItem := lcl.NewMenuItem(m.tab)
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
func (m *Designer) addDesignerFormTab(defaultId ...int) *FormTab {
	m.hideFormTabs()
	form := new(FormTab)
	form.componentName = make(map[string]int)
	// 组件树
	form.tree = lcl.NewTreeView(inspector.componentTree.treeComponentTree)
	form.tree.SetAutoExpand(true)
	form.tree.SetReadOnly(true)
	form.tree.SetDoubleBuffered(true)
	//m.tree.SetMultiSelect(true) // 多选控制
	form.tree.SetAlign(types.AlClient)
	form.tree.SetVisible(false)
	SetComponentDefaultColor(form.tree)
	form.tree.SetBorderStyleToBorderStyle(types.BsNone)
	form.tree.SetImages(imageComponents.ImageList100())
	form.tree.SetOnGetSelectedIndex(form.TreeOnGetSelectedIndex)
	form.tree.SetOnMouseDown(form.TreeOnMouseDown)
	form.tree.SetOnContextPopup(form.TreeOnContextPopup)
	// 组件树右键菜单
	form.CreateComponentMenu()
	form.tree.SetPopupMenu(form.componentMenu.treePopupMenu)
	form.tree.SetParent(inspector.componentTree.treeComponentTree)

	// 默认名
	if len(defaultId) > 0 {
		form.Id = defaultId[0]
	} else {
		form.Id = len(m.designerForms) + 1
	}
	form.name = fmt.Sprintf("Form%v", form.Id)
	// 窗体ID
	m.designerForms[form.Id] = form

	//form.sheet = lcl.NewTabSheet(m.page)
	form.sheet = m.tab.NewPage()
	form.sheet.Button().SetIconFavoriteFormBytes(resources.Images("components/tform.png"))
	form.sheet.Button().SetIconCloseFormBytes(resources.Images("button/close.png"))
	form.sheet.Button().SetIconCloseHighlightFormBytes(resources.Images("button/close_highlight.png"))
	form.sheet.Button().SetCloseHintText("关闭设计窗体")
	form.sheet.Button().SetBorderDirections(types.NewSet(wg.BbdTop))
	form.sheet.Button().SetCaption(form.name)
	form.sheet.Button().Font().SetColor(colors.ClBlack)
	form.sheet.Button().SetColorGradient(bgLightColor, bgLightColor) // 设置标签按钮过度颜色
	form.sheet.SetDefaultColor(bgLightColor)                         // 设置默认颜色
	form.sheet.SetActiveColor(wg.DarkenColor(bgLightColor, 0.1))     // 设置激活颜色
	form.sheet.SetColor(wg.DarkenColor(bgLightColor, 0.1))           // 设置背景色
	form.sheet.SetOnHide(form.tabSheetOnHide)
	form.sheet.SetOnShow(form.tabSheetOnShow)
	form.sheet.SetOnClose(form.tabSheetOnClose)
	SetComponentDefaultColor(form.sheet) // 设置背景色
	//form.sheet.SetAlign(types.AlClient)
	form.sheet.SetParent(m.tab)

	form.scroll = lcl.NewScrollBox(form.sheet)
	form.scroll.SetAlign(types.AlClient)
	form.scroll.SetAutoScroll(true)
	form.scroll.SetBorderStyleToBorderStyle(types.BsNone)
	form.scroll.SetDoubleBuffered(true)
	form.scroll.SetParent(form.sheet)

	//newStatusBar(form.scroll)

	// 创建设计窗体
	form.NewFormDesigner()

	m.tab.EnableScrollButton(true)
	return form
}

// 激活指定的 tab
// 触发 tab 的 onshow 事件
func (m *Designer) ActiveFormTab(tab *FormTab) {
	tab.sheet.SetActive(true)
	for _, form := range m.designerForms {
		form.isDesigner = false
	}
	tab.isDesigner = true
}

// GetFormTab 获取指定窗体
//
//	formId - 窗体ID
func (m *Designer) GetFormTab(formId int) *FormTab {
	return m.designerForms[formId]
}

// 绘制刻度尺, 在外层 scroll 上
//
//	func (m *FormTab) scrollDrawRuler() {
//		gridSize := 5 // 小刻度
//		//canvas := m.bg.Canvas()
//		canvas := m.scroll.Canvas()
//		canvas.PenToPen().SetColor(colors.ClBlack)
//		width, height := m.FormRoot.Width(), m.FormRoot.Height()
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
