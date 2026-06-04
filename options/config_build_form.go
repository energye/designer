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
	"context"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/rtl/version"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"path/filepath"
	"strings"
)

var (
	buildFormWidth  = int32(555)
	buildFormHeight = int32(600)
)

func NewBuildForm() *TBuildForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TBuildForm{TEngForm: *newEngForm.(*lcl.TEngForm)}
	newForm.FormCreate(newEngForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TBuildForm struct {
	lcl.TEngForm
	closing      bool
	font         lcl.IFont
	titleFont    lcl.IFont
	titleFontTwo lcl.IFont
	selectDir    lcl.ISelectDirectoryDialog
	//openFile     lcl.IOpenDialog

	buildTab            *wg.TTab
	buildTabPageConfig  *wg.TPage
	buildTabPagePackage *wg.TPage

	// 基础配置
	windowsCheckBox         lcl.ICheckBox
	macOSCheckBox           lcl.ICheckBox
	linuxCheckBox           lcl.ICheckBox
	amd64CheckBox           lcl.ICheckBox
	arm64CheckBox           lcl.ICheckBox
	i386CheckBox            lcl.ICheckBox
	armCheckBox             lcl.ICheckBox
	loong64CheckBox         lcl.ICheckBox
	uiWin32Box              lcl.ICheckBox
	uiCocoaBox              lcl.ICheckBox
	uiGtk2Box               lcl.ICheckBox
	uiGtk3Box               lcl.ICheckBox
	outputEdit              lcl.IEdit
	selectOutputDirBtn      *wg.TButton
	buildFileNameEdit       lcl.IEdit
	buildModeDebugRdo       lcl.IRadioButton
	buildModeReleaseRdo     lcl.IRadioButton
	buildCGOEnabledChk      lcl.ICheckBox
	buildOtherPlatformChk   lcl.ICheckBox
	buildArgsEdit           lcl.IEdit
	codeObfuscationCheckBox lcl.ICheckBox
	disableDebugCheckBox    lcl.ICheckBox

	// 构建打包
	platformTab            *wg.TTab
	platformTabPageWindows *wg.TPage
	platformTabPageMacOS   *wg.TPage
	platformTabPageLinux   *wg.TPage

	packageNameEdit lcl.ILabeledEdit

	// 打包配置
	winPackConfigTab                              *wg.TTab
	winPackConfigTabPageBinSign                   *wg.TPage
	winPackConfigTabPageAssociateFiles            *wg.TPage
	winPackConfigTabPageAssociateProtocols        *wg.TPage
	winPackConfigTabPageAppxAssets                *wg.TPage
	winPackConfigTabPageNSISAssets                *wg.TPage
	winPackConfigTabPageNSISLicense               *wg.TPage
	winPackConfigTabPageBinSignMemoBox            *TCommonMemoBox
	winPackConfigTabPageAssociateFilesMemoBox     *TCommonMemoBox
	winPackConfigTabPageAssociateProtocolsMemoBox *TCommonMemoBox
	winPackConfigTabPageAppxAssetsMemoBox         *TCommonMemoBox
	winPackConfigTabPageNSISAssetsMemoBox         *TCommonMemoBox
	winPackConfigTabPageNSISLicenseMemoBox        *TCommonMemoBox

	winMsiCheckBox        lcl.ICheckBox
	winExeCheckBox        lcl.ICheckBox
	winDefaultInstallEdit lcl.ILabeledEdit
	winSignEnable         *TEnableButton

	macDMGCheckBox       lcl.ICheckBox
	macPKGCheckBox       lcl.ICheckBox
	macCommonLibCheckBox lcl.ICheckBox

	// 打包配置
	macPackConfigTab                          *wg.TTab
	macPackConfigTabPageBinSign               *wg.TPage
	macPackConfigTabPageAssociateFiles        *wg.TPage
	macPackConfigTabPageAssociateProtocols    *wg.TPage
	macPackConfigTabPageUniversalLink         *wg.TPage
	macPackConfigTabPageBinSignMemoBox        *TCommonMemoBox
	macPackConfigTabPageAssociateFilesBox     *TCommonMemoBox
	macPackConfigTabPageAssociateProtocolsBox *TCommonMemoBox
	macPackConfigTabPageUniversalLinkBox      *TCommonMemoBox
	macSignEnable                             *TEnableButton

	linuxDEBCheckBox      lcl.ICheckBox
	linuxRPMCheckBox      lcl.ICheckBox
	linuxAppImageCheckBox lcl.ICheckBox
	dependsEdit           lcl.IEdit
	categoriesEdit        lcl.IEdit
	homepageEdit          lcl.IEdit
	maintainerEdit        lcl.IEdit
	licenseEdit           lcl.IEdit

	// 操作按钮
	saveBtn       *wg.TButton
	packageBtn    *wg.TButton
	packageing    bool
	packageCancel func()

	statusBar lcl.IStatusBar
}

func (m *TBuildForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TBuildForm FormCreate")
	fontSize := int32(10)
	if tool.IsLinux {
		fontSize = 10
	}
	m.SetName("BuildForm")
	m.SetCaption(metadata.GI18n.Dict("BuildForm.Caption"))
	m.SetWidth(buildFormWidth)
	m.SetHeight(buildFormHeight)
	//constr := m.Constraints()
	//constr.SetMaxWidth(buildFormWidth)
	//constr.SetMaxHeight(buildFormHeight)
	//constr.SetMinWidth(buildFormWidth)
	//constr.SetMinHeight(buildFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	SetWindowCenterByMainWindow(m)

	m.font = lcl.NewFont()
	m.font.SetName("微软雅黑")
	m.font.SetSize(fontSize)
	m.SetColor(colors.ClWhite)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)
	m.selectDir.SetName("BuildFormSelectDir")

	{
		tabBR := types.TRect{Left: 0, Top: 5}
		tabBR.SetWidth(m.Width())
		tabBR.SetHeight(m.Height() - tabBR.Top)
		m.buildTab = wg.NewTab(m)
		m.buildTab.Margin = 0
		m.buildTab.SetBoundsRect(tabBR)
		m.buildTab.SetColor(colors.ClWhite)
		m.buildTab.EnableScrollButton(false)
		m.buildTab.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.buildTab.SetParent(m)
		m.buildTab.SetOnChange(func(sender lcl.IObject) {
			for _, page := range m.buildTab.Pages() {
				if page.Active() {
					page.Button().SetColor(tabActiveBgColor)
					page.Button().Font().SetColor(tabActiveTextColor)
					page.Button().SetBorderColor(wg.BbdNone, tabActiveBorderColor)
				} else {
					page.Button().SetColor(tabNoActiveBgColor)
					page.Button().Font().SetColor(tabNoActiveTextColor)
					page.Button().SetBorderColor(wg.BbdNone, tabNoActiveBorderColor)
				}
			}
		})

		m.buildTabPageConfig = m.buildTab.NewPage()
		m.buildTabPageConfig.SetName("BuildFormConfig")
		m.buildTabPageConfig.ICustomPanel.SetCaption("")
		m.buildTabPageConfig.SetCaption(metadata.GI18n.Dict("BuildFormConfig.Caption"))
		//m.buildTabPageConfig.Button().SetAutoSize(true)
		m.buildTabPageConfig.Button().SetWidth(150)
		m.buildTabPageConfig.Button().SetCursor(types.CrHandPoint)

		m.buildTabPagePackage = m.buildTab.NewPage()
		m.buildTabPagePackage.SetName("BuildFormPackage")
		m.buildTabPagePackage.ICustomPanel.SetCaption("")
		m.buildTabPagePackage.SetCaption(metadata.GI18n.Dict("BuildFormPackage.Caption"))
		//m.buildTabPagePackage.Button().SetAutoSize(true)
		m.buildTabPagePackage.Button().SetWidth(150)
		m.buildTabPagePackage.Button().SetCursor(types.CrHandPoint)

		m.buildTabPageConfig.SetActive(true)

	}
	// 初始化创建配置组件
	m.initConfigComponent()
	// 初始化构建打包组件
	m.initBuildComponent()

	{
		btnFont := lcl.NewFont()
		btnFont.SetSize(10)

		saveBtnRect := types.TRect{Left: 390, Top: 0}
		saveBtnRect.SetWidth(75)
		saveBtnRect.SetHeight(25)
		m.saveBtn = wg.NewButton(m)
		m.saveBtn.SetName("BuildFormSaveBtn")
		m.saveBtn.SetText(metadata.GI18n.Dict("BuildFormSaveBtn.Caption"))
		m.saveBtn.SetFont(btnFont)
		m.saveBtn.SetRadius(3)
		m.saveBtn.SetCursor(types.CrHandPoint)
		m.saveBtn.SetBoundsRect(saveBtnRect)
		m.saveBtn.SetColor(colors.RGBToColor(243, 244, 246))
		m.saveBtn.Font().SetColor(colors.RGBToColor(55, 65, 81))
		m.saveBtn.SetAnchors(types.NewSet(types.AkRight, types.AkTop))
		m.saveBtn.SetParent(m.buildTab)
		m.saveBtn.SetOnClick(m.saveClick)

		buildBtnRect := types.TRect{Left: saveBtnRect.Left + saveBtnRect.Width() + 10, Top: saveBtnRect.Top}
		buildBtnRect.SetWidth(75)
		buildBtnRect.SetHeight(25)
		m.packageBtn = wg.NewButton(m)
		m.packageBtn.SetName("BuildFormPackageBtn")
		m.packageBtn.SetText(metadata.GI18n.Dict("BuildFormPackageBtn.Caption"))
		m.packageBtn.SetFont(btnFont)
		m.packageBtn.SetCursor(types.CrHandPoint)
		m.packageBtn.Font().SetColor(colors.ClWhite)
		m.packageBtn.SetRadius(3)
		m.packageBtn.SetBoundsRect(buildBtnRect)
		m.packageBtn.SetColor(colors.RGBToColor(34, 197, 94))
		m.packageBtn.SetAnchors(types.NewSet(types.AkRight, types.AkTop))
		m.packageBtn.SetParent(m.buildTab)
		m.packageBtn.SetOnClick(m.packageClick)
	}
	//(&hook.TWindowHook{Form: m}).Hook()

	m.statusBar = lcl.NewStatusBar(m)
	m.statusBar.SetParent(m)
	m.statusBar.SetAutoHint(true)
}

