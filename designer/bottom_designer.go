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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"path/filepath"
	"strings"
)

// 窗体设计功能

var (
	designer                                    *Designer
	designerDefaultWidth, designerDefaultHeight int32 = 600, 400
	canvasBorderSpacing                         int32 = 8
)

// 主设计器
type Designer struct {
	tab            *wg.TTab                  // 设计器 tabs
	tabMenu        lcl.IPopupMenu            // tab 菜单
	designerForms  map[int]*FormTab          // 设计器窗体列表
	codeEditorTabs map[string]*CodeEditorTab // 代码编辑器标签列表
	defaultTip     *wg.TButton               // Home 默认提示
}

func (m *Designer) IsDuplicateName(currComp *TDesigningComponent, formName string) bool {
	for _, formTab := range m.designerForms {
		if formTab.FormRoot != currComp && formTab.FormRoot.Name() == formName {
			return true
		}
	}
	return false
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

// ResetDesigner 重置设计器
// 新建项目时/重新打开项目时调用
// 打开设计窗体不调用
func ResetDesigner() {
	if designer == nil {
		return
	}
	// 关闭所有已打开的代码编辑器标签
	for filePath, tab := range designer.codeEditorTabs {
		editor.CloseFileInEditor(filePath)
		tab.mainPage.Close()
	}
	designer.codeEditorTabs = make(map[string]*CodeEditorTab)
	// 关闭所有已打开的设计窗体
	tempForms := designer.designerForms
	// 关闭之前打开的所有设计窗体
	for _, form := range tempForms {
		if form == nil {
			continue
		}
		form.Remove()
	}
	designer.designerForms = make(map[int]*FormTab) // 清空设计窗体
	ProjectTreeClearComponentTreeNode()
	ProjectTreeClearSrcTreeNode()
	gFromEditor = nil
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

func setDesignerFormTabPageStyle(page *wg.TPage) {
	tabColor := colors.ClWhite
	btnColor := colors.RGBToColor(234, 239, 249)
	page.Button().SetIconFavoriteFormBytes(resources.Images("components/design.png"))
	page.Button().SetIconCloseFormBytes(resources.Images("button/close.png"))
	page.Button().SetIconCloseHighlightFormBytes(resources.Images("button/close_highlight.png"))
	page.Button().SetCloseHintText("关闭设计窗体")
	page.Button().Font().SetColor(colors.ClBlack)
	page.Button().Font().SetStyle(types.NewSet(types.FsBold))
	page.Button().RoundedCorner = types.NewSet(wg.RcLeftTop, wg.RcRightTop)
	page.Button().SetBorderColor(wg.BbdNone, wg.LightenColor(btnColor, 0.8))
	page.Button().SetRadius(5)
	page.Button().SetColor(tabColor)
	page.Button().SetDownColor(wg.LightenColor(btnColor, 0.3), wg.LightenColor(btnColor, 0.5))
	page.Button().SetEnterColor(wg.LightenColor(btnColor, 0.1), wg.LightenColor(btnColor, 0.3))
	page.SetDefaultColor(tabColor)
	page.SetActiveColor(btnColor)
	page.Button().SetCursor(types.CrHandPoint)
}

// GetNewDesignFormName 获取新的设计表单名称和ID
// 该方法用于生成一个唯一的表单ID和表单名称，确保与现有的表单不冲突。
// 表单名称采用"Form{数字}"的格式，其中数字从1开始递增查找第一个可用的编号。
func (m *Designer) GetNewDesignFormName() (int, string) {
	//	matchFormId 检查给定的ID和名称是否与现有表单冲突
	//	参数:
	//	    - newId: 待检查的表单ID
	//	    - newName: 待检查的表单名称
	//	返回值: 如果ID和名称都未被使用则返回true，否则返回false
	matchFormId := func(newId int, newName string) bool {
		ok := true
		for id, form := range m.designerForms {
			if newId == id || form.FormRoot.Name() == newName {
				ok = false
				break
			}
		}
		return ok
	}

	// 从1开始查找第一个可用的表单ID和名称
	// 循环检查每个递增的数字，直到找到一个未被使用的ID和对应的表单名称，找到后立即返回
	nextId := 1
	for {
		newTmpName := fmt.Sprintf("Form%v", nextId)
		if matchFormId(nextId, newTmpName) {
			return nextId, newTmpName
		}
		nextId++
	}
}

// 添加一个窗体设计器 form tab
func (m *Designer) addDesignFormTab(uiForm *bean.TUIForm) *FormTab {
	form := new(FormTab)
	form.FormRoot = new(TDesigningComponent)
	form.FormRoot.ComponentType = consts.CtForm

	// 默认名
	if uiForm != nil {
		form.Id = uiForm.Id
		form.FormRoot.SetName(uiForm.Name)
	} else {
		newId, newName := m.GetNewDesignFormName()
		form.Id = newId
		form.FormRoot.SetName(newName)
	}

	// 窗体ID
	m.designerForms[form.Id] = form

	form.mainPage = m.tab.NewPage()
	form.mainPage.Button().SetCaption(form.FormRoot.Name())
	form.mainPage.SetOnHide(form.tabSheetOnHide)
	form.mainPage.SetOnShow(form.tabSheetOnShow)
	form.mainPage.SetOnClose(form.tabSheetOnClose)
	SetComponentDefaultColor(form.mainPage) // 设置背景色
	setDesignerFormTabPageStyle(form.mainPage)

	// 窗体和代码 tab
	form.formDesignPage = NewFormDesignPage(form)

	// 创建设计窗体
	form.NewFormDesigner()

	m.tab.EnableScrollButton(true)
	return form
}

// 激活指定的 tab
// 触发 tab 的 onshow 事件
func (m *Designer) ActiveFormTab(tab *FormTab) {
	// 将所有窗体设为非设计窗体
	for _, form := range m.designerForms {
		if form == nil {
			continue
		}
		form.IsDesigner = false
	}
	// 设为设计窗体
	tab.IsDesigner = true
	// 激活当前 tab
	tab.mainPage.SetActive(true)
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
//
// FindFormTabByFile 检查给定文件路径是否属于某个设计窗体
// 返回对应的 FormTab 和文件类型 ("ui", "uigo", "go"), 如果不属于任何窗体返回 nil
func (m *Designer) FindFormTabByFile(filePath string) (*FormTab, string) {
	if m == nil || m.designerForms == nil {
		return nil, ""
	}
	filePath = filepath.Clean(filePath)
	for _, formTab := range m.designerForms {
		goPath := filepath.Clean(filepath.Join(bean.CodePath(), formTab.GOFile()))
		goUserPath := filepath.Clean(filepath.Join(bean.CodePath(), formTab.GOUserFile()))
		uiPath := filepath.Clean(filepath.Join(bean.CodePath(), formTab.UIFile()))
		if filePath == goPath {
			return formTab, "uigo"
		} else if filePath == goUserPath {
			return formTab, "go"
		} else if filePath == uiPath {
			return formTab, "ui"
		}
	}
	// 检查 .ui 文件是否在 layouts 目录下
	for _, formTab := range m.designerForms {
		uiName := strings.ToLower(formTab.FormRoot.Name()) + consts.UIExt
		if filepath.Base(filePath) == uiName {
			return formTab, "ui"
		}
	}
	return nil, ""
}

// openFileInAppropriateTab 根据文件类型在合适的位置打开文件
// 如果文件属于设计窗体, 切换到对应窗体的子标签
// 如果不属于任何窗体, 创建或激活代码编辑器标签
// 非文本文件(图片、二进制等)不会在编辑器中打开
func openFileInAppropriateTab(filePath string) {
	if designer == nil {
		return
	}
	if !editor.IsTextFile(filePath) {
		return
	}
	formTab, fileType := designer.FindFormTabByFile(filePath)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if formTab != nil {
			switch fileType {
			case "go":
				formTab.formDesignPage.formPageControl.SetActivePageIndex(1)
			case "uigo":
				formTab.formDesignPage.formPageControl.SetActivePageIndex(2)
			case "ui":
				formTab.formDesignPage.formPageControl.SetActivePageIndex(0)
			}
			designer.tab.HideAllActivated()
			designer.ActiveFormTab(formTab)
		} else {
			tab := designer.addCodeEditorTab(filePath)
			designer.ActivateCodeEditorTab(tab)
			designer.tab.RecalculatePosition()
		}
	})
}

// goToDefinition Ctrl+Click 跳转定义处理
// xxx.ui.go 在当前窗体的 UI代码 子标签打开
// xxx.go 在当前窗体的 代码 子标签打开
// 其它文件在主 tab 标签打开
func goToDefinition(filePath string, line, character int) {
	openFileInAppropriateTab(filePath)
	// 定位到定义的具体行和列
	editor.OpenFileInEditor(filePath)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		ipc.Emit("goto-position", filePath, line, character)
	})
}
