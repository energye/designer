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
	"github.com/energye/designer/options/bean"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
	"os"
	"path/filepath"
	"strings"
)

func (m *TBuildForm) initWindowsOptions() {
	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	windowsPackageFmtTitle := lcl.NewLabel(m)
	windowsPackageFmtTitle.SetCaption("打包格式")
	windowsPackageFmtTitle.SetLeft(10)
	windowsPackageFmtTitle.SetTop(nextTop(5))
	windowsPackageFmtTitle.SetFont(m.titleFontTwo)
	windowsPackageFmtTitle.SetParent(m.platformTabPageWindows)

	m.winMsiCheckBox = lcl.NewCheckBox(m)
	m.winMsiCheckBox.SetCaption("MSIX 安装包(MakeAppx)")
	m.winMsiCheckBox.SetLeft(20)
	m.winMsiCheckBox.SetTop(nextTop(25))
	m.winMsiCheckBox.SetFont(m.font)
	m.winMsiCheckBox.SetChecked(bean.GProject.BuildOption.WinMsi)
	m.winMsiCheckBox.SetParent(m.platformTabPageWindows)
	m.winExeCheckBox = lcl.NewCheckBox(m)
	m.winExeCheckBox.SetCaption("EXE 安装包(MakeNsis)")
	m.winExeCheckBox.SetLeft(210)
	m.winExeCheckBox.SetTop(m.winMsiCheckBox.Top())
	m.winExeCheckBox.SetFont(m.font)
	m.winExeCheckBox.SetChecked(bean.GProject.BuildOption.WinExe)
	m.winExeCheckBox.SetParent(m.platformTabPageWindows)
	m.winExeCheckBox.SetShowHint(true)
	m.winExeCheckBox.SetHint(`Enable the makensis command to create an installation package program`)

	m.winDefaultInstallEdit = lcl.NewLabeledEdit(m)
	m.winDefaultInstallEdit.SetBounds(90, nextTop(30), buildFormWidth-110, 30)
	m.winDefaultInstallEdit.SetFont(m.font)
	m.winDefaultInstallEdit.SetTextHint("Default installation path of the application, e.g: C:\\Program Files")
	m.winDefaultInstallEdit.SetText(bean.GProject.BuildOption.WinDefaultInstall)
	m.winDefaultInstallEdit.SetAnchors(types.NewSet(types.AkLeft, types.AkRight, types.AkTop))
	m.winDefaultInstallEdit.EditLabel().SetCaption("默认安装路径")
	m.winDefaultInstallEdit.SetLabelPosition(types.LpLeft)
	m.winDefaultInstallEdit.SetParent(m.platformTabPageWindows)

	winConfigTitle := lcl.NewLabel(m)
	winConfigTitle.SetCaption("配置选项")
	winConfigTitle.SetLeft(10)
	winConfigTitle.SetTop(nextTop(35))
	winConfigTitle.SetFont(m.titleFontTwo)
	winConfigTitle.SetParent(m.platformTabPageWindows)

	winPackConfigBR := types.TRect{Left: 0, Top: nextTop(25)}
	winPackConfigBR.SetWidth(m.platformTabPageWindows.Width())
	winPackConfigBR.SetHeight(m.platformTabPageWindows.Height() - winPackConfigBR.Top)
	tabColor := colors.ClWhite //colors.TColor(0xF3F4F6)
	btnColor := colors.RGBToColor(234, 239, 249)
	setWinPackConfigTabPageStyle := func(page *wg.TPage) {
		page.SetHeight(m.winPackConfigTab.Height() - page.Top())
		page.SetColor(btnColor) // 设置背景色
		page.Button().SetWidth(80)
		page.Button().SetHeight(25)
		page.Button().SetLeft(0)
		page.Button().RoundedCorner = types.NewSet(wg.RcLeftTop, wg.RcRightTop)
		page.Button().Font().SetColor(colors.ClBlack)
		page.Button().SetBorderColor(wg.BbdNone, wg.LightenColor(btnColor, 0.8))
		page.Button().SetRadius(5)
		page.Button().SetColor(tabColor)
		page.Button().SetDownColor(wg.LightenColor(btnColor, 0.3), wg.LightenColor(btnColor, 0.5))
		page.Button().SetEnterColor(wg.LightenColor(btnColor, 0.1), wg.LightenColor(btnColor, 0.3))
		page.SetDefaultColor(tabColor)
		page.SetActiveColor(btnColor)
		page.Button().SetCursor(types.CrHandPoint)
	}
	m.winPackConfigTab = wg.NewTab(m)
	m.winPackConfigTab.Margin = 0
	m.winPackConfigTab.SetBoundsRect(winPackConfigBR)
	m.winPackConfigTab.SetColor(colors.ClWhite)
	m.winPackConfigTab.EnableScrollButton(false)
	m.winPackConfigTab.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
	m.winPackConfigTab.SetParent(m.platformTabPageWindows)
	m.winPackConfigTab.SetOnChange(func(sender lcl.IObject) {
		page := sender.(*wg.TPage)
		if page == m.winPackConfigTabPageBinSign {
			m.createWinSignCommandList()
		} else if page == m.winPackConfigTabPageAssociateFiles {
			m.createWinAssociateFiles()
		} else if page == m.winPackConfigTabPageAssociateProtocols {
			m.createWinAssociateProtocols()
		} else if page == m.winPackConfigTabPageAppxAssets {
			m.createWinAppxAssets()
		} else if page == m.winPackConfigTabPageNSISAssets {
			m.createWinNSISAssets()
		} else if page == m.winPackConfigTabPageNSISLicense {
			m.createWinNSISLicense()
		}
	})
	{
		m.winPackConfigTabPageBinSign = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageBinSign.SetCaption("签 名")
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageBinSign)
		m.winPackConfigTabPageBinSign.SetActive(true)
		winSignEnableRect := types.TRect{Left: 10, Top: 5}
		winSignEnableRect.SetWidth(50)
		winSignEnableRect.SetHeight(35)
		winSignEnableEnableColor := colors.RGBToColor(66, 133, 244)
		winSignEnableDisableColor := colors.RGBToColor(224, 224, 224)
		winSignEnableEnableFont := lcl.NewFont()
		winSignEnableEnableFont.SetName("微软雅黑")
		winSignEnableEnableFont.SetSize(10)
		winSignEnableEnableFont.SetStyle(types.NewSet(types.FsBold))
		winSignEnableEnableFont.SetColor(colors.RGBToColor(255, 255, 255))
		winSignEnableDisableFont := lcl.NewFont()
		winSignEnableDisableFont.SetName("微软雅黑")
		winSignEnableDisableFont.SetSize(10)
		winSignEnableDisableFont.SetStyle(types.NewSet(types.FsBold))
		winSignEnableDisableFont.SetColor(colors.RGBToColor(158, 158, 158))
		m.winSignEnable = wg.NewButton(m)
		m.winSignEnable.Font().SetColor(colors.ClWhite)
		m.winSignEnable.SetRadius(10)
		m.winSignEnable.SetCursor(types.CrHandPoint)
		m.winSignEnable.SetBoundsRect(winSignEnableRect)
		m.winSignEnable.SetColor(winSignEnableEnableColor)
		m.winSignEnable.SetDownColor(winSignEnableEnableColor, winSignEnableEnableColor)
		m.winSignEnable.SetEnterColor(winSignEnableEnableColor, winSignEnableEnableColor)
		m.winSignEnable.SetDisabledColor(winSignEnableDisableColor, winSignEnableDisableColor)
		m.winSignEnable.SetParent(m.winPackConfigTabPageBinSign)
		m.winSignEnable.SetOnClick(func(sender lcl.IObject) {
			if m.winSignEnable.Disable() {
				m.winSignEnable.SetText("已启用")
				m.winSignEnable.SetDisable(false)
				m.winSignEnable.SetFont(winSignEnableEnableFont)
			} else {
				m.winSignEnable.SetText("已禁用")
				m.winSignEnable.SetDisable(true)
				m.winSignEnable.SetFont(winSignEnableDisableFont)
			}
		})
		m.winSignEnable.SetDisable(!bean.GProject.BuildOption.WinSign.Enable)
		if m.winSignEnable.Disable() {
			m.winSignEnable.SetText("已禁用")
		} else {
			m.winSignEnable.SetText("已启用")
		}
	}
	{
		m.winPackConfigTabPageAssociateFiles = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageAssociateFiles.SetCaption("关联文件")
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageAssociateFiles)
	}
	{
		m.winPackConfigTabPageAssociateProtocols = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageAssociateProtocols.SetCaption("关联协议")
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageAssociateProtocols)
	}
	{
		m.winPackConfigTabPageAppxAssets = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageAppxAssets.SetCaption("Appx Assets")
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageAppxAssets)
	}
	{
		m.winPackConfigTabPageNSISAssets = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageNSISAssets.SetCaption("NSIS Assets")
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageNSISAssets)
	}
	{
		m.winPackConfigTabPageNSISLicense = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageNSISLicense.SetCaption("NSIS License")
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageNSISLicense)
	}
}

