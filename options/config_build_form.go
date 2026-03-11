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
	"errors"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/rtl/version"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
	"os"
	"path/filepath"
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
	openFile  lcl.IOpenDialog

	buildTab            *wg.TTab
	buildTabPageConfig  *wg.TPage
	buildTabPagePackage *wg.TPage

	// 基础配置
	windowsCheckBox         lcl.ICheckBox
	macOSCheckBox           lcl.ICheckBox
	linuxCheckBox           lcl.ICheckBox
	x86_64CheckBox          lcl.ICheckBox
	aarch64CheckBox         lcl.ICheckBox
	i386CheckBox            lcl.ICheckBox
	armCheckBox             lcl.ICheckBox
	loongarch64CheckBox     lcl.ICheckBox
	uiWin32Box              lcl.ICheckBox
	uiCocoaBox              lcl.ICheckBox
	uiGtk2Box               lcl.ICheckBox
	uiGtk3Box               lcl.ICheckBox
	outputEdit              lcl.IEdit
	selectOutputDirBtn      *wg.TButton
	buildFileNameEdit       lcl.IEdit
	buildModeDebugRdo       lcl.IRadioButton
	buildModeReleaseRdo     lcl.IRadioButton
	buildArgsEdit           lcl.IEdit
	codeObfuscationCheckBox lcl.ICheckBox
	disableDebugCheckBox    lcl.ICheckBox

	// 构建打包
	packageNameEdit lcl.IEdit

	winMsiCheckBox             lcl.ICheckBox
	winExeCheckBox             lcl.ICheckBox
	winDefaultInstallEdit      lcl.IEdit
	winDesktopShortcutCheckBox lcl.ICheckBox
	winAddStartMenuCheckBox    lcl.ICheckBox

	macDMGCheckBox  lcl.ICheckBox
	macPKGCheckBox  lcl.ICheckBox
	macCertCheckBox lcl.ICheckBox
	macCertComboBox lcl.IComboBox

	macCommonLibCheckBox lcl.ICheckBox

	linuxDEBCheckBox lcl.ICheckBox
	dependsEdit      lcl.IEdit

	// 操作按钮
	saveBtn    *wg.TButton
	packageBtn *wg.TButton
	packageing bool
}

