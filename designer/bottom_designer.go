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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
)

// 窗体设计功能

var (
	designer                    *Designer
	margin                      int32 = 0
	borderWidth                 int32 = 8
	defaultWidth, defaultHeight int32 = 600, 400
)

// 主设计器
type Designer struct {
	tab           *wg.TTab         // 设计器 tabs
	tabMenu       lcl.IPopupMenu   // tab 菜单
	designerForms map[int]*FormTab // 设计器窗体列表
	designerCount int              // 设计器窗总数，值动态更新
	// 默认提示
	defaultTip *wg.TButton
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

// 隐藏所有组件树
func (m *Designer) hideAllComponentTrees() {
	for _, formTab := range m.designerForms {
		if formTab == nil {
			continue
		}
		formTab.tree.SetVisible(false)
	}
}

// ResetDesigner 重置设计器
// 新建项目时/重新打开项目时调用
// 打开设计窗体不调用
func ResetDesigner() {
	if designer == nil {
		return
	}
	// 关闭所有已打开的设计窗体
	tempForms := designer.designerForms
	// 关闭之前打开的所有设计窗体
	for _, form := range tempForms {
		if form == nil {
			continue
		}
		form.Close()
	}
	designer.designerForms = make(map[int]*FormTab) // 清空设计窗体
	SetDesignerCount(0)
}

func UpdateHistoryProject(egpFilePath string) {
	// 更新打开项目历史记录
	config.UpdateHistoryProject(egpFilePath)
	config.UpdateConfig()
	// 更新设计器菜单-文件-历史项目
	lcl.RunOnMainThreadAsync(func(id uint32) {
		MainWindow.mainMenu.fileHistoryProjectMenu()
	})
}

// SetDesignerCount
// 设置当前设计器的窗体总数
// 当前设计器添加窗体时 +1
// 当前设计器删除窗体时 -1
// 当前设计器重置时 =0
func SetDesignerCount(count int) {
	if count < 0 {
		count = 0
	}
	designer.designerCount = count
}

// SetRecvMethods 设置指定表单的接收方法列表
// formId: 表单ID，用于定位需要设置接收方法的表单
// methods: 接收方法信息列表，包含方法的详细信息
func SetRecvMethods(formId int, methods []*dast.TFuncInfo) {
	if designer == nil {
		return
	}
	if form, ok := designer.designerForms[formId]; ok {
		form.SetRecvMethods(methods)
	}
}

// 添加一个窗体设计器 form tab
func (m *Designer) addDesignerFormTab(defaultId ...int) *FormTab {
	SetDesignerCount(m.designerCount + 1)
	form := new(FormTab)
	form.componentName = make(map[string]int)
	// 组件树
	form.tree = lcl.NewTreeView(MainWindow.contentLayout.layoutProject.box)
	form.tree.SetAutoExpand(true)
	form.tree.SetReadOnly(true)
	form.tree.SetDoubleBuffered(true)
	form.tree.SetTreeLineColor(colors.RGBToColor(128, 128, 128))
	form.tree.SetTreeLinePenStyle(types.PsSolid)
	//m.tree.SetMultiSelect(true) // 多选控制
	form.tree.SetAlign(types.AlClient)
	form.tree.SetVisible(false)
	SetComponentDefaultColor(form.tree)
	form.tree.SetBorderStyleToBorderStyle(types.BsNone)
	form.tree.SetImages(imageComponents.ImageList100())
	form.tree.SetOnGetSelectedIndex(form.TreeOnGetSelectedIndex)
	form.tree.SetOnMouseDown(form.TreeOnMouseDown)
	form.tree.SetOnContextPopup(form.TreeOnContextPopup)
	form.tree.Font().SetHeight(-11)
	// 组件树右键菜单
	form.CreateComponentMenu()
	form.tree.SetPopupMenu(form.componentMenu.treePopupMenu)
	form.tree.SetParent(MainWindow.contentLayout.layoutProject.box)

	// 默认名
	if len(defaultId) > 0 {
		form.Id = defaultId[0]
	} else {
		form.Id = m.designerCount
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
	form.sheet.Button().SetCaption(form.name)
	form.sheet.Button().Font().SetColor(colors.ClBlack)
	//form.sheet.Button().SetBorderDirections(types.NewSet(wg.BbdTop))
	//form.sheet.Button().SetBorderColor(wg.BbdTop, colors.ClBlue)
	//form.sheet.Button().SetColorGradient(bgLightColor, bgLightColor) // 设置标签按钮过度颜色
	//form.sheet.SetDefaultColor(bgLightColor)                         // 设置默认颜色
	//form.sheet.SetActiveColor(bgLightColor)                          // 设置激活颜色
	//form.sheet.SetColor(bgLightColor)                                // 设置背景色
	form.sheet.SetOnHide(form.tabSheetOnHide)
	form.sheet.SetOnShow(form.tabSheetOnShow)
	form.sheet.SetOnClose(form.tabSheetOnClose)
	SetComponentDefaultColor(form.sheet) // 设置背景色
	//form.sheet.SetAlign(types.AlClient)
	form.sheet.SetParent(m.tab)

	form.scroll = lcl.NewScrollBox(form.sheet)
	form.scroll.SetAlign(types.AlClient)
	form.scroll.SetAutoScroll(true)
	if tool.IsDarwin {
		// fix: laz MacOS bug 默认隐藏滚动条, 手动控制显示
		// 该bug体现为当同时出现横坚滚动条时, UI 锁死崩溃, 在laz 4.6 复现
		hBar := form.scroll.HorzScrollBar()
		vBar := form.scroll.VertScrollBar()
		hBar.SetVisible(false)
		vBar.SetVisible(false)
		form.scroll.SetOnResize(func(sender lcl.IObject) {
			// 这里使用 < 小于做判断
			// 当滚动box宽或高小于设计窗体大小时显示或隐藏滚动条
			vBar.SetVisible(form.scroll.Height() < form.formDesigner.Form.Height())
			hBar.SetVisible(form.scroll.Width() < form.formDesigner.Form.Width())
		})
	}
	form.scroll.SetBorderStyleToBorderStyle(types.BsNone)
	form.scroll.SetDoubleBuffered(true)
	form.scroll.SetParent(form.sheet)

	// 创建设计窗体
	form.NewFormDesigner()

	m.tab.EnableScrollButton(true)
	return form
}

// 激活指定的 tab
// 触发 tab 的 onshow 事件
func (m *Designer) ActiveFormTab(tab *FormTab) {
	for _, form := range m.designerForms {
		if form == nil {
			continue
		}
		form.IsDesigner = false
	}
	tab.IsDesigner = true
	tab.sheet.SetActive(true)
}

// GetFormTab 获取指定窗体
//
//	formId - 窗体ID
func (m *Designer) GetFormTab(formId int) *FormTab {
	if m == nil || m.designerForms == nil {
		return nil
	}
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