func (m *TBuildForm) createWinSignCommandList() {
	if m.winPackConfigTabPageBinSignMemoBox == nil {
		rect := types.TRect{Top: 40, Left: 0}
		rect.SetWidth(m.winPackConfigTabPageBinSign.Width())
		rect.SetHeight(m.winPackConfigTabPageBinSign.Height())
		m.winPackConfigTabPageBinSignMemoBox = NewCommonMemoBox(rect, "配置 Windows SDK signtool 签名命令", m.winPackConfigTabPageBinSign)
		m.winPackConfigTabPageBinSignMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageBinSignMemoBox.SetMultipleLine(false)
		m.winPackConfigTabPageBinSignMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.WinSign.Cert, "\n"))
		m.winPackConfigTabPageBinSignMemoBox.SetDemoText(`使用说明: signtool.exe 工具签名. 自动(auto=cmd)或指定证书签名(file=cmd)
备注: 证书相对路径需放在 resources/assets 目录
auto=signtool sign /a /fd SHA256 /tr http://timestamp.digicert.com /td SHA256
file=signtool sign /f cert.pfx /p 密码 /fd SHA256`)
		m.winPackConfigTabPageBinSignMemoBox.Show()
	}
}

func (m *TBuildForm) createWinAssociateFiles() {
	if m.winPackConfigTabPageAssociateFilesMemoBox == nil {
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.winPackConfigTabPageAssociateFiles.Width())
		rect.SetHeight(m.winPackConfigTabPageAssociateFiles.Height())
		m.winPackConfigTabPageAssociateFilesMemoBox = NewCommonMemoBox(rect, "配置应用关联文件", m.winPackConfigTabPageAssociateFiles)
		m.winPackConfigTabPageAssociateFilesMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageAssociateFilesMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.WinAssociateFileList, "\n"))
		m.winPackConfigTabPageAssociateFilesMemoBox.SetDemoText(`使用说明: 多个换行, 字段使用 | 分割
字段说明: EXT(后缀名) | FILECLASS(唯一类名) | DESCRIPTION(类型描述)  | ICON(图标路径) | COMMANDTEXT(右键菜单显示文字)
txt | AppTxtFile | My Project File | MyFile.ico | Open with Your App
eng | MyProductEngFile | Custom Config File | YourFile.ico | Open with energy project`)
		m.winPackConfigTabPageAssociateFilesMemoBox.Show()
	}
}