func (m *TBuildForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TBuildForm FormCreate")
	fontSize := int32(10)
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
	m.openFile = lcl.NewOpenDialog(m)

	{
		m.openFile.SetTitle("打开证书")
		m.openFile.SetFilter(config.DialogFilter.MacCertFilter())
		m.openFile.SetFilterIndex(1)
	}

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

		m.packageBtn = wg.NewButton(m)
		m.packageBtn.SetText("开始打包")
		m.packageBtn.SetFont(btnFont)
		m.packageBtn.Font().SetColor(colors.ClWhite)
		m.packageBtn.SetRadius(3)
		buildBtnRect := types.TRect{Left: saveBtnRect.Left + saveBtnRect.Width() + 10, Top: saveBtnRect.Top}
		buildBtnRect.SetWidth(75)
		buildBtnRect.SetHeight(25)
		m.packageBtn.SetBoundsRect(buildBtnRect)
		m.packageBtn.SetColor(colors.RGBToColor(46, 204, 113))
		m.packageBtn.SetParent(m.buildTab)
		m.packageBtn.SetOnClick(m.packageClick)
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

	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	targetPlatformTitle := lcl.NewLabel(m)
	targetPlatformTitle.SetFont(titleFont)
	targetPlatformTitle.SetCaption("平台与架构")
	targetPlatformTitle.SetTop(nextTop(5))
	targetPlatformTitle.SetLeft(10)
	targetPlatformTitle.SetParent(m.buildTabPageConfig)

	platformTitle := lcl.NewLabel(m)
	platformTitle.SetCaption("平台")
	platformTitle.SetLeft(10)
	platformTitle.SetTop(nextTop(30))
	platformTitle.SetFont(titleFontTwo)
	platformTitle.SetParent(m.buildTabPageConfig)

	m.windowsCheckBox = lcl.NewCheckBox(m)
	m.windowsCheckBox.SetCaption("Windows")
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
	m.linuxCheckBox.SetCaption("Linux")
	m.linuxCheckBox.SetLeft(210)
	m.linuxCheckBox.SetTop(m.windowsCheckBox.Top())
	m.linuxCheckBox.SetFont(m.font)
	m.linuxCheckBox.SetChecked(bean.GProject.BuildOption.PlatformLinux)
	m.linuxCheckBox.SetParent(m.buildTabPageConfig)

	archTitle := lcl.NewLabel(m)
	archTitle.SetCaption("架构")
	archTitle.SetLeft(10)
	archTitle.SetTop(nextTop(30))
	archTitle.SetFont(titleFontTwo)
	archTitle.SetParent(m.buildTabPageConfig)

	m.x86_64CheckBox = lcl.NewCheckBox(m)
	m.x86_64CheckBox.SetCaption("x86_64")
	m.x86_64CheckBox.SetLeft(20)
	m.x86_64CheckBox.SetTop(nextTop(25))
	m.x86_64CheckBox.SetFont(m.font)
	m.x86_64CheckBox.SetChecked(bean.GProject.BuildOption.ArchX86_64)
	m.x86_64CheckBox.SetParent(m.buildTabPageConfig)
	m.i386CheckBox = lcl.NewCheckBox(m)
	m.i386CheckBox.SetCaption("i386")
	m.i386CheckBox.SetLeft(120)
	m.i386CheckBox.SetTop(m.x86_64CheckBox.Top())
	m.i386CheckBox.SetFont(m.font)
	m.i386CheckBox.SetChecked(bean.GProject.BuildOption.ArchI386)
	m.i386CheckBox.SetParent(m.buildTabPageConfig)
	m.aarch64CheckBox = lcl.NewCheckBox(m)
	m.aarch64CheckBox.SetCaption("aarch64")
	m.aarch64CheckBox.SetLeft(210)
	m.aarch64CheckBox.SetTop(m.x86_64CheckBox.Top())
	m.aarch64CheckBox.SetFont(m.font)
	m.aarch64CheckBox.SetChecked(bean.GProject.BuildOption.ArchAarch64)
	m.aarch64CheckBox.SetParent(m.buildTabPageConfig)
	m.armCheckBox = lcl.NewCheckBox(m)
	m.armCheckBox.SetCaption("arm")
	m.armCheckBox.SetLeft(300)
	m.armCheckBox.SetTop(m.x86_64CheckBox.Top())
	m.armCheckBox.SetFont(m.font)
	m.armCheckBox.SetChecked(bean.GProject.BuildOption.ArchArm)
	m.armCheckBox.SetParent(m.buildTabPageConfig)
	m.loongarch64CheckBox = lcl.NewCheckBox(m)
	m.loongarch64CheckBox.SetCaption("loongarch64")
	m.loongarch64CheckBox.SetLeft(390)
	m.loongarch64CheckBox.SetTop(m.x86_64CheckBox.Top())
	m.loongarch64CheckBox.SetFont(m.font)
	m.loongarch64CheckBox.SetEnabled(false)
	m.loongarch64CheckBox.SetChecked(bean.GProject.BuildOption.ArchLoongarch64)
	m.loongarch64CheckBox.SetParent(m.buildTabPageConfig)

	widgetTitle := lcl.NewLabel(m)
	widgetTitle.SetCaption("UI")
	widgetTitle.SetLeft(10)
	widgetTitle.SetTop(nextTop(30))
	widgetTitle.SetFont(titleFontTwo)
	widgetTitle.SetParent(m.buildTabPageConfig)

	m.uiWin32Box = lcl.NewCheckBox(m)
	m.uiWin32Box.SetCaption("✱Win32/64")
	m.uiWin32Box.SetLeft(20)
	m.uiWin32Box.SetTop(nextTop(25))
	m.uiWin32Box.SetFont(m.font)
	//m.uiWin32Box.SetChecked(gProject.BuildOption.UIWin32_64)
	m.uiWin32Box.SetChecked(true)
	m.uiWin32Box.SetParent(m.buildTabPageConfig)

	m.uiCocoaBox = lcl.NewCheckBox(m)
	m.uiCocoaBox.SetCaption("✱Cocoa")
	m.uiCocoaBox.SetLeft(120)
	m.uiCocoaBox.SetTop(m.uiWin32Box.Top())
	m.uiCocoaBox.SetFont(m.font)
	//m.uiCocoaBox.SetChecked(gProject.BuildOption.UICocoa)
	m.uiCocoaBox.SetChecked(true)
	m.uiCocoaBox.SetParent(m.buildTabPageConfig)

	m.uiGtk2Box = lcl.NewCheckBox(m)
	m.uiGtk2Box.SetCaption("GTK2")
	m.uiGtk2Box.SetLeft(210)
	m.uiGtk2Box.SetTop(m.uiWin32Box.Top())
	m.uiGtk2Box.SetFont(m.font)
	m.uiGtk2Box.SetChecked(bean.GProject.BuildOption.UIGtk2)
	m.uiGtk2Box.SetParent(m.buildTabPageConfig)

	m.uiGtk3Box = lcl.NewCheckBox(m)
	m.uiGtk3Box.SetCaption("GTK3")
	m.uiGtk3Box.SetLeft(300)
	m.uiGtk3Box.SetTop(m.uiWin32Box.Top())
	m.uiGtk3Box.SetFont(m.font)
	m.uiGtk3Box.SetChecked(bean.GProject.BuildOption.UIGtk3)
	m.uiGtk3Box.SetParent(m.buildTabPageConfig)

	outputTitle := lcl.NewLabel(m)
	outputTitle.SetCaption("输出目录")
	outputTitle.SetLeft(10)
	outputTitle.SetTop(nextTop(30))
	outputTitle.SetFont(titleFontTwo)
	outputTitle.SetParent(m.buildTabPageConfig)

	m.outputEdit = lcl.NewEdit(m)
	m.outputEdit.SetBounds(20, nextTop(30), 260, 30)
	m.outputEdit.SetFont(m.font)
	m.outputEdit.SetTextHint("构建二进制输出目录, 默认: ./build")
	m.outputEdit.SetText(bean.GProject.BuildOption.Output)
	m.outputEdit.SetParent(m.buildTabPageConfig)

	m.selectOutputDirBtn = wg.NewButton(m)
	m.selectOutputDirBtn.SetIconFormBytes(resources.Images("actions/add.png"))
	m.selectOutputDirBtn.SetRadius(3)
	selectOutputDirRect := types.TRect{Left: m.outputEdit.Left() + m.outputEdit.Width() + 5, Top: m.outputEdit.Top()}
	selectOutputDirRect.SetWidth(30)
	if tool.IsLinux {
		selectOutputDirRect.SetHeight(35)
	} else {
		selectOutputDirRect.SetHeight(30)
	}
	m.selectOutputDirBtn.SetBoundsRect(selectOutputDirRect)
	m.selectOutputDirBtn.SetParent(m.buildTabPageConfig)
	m.selectOutputDirBtn.SetOnClick(m.selectOutputDirClick)

	buildFileNameTitle := lcl.NewLabel(m)
	buildFileNameTitle.SetCaption("可执行文件名称")
	buildFileNameTitle.SetLeft(m.outputEdit.Left() + m.outputEdit.Width() + selectOutputDirRect.Width() + 20)
	buildFileNameTitle.SetTop(outputTitle.Top())
	buildFileNameTitle.SetFont(titleFontTwo)
	buildFileNameTitle.SetParent(m.buildTabPageConfig)

	m.buildFileNameEdit = lcl.NewEdit(m)
	m.buildFileNameEdit.SetBounds(buildFileNameTitle.Left()+10, m.outputEdit.Top(), 195, 30)
	m.buildFileNameEdit.SetFont(m.font)
	m.buildFileNameEdit.SetTextHint(`构建的二进制文件名`)
	m.buildFileNameEdit.SetText(bean.GProject.BuildOption.BuildFileName)
	m.buildFileNameEdit.SetParent(m.buildTabPageConfig)

	compileArgsTitle := lcl.NewLabel(m)
	compileArgsTitle.SetFont(titleFont)
	compileArgsTitle.SetCaption("编译")
	compileArgsTitle.SetTop(nextTop(50))
	compileArgsTitle.SetLeft(10)
	compileArgsTitle.SetParent(m.buildTabPageConfig)

	buildModeTitle := lcl.NewLabel(m)
	buildModeTitle.SetFont(titleFontTwo)
	buildModeTitle.SetCaption("编译模式")
	buildModeTitle.SetTop(nextTop(35))
	buildModeTitle.SetLeft(10)
	buildModeTitle.SetParent(m.buildTabPageConfig)

	m.buildModeDebugRdo = lcl.NewRadioButton(m)
	m.buildModeDebugRdo.SetCaption("调试模式")
	m.buildModeDebugRdo.SetLeft(20)
	m.buildModeDebugRdo.SetTop(nextTop(25))
	m.buildModeDebugRdo.SetFont(m.font)
	m.buildModeDebugRdo.SetChecked(true)
	m.buildModeDebugRdo.SetShowHint(true)
	m.buildModeDebugRdo.SetHint("调试模式保留调试信息")
	m.buildModeDebugRdo.SetChecked(bean.GProject.BuildOption.BuildModeDebug)
	m.buildModeDebugRdo.SetParent(m.buildTabPageConfig)
	m.buildModeReleaseRdo = lcl.NewRadioButton(m)
	m.buildModeReleaseRdo.SetCaption("发布模式")
	m.buildModeReleaseRdo.SetLeft(210)
	m.buildModeReleaseRdo.SetTop(m.buildModeDebugRdo.Top())
	m.buildModeReleaseRdo.SetFont(m.font)
	m.buildModeReleaseRdo.SetShowHint(true)
	m.buildModeReleaseRdo.SetHint("发布模式优化体积, 去除调试信息和符号")
	m.buildModeReleaseRdo.SetChecked(bean.GProject.BuildOption.BuildModeRelease)
	m.buildModeReleaseRdo.SetParent(m.buildTabPageConfig)

	buildArgsTitle := lcl.NewLabel(m)
	buildArgsTitle.SetFont(titleFontTwo)
	buildArgsTitle.SetCaption("构建参数")
	buildArgsTitle.SetTop(nextTop(35))
	buildArgsTitle.SetLeft(10)
	buildArgsTitle.SetParent(m.buildTabPageConfig)

	m.buildArgsEdit = lcl.NewEdit(m)
	m.buildArgsEdit.SetBounds(20, nextTop(30), 515, 30)
	m.buildArgsEdit.SetFont(m.font)
	m.buildArgsEdit.SetTextHint(`传递给 go build 额外参数 如: -tags xxx -ldflags "-xxx"`)
	m.buildArgsEdit.SetText(bean.GProject.BuildOption.GoArgs)
	m.buildArgsEdit.SetParent(m.buildTabPageConfig)

	decompileTitle := lcl.NewLabel(m)
	decompileTitle.SetFont(titleFontTwo)
	decompileTitle.SetCaption("反编译防护")
	decompileTitle.SetTop(nextTop(45))
	decompileTitle.SetLeft(10)
	decompileTitle.SetParent(m.buildTabPageConfig)

	m.codeObfuscationCheckBox = lcl.NewCheckBox(m)
	m.codeObfuscationCheckBox.SetCaption("代码混淆")
	m.codeObfuscationCheckBox.SetLeft(20)
	m.codeObfuscationCheckBox.SetTop(nextTop(25))
	m.codeObfuscationCheckBox.SetFont(m.font)
	m.codeObfuscationCheckBox.SetShowHint(true)
	m.codeObfuscationCheckBox.SetHint("对 Go 代码进行简单混淆")
	m.codeObfuscationCheckBox.SetEnabled(false)
	m.codeObfuscationCheckBox.SetChecked(bean.GProject.BuildOption.CodeObfuscation)
	m.codeObfuscationCheckBox.SetParent(m.buildTabPageConfig)
	m.disableDebugCheckBox = lcl.NewCheckBox(m)
	m.disableDebugCheckBox.SetCaption("禁止调试")
	m.disableDebugCheckBox.SetLeft(210)
	m.disableDebugCheckBox.SetTop(m.codeObfuscationCheckBox.Top())
	m.disableDebugCheckBox.SetFont(m.font)
	m.disableDebugCheckBox.SetShowHint(true)
	m.disableDebugCheckBox.SetHint("提高二进制文件的反编译难度")
	m.disableDebugCheckBox.SetEnabled(false)
	m.disableDebugCheckBox.SetChecked(bean.GProject.BuildOption.DisableDebug)
	m.disableDebugCheckBox.SetParent(m.buildTabPageConfig)

	m.initConfigData()
}

