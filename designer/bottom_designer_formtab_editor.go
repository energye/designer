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
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/webview"
	projBean "github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"path/filepath"
	"time"
)

var gFromEditor editor.IEditor

type IFormDesignCode interface {
	UIFile() string
	GOFile() string
	GOUserFile() string
}

// 窗体设计页, 只针对设计窗体和对应的代码 tab page
// 如果非设计窗口代码, 常规代码文件直接在 mainPage 初始化编辑器
type TFormDesignPage struct {
	//formTab            *wg.TTab  // 设计窗体Tab
	//formDesignPage     *wg.TPage // 设计窗体
	//formUserEditorPage *wg.TPage // 用户代码 tab page
	////formUIEditorPage   *wg.TPage      // UI 代码 tab page

	formDesignScroll lcl.IScrollBox // 设计窗体滚动条

	//
	formPageControl    lcl.IPageControl
	formDesignPage     lcl.ITabSheet // "窗体"
	formUserEditorPage lcl.ITabSheet // "代码"
	formUIEditorPage   lcl.ITabSheet // "UI代码" - 只读

	wvWindowParent wv.IWindowParent

	formDesignCode IFormDesignCode
}

func NewFormDesignPage(formTab *FormTab) *TFormDesignPage {
	m := &TFormDesignPage{formDesignCode: formTab}
	//m.formTab = wg.NewTab(formTab.mainPage)
	//m.formTab.SetBounds(0, 0, formTab.mainPage.Width(), formTab.mainPage.Height())
	//m.formTab.SetAlign(types.AlClient)
	//m.formTab.EnableScrollButton(false)
	//m.formTab.SetParent(formTab.mainPage)
	//
	//m.formDesignPage = m.formTab.NewPage()
	//m.formDesignPage.SetDoubleBuffered(true)
	//m.formDesignPage.Button().SetCaption("窗体")
	//m.formDesignPage.SetOnHide(func(sender lcl.IObject) {
	//	m.formDesignScroll.SetVisible(false)
	//})
	//m.formDesignPage.SetOnShow(func(sender lcl.IObject) {
	//	m.formDesignScroll.SetVisible(true)
	//})
	//setFormDesignPageStyle(m.formDesignPage, resources.Images("components/tform.png"))
	//
	//m.formUserEditorPage = m.formTab.NewPage()
	//m.formUserEditorPage.SetDoubleBuffered(true)
	//m.formUserEditorPage.Button().SetCaption("代码")
	//m.formUserEditorPage.SetOnShow(m.UserEditorPageOnShow)
	//name := tool.GetSVGIconPath("go")
	//codeIconPngData, _ := tool.SVGToPNG(resources.Images(name), 24, 24)
	//setFormDesignPageStyle(m.formUserEditorPage, codeIconPngData)

	m.formPageControl = lcl.NewPageControl(formTab.mainPage)
	m.formPageControl.SetAlign(types.AlClient)
	m.formPageControl.SetParent(formTab.mainPage)

	m.formDesignPage = lcl.NewTabSheet(formTab.mainPage)
	m.formDesignPage.SetPageControl(m.formPageControl)
	m.formDesignPage.SetCaption("窗体")

	m.formUserEditorPage = lcl.NewTabSheet(formTab.mainPage)
	m.formUserEditorPage.SetPageControl(m.formPageControl)
	m.formUserEditorPage.SetCaption("代码")
	m.formUserEditorPage.SetOnShow(m.UserEditorPageOnShow)

	//m.formUIEditorPage = lcl.NewTabSheet(formTab.mainPage)
	//m.formUIEditorPage.SetPageControl(m.formPageControl)
	//m.formUIEditorPage.SetCaption("UI代码")
	//m.formUIEditorPage.SetOnShow(m.UIEditorPageOnShow)

	m.formDesignScroll = lcl.NewScrollBox(formTab.mainPage)
	m.formDesignScroll.SetAlign(types.AlClient)
	m.formDesignScroll.SetAutoScroll(true)
	if tool.IsDarwin {
		// fix: laz MacOS bug 默认隐藏滚动条, 手动控制显示
		// 该bug体现为当同时出现横竖滚动条时, UI 锁死崩溃, 在laz 4.6 复现
		hBar := m.formDesignScroll.HorzScrollBar()
		vBar := m.formDesignScroll.VertScrollBar()
		hBar.SetVisible(false)
		vBar.SetVisible(false)
		m.formDesignScroll.SetOnResize(func(sender lcl.IObject) {
			// 这里使用 < 小于做判断
			// 当滚动box宽或高小于设计窗体大小时显示或隐藏滚动条
			vBar.SetVisible(m.formDesignScroll.Height() < formTab.formDesigner.Form.Height())
			hBar.SetVisible(m.formDesignScroll.Width() < formTab.formDesigner.Form.Width())
		})
	}
	m.formDesignScroll.SetBorderStyleToBorderStyle(types.BsNone)
	m.formDesignScroll.SetDoubleBuffered(true)
	m.formDesignScroll.SetParent(m.formDesignPage)

	//m.formTab.HideAllActivated()
	//m.formDesignPage.SetActive(true)
	//m.formTab.RecalculatePosition()

	return m
}

