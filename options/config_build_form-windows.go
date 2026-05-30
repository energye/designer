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
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
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
	windowsPackageFmtTitle.SetName("BuildFormWindowsPackageFmtTitle")
	windowsPackageFmtTitle.SetCaption(metadata.GI18n.Dict("BuildFormWindowsPackageFmtTitle.Caption"))
	windowsPackageFmtTitle.SetLeft(10)
	windowsPackageFmtTitle.SetTop(nextTop(5))
	windowsPackageFmtTitle.SetFont(m.titleFontTwo)
	windowsPackageFmtTitle.SetParent(m.platformTabPageWindows)

	m.winMsiCheckBox = lcl.NewCheckBox(m)
	m.winMsiCheckBox.SetName("BuildFormWinMsiCheckBox")
	m.winMsiCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormWinMsiCheckBox.Caption"))
	m.winMsiCheckBox.SetLeft(20)
	m.winMsiCheckBox.SetTop(nextTop(25))
	m.winMsiCheckBox.SetFont(m.font)
	m.winMsiCheckBox.SetChecked(bean.GProject.BuildOption.WinMsi)
	m.winMsiCheckBox.SetParent(m.platformTabPageWindows)

	m.winExeCheckBox = lcl.NewCheckBox(m)
	m.winExeCheckBox.SetName("BuildFormWinExeCheckBox")
	m.winExeCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormWinExeCheckBox.Caption"))
	m.winExeCheckBox.SetLeft(210)
	m.winExeCheckBox.SetTop(m.winMsiCheckBox.Top())
	m.winExeCheckBox.SetFont(m.font)
	m.winExeCheckBox.SetChecked(bean.GProject.BuildOption.WinExe)
	m.winExeCheckBox.SetParent(m.platformTabPageWindows)
	m.winExeCheckBox.SetShowHint(true)
	m.winExeCheckBox.SetHint(metadata.GI18n.Dict("BuildFormWinExeCheckBox.Hint"))

	m.winDefaultInstallEdit = lcl.NewLabeledEdit(m)
	m.winDefaultInstallEdit.SetName("BuildFormWinDefaultInstallEdit")
	m.winDefaultInstallEdit.SetBounds(90, nextTop(30), buildFormWidth-110, 30)
	m.winDefaultInstallEdit.SetFont(m.font)
	m.winDefaultInstallEdit.SetTextHint(metadata.GI18n.Dict("BuildFormWinDefaultInstallEdit.TextHint"))
	m.winDefaultInstallEdit.SetText(bean.GProject.BuildOption.WinDefaultInstall)
	m.winDefaultInstallEdit.SetAnchors(types.NewSet(types.AkLeft, types.AkRight, types.AkTop))
	m.winDefaultInstallEdit.EditLabel().SetCaption(metadata.GI18n.Dict("BuildFormWinDefaultInstallEdit.EditLabel.Caption"))
	m.winDefaultInstallEdit.SetLabelPosition(types.LpLeft)
	m.winDefaultInstallEdit.SetParent(m.platformTabPageWindows)

	winConfigTitle := lcl.NewLabel(m)
	winConfigTitle.SetName("BuildFormWinConfigTitle")
	winConfigTitle.SetCaption(metadata.GI18n.Dict("BuildFormWinConfigTitle.Caption"))
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
		m.winPackConfigTabPageBinSign.SetName("BuildFormWinPackConfigTabPageBinSign")
		m.winPackConfigTabPageBinSign.SetCaption(metadata.GI18n.Dict("BuildFormWinPackConfigTabPageBinSign.Caption"))
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageBinSign)
		m.winPackConfigTabPageBinSign.SetActive(true)
		m.winSignEnable = NewEnableButton(m.winPackConfigTabPageBinSign)
		m.winSignEnable.DisableText = metadata.GI18n.Dict("BuildFormWinSignEnable.DisableText")
		m.winSignEnable.EnableText = metadata.GI18n.Dict("BuildFormWinSignEnable.EnableText")
		m.winSignEnable.SetDisable(!bean.GProject.BuildOption.WinSign.Enable)
	}
	{
		m.winPackConfigTabPageAssociateFiles = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageAssociateFiles.SetName("BuildFormWinPackConfigTabPageAssociateFiles")
		m.winPackConfigTabPageAssociateFiles.SetCaption(metadata.GI18n.Dict("BuildFormWinPackConfigTabPageAssociateFiles.Caption"))
		setWinPackConfigTabPageStyle(m.winPackConfigTabPageAssociateFiles)
	}
	{
		m.winPackConfigTabPageAssociateProtocols = m.winPackConfigTab.NewPage()
		m.winPackConfigTabPageAssociateProtocols.SetName("BuildFormWinPackConfigTabPageAssociateProtocols")
		m.winPackConfigTabPageAssociateProtocols.SetCaption(metadata.GI18n.Dict("BuildFormWinPackConfigTabPageAssociateProtocols.Caption"))
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
		m.winPackConfigTabPageBinSignMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormWinSignCommandList.Title"), m.winPackConfigTabPageBinSign)
		m.winPackConfigTabPageBinSignMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageBinSignMemoBox.SetMultipleLine(false)
		m.winPackConfigTabPageBinSignMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.WinSign.Cert, "\n"))
		m.winPackConfigTabPageBinSignMemoBox.SetDemoText(`Use Example: signtool.exe tool sign. auto(auto=cmd) or specified cert(file=cmd)
remark: cert relative path resources/assets
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
		m.winPackConfigTabPageAssociateFilesMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormWinAssociateFiles.Title"), m.winPackConfigTabPageAssociateFiles)
		m.winPackConfigTabPageAssociateFilesMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageAssociateFilesMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.WinAssociateFileList, "\n"))
		m.winPackConfigTabPageAssociateFilesMemoBox.SetDemoText(`Use Example: Multiple lines, Separate fields with |
Field Desc: Extension | Unique Class Name | Type Description | Icon | Right-Click Menu Text
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
		m.winPackConfigTabPageAssociateProtocolsMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormWinAssociateProtocols.Title"), m.winPackConfigTabPageAssociateProtocols)
		m.winPackConfigTabPageAssociateProtocolsMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageAssociateProtocolsMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.WinAssociateProtocolList, "\n"))
		m.winPackConfigTabPageAssociateProtocolsMemoBox.SetDemoText(`Use Example: Multiple lines, Separate fields with |
Field Desc: Protocol Name | Protocol Description
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
		m.winPackConfigTabPageAppxAssetsMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormWinAppxAssets.Title"), m.winPackConfigTabPageAppxAssets)
		m.winPackConfigTabPageAppxAssetsMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageAppxAssetsMemoBox.SetDefaultText(strings.Join(winAppxAssets, "\n"))
		m.winPackConfigTabPageAppxAssetsMemoBox.SetDemoText(`Instructions: Custom Assets images | resources/assets directory
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
		m.winPackConfigTabPageNSISAssetsMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormWinNSISAssets.Title"), m.winPackConfigTabPageNSISAssets)
		m.winPackConfigTabPageNSISAssetsMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.winPackConfigTabPageNSISAssetsMemoBox.SetDefaultText(strings.Join(winNsisBanner, "\n"))
		m.winPackConfigTabPageNSISAssetsMemoBox.SetDemoText(`Use Example: welcome and header (.png .bmp), icon(.png .ico) in resources/assets directory
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
		m.winPackConfigTabPageNSISLicenseMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormWinNSISLicense.Title"), m.winPackConfigTabPageNSISLicense)
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