func (m *TBuildForm) initBuildComponent() {
	titleFont := lcl.NewFont()
	titleFont.SetName("微软雅黑")
	titleFont.SetSize(12)
	titleFont.SetStyle(types.NewSet(types.FsBold))

	titleFontTwo := lcl.NewFont()
	titleFontTwo.SetName("微软雅黑")
	titleFontTwo.SetSize(10)
	titleFontTwo.SetStyle(types.NewSet(types.FsBold))

	baseTop := int32(0)
	nextTop := func(top int32) int32 {
		baseTop += top
		return baseTop
	}

	packageNameTitle := lcl.NewLabel(m)
	packageNameTitle.SetFont(titleFontTwo)
	packageNameTitle.SetCaption("安装包名称")
	packageNameTitle.SetTop(nextTop(10))
	packageNameTitle.SetLeft(10)
	packageNameTitle.SetParent(m.buildTabPagePackage)
	m.packageNameEdit = lcl.NewEdit(m)
	m.packageNameEdit.SetBounds(80, packageNameTitle.Top()-5, 435, 30)
	m.packageNameEdit.SetFont(m.font)
	m.packageNameEdit.SetTextHint("安装包名称, 默认可执行文件名称")
	m.packageNameEdit.SetText(bean.GProject.BuildOption.PackageName)
	m.packageNameEdit.SetParent(m.buildTabPagePackage)

	windowsPackageTitle := lcl.NewLabel(m)
	windowsPackageTitle.SetFont(titleFont)
	windowsPackageTitle.SetCaption("Windows 打包配置")
	windowsPackageTitle.SetTop(nextTop(35))
	windowsPackageTitle.SetLeft(10)
	windowsPackageTitle.SetParent(m.buildTabPagePackage)

	windowsPackageFmtTitle := lcl.NewLabel(m)
	windowsPackageFmtTitle.SetCaption("打包格式")
	windowsPackageFmtTitle.SetLeft(10)
	windowsPackageFmtTitle.SetTop(nextTop(30))
	windowsPackageFmtTitle.SetFont(titleFontTwo)
	windowsPackageFmtTitle.SetParent(m.buildTabPagePackage)

	m.winMsiCheckBox = lcl.NewCheckBox(m)
	m.winMsiCheckBox.SetCaption("MSI 安装包")
	m.winMsiCheckBox.SetLeft(20)
	m.winMsiCheckBox.SetTop(nextTop(25))
	m.winMsiCheckBox.SetFont(m.font)
	m.winMsiCheckBox.SetChecked(bean.GProject.BuildOption.WinMsi)
	m.winMsiCheckBox.SetParent(m.buildTabPagePackage)
	m.winExeCheckBox = lcl.NewCheckBox(m)
	m.winExeCheckBox.SetCaption("EXE 安装包")
	m.winExeCheckBox.SetLeft(210)
	m.winExeCheckBox.SetTop(m.winMsiCheckBox.Top())
	m.winExeCheckBox.SetFont(m.font)
	m.winExeCheckBox.SetChecked(bean.GProject.BuildOption.WinExe)
	m.winExeCheckBox.SetParent(m.buildTabPagePackage)

	winDefaultInstallTitle := lcl.NewLabel(m)
	winDefaultInstallTitle.SetCaption("默认安装路径")
	winDefaultInstallTitle.SetLeft(10)
	winDefaultInstallTitle.SetTop(nextTop(30))
	winDefaultInstallTitle.SetFont(titleFontTwo)
	winDefaultInstallTitle.SetParent(m.buildTabPagePackage)

	m.winDefaultInstallEdit = lcl.NewEdit(m)
	m.winDefaultInstallEdit.SetBounds(20, nextTop(25), 515, 30)
	m.winDefaultInstallEdit.SetFont(m.font)
	m.winDefaultInstallEdit.SetTextHint("Windows 应用的默认安装路径 如: C:\\Program Files")
	m.winDefaultInstallEdit.SetText(bean.GProject.BuildOption.WinDefaultInstall)
	m.winDefaultInstallEdit.SetParent(m.buildTabPagePackage)

	m.winDesktopShortcutCheckBox = lcl.NewCheckBox(m)
	m.winDesktopShortcutCheckBox.SetCaption("创建桌面快捷方式")
	m.winDesktopShortcutCheckBox.SetLeft(20)
	m.winDesktopShortcutCheckBox.SetTop(nextTop(40))
	m.winDesktopShortcutCheckBox.SetFont(m.font)
	m.winDesktopShortcutCheckBox.SetChecked(bean.GProject.BuildOption.WinDesktopShortcut)
	m.winDesktopShortcutCheckBox.SetParent(m.buildTabPagePackage)

	m.winAddStartMenuCheckBox = lcl.NewCheckBox(m)
	m.winAddStartMenuCheckBox.SetCaption("添加到开始菜单")
	m.winAddStartMenuCheckBox.SetLeft(210)
	m.winAddStartMenuCheckBox.SetTop(m.winDesktopShortcutCheckBox.Top())
	m.winAddStartMenuCheckBox.SetFont(m.font)
	m.winAddStartMenuCheckBox.SetChecked(bean.GProject.BuildOption.WinAddStartMenu)
	m.winAddStartMenuCheckBox.SetParent(m.buildTabPagePackage)

	macOSPackageTitle := lcl.NewLabel(m)
	macOSPackageTitle.SetFont(titleFont)
	macOSPackageTitle.SetCaption("macOS 打包配置")
	macOSPackageTitle.SetTop(nextTop(35))
	macOSPackageTitle.SetLeft(10)
	macOSPackageTitle.SetParent(m.buildTabPagePackage)

	macOSPackageFmtTitle := lcl.NewLabel(m)
	macOSPackageFmtTitle.SetCaption("打包格式")
	macOSPackageFmtTitle.SetLeft(10)
	macOSPackageFmtTitle.SetTop(nextTop(30))
	macOSPackageFmtTitle.SetFont(titleFontTwo)
	macOSPackageFmtTitle.SetParent(m.buildTabPagePackage)

	m.macDMGCheckBox = lcl.NewCheckBox(m)
	m.macDMGCheckBox.SetCaption("DMG 镜像")
	m.macDMGCheckBox.SetLeft(20)
	m.macDMGCheckBox.SetTop(nextTop(30))
	m.macDMGCheckBox.SetFont(m.font)
	m.macDMGCheckBox.SetChecked(bean.GProject.BuildOption.MacDMG)
	m.macDMGCheckBox.SetParent(m.buildTabPagePackage)

	m.macPKGCheckBox = lcl.NewCheckBox(m)
	m.macPKGCheckBox.SetCaption("PKG 安装包")
	m.macPKGCheckBox.SetLeft(210)
	m.macPKGCheckBox.SetTop(m.macDMGCheckBox.Top())
	m.macPKGCheckBox.SetFont(m.font)
	m.macPKGCheckBox.SetChecked(bean.GProject.BuildOption.MacPKG)
	m.macPKGCheckBox.SetParent(m.buildTabPagePackage)

	m.macCertCheckBox = lcl.NewCheckBox(m)
	m.macCertCheckBox.SetCaption("签名")
	m.macCertCheckBox.SetLeft(20)
	m.macCertCheckBox.SetTop(nextTop(30))
	m.macCertCheckBox.SetFont(m.font)
	m.macCertCheckBox.SetChecked(bean.GProject.BuildOption.MacCert)
	m.macCertCheckBox.SetParent(m.buildTabPagePackage)
	m.macCertCheckBox.SetOnChange(func(sender lcl.IObject) {
		m.macCertComboBox.SetVisible(m.macCertCheckBox.Checked())
	})

	m.macCertComboBox = lcl.NewComboBox(m)
	m.macCertComboBox.SetBounds(85, m.macCertCheckBox.Top(), 450, 30)
	m.macCertComboBox.SetFont(m.font)
	m.macCertComboBox.SetTextHint(`选择用于签名的证书`)
	m.macCertComboBox.SetShowHint(true)
	m.macCertComboBox.SetDoubleBuffered(true)
	m.macCertComboBox.SetStyle(types.CsDropDownList)
	m.macCertComboBox.SetBorderStyle(types.BsSingle)
	m.macCertComboBox.SetVisible(m.macCertCheckBox.Checked())
	m.macCertComboBox.Items().Add("-- 选择证书 --")
	for _, item := range bean.GProject.BuildOption.MacCertList {
		m.macCertComboBox.Items().Add(item)
	}
	m.macCertComboBox.SetItemIndex(bean.GProject.BuildOption.MacCertListIndex)
	m.macCertComboBox.SetOnChange(m.macCertComboBoxChange)
	m.macCertComboBox.SetParent(m.buildTabPagePackage)

	m.macCommonLibCheckBox = lcl.NewCheckBox(m)
	m.macCommonLibCheckBox.SetCaption("‌通用二进制文件(Universal Binary)")
	m.macCommonLibCheckBox.SetLeft(210)
	m.macCommonLibCheckBox.SetTop(m.macCertCheckBox.Top())
	m.macCommonLibCheckBox.SetFont(m.font)
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
		bean.GProject.BuildOption.MacCommonLib = false
		m.macCommonLibCheckBox.SetEnabled(false)
	}
	m.macCommonLibCheckBox.SetChecked(bean.GProject.BuildOption.MacCommonLib)
	m.macCommonLibCheckBox.SetParent(m.buildTabPagePackage)

	linuxPackageTitle := lcl.NewLabel(m)
	linuxPackageTitle.SetFont(titleFont)
	linuxPackageTitle.SetCaption("Linux 打包配置")
	linuxPackageTitle.SetTop(nextTop(35))
	linuxPackageTitle.SetLeft(10)
	linuxPackageTitle.SetParent(m.buildTabPagePackage)

	linuxPackageFmtTitle := lcl.NewLabel(m)
	linuxPackageFmtTitle.SetCaption("打包格式")
	linuxPackageFmtTitle.SetLeft(10)
	linuxPackageFmtTitle.SetTop(nextTop(30))
	linuxPackageFmtTitle.SetFont(titleFontTwo)
	linuxPackageFmtTitle.SetParent(m.buildTabPagePackage)

	m.linuxDEBCheckBox = lcl.NewCheckBox(m)
	m.linuxDEBCheckBox.SetCaption("DEB 包")
	m.linuxDEBCheckBox.SetLeft(20)
	m.linuxDEBCheckBox.SetTop(nextTop(30))
	m.linuxDEBCheckBox.SetFont(m.font)
	m.linuxDEBCheckBox.SetChecked(bean.GProject.BuildOption.LinuxDEB)
	m.linuxDEBCheckBox.SetParent(m.buildTabPagePackage)

	dependsTitle := lcl.NewLabel(m)
	dependsTitle.SetCaption("依赖项")
	dependsTitle.SetLeft(10)
	dependsTitle.SetTop(nextTop(35))
	dependsTitle.SetFont(titleFontTwo)
	dependsTitle.SetParent(m.buildTabPagePackage)

	m.dependsEdit = lcl.NewEdit(m)
	m.dependsEdit.SetBounds(20, nextTop(30), 515, 30)
	m.dependsEdit.SetFont(m.font)
	m.dependsEdit.SetTextHint("用逗号分隔的依赖项列表, 如: libc6 (>= 2.14)")
	m.dependsEdit.SetText(bean.GProject.BuildOption.Depends)
	m.dependsEdit.SetParent(m.buildTabPagePackage)

	m.initBuildData()
}

