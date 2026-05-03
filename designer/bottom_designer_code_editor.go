// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"fmt"
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/webview"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"path/filepath"
)

// 代码编辑器标签 - 非设计窗体的代码文件
// 在设计器主 tab 同级创建代码编辑标签, 复用现有的 Monaco 编辑器

type CodeEditorTab struct {
	filePath       string           // 绝对文件路径
	mainPage       *wg.TPage        // Level 1 tab page in Designer.tab
	wvWindowParent wv.IWindowParent // WebView2 窗口父容器
}

// addCodeEditorTab 创建或获取一个代码编辑器标签
// 如果文件已经打开, 直接激活已有标签
func (m *Designer) addCodeEditorTab(filePath string) *CodeEditorTab {
	if m.codeEditorTabs == nil {
		m.codeEditorTabs = make(map[string]*CodeEditorTab)
	}
	// 如果已打开, 直接激活
	if existing, ok := m.codeEditorTabs[filePath]; ok {
		m.tab.HideAllActivated()
		existing.mainPage.SetActive(true)
		return existing
	}

	tab := &CodeEditorTab{
		filePath: filePath,
	}

	// 在 Designer.tab 中创建顶级标签页
	tab.mainPage = m.tab.NewPage()
	tab.mainPage.Button().SetCaption(filepath.Base(filePath))
	tab.mainPage.SetOnShow(tab.onShow)
	tab.mainPage.SetOnHide(tab.onHide)
	tab.mainPage.SetOnClose(tab.onClose)
	setCodeEditorTabPageStyle(tab.mainPage)

	m.codeEditorTabs[filePath] = tab
	m.tab.EnableScrollButton(true)
	return tab
}

// ActivateCodeEditorTab 激活指定的代码编辑器标签
func (m *Designer) ActivateCodeEditorTab(tab *CodeEditorTab) {
	m.tab.HideAllActivated()
	tab.mainPage.SetActive(true)
}

// onShow 代码编辑器标签显示事件
func (m *CodeEditorTab) onShow(sender lcl.IObject) {
	fmt.Println("CodeEditorTab onShow:", m.filePath)
	m.initEditor()
	m.switchEditorToThisTab()
}

// onHide 代码编辑器标签隐藏事件
func (m *CodeEditorTab) onHide(sender lcl.IObject) {
	if m.mainPage.IsEnterClose() {
		return
	}
}

// onClose 代码编辑器标签关闭事件
func (m *CodeEditorTab) onClose(page *wg.TPage, canClose *bool) {
	*canClose = true
	if gFromEditor != nil {
		if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
			// Only reparent the browser if this tab currently holds it,
			// otherwise we'd detach the browser from the active tab.
			if m.mainPage.Active() {
				wvEditor.SwitchTabPage(designer.tab, wvEditor.Webview().WindowParent())
			}
		}
	}
	// 在 Monaco 前端关闭文件
	editor.CloseFileInEditor(m.filePath)
	// 从列表中移除
	delete(designer.codeEditorTabs, m.filePath)
	if len(designer.tab.Pages()) == 0 {
		designer.tab.EnableScrollButton(false)
	}
}

// initEditor 确保共享编辑器实例存在, 并创建窗口父容器
func (m *CodeEditorTab) initEditor() {
	if gFromEditor == nil {
		gFromEditor = editor.NewEditor(designer.tab)
		if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
			m.wvWindowParent = wvEditor.Webview().WindowParent()
		}
		webview.SetOnGoToDefinition(goToDefinition)
	} else {
		if gFromEditor.Type() == editor.EtWebview && m.wvWindowParent == nil {
			m.wvWindowParent = webview.NewWebviewWindowParent(designer.tab)
		}
	}
}

// switchEditorToThisTab 将共享编辑器切换到当前代码编辑器标签
func (m *CodeEditorTab) switchEditorToThisTab() {
	if gFromEditor != nil && m.wvWindowParent != nil {
		if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
			canLoad := make(chan error, 1)
			wvEditor.SetCanLoadChan(canLoad)
			go func() {
				err := <-canLoad
				wvEditor.SetCanLoadChan(nil)
				close(canLoad)
				fmt.Println("CodeEditorTab canLoad", err, wvEditor.Initialized(), "filePath:", m.filePath)
				if wvEditor.Initialized() {
					lcl.RunOnMainThreadAsync(func(id uint32) {
						editor.OpenFileInEditor(m.filePath, false)
					})
				}
			}()
			wvEditor.SwitchTabPage(m.mainPage, m.wvWindowParent)
			wvEditor.CreateBrowser()
		}
	}
}

// setCodeEditorTabPageStyle 设置代码编辑器标签样式
func setCodeEditorTabPageStyle(page *wg.TPage) {
	tabColor := colors.ClWhite
	btnColor := colors.RGBToColor(234, 239, 249)
	page.Button().SetIconCloseFormBytes(resources.Images("button/close.png"))
	page.Button().SetIconCloseHighlightFormBytes(resources.Images("button/close_highlight.png"))
	page.Button().SetCloseHintText("关闭")
	page.Button().Font().SetColor(colors.ClBlack)
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
