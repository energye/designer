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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 窗体设计页, 只针对设计窗体和对应的代码 tab page
// 如果非设计窗口代码, 常规代码文件直接在 mainPage 初始化编辑器
type TFormDesignPage struct {
	formTab            *wg.TTab       // 设计窗体Tab
	formDesignPage     *wg.TPage      // 设计窗体
	formDesignScroll   lcl.IScrollBox // 设计窗体滚动条
	formUserEditorPage *wg.TPage      // 用户代码 tab page
	//formUIEditorPage   *wg.TPage      // UI 代码 tab page
	editor editor.IEditor
}

func NewFormDesignPage(formTab *FormTab) *TFormDesignPage {
	m := &TFormDesignPage{}
	m.formTab = wg.NewTab(formTab.mainPage)
	m.formTab.SetBounds(0, 0, formTab.mainPage.Width(), formTab.mainPage.Height())
	m.formTab.SetAlign(types.AlClient)
	m.formTab.EnableScrollButton(false)
	m.formTab.SetParent(formTab.mainPage)

	m.formDesignPage = m.formTab.NewPage()
	m.formDesignPage.Button().SetCaption("窗体")
	m.formUserEditorPage = m.formTab.NewPage()
	m.formUserEditorPage.Button().SetCaption("代码")
	m.formUserEditorPage.SetOnShow(m.UserEditorPageOnShow)
	//m.formUIEditorPage = m.formTab.NewPage()
	//m.formUIEditorPage.Button().SetCaption("UI代码")

	setFormDesignPageStyle(m.formDesignPage, resources.Images("components/tform.png"))
	name := tool.GetSVGIconPath("go")
	codeIconPngData, _ := tool.SVGToPNG(resources.Images(name), 24, 24)
	setFormDesignPageStyle(m.formUserEditorPage, codeIconPngData)
	//setFormDesignPageStyle(m.formUIEditorPage)

	m.formDesignScroll = lcl.NewScrollBox(formTab.mainPage)
	m.formDesignScroll.SetAlign(types.AlClient)
	m.formDesignScroll.SetAutoScroll(true)
	if tool.IsDarwin {
		// fix: laz MacOS bug 默认隐藏滚动条, 手动控制显示
		// 该bug体现为当同时出现横坚滚动条时, UI 锁死崩溃, 在laz 4.6 复现
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

	m.formTab.HideAllActivated()
	m.formDesignPage.SetActive(true)
	m.formTab.RecalculatePosition()

	if m.editor == nil {
		m.editor = editor.NewEditor(m.formUserEditorPage)
	}
	return m
}

func setFormDesignPageStyle(page *wg.TPage, icon []byte) {
	tabColor := colors.ClWhite
	btnColor := colors.RGBToColor(234, 239, 249)
	if icon != nil && len(icon) > 0 {
		page.Button().SetIconFavoriteFormBytes(icon)
	}
	page.Button().SetWidth(65)
	page.Button().Font().SetColor(colors.ClBlack)
	page.Button().RoundedCorner = types.NewSet(wg.RcLeftTop, wg.RcRightTop)
	page.Button().TextOffSetX = 10
	page.Button().SetBorderColor(wg.BbdNone, wg.LightenColor(btnColor, 0.8))
	page.Button().SetRadius(5)
	page.Button().SetColor(tabColor)
	page.Button().SetDownColor(wg.LightenColor(btnColor, 0.3), wg.LightenColor(btnColor, 0.5))
	page.Button().SetEnterColor(wg.LightenColor(btnColor, 0.1), wg.LightenColor(btnColor, 0.3))
	page.SetDefaultColor(tabColor)
	page.SetActiveColor(btnColor)
	page.Button().SetCursor(types.CrHandPoint)
	page.SetColor(tabColor)
}

func (m *TFormDesignPage) ActiveDesignPage() {
	m.formTab.HideAllActivated()
	m.formDesignPage.SetActive(true)
}
func (m *TFormDesignPage) ActiveEditorPage() {
	m.formTab.HideAllActivated()
	m.formUserEditorPage.SetActive(true)
}

func (m *TFormDesignPage) UserEditorPageOnShow(sender lcl.IObject) {
	fmt.Println("UserEditorPageOnShow IsMainThread:", tool.IsMainThread(), m.formUserEditorPage.BoundsRect())
	if m.editor != nil {
		if wvEditor, ok := m.editor.(editor.IWebviewEditor); ok {
			wvEditor.LoadURL("energy://designer/index.html")
			wvEditor.CreateBrowser()
		}
	}
}