func (m *TBuildForm) initConfigData() {

}

func (m *TBuildForm) initBuildData() {

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

func (m *TBuildForm) selectOutputDirClick(sender lcl.IObject) {
	m.selectDir.SetTitle("选择输出目录")
	if m.selectDir.Execute() {
		output := m.selectDir.FileName()
		m.outputEdit.SetText(output)
	}
}

func (m *TBuildForm) macCertComboBoxChange(sender lcl.IObject) {
	if m.macCertComboBox.ItemIndex() == 0 {
		if m.openFile.Execute() {
			output := m.openFile.FileName()
			isAdd := true
			items := m.macCertComboBox.Items()
			for i := int32(0); i < items.Count(); i++ {
				if items.Strings(i) == output {
					isAdd = false
					break
				}
			}
			if isAdd {
				m.macCertComboBox.Items().Add(output)
				m.macCertComboBox.SetText(output)
			}
		}
	}
}

func (m *TBuildForm) saveClick(sender lcl.IObject) {
	event.ConsoleWriteInfo("构建配置-保存")
	// 基础配置
	bean.GProject.BuildOption.PlatformWindows = m.windowsCheckBox.Checked()
	bean.GProject.BuildOption.PlatformMacOS = m.macOSCheckBox.Checked()
	bean.GProject.BuildOption.PlatformLinux = m.linuxCheckBox.Checked()
	bean.GProject.BuildOption.ArchX86_64 = m.x86_64CheckBox.Checked()
	bean.GProject.BuildOption.ArchI386 = m.i386CheckBox.Checked()
	bean.GProject.BuildOption.ArchAarch64 = m.aarch64CheckBox.Checked()
	bean.GProject.BuildOption.ArchArm = m.armCheckBox.Checked()
	bean.GProject.BuildOption.ArchLoongarch64 = m.loongarch64CheckBox.Checked()
	bean.GProject.BuildOption.UIWin32_64 = m.uiWin32Box.Checked()
	bean.GProject.BuildOption.UICocoa = m.uiCocoaBox.Checked()
	bean.GProject.BuildOption.UIGtk2 = m.uiGtk2Box.Checked()
	bean.GProject.BuildOption.UIGtk3 = m.uiGtk3Box.Checked()
	bean.GProject.BuildOption.Output = m.outputEdit.Text()
	bean.GProject.BuildOption.BuildFileName = m.buildFileNameEdit.Text()
	bean.GProject.BuildOption.BuildModeDebug = m.buildModeDebugRdo.Checked()
	bean.GProject.BuildOption.BuildModeRelease = m.buildModeReleaseRdo.Checked()
	bean.GProject.BuildOption.GoArgs = m.buildArgsEdit.Text()
	bean.GProject.BuildOption.CodeObfuscation = m.codeObfuscationCheckBox.Checked()
	bean.GProject.BuildOption.DisableDebug = m.disableDebugCheckBox.Checked()
	// 打包配置
	bean.GProject.BuildOption.PackageName = m.packageNameEdit.Text()
	bean.GProject.BuildOption.WinMsi = m.winMsiCheckBox.Checked()
	bean.GProject.BuildOption.WinExe = m.winExeCheckBox.Checked()
	bean.GProject.BuildOption.WinDefaultInstall = m.winDefaultInstallEdit.Text()
	bean.GProject.BuildOption.WinDesktopShortcut = m.winDesktopShortcutCheckBox.Checked()
	bean.GProject.BuildOption.WinAddStartMenu = m.winAddStartMenuCheckBox.Checked()
	bean.GProject.BuildOption.MacDMG = m.macDMGCheckBox.Checked()
	bean.GProject.BuildOption.MacPKG = m.macPKGCheckBox.Checked()
	bean.GProject.BuildOption.MacCert = m.macCertCheckBox.Checked()
	var macCertList []string
	for i := int32(1); i < m.macCertComboBox.Items().Count(); i++ {
		macCertList = append(macCertList, m.macCertComboBox.Items().Strings(i))
	}
	bean.GProject.BuildOption.MacCertList = macCertList
	bean.GProject.BuildOption.MacCertListIndex = m.macCertComboBox.ItemIndex()
	bean.GProject.BuildOption.MacCommonLib = m.macCommonLibCheckBox.Checked()
	bean.GProject.BuildOption.LinuxDEB = m.linuxDEBCheckBox.Checked()
	bean.GProject.BuildOption.Depends = m.dependsEdit.Text()
	go func() {
		// 更新项目配置文件
		if err := WriteEGPConfig(bean.GPath, bean.GProject); err != nil {
			event.ConsoleWriteError("保存-写入项目配置文件失败", err.Error())
			return
		}
		if err := updateAppGoFile(); err != nil {
			event.ConsoleWriteError("创建/更新 app.go 文件失败:", err.Error())
			return
		}
		if err := m.mergeMacOSUniversalBinary(); err != nil {
			event.ConsoleWriteError("合并通用二进制文件失败:", err.Error())
			return
		}
		event.ConsoleWriteInfo("构建配置-保存-完成")
	}()
}

func (m *TBuildForm) packageClick(sender lcl.IObject) {
	if m.packageing {
		return
	}
	m.packageing = true
	m.packageBtn.SetDisable(true)
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
		bean.GProject.BuildOption.MacCommonLib = false
	}
	go func() {
		configBuildPackage()
		m.packageBtn.SetDisable(false)
		m.packageing = false
	}()
}