func (m *TBuildForm) createWinAssociateProtocols() {
	if m.winPackConfigTabPageAssociateProtocolsMemoBox == nil {
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.winPackConfigTabPageAssociateProtocols.Width())
		rect.SetHeight(m.winPackConfigTabPageAssociateProtocols.Height())
		m.winPackConfigTabPageAssociateProtocolsMemoBox = NewCommonMemoBox(rect, "配置应用关联协议", m.winPackConfigTabPageAssociateProtocols)
		m.winPackConfigTabPageAssociateProtocolsMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageAssociateProtocolsMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.WinAssociateProtocolList, "\n"))
		m.winPackConfigTabPageAssociateProtocolsMemoBox.SetDemoText(`使用说明: 多个换行, 字段使用 | 分割
字段说明: Scheme(协议头) | DESCRIPTION(协议描述)
myapp | Open My App
fs | fs soft scheme`)
		m.winPackConfigTabPageAssociateProtocolsMemoBox.Show()
	}
}

func (m *TBuildForm) createWinAppxAssets() {
	if m.winPackConfigTabPageAppxAssetsMemoBox == nil {
		var winAppxAssets []string
		winAppxAssets = append(winAppxAssets, "propertiesLogo="+bean.GProject.BuildOption.WinAppx.PropertiesLogo)
		winAppxAssets = append(winAppxAssets, "square44x44Logo="+bean.GProject.BuildOption.WinAppx.Square44x44Logo)
		winAppxAssets = append(winAppxAssets, "square150x150Logo="+bean.GProject.BuildOption.WinAppx.Square150x150Logo)
		winAppxAssets = append(winAppxAssets, "wide310x150Logo="+bean.GProject.BuildOption.WinAppx.Wide310x150Logo)
		winAppxAssets = append(winAppxAssets, "splashScreen="+bean.GProject.BuildOption.WinAppx.SplashScreen)
		winAppxAssets = append(winAppxAssets, "associateFileIcon="+bean.GProject.BuildOption.WinAppx.AssociateFileIcon)
		winAppxAssets = append(winAppxAssets, "associateProtocolLogo="+bean.GProject.BuildOption.WinAppx.AssociateProtocolLogo)
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.winPackConfigTabPageAppxAssets.Width())
		rect.SetHeight(m.winPackConfigTabPageAppxAssets.Height())
		m.winPackConfigTabPageAppxAssetsMemoBox = NewCommonMemoBox(rect, "配置 Appx Assets 图片资源", m.winPackConfigTabPageAppxAssets)
		m.winPackConfigTabPageAppxAssetsMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageAppxAssetsMemoBox.SetDefaultText(strings.Join(winAppxAssets, "\n"))
		m.winPackConfigTabPageAppxAssetsMemoBox.SetDemoText(`使用说明: 自定义 Assets 图片, 需放在 resources/assets 目录
propertiesLogo=PropertiesLogo.png
square44x44Logo=Square44x44Logo.png
square150x150Logo=Square150x150Logo.png
wide310x150Logo=Wide310x150Logo.png
splashScreen=SplashScreen.png
associateFileIcon=AssociateFileIcon.png
associateProtocolLogo=AssociateProtocolLogo.png`)
		m.winPackConfigTabPageAppxAssetsMemoBox.Show()
	}
}