//func (m *TFormDesignPage) ActiveDesignPage() {
//	m.formTab.HideAllActivated()
//	m.formDesignPage.SetActive(true)
//}
//func (m *TFormDesignPage) ActiveEditorPage() {
//	m.formTab.HideAllActivated()
//	m.formUserEditorPage.SetActive(true)
//}

func (m *TFormDesignPage) UserEditorPageOnShow(sender lcl.IObject) {
	m.initEditor()
	m.SwitchTabPageEditor(false)
}

func (m *TFormDesignPage) UIEditorPageOnShow(sender lcl.IObject) {
	m.initEditor()
	m.SwitchTabPageEditor(true)
}

func (m *TFormDesignPage) initEditor() {
	if gFromEditor == nil {
		gFromEditor = editor.NewEditor(designer.tab)
		if gFromEditor.Type() == editor.EtWebview {
			if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
				m.wvWindowParent = wvEditor.Webview().WindowParent()
			}
		}
	} else {
		if gFromEditor.Type() == editor.EtWebview && m.wvWindowParent == nil {
			m.wvWindowParent = webview.NewWebviewWindowParent(designer.tab)
		}
	}
	webview.SetOnGoToDefinition(goToDefinition)
}

func (m *TFormDesignPage) SwitchTabPageEditor(uiCode bool) {
	if gFromEditor != nil && m.wvWindowParent != nil {
		if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
			canLoad := make(chan error, 1)
			wvEditor.SetCanLoadChan(canLoad)
			go func() {
				_ = <-canLoad
				wvEditor.SetCanLoadChan(nil)
				close(canLoad)
				var (
					filePath string
					readOnly bool
				)
				if uiCode {
					filePath = filepath.Join(projBean.CodePath(), m.formDesignCode.GOFile())
					readOnly = true
				} else {
					filePath = filepath.Join(projBean.CodePath(), m.formDesignCode.GOUserFile())
					readOnly = false
				}
				if wvEditor.Initialized() {
					lcl.RunOnMainThreadAsync(func(id uint32) {
						editor.OpenFileInEditor(filePath, readOnly)
					})
				}
			}()
			var targetOwner lcl.IWinControl
			if uiCode && m.formUIEditorPage != nil {
				targetOwner = m.formUIEditorPage
			} else {
				targetOwner = m.formUserEditorPage
			}
			wvEditor.SwitchTabPage(targetOwner, m.wvWindowParent)
			wvEditor.CreateBrowser()
		}
	}
}

func (m *TFormDesignPage) ActiveCodeEditorTab() {
	activeIdx := m.formPageControl.ActivePageIndex()
	m.initEditor()
	m.SwitchTabPageEditor(activeIdx == 2)
	if activeIdx >= 1 {
		go func() {
			time.AfterFunc(time.Millisecond*200, func() {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					m.formPageControl.SetActivePageIndex(activeIdx)
				})
			})
		}()
	}
}
