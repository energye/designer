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

package options

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
)

var (
	buildFormWidth  = int32(555)
	buildFormHeight = int32(555)
)

func NewBuildForm() *TBuildForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TBuildForm{TEngForm: newEngForm.(*lcl.TEngForm)}
	newForm.FormCreate(newEngForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TBuildForm struct {
	*lcl.TEngForm
	closing   bool
	font      lcl.IFont
	selectDir lcl.ISelectDirectoryDialog

	buildTab            *wg.TTab
	buildTabPageConfig  *wg.TPage
	buildTabPagePackage *wg.TPage

	// 基础配置

	// 构建打包

	// 操作按钮
	saveBtn  *wg.TButton
	buildBtn *wg.TButton
}

func (m *TBuildForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TBuildForm FormCreate")
	fontSize := int32(12)
	if tool.IsLinux {
		fontSize = 10
	}
	m.SetCaption("构建配置")
	m.SetWidth(buildFormWidth)
	m.SetHeight(buildFormHeight)
	constr := m.Constraints()
	constr.SetMaxWidth(buildFormWidth)
	constr.SetMaxHeight(buildFormHeight)
	constr.SetMinWidth(buildFormWidth)
	constr.SetMinHeight(buildFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.WorkAreaCenter()
	m.font = lcl.NewFont()
	m.font.SetName("微软雅黑")
	m.font.SetSize(fontSize)
	m.SetColor(colors.ClWhite)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)

	{
		m.buildTab = wg.NewTab(m)
		m.buildTab.Margin = 5
		tabBR := types.TRect{Left: 0, Top: 5}
		tabBR.SetWidth(m.Width())
		tabBR.SetHeight(m.Height() - tabBR.Top)
		m.buildTab.SetBoundsRect(tabBR)
		m.buildTab.SetColor(colors.ClWhite)
		m.buildTab.EnableScrollButton(false)
		m.buildTab.SetParent(m)
		m.buildTab.SetOnChange(func(sender lcl.IObject) {
			//for _, page := range m.buildTab.Pages() {
			//	if page.Active() {
			//		page.Button().SetBorderDirections(0)
			//		page.Button().Font().SetColor(colors.ClWhite)
			//		page.Button().SetIconFavoriteFormBytes(buttons.Get(page.Button().Text()).iconActive)
			//	} else {
			//		page.Button().SetBorderDirections(types.NewSet(wg.BbdBottom, wg.BbdLeft, wg.BbdTop, wg.BbdRight))
			//		page.Button().Font().SetColor(colors.ClBlack)
			//		page.Button().SetIconFavoriteFormBytes(buttons.Get(page.Button().Text()).iconDefault)
			//	}
			//}
		})

		m.buildTabPageConfig = m.buildTab.NewPage()
		m.buildTabPageConfig.SetCaption("基础配置")
		//m.buildTabPageConfig.Button().SetIconFavoriteFormBytes(buttons.Get("Windows").iconDefault)
		//setTabPageStyle(m.platformTabPageWindows)

		m.buildTabPagePackage = m.buildTab.NewPage()
		m.buildTabPagePackage.SetCaption("构建打包")
		//m.buildTabPageConfig.Button().SetIconFavoriteFormBytes(buttons.Get("Windows").iconDefault)
		//setTabPageStyle(m.platformTabPageWindows)

		m.buildTabPageConfig.SetActive(true)
	}

	m.initConfigComponent()
	m.initBuildComponent()

	{
		btnFont := lcl.NewFont()
		btnFont.SetName("微软雅黑")
		btnFont.SetSize(10)
		btnFont.SetStyle(types.NewSet(types.FsBold))

		m.saveBtn = wg.NewButton(m)
		m.saveBtn.SetText("保存配置")
		m.saveBtn.SetFont(btnFont)
		m.saveBtn.Font().SetColor(colors.ClWhite)
		m.saveBtn.SetRadius(3)
		saveBtnRect := types.TRect{Left: 390, Top: 0}
		saveBtnRect.SetWidth(75)
		saveBtnRect.SetHeight(25)
		m.saveBtn.SetBoundsRect(saveBtnRect)
		m.saveBtn.SetColor(colors.RGBToColor(59, 130, 246))
		m.saveBtn.SetParent(m.buildTab)
		m.saveBtn.SetOnClick(m.saveClick)

		m.buildBtn = wg.NewButton(m)
		m.buildBtn.SetText("开始构建")
		m.buildBtn.SetFont(btnFont)
		m.buildBtn.Font().SetColor(colors.ClWhite)
		m.buildBtn.SetRadius(3)
		buildBtnRect := types.TRect{Left: saveBtnRect.Left + saveBtnRect.Width() + 10, Top: saveBtnRect.Top}
		buildBtnRect.SetWidth(75)
		buildBtnRect.SetHeight(25)
		m.buildBtn.SetBoundsRect(buildBtnRect)
		m.buildBtn.SetColor(colors.RGBToColor(46, 204, 113))
		m.buildBtn.SetParent(m.buildTab)
		m.buildBtn.SetOnClick(m.saveClick)
	}
	//(&hook.TWindowHook{Form: m}).Hook()
}

func (m *TBuildForm) initConfigComponent() {
	titleFont := lcl.NewFont()
	titleFont.SetName("微软雅黑")
	titleFont.SetSize(12)
	titleFont.SetStyle(types.NewSet(types.FsBold))

	titleFontTwo := lcl.NewFont()
	titleFontTwo.SetName("微软雅黑")
	titleFontTwo.SetSize(10)
	titleFontTwo.SetStyle(types.NewSet(types.FsBold))

	targetPlatformTitle := lcl.NewLabel(m)
	targetPlatformTitle.SetFont(titleFont)
	targetPlatformTitle.SetCaption("目标平台与架构")
	targetPlatformTitle.SetTop(5)
	targetPlatformTitle.SetLeft(10)
	targetPlatformTitle.SetParent(m.buildTabPageConfig)

	platformTitle := lcl.NewLabel(m)
	platformTitle.SetCaption("平台")
	platformTitle.SetLeft(10)
	platformTitle.SetTop(35)
	platformTitle.SetFont(titleFontTwo)
	platformTitle.SetParent(m.buildTabPageConfig)
}

func (m *TBuildForm) initBuildComponent() {

}

func (m *TBuildForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TBuildForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TBuildForm) closeClick(sender lcl.IObject) {
	m.Close()
}

func (m *TBuildForm) saveClick(sender lcl.IObject) {

}
