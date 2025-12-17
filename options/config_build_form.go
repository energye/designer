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
	"github.com/energye/designer/resources"
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
	windowsCheckBox     lcl.ICheckBox
	macOSCheckBox       lcl.ICheckBox
	linuxCheckBox       lcl.ICheckBox
	x86_64CheckBox      lcl.ICheckBox
	aarch64CheckBox     lcl.ICheckBox
	i386CheckBox        lcl.ICheckBox
	armCheckBox         lcl.ICheckBox
	loongarch64CheckBox lcl.ICheckBox
	outputEdit          lcl.IEdit
	selectOutputDirBtn  *wg.TButton

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

	m.windowsCheckBox = lcl.NewCheckBox(m)
	m.windowsCheckBox.SetCaption("Windows")
	m.windowsCheckBox.SetLeft(20)
	m.windowsCheckBox.SetTop(60)
	m.windowsCheckBox.SetFont(m.font)
	m.windowsCheckBox.SetParent(m.buildTabPageConfig)
	m.macOSCheckBox = lcl.NewCheckBox(m)
	m.macOSCheckBox.SetCaption("macOS")
	m.macOSCheckBox.SetLeft(120)
	m.macOSCheckBox.SetTop(60)
	m.macOSCheckBox.SetFont(m.font)
	m.macOSCheckBox.SetParent(m.buildTabPageConfig)
	m.linuxCheckBox = lcl.NewCheckBox(m)
	m.linuxCheckBox.SetCaption("Linux")
	m.linuxCheckBox.SetLeft(210)
	m.linuxCheckBox.SetTop(60)
	m.linuxCheckBox.SetFont(m.font)
	m.linuxCheckBox.SetParent(m.buildTabPageConfig)

	archTitle := lcl.NewLabel(m)
	archTitle.SetCaption("架构")
	archTitle.SetLeft(10)
	archTitle.SetTop(95)
	archTitle.SetFont(titleFontTwo)
	archTitle.SetParent(m.buildTabPageConfig)

	m.x86_64CheckBox = lcl.NewCheckBox(m)
	m.x86_64CheckBox.SetCaption("x86_64")
	m.x86_64CheckBox.SetLeft(20)
	m.x86_64CheckBox.SetTop(120)
	m.x86_64CheckBox.SetFont(m.font)
	m.x86_64CheckBox.SetParent(m.buildTabPageConfig)
	m.aarch64CheckBox = lcl.NewCheckBox(m)
	m.aarch64CheckBox.SetCaption("aarch64")
	m.aarch64CheckBox.SetLeft(120)
	m.aarch64CheckBox.SetTop(120)
	m.aarch64CheckBox.SetFont(m.font)
	m.aarch64CheckBox.SetParent(m.buildTabPageConfig)
	m.i386CheckBox = lcl.NewCheckBox(m)
	m.i386CheckBox.SetCaption("i386")
	m.i386CheckBox.SetLeft(210)
	m.i386CheckBox.SetTop(120)
	m.i386CheckBox.SetFont(m.font)
	m.i386CheckBox.SetParent(m.buildTabPageConfig)
	m.armCheckBox = lcl.NewCheckBox(m)
	m.armCheckBox.SetCaption("arm")
	m.armCheckBox.SetLeft(300)
	m.armCheckBox.SetTop(120)
	m.armCheckBox.SetFont(m.font)
	m.armCheckBox.SetParent(m.buildTabPageConfig)
	m.loongarch64CheckBox = lcl.NewCheckBox(m)
	m.loongarch64CheckBox.SetCaption("loongarch64")
	m.loongarch64CheckBox.SetLeft(390)
	m.loongarch64CheckBox.SetTop(120)
	m.loongarch64CheckBox.SetFont(m.font)
	m.loongarch64CheckBox.SetEnabled(false)
	m.loongarch64CheckBox.SetParent(m.buildTabPageConfig)

	outputTitle := lcl.NewLabel(m)
	outputTitle.SetCaption("输出目录")
	outputTitle.SetLeft(10)
	outputTitle.SetTop(155)
	outputTitle.SetFont(titleFontTwo)
	outputTitle.SetParent(m.buildTabPageConfig)

	m.outputEdit = lcl.NewEdit(m)
	m.outputEdit.SetBounds(20, 185, 480, 30)
	m.outputEdit.SetFont(m.font)
	m.outputEdit.SetTextHint("构建二进制输出目录, 默认: ./build")
	m.outputEdit.SetText("./build")
	m.outputEdit.SetParent(m.buildTabPageConfig)

	m.selectOutputDirBtn = wg.NewButton(m)
	m.selectOutputDirBtn.SetIconFormBytes(resources.Images("actions/add.png"))
	m.selectOutputDirBtn.SetRadius(3)
	selectOutputDirRect := types.TRect{Left: m.outputEdit.Left() + m.outputEdit.Width() + 5, Top: 185}
	selectOutputDirRect.SetWidth(30)
	if tool.IsLinux {
		selectOutputDirRect.SetHeight(35)
	} else {
		selectOutputDirRect.SetHeight(30)
	}
	m.selectOutputDirBtn.SetBoundsRect(selectOutputDirRect)
	m.selectOutputDirBtn.SetParent(m.buildTabPageConfig)
	m.selectOutputDirBtn.SetOnClick(m.selectOutputDirClick)
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

func (m *TBuildForm) selectOutputDirClick(sender lcl.IObject) {
	m.selectDir.SetTitle("选择输出目录")
	if m.selectDir.Execute() {
		output := m.selectDir.FileName()
		m.outputEdit.SetText(output)
	}
}