func (m *TBuildForm) initConfigComponent() {
	m.titleFont = lcl.NewFont()
	m.titleFont.SetName("微软雅黑")
	m.titleFont.SetSize(12)
	m.titleFont.SetStyle(types.NewSet(types.FsBold))

	m.titleFontTwo = lcl.NewFont()
	m.titleFontTwo.SetName("微软雅黑")
	m.titleFontTwo.SetSize(10)
	m.titleFontTwo.SetStyle(types.NewSet(types.FsBold))

	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	targetPlatformTitle := lcl.NewLabel(m)
	targetPlatformTitle.SetFont(m.titleFont)
	targetPlatformTitle.SetName("BuildFormConfigTargetPlatformTitle")
	targetPlatformTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigTargetPlatformTitle.Caption"))
	targetPlatformTitle.SetTop(nextTop(5))
	targetPlatformTitle.SetLeft(10)
	targetPlatformTitle.SetParent(m.buildTabPageConfig)

	platformTitle := lcl.NewLabel(m)
	platformTitle.SetName("BuildFormConfigPlatformTitle")
	platformTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigPlatformTitle.Caption"))
	platformTitle.SetLeft(10)
	platformTitle.SetTop(nextTop(30))
	platformTitle.SetFont(m.titleFontTwo)
	platformTitle.SetParent(m.buildTabPageConfig)

	m.windowsCheckBox = lcl.NewCheckBox(m)
	m.windowsCheckBox.SetCaption("windows")
	m.windowsCheckBox.SetLeft(20)
	m.windowsCheckBox.SetTop(nextTop(25))
	m.windowsCheckBox.SetFont(m.font)
	m.windowsCheckBox.SetChecked(bean.GProject.BuildOption.PlatformWindows)
	m.windowsCheckBox.SetParent(m.buildTabPageConfig)
	m.macOSCheckBox = lcl.NewCheckBox(m)
	m.macOSCheckBox.SetCaption("macOS")
	m.macOSCheckBox.SetLeft(120)
	m.macOSCheckBox.SetTop(m.windowsCheckBox.Top())
	m.macOSCheckBox.SetFont(m.font)
	m.macOSCheckBox.SetChecked(bean.GProject.BuildOption.PlatformMacOS)
	m.macOSCheckBox.SetParent(m.buildTabPageConfig)
	m.linuxCheckBox = lcl.NewCheckBox(m)
	m.linuxCheckBox.SetCaption("linux")
	m.linuxCheckBox.SetLeft(210)
	m.linuxCheckBox.SetTop(m.windowsCheckBox.Top())
	m.linuxCheckBox.SetFont(m.font)
	m.linuxCheckBox.SetChecked(bean.GProject.BuildOption.PlatformLinux)
	m.linuxCheckBox.SetParent(m.buildTabPageConfig)

	archTitle := lcl.NewLabel(m)
	archTitle.SetName("BuildFormConfigArchTitle")
	archTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigArchTitle.Caption"))
	archTitle.SetLeft(10)
	archTitle.SetTop(nextTop(30))
	archTitle.SetFont(m.titleFontTwo)
	archTitle.SetParent(m.buildTabPageConfig)

	m.amd64CheckBox = lcl.NewCheckBox(m)
	m.amd64CheckBox.SetCaption("amd64")
	m.amd64CheckBox.SetLeft(20)
	m.amd64CheckBox.SetTop(nextTop(25))
	m.amd64CheckBox.SetFont(m.font)
	m.amd64CheckBox.SetChecked(bean.GProject.BuildOption.ArchAmd64)
	m.amd64CheckBox.SetParent(m.buildTabPageConfig)

	m.i386CheckBox = lcl.NewCheckBox(m)
	m.i386CheckBox.SetCaption("386")
	m.i386CheckBox.SetLeft(120)
	m.i386CheckBox.SetTop(m.amd64CheckBox.Top())
	m.i386CheckBox.SetFont(m.font)
	m.i386CheckBox.SetChecked(bean.GProject.BuildOption.Arch386)
	m.i386CheckBox.SetShowHint(true)
	m.i386CheckBox.SetHint("macOS Not supported")
	m.i386CheckBox.SetParent(m.buildTabPageConfig)

	m.arm64CheckBox = lcl.NewCheckBox(m)
	m.arm64CheckBox.SetCaption("arm64")
	m.arm64CheckBox.SetLeft(210)
	m.arm64CheckBox.SetTop(m.amd64CheckBox.Top())
	m.arm64CheckBox.SetFont(m.font)
	m.arm64CheckBox.SetChecked(bean.GProject.BuildOption.ArchArm64)
	m.arm64CheckBox.SetShowHint(true)
	m.arm64CheckBox.SetHint("windows Not supported")
	m.arm64CheckBox.SetParent(m.buildTabPageConfig)

	m.armCheckBox = lcl.NewCheckBox(m)
	m.armCheckBox.SetCaption("arm")
	m.armCheckBox.SetLeft(300)
	m.armCheckBox.SetTop(m.amd64CheckBox.Top())
	m.armCheckBox.SetFont(m.font)
	m.armCheckBox.SetChecked(bean.GProject.BuildOption.ArchArm)
	m.armCheckBox.SetShowHint(true)
	m.armCheckBox.SetHint("windows and macOS Not supported")
	m.armCheckBox.SetParent(m.buildTabPageConfig)

	m.loong64CheckBox = lcl.NewCheckBox(m)
	m.loong64CheckBox.SetCaption("loong64")
	m.loong64CheckBox.SetLeft(390)
	m.loong64CheckBox.SetTop(m.amd64CheckBox.Top())
	m.loong64CheckBox.SetFont(m.font)
	m.loong64CheckBox.SetEnabled(false)
	m.loong64CheckBox.SetChecked(bean.GProject.BuildOption.ArchLoong64)
	m.loong64CheckBox.SetShowHint(true)
	m.loong64CheckBox.SetHint("Not supported temporarily")
	m.loong64CheckBox.SetParent(m.buildTabPageConfig)

	widgetTitle := lcl.NewLabel(m)
	widgetTitle.SetCaption("UI")
	widgetTitle.SetLeft(10)
	widgetTitle.SetTop(nextTop(30))
	widgetTitle.SetFont(m.titleFontTwo)
	widgetTitle.SetParent(m.buildTabPageConfig)

	m.uiWin32Box = lcl.NewCheckBox(m)
	m.uiWin32Box.SetCaption("✱Win32")
	m.uiWin32Box.SetLeft(20)
	m.uiWin32Box.SetTop(nextTop(25))
	m.uiWin32Box.SetFont(m.font)
	//m.uiWin32Box.SetChecked(gProject.BuildOption.UIWin32)
	m.uiWin32Box.SetChecked(true)
	m.uiWin32Box.SetShowHint(true)
	m.uiWin32Box.SetHint("Windows")
	m.uiWin32Box.SetParent(m.buildTabPageConfig)

	m.uiCocoaBox = lcl.NewCheckBox(m)
	m.uiCocoaBox.SetCaption("✱Cocoa")
	m.uiCocoaBox.SetLeft(120)
	m.uiCocoaBox.SetTop(m.uiWin32Box.Top())
	m.uiCocoaBox.SetFont(m.font)
	//m.uiCocoaBox.SetChecked(gProject.BuildOption.UICocoa)
	m.uiCocoaBox.SetChecked(true)
	m.uiCocoaBox.SetShowHint(true)
	m.uiCocoaBox.SetHint("macOS")
	m.uiCocoaBox.SetParent(m.buildTabPageConfig)

	m.uiGtk2Box = lcl.NewCheckBox(m)
	m.uiGtk2Box.SetCaption("GTK2")
	m.uiGtk2Box.SetLeft(210)
	m.uiGtk2Box.SetTop(m.uiWin32Box.Top())
	m.uiGtk2Box.SetFont(m.font)
	m.uiGtk2Box.SetChecked(bean.GProject.BuildOption.UIGtk2)
	m.uiGtk2Box.SetShowHint(true)
	m.uiGtk2Box.SetHint("Linux: UI is LCL")
	m.uiGtk2Box.SetParent(m.buildTabPageConfig)

	m.uiGtk3Box = lcl.NewCheckBox(m)
	m.uiGtk3Box.SetCaption("GTK3")
	m.uiGtk3Box.SetLeft(300)
	m.uiGtk3Box.SetTop(m.uiWin32Box.Top())
	m.uiGtk3Box.SetFont(m.font)
	m.uiGtk3Box.SetChecked(bean.GProject.BuildOption.UIGtk3)
	m.uiGtk3Box.SetShowHint(true)
	m.uiGtk3Box.SetHint("Linux: UI is LCL + (WV or CEF)")
	m.uiGtk3Box.SetParent(m.buildTabPageConfig)

	outputTitle := lcl.NewLabel(m)
	outputTitle.SetName("BuildFormConfigOutputTitle")
	outputTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigOutputTitle.Caption"))
	outputTitle.SetLeft(10)
	outputTitle.SetTop(nextTop(30))
	outputTitle.SetFont(m.titleFontTwo)
	outputTitle.SetParent(m.buildTabPageConfig)

	m.outputEdit = lcl.NewEdit(m)
	m.outputEdit.SetName("BuildFormConfigOutputEdit")
	m.outputEdit.SetBounds(20, nextTop(30), 260, 30)
	m.outputEdit.SetFont(m.font)
	m.outputEdit.SetShowHint(true)
	m.outputEdit.SetHint(metadata.GI18n.Dict("BuildFormConfigOutputEdit.Hint"))
	m.outputEdit.SetTextHint(metadata.GI18n.Dict("BuildFormConfigOutputEdit.TextHint"))
	m.outputEdit.SetText(bean.GProject.BuildOption.Output)
	m.outputEdit.SetParent(m.buildTabPageConfig)

	m.selectOutputDirBtn = wg.NewButton(m)
	m.selectOutputDirBtn.SetIconFormBytes(resources.Images("menu/menu_project_open.png"))
	m.selectOutputDirBtn.SetRadius(3)
	m.selectOutputDirBtn.SetCursor(types.CrHandPoint)
	m.selectOutputDirBtn.SetColor(grayBtnColor)
	selectOutputDirRect := types.TRect{Left: m.outputEdit.Left() + m.outputEdit.Width() + 5, Top: m.outputEdit.Top() - 2}
	selectOutputDirRect.SetWidth(30)
	if tool.IsLinux {
		selectOutputDirRect.SetHeight(35)
	} else {
		selectOutputDirRect.SetHeight(25)
	}
	m.selectOutputDirBtn.SetBoundsRect(selectOutputDirRect)
	m.selectOutputDirBtn.SetParent(m.buildTabPageConfig)
	m.selectOutputDirBtn.SetOnClick(m.selectOutputDirClick)

	buildFileNameTitle := lcl.NewLabel(m)
	buildFileNameTitle.SetName("BuildFormConfigExeNameTitle")
	buildFileNameTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigExeNameTitle.Caption"))
	buildFileNameTitle.SetLeft(m.outputEdit.Left() + m.outputEdit.Width() + selectOutputDirRect.Width() + 20)
	buildFileNameTitle.SetTop(outputTitle.Top())
	buildFileNameTitle.SetFont(m.titleFontTwo)
	buildFileNameTitle.SetParent(m.buildTabPageConfig)

	m.buildFileNameEdit = lcl.NewEdit(m)
	m.buildFileNameEdit.SetName("BuildFormConfigExeNameEdit")
	m.buildFileNameEdit.SetBounds(buildFileNameTitle.Left()+10, m.outputEdit.Top(), 195, 30)
	m.buildFileNameEdit.SetFont(m.font)
	m.buildFileNameEdit.SetTextHint(metadata.GI18n.Dict("BuildFormConfigExeNameEdit.TextHint"))
	m.buildFileNameEdit.SetText(bean.GProject.BuildOption.BuildFileName)
	m.buildFileNameEdit.SetParent(m.buildTabPageConfig)

	compileArgsTitle := lcl.NewLabel(m)
	compileArgsTitle.SetName("BuildFormConfigCompileArgsTitle")
	compileArgsTitle.SetFont(m.titleFont)
	compileArgsTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigCompileArgsTitle.Caption"))
	compileArgsTitle.SetTop(nextTop(50))
	compileArgsTitle.SetLeft(10)
	compileArgsTitle.SetParent(m.buildTabPageConfig)

	buildModeTitle := lcl.NewLabel(m)
	buildModeTitle.SetName("BuildFormConfigBuildModeTitle")
	buildModeTitle.SetFont(m.titleFontTwo)
	buildModeTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildModeTitle.Caption"))
	buildModeTitle.SetTop(nextTop(35))
	buildModeTitle.SetLeft(10)
	buildModeTitle.SetParent(m.buildTabPageConfig)

	m.buildModeDebugRdo = lcl.NewRadioButton(m)
	m.buildModeDebugRdo.SetName("BuildFormConfigBuildModeDebugRdo")
	m.buildModeDebugRdo.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildModeDebugRdo.Caption"))
	m.buildModeDebugRdo.SetLeft(20)
	m.buildModeDebugRdo.SetTop(nextTop(25))
	m.buildModeDebugRdo.SetFont(m.font)
	m.buildModeDebugRdo.SetChecked(true)
	m.buildModeDebugRdo.SetShowHint(true)
	m.buildModeDebugRdo.SetHint(metadata.GI18n.Dict("BuildFormConfigBuildModeDebugRdo.Hint"))
	m.buildModeDebugRdo.SetChecked(bean.GProject.BuildOption.BuildModeDebug)
	m.buildModeDebugRdo.SetParent(m.buildTabPageConfig)

	m.buildModeReleaseRdo = lcl.NewRadioButton(m)
	m.buildModeReleaseRdo.SetName("BuildFormConfigBuildModeReleaseRdo")
	m.buildModeReleaseRdo.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildModeReleaseRdo.Caption"))
	m.buildModeReleaseRdo.SetLeft(120)
	m.buildModeReleaseRdo.SetTop(m.buildModeDebugRdo.Top())
	m.buildModeReleaseRdo.SetFont(m.font)
	m.buildModeReleaseRdo.SetShowHint(true)
	m.buildModeReleaseRdo.SetHint(metadata.GI18n.Dict("BuildFormConfigBuildModeReleaseRdo.Hint"))
	m.buildModeReleaseRdo.SetChecked(bean.GProject.BuildOption.BuildModeRelease)
	m.buildModeReleaseRdo.SetParent(m.buildTabPageConfig)

	m.buildCGOEnabledChk = lcl.NewCheckBox(m)
	m.buildCGOEnabledChk.SetName("BuildFormConfigBuildCGOEnabledChk")
	m.buildCGOEnabledChk.SetShowHint(true)
	m.buildCGOEnabledChk.SetHint(metadata.GI18n.Dict("BuildFormConfigBuildCGOEnabledChk.Hint"))
	m.buildCGOEnabledChk.SetLeft(210)
	m.buildCGOEnabledChk.SetTop(m.buildModeDebugRdo.Top())
	m.buildCGOEnabledChk.SetFont(m.font)
	m.buildCGOEnabledChk.SetChecked(bean.GProject.BuildOption.BuildCGOEnabled)
	m.buildCGOEnabledChk.SetParent(m.buildTabPageConfig)
	setBuildCGOEnabledChkCaption := func() {
		if m.buildCGOEnabledChk.Checked() {
			m.buildCGOEnabledChk.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildCGOEnabledChk.Enable"))
			m.buildOtherPlatformChk.SetEnabled(false)
		} else {
			m.buildCGOEnabledChk.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildCGOEnabledChk.Disable"))
			m.buildOtherPlatformChk.SetEnabled(true)
		}
	}
	m.buildCGOEnabledChk.SetOnChange(func(sender lcl.IObject) {
		setBuildCGOEnabledChkCaption()
	})

	m.buildOtherPlatformChk = lcl.NewCheckBox(m)
	m.buildOtherPlatformChk.SetName("BuildFormConfigBuildOtherPlatformChk")
	m.buildOtherPlatformChk.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildOtherPlatformChk.Caption"))
	m.buildOtherPlatformChk.SetShowHint(true)
	m.buildOtherPlatformChk.SetHint(metadata.GI18n.Dict("BuildFormConfigBuildOtherPlatformChk.Hint"))
	m.buildOtherPlatformChk.SetLeft(310)
	m.buildOtherPlatformChk.SetTop(m.buildCGOEnabledChk.Top())
	m.buildOtherPlatformChk.SetFont(m.font)
	m.buildOtherPlatformChk.SetChecked(bean.GProject.BuildOption.BuildOtherPlatform)
	m.buildOtherPlatformChk.SetParent(m.buildTabPageConfig)

	setBuildCGOEnabledChkCaption()

	buildArgsTitle := lcl.NewLabel(m)
	buildArgsTitle.SetName("BuildFormConfigBuildArgsTitle")
	buildArgsTitle.SetFont(m.titleFontTwo)
	buildArgsTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigBuildArgsTitle.Caption"))
	buildArgsTitle.SetTop(nextTop(35))
	buildArgsTitle.SetLeft(10)
	buildArgsTitle.SetParent(m.buildTabPageConfig)

	m.buildArgsEdit = lcl.NewEdit(m)
	m.buildArgsEdit.SetName("BuildFormConfigBuildArgsEdit")
	m.buildArgsEdit.SetBounds(20, nextTop(30), 515, 30)
	m.buildArgsEdit.SetFont(m.font)
	m.buildArgsEdit.SetTextHint(metadata.GI18n.Dict("BuildFormConfigBuildArgsEdit.TextHint"))
	m.buildArgsEdit.SetText(bean.GProject.BuildOption.GoArgs)
	m.buildArgsEdit.SetParent(m.buildTabPageConfig)

	decompileTitle := lcl.NewLabel(m)
	decompileTitle.SetName("BuildFormConfigDecompileTitle")
	decompileTitle.SetFont(m.titleFontTwo)
	decompileTitle.SetCaption(metadata.GI18n.Dict("BuildFormConfigDecompileTitle.Caption"))
	decompileTitle.SetTop(nextTop(45))
	decompileTitle.SetLeft(10)
	decompileTitle.SetParent(m.buildTabPageConfig)

	m.codeObfuscationCheckBox = lcl.NewCheckBox(m)
	m.codeObfuscationCheckBox.SetName("BuildFormConfigCodeObfuscationCheckBox")
	m.codeObfuscationCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormConfigCodeObfuscationCheckBox.Caption"))
	m.codeObfuscationCheckBox.SetLeft(20)
	m.codeObfuscationCheckBox.SetTop(nextTop(25))
	m.codeObfuscationCheckBox.SetFont(m.font)
	m.codeObfuscationCheckBox.SetShowHint(true)
	m.codeObfuscationCheckBox.SetHint(metadata.GI18n.Dict("BuildFormConfigCodeObfuscationCheckBox.Hint"))
	m.codeObfuscationCheckBox.SetEnabled(false)
	m.codeObfuscationCheckBox.SetChecked(bean.GProject.BuildOption.CodeObfuscation)
	m.codeObfuscationCheckBox.SetParent(m.buildTabPageConfig)

	m.disableDebugCheckBox = lcl.NewCheckBox(m)
	m.disableDebugCheckBox.SetCaption("BuildFormConfigDisableDebugCheckBox")
	m.disableDebugCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormConfigDisableDebugCheckBox.Caption"))
	m.disableDebugCheckBox.SetLeft(210)
	m.disableDebugCheckBox.SetTop(m.codeObfuscationCheckBox.Top())
	m.disableDebugCheckBox.SetFont(m.font)
	m.disableDebugCheckBox.SetShowHint(true)
	m.disableDebugCheckBox.SetHint(metadata.GI18n.Dict("BuildFormConfigDisableDebugCheckBox.Hint"))
	m.disableDebugCheckBox.SetEnabled(false)
	m.disableDebugCheckBox.SetChecked(bean.GProject.BuildOption.DisableDebug)
	m.disableDebugCheckBox.SetParent(m.buildTabPageConfig)

	m.initConfigData()
}