func (m *TBuildForm) createWinNSISAssets() {
	if m.winPackConfigTabPageNSISAssetsMemoBox == nil {
		var winNsisBanner []string
		winNsisBanner = append(winNsisBanner, "welcome="+bean.GProject.BuildOption.NSIS.WelcomeBanner)
		winNsisBanner = append(winNsisBanner, "header="+bean.GProject.BuildOption.NSIS.HeaderBanner)
		winNsisBanner = append(winNsisBanner, "icon="+bean.GProject.BuildOption.NSIS.ICON)
		winNsisBanner = append(winNsisBanner, "unicon="+bean.GProject.BuildOption.NSIS.UnICON)
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.winPackConfigTabPageNSISAssets.Width())
		rect.SetHeight(m.winPackConfigTabPageNSISAssets.Height())
		m.winPackConfigTabPageNSISAssetsMemoBox = NewCommonMemoBox(rect, "配置 NSIS Banner/ICON 图片资源", m.winPackConfigTabPageNSISAssets)
		m.winPackConfigTabPageNSISAssetsMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageNSISAssetsMemoBox.SetDefaultText(strings.Join(winNsisBanner, "\n"))
		m.winPackConfigTabPageNSISAssetsMemoBox.SetDemoText(`使用说明: welcome 和 header (.png .bmp), icon(.png .ico) 需放在 resources/assets 目录
welcome=welcome.png
header=header.png
icon=nsis_icon.ico
unicon=nsis_unicon.ico`)
		m.winPackConfigTabPageNSISAssetsMemoBox.Show()
	}
}

func (m *TBuildForm) createWinNSISLicense() {
	if m.winPackConfigTabPageNSISLicenseMemoBox == nil {
		licenseName := bean.GProject.BuildOption.NSIS.License
		if licenseName == "" {
			licenseName = bean.GProject.Name + "-license.txt"
		}
		// 读取保存到 resource/xxx-license.txt 的内容
		licensePath := filepath.Join(bean.ResourcePath(), licenseName)
		licenseText := ""
		if data, err := os.ReadFile(licensePath); err == nil {
			licenseText = string(data)
		}
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.winPackConfigTabPageNSISLicense.Width())
		rect.SetHeight(m.winPackConfigTabPageNSISLicense.Height())
		m.winPackConfigTabPageNSISLicenseMemoBox = NewCommonMemoBox(rect, "配置 NSIS 许可证内容", m.winPackConfigTabPageNSISLicense)
		m.winPackConfigTabPageNSISLicenseMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageNSISLicenseMemoBox.SetDefaultText(licenseText)
		m.winPackConfigTabPageNSISLicenseMemoBox.Show()
	}
	// 文本保存到临时文件
	//_ = os.Remove(licensePath)
	//if len(lines) > 0 {
	//	data := strings.Join(lines, "\n")
	//	utf8Bom := []byte{0xEF, 0xBB, 0xBF}
	//	licenseData := append(utf8Bom, data...)
	//	_ = os.WriteFile(licensePath, licenseData, 0644)
	//	m.license = licensePath
	//} else {
	//	m.license = ""
	//}
}

func (m *TBuildForm) TemplateVariablesClick(sender lcl.IObject) {
}