func (m *TBuildForm) mergeMacOSUniversalBinary() error {
	if !tool.IsDarwin {
		return nil
	}
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
		bean.GProject.BuildOption.MacCommonLib = false
	}
	event.ConsoleWriteInfo("Merge macOS UniversalBinary, MacCommonLib:", tool.BoolToString(bean.GProject.BuildOption.MacCommonLib))
	if bean.GProject.BuildOption.MacCommonLib {
		// 启用通用二进制, 保存到 designer frameworks/runtime 目录
		libArm64 := lib.Libs().Get(lib.PathARM64Cocoa)
		if libArm64 == nil {
			return errors.New("libArm64 is nil")
		}
		libAmd64 := lib.Libs().Get(lib.PathAMD64Cocoa)
		if libAmd64 == nil {
			return errors.New("libAmd64 is nil")
		}
		outputLibPath := filepath.Join(config.Config.FrameworkDir, "runtime")
		tempArm64LibName := libArm64.OutputFilename
		tempAmd64LibName := libAmd64.OutputFilename
		arm64LibFilePath := filepath.Join(outputLibPath, tempArm64LibName)
		amd64LibFilePath := filepath.Join(outputLibPath, tempAmd64LibName)
		universalLibFilePath := filepath.Join(outputLibPath, "libenergy-darwin-universal-cocoa.dylib")
		event.ConsoleWriteInfo("Merge macOS UniversalBinary, arm64LibFilePath:", arm64LibFilePath)
		event.ConsoleWriteInfo("Merge macOS UniversalBinary, amd64LibFilePath:", amd64LibFilePath)
		_ = os.Remove(universalLibFilePath)
		cmd := command.NewCMD()
		cmd.Command("lipo", "-create", amd64LibFilePath, arm64LibFilePath, "-output", universalLibFilePath)
		event.ConsoleWriteInfo("Merge macOS UniversalBinary, universalLibFilePath:", universalLibFilePath)
	}
	return nil
}