func (m *TBuildForm) initBuildComponent() {
	m.packageNameEdit = lcl.NewLabeledEdit(m)
	m.packageNameEdit.SetBounds(80, 5, buildFormWidth-100, 30)
	m.packageNameEdit.SetFont(m.font)
	m.packageNameEdit.SetTextHint("Installer package name")
	m.packageNameEdit.SetAnchors(types.NewSet(types.AkLeft, types.AkRight, types.AkTop))
	m.packageNameEdit.SetText(bean.GProject.BuildOption.PackageName)
	editLabel := m.packageNameEdit.EditLabel()
	editLabel.SetName("BuildFormPackageNameEdit")
	editLabel.SetCaption(metadata.GI18n.Dict("BuildFormPackageNameEdit.EditLabel.Caption"))
	m.packageNameEdit.SetLabelPosition(types.LpLeft)
	m.packageNameEdit.SetParent(m.buildTabPagePackage)

	{
		tabBR := types.TRect{Left: 0, Top: 40}
		tabBR.SetWidth(m.buildTabPagePackage.Width())
		tabBR.SetHeight(m.buildTabPagePackage.Height() - tabBR.Top)
		m.platformTab = wg.NewTab(m)
		m.platformTab.Margin = 0
		m.platformTab.SetBoundsRect(tabBR)
		m.platformTab.SetColor(colors.ClWhite)
		m.platformTab.EnableScrollButton(false)
		m.platformTab.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.platformTab.SetParent(m.buildTabPagePackage)
		m.platformTab.SetOnChange(func(sender lcl.IObject) {
			for _, page := range m.platformTab.Pages() {
				if page.Active() {
					page.Button().SetColor(colors.RGBToColor(219, 234, 254))
					page.Button().Font().SetColor(colors.RGBToColor(29, 78, 216))
					page.Button().SetBorderColor(wg.BbdNone, colors.RGBToColor(219, 234, 254))
				} else {
					page.Button().SetColor(tabNoActiveBgColor)
					page.Button().Font().SetColor(tabNoActiveTextColor)
					page.Button().SetBorderColor(wg.BbdNone, tabNoActiveBorderColor)
				}
			}
		})
		// 设置标签按钮样式
		setTabPageStyle := func(page *wg.TPage) {
			page.SetTop(30)
			page.SetColor(m.platformTab.Color()) // 设置背景色
			page.SetHeight(m.platformTab.Height() - page.Top())
			page.Button().SetRadius(0)
			page.Button().SetCursor(types.CrHandPoint)
			page.Button().SetWidth(80)
			page.Button().SetHeight(25)
		}

		m.platformTabPageWindows = m.platformTab.NewPage()
		m.platformTabPageWindows.SetCaption("Windows")
		m.platformTabPageWindows.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		setTabPageStyle(m.platformTabPageWindows)
		m.initWindowsOptions()

		m.platformTabPageMacOS = m.platformTab.NewPage()
		m.platformTabPageMacOS.SetCaption("MacOS")
		setTabPageStyle(m.platformTabPageMacOS)
		m.initMacOSOptions()

		m.platformTabPageLinux = m.platformTab.NewPage()
		m.platformTabPageLinux.SetCaption("Linux")
		setTabPageStyle(m.platformTabPageLinux)
		m.initLinuxOptions()

		if tool.IsWindows {
			m.platformTabPageWindows.SetActive(true)
		} else if tool.IsDarwin {
			m.platformTabPageMacOS.SetActive(true)
		} else if tool.IsLinux {
			m.platformTabPageLinux.SetActive(true)
		}
	}

	m.initBuildData()
}

func (m *TBuildForm) initConfigData() {

}

func (m *TBuildForm) initBuildData() {

}

func (m *TBuildForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	if m.packageing {
		*canClose = false
		return
	}
	m.closing = true
}

func (m *TBuildForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TBuildForm) closeClick(sender lcl.IObject) {
	m.Close()
}

func (m *TBuildForm) selectOutputDirClick(sender lcl.IObject) {
	m.selectDir.SetTitle(metadata.GI18n.Dict("BuildFormSelectDir.Title"))
	if m.selectDir.Execute() {
		output := m.selectDir.FileName()
		m.outputEdit.SetText(output)
	}
}

func (m *TBuildForm) saveClick(sender lcl.IObject) {
	event.ConsoleWriteInfo("Build Configuration - Save")
	// 基础配置
	bean.GProject.BuildOption.PlatformWindows = m.windowsCheckBox.Checked()
	bean.GProject.BuildOption.PlatformMacOS = m.macOSCheckBox.Checked()
	bean.GProject.BuildOption.PlatformLinux = m.linuxCheckBox.Checked()
	bean.GProject.BuildOption.ArchAmd64 = m.amd64CheckBox.Checked()
	bean.GProject.BuildOption.Arch386 = m.i386CheckBox.Checked()
	bean.GProject.BuildOption.ArchArm64 = m.arm64CheckBox.Checked()
	bean.GProject.BuildOption.ArchArm = m.armCheckBox.Checked()
	bean.GProject.BuildOption.ArchLoong64 = m.loong64CheckBox.Checked()
	bean.GProject.BuildOption.UIWin32 = m.uiWin32Box.Checked()
	bean.GProject.BuildOption.UICocoa = m.uiCocoaBox.Checked()
	bean.GProject.BuildOption.UIGtk2 = m.uiGtk2Box.Checked()
	bean.GProject.BuildOption.UIGtk3 = m.uiGtk3Box.Checked()
	bean.GProject.BuildOption.Output = m.outputEdit.Text()
	bean.GProject.BuildOption.BuildFileName = m.buildFileNameEdit.Text()
	bean.GProject.BuildOption.BuildModeDebug = m.buildModeDebugRdo.Checked()
	bean.GProject.BuildOption.BuildModeRelease = m.buildModeReleaseRdo.Checked()
	bean.GProject.BuildOption.BuildCGOEnabled = m.buildCGOEnabledChk.Checked()
	bean.GProject.BuildOption.BuildOtherPlatform = m.buildOtherPlatformChk.Checked()
	bean.GProject.BuildOption.GoArgs = m.buildArgsEdit.Text()
	bean.GProject.BuildOption.CodeObfuscation = m.codeObfuscationCheckBox.Checked()
	bean.GProject.BuildOption.DisableDebug = m.disableDebugCheckBox.Checked()
	// 打包配置
	packageName := strings.TrimSpace(m.packageNameEdit.Text())
	if packageName == "" {
		packageName = bean.GProject.Name
	}
	bean.GProject.BuildOption.PackageName = packageName
	// windows
	bean.GProject.BuildOption.WinMsi = m.winMsiCheckBox.Checked()
	bean.GProject.BuildOption.WinExe = m.winExeCheckBox.Checked()
	bean.GProject.BuildOption.WinDefaultInstall = m.winDefaultInstallEdit.Text()
	if m.winPackConfigTabPageAssociateFilesMemoBox != nil {
		bean.GProject.BuildOption.WinAssociateFileList = m.winPackConfigTabPageAssociateFilesMemoBox.Lines()
	}
	if m.winPackConfigTabPageAssociateProtocolsMemoBox != nil {
		bean.GProject.BuildOption.WinAssociateProtocolList = m.winPackConfigTabPageAssociateProtocolsMemoBox.Lines()
	}
	if m.winPackConfigTabPageNSISLicenseMemoBox != nil {
		licensePath := filepath.Join(bean.ResourcePath(), bean.GProject.Name+"-license.txt")
		licenseFileName := ""
		_ = os.Remove(licensePath)
		lines := m.winPackConfigTabPageNSISLicenseMemoBox.Lines()
		if len(lines) > 0 {
			data := strings.Join(lines, "\n")
			if strings.TrimSpace(data) != "" {
				utf8Bom := []byte{0xEF, 0xBB, 0xBF}
				licenseData := append(utf8Bom, data...)
				_ = os.WriteFile(licensePath, licenseData, 0644)
				_, licenseFileName = filepath.Split(licensePath)
			}
		}
		bean.GProject.BuildOption.NSIS.License = licenseFileName
	}
	if m.winPackConfigTabPageAppxAssetsMemoBox != nil {
		var PropertiesLogo, Square44x44Logo, Square150x150Logo, Wide310x150Logo, SplashScreen, AssociateFileIcon, AssociateProtocolLogo string
		for _, line := range m.winPackConfigTabPageAppxAssetsMemoBox.Lines() {
			line = strings.TrimSpace(line)
			assets := strings.Split(line, "=")
			if len(assets) == 2 {
				name := strings.ToLower(strings.TrimSpace(assets[0]))
				image := strings.TrimSpace(assets[1])
				switch name {
				case "propertieslogo":
					PropertiesLogo = image
				case "square44x44logo":
					Square44x44Logo = image
				case "square150x150logo":
					Square150x150Logo = image
				case "wide310x150logo":
					Wide310x150Logo = image
				case "splashscreen":
					SplashScreen = image
				case "associatefileicon":
					AssociateFileIcon = image
				case "associateprotocollogo":
					AssociateProtocolLogo = image
				}
			}
		}
		bean.GProject.BuildOption.WinAppx.PropertiesLogo = PropertiesLogo
		bean.GProject.BuildOption.WinAppx.Square44x44Logo = Square44x44Logo
		bean.GProject.BuildOption.WinAppx.Square150x150Logo = Square150x150Logo
		bean.GProject.BuildOption.WinAppx.Wide310x150Logo = Wide310x150Logo
		bean.GProject.BuildOption.WinAppx.SplashScreen = SplashScreen
		bean.GProject.BuildOption.WinAppx.AssociateFileIcon = AssociateFileIcon
		bean.GProject.BuildOption.WinAppx.AssociateProtocolLogo = AssociateProtocolLogo
	}
	if m.winPackConfigTabPageNSISAssetsMemoBox != nil {
		var winNsisAssets []string
		var nsisWelcomeBanner, nsisHeaderBanner, icon, unIcon string
		for _, line := range m.winPackConfigTabPageNSISAssetsMemoBox.Lines() {
			banner := strings.Split(line, "=")
			if len(banner) == 2 {
				winNsisAssets = append(winNsisAssets, line)
			}
		}
		for _, line := range winNsisAssets {
			banner := strings.Split(line, "=")
			if len(banner) == 2 {
				name := strings.ToLower(strings.TrimSpace(banner[0]))
				image := strings.TrimSpace(banner[1])
				if name == "welcome" && image != "" {
					nsisWelcomeBanner = image
				} else if name == "header" && image != "" {
					nsisHeaderBanner = image
				} else if name == "icon" && image != "" {
					icon = image
				} else if name == "unicon" && image != "" {
					unIcon = image
				}
			}
		}
		bean.GProject.BuildOption.NSIS.WelcomeBanner = nsisWelcomeBanner
		bean.GProject.BuildOption.NSIS.HeaderBanner = nsisHeaderBanner
		bean.GProject.BuildOption.NSIS.ICON = icon
		bean.GProject.BuildOption.NSIS.UnICON = unIcon
	}

	bean.GProject.BuildOption.WinSign.Enable = !m.winSignEnable.Disable()
	if m.winPackConfigTabPageBinSignMemoBox != nil {
		var signCMD []string
		for _, line := range m.winPackConfigTabPageBinSignMemoBox.Lines() {
			banner := strings.Split(line, "=")
			if len(banner) == 2 {
				signCMD = append(signCMD, line)
			}
		}
		bean.GProject.BuildOption.WinSign.Cert = signCMD
	}

	// mac
	bean.GProject.BuildOption.MacDMG = m.macDMGCheckBox.Checked()
	bean.GProject.BuildOption.MacPKG = m.macPKGCheckBox.Checked()
	bean.GProject.BuildOption.MacSign.Enable = !m.macSignEnable.Disable()
	if m.macPackConfigTabPageBinSignMemoBox != nil {
		var signCMD []string
		for _, line := range m.macPackConfigTabPageBinSignMemoBox.Lines() {
			signCMD = append(signCMD, line)
		}
		bean.GProject.BuildOption.MacSign.Cert = signCMD
	}
	if m.macPackConfigTabPageAssociateFilesBox != nil {
		bean.GProject.BuildOption.MacAssociateFileList = m.macPackConfigTabPageAssociateFilesBox.Lines()
	}
	if m.macPackConfigTabPageAssociateProtocolsBox != nil {
		bean.GProject.BuildOption.MacAssociateProtocolList = m.macPackConfigTabPageAssociateProtocolsBox.Lines()
	}
	if m.macPackConfigTabPageUniversalLinkBox != nil {
		bean.GProject.BuildOption.MacUniversalLink = m.macPackConfigTabPageUniversalLinkBox.Lines()
	}
	bean.GProject.BuildOption.MacCommonLib = m.macCommonLibCheckBox.Checked()
	bean.GProject.BuildOption.LinuxDEB = m.linuxDEBCheckBox.Checked()
	bean.GProject.BuildOption.LinuxRPM = m.linuxRPMCheckBox.Checked()
	bean.GProject.BuildOption.LinuxAppImage = m.linuxAppImageCheckBox.Checked()
	bean.GProject.AppOption.Linux.Depends = m.dependsEdit.Text()
	bean.GProject.AppOption.Linux.Categories = m.categoriesEdit.Text()
	bean.GProject.AppOption.Linux.Homepage = m.homepageEdit.Text()
	bean.GProject.AppOption.Linux.Maintainer = m.maintainerEdit.Text()
	bean.GProject.AppOption.Linux.License = m.licenseEdit.Text()
	go func() {
		// 更新项目配置文件
		if err := bean.GProject.Save(); err != nil {
			event.ConsoleWriteError("Save failed: Error writing project configuration file", err.Error())
			return
		}
		if err := updateAppGoFile(); err != nil {
			event.ConsoleWriteError("Failed to create/update app.go file:", err.Error())
			return
		}
		if err := m.mergeMacOSUniversalBinary(); err != nil {
			event.ConsoleWriteError("Failed to merge Universal Binary files:", err.Error())
			return
		}
		event.ConsoleWriteInfo("Build Configuration - Save Completed")
	}()
}

func (m *TBuildForm) packageClick(sender lcl.IObject) {
	if m.packageing {
		isOK := api.MessageDlg(metadata.GI18n.Dict("BuildFormPackageBtn.PackageStopMsg"), types.MtInformation,
			types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.MrYes
		if isOK && m.packageCancel != nil {
			m.packageCancel()
		}
		return
	}
	m.packageing = true
	m.packageBtn.SetText(metadata.GI18n.Dict("BuildFormPackageBtn.Packaging"))
	m.packageBtn.SetDisable(true)
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
		bean.GProject.BuildOption.MacCommonLib = false
	}
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		m.packageCancel = cancel
		configBuildPackage(ctx)
		m.packageCancel = nil
		if m.closing {
			return
		}

		// 恢复按钮状态
		m.packageBtn.SetText(metadata.GI18n.Dict("BuildFormPackageBtn.Caption"))
		m.packageBtn.SetDisable(false)
		m.packageing = false
	}()
}
