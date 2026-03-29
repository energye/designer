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
	m.winMsiCheckBox.SetEnabled(false) // TODO 先禁用, 还未实现
	m.winMsiCheckBox.SetParent(m.platformTabPageWindows)
	m.winExeCheckBox = lcl.NewCheckBox(m)
	m.winExeCheckBox.SetCaption("EXE 安装包(MakeNsis)")
	m.winExeCheckBox.SetLeft(210)
	m.winExeCheckBox.SetTop(m.winMsiCheckBox.Top())
	m.winExeCheckBox.SetFont(m.font)
	m.winExeCheckBox.SetChecked(bean.GProject.BuildOption.WinExe)
	m.winExeCheckBox.SetParent(m.platformTabPageWindows)

	winDefaultInstallTitle := lcl.NewLabel(m)
	winDefaultInstallTitle.SetCaption("默认安装路径")
	winDefaultInstallTitle.SetLeft(10)
	winDefaultInstallTitle.SetTop(nextTop(30))
	winDefaultInstallTitle.SetFont(m.titleFontTwo)
	winDefaultInstallTitle.SetParent(m.platformTabPageWindows)

	m.winDefaultInstallEdit = lcl.NewEdit(m)
	m.winDefaultInstallEdit.SetBounds(20, nextTop(25), 515, 30)
	m.winDefaultInstallEdit.SetFont(m.font)
	m.winDefaultInstallEdit.SetTextHint("Windows 应用的默认安装路径 如: C:\\Program Files")
	m.winDefaultInstallEdit.SetText(bean.GProject.BuildOption.WinDefaultInstall)
	m.winDefaultInstallEdit.SetParent(m.platformTabPageWindows)

	m.winDesktopShortcutCheckBox = lcl.NewCheckBox(m)
	m.winDesktopShortcutCheckBox.SetCaption("创建桌面快捷方式")
	m.winDesktopShortcutCheckBox.SetLeft(20)
	m.winDesktopShortcutCheckBox.SetTop(nextTop(40))
	m.winDesktopShortcutCheckBox.SetFont(m.font)
	m.winDesktopShortcutCheckBox.SetChecked(bean.GProject.BuildOption.WinDesktopShortcut)
	m.winDesktopShortcutCheckBox.SetParent(m.platformTabPageWindows)

	m.winAddStartMenuCheckBox = lcl.NewCheckBox(m)
	m.winAddStartMenuCheckBox.SetCaption("添加到开始菜单")
	m.winAddStartMenuCheckBox.SetLeft(210)
	m.winAddStartMenuCheckBox.SetTop(m.winDesktopShortcutCheckBox.Top())
	m.winAddStartMenuCheckBox.SetFont(m.font)
	m.winAddStartMenuCheckBox.SetChecked(bean.GProject.BuildOption.WinAddStartMenu)
	m.winAddStartMenuCheckBox.SetParent(m.platformTabPageWindows)

	winAssociateFilesRect := types.TRect{Left: 20, Top: nextTop(30)}
	winAssociateFilesRect.SetWidth(90)
	winAssociateFilesRect.SetHeight(25)
	m.winAssociateFilesBtn = wg.NewButton(m)
	m.winAssociateFilesBtn.SetBoundsRect(winAssociateFilesRect)
	m.winAssociateFilesBtn.SetText("设置关联文件")
	m.winAssociateFilesBtn.Font().SetColor(colors.ClWhite)
	m.winAssociateFilesBtn.SetRadius(3)
	m.winAssociateFilesBtn.SetCursor(types.CrHandPoint)
	m.winAssociateFilesBtn.SetColor(colors.RGBToColor(59, 130, 246))
	m.winAssociateFilesBtn.SetParent(m.platformTabPageWindows)
	m.winAssociateFilesBtn.SetOnClick(m.AssociateFilesClick)
	m.winAssociateFileArray = bean.GProject.BuildOption.WinAssociateFileList

	winAssociateProtocolsRect := types.TRect{Left: winAssociateFilesRect.Left + winAssociateFilesRect.Width() + 20,
		Top: winAssociateFilesRect.Top}
	winAssociateProtocolsRect.SetWidth(90)
	winAssociateProtocolsRect.SetHeight(25)
	m.winAssociateProtocolsBtn = wg.NewButton(m)
	m.winAssociateProtocolsBtn.SetBoundsRect(winAssociateProtocolsRect)
	m.winAssociateProtocolsBtn.SetText("设置关联协议")
	m.winAssociateProtocolsBtn.Font().SetColor(colors.ClWhite)
	m.winAssociateProtocolsBtn.SetRadius(3)
	m.winAssociateProtocolsBtn.SetCursor(types.CrHandPoint)
	m.winAssociateProtocolsBtn.SetColor(colors.RGBToColor(59, 130, 246))
	m.winAssociateProtocolsBtn.SetParent(m.platformTabPageWindows)
	m.winAssociateProtocolsBtn.SetOnClick(m.AssociateProtocolsClick)
	m.winAssociateProtocolArray = bean.GProject.BuildOption.WinAssociateProtocolList

	bannerRect := types.TRect{Left: winAssociateProtocolsRect.Left + winAssociateProtocolsRect.Width() + 20,
		Top: winAssociateFilesRect.Top}
	bannerRect.SetWidth(90)
	bannerRect.SetHeight(25)
	m.bannerBtn = wg.NewButton(m)
	m.bannerBtn.SetBoundsRect(bannerRect)
	m.bannerBtn.SetText("设置 Banner")
	m.bannerBtn.Font().SetColor(colors.ClWhite)
	m.bannerBtn.SetRadius(3)
	m.bannerBtn.SetCursor(types.CrHandPoint)
	m.bannerBtn.SetColor(colors.RGBToColor(59, 130, 246))
	m.bannerBtn.SetParent(m.platformTabPageWindows)
	m.bannerBtn.SetOnClick(m.BannerClick)

	licenseRect := types.TRect{Left: bannerRect.Left + bannerRect.Width() + 20,
		Top: bannerRect.Top}
	licenseRect.SetWidth(90)
	licenseRect.SetHeight(25)
	m.licenseBtn = wg.NewButton(m)
	m.licenseBtn.SetBoundsRect(licenseRect)
	m.licenseBtn.SetText("设置许可证")
	m.licenseBtn.Font().SetColor(colors.ClWhite)
	m.licenseBtn.SetRadius(3)
	m.licenseBtn.SetCursor(types.CrHandPoint)
	m.licenseBtn.SetColor(colors.RGBToColor(59, 130, 246))
	m.licenseBtn.SetParent(m.platformTabPageWindows)
	m.licenseBtn.SetOnClick(m.LicenseClick)

	// 签名 signtool
	// 授权协议
	// banner
	//!define MUI_WELCOMEFINISHPAGE_BITMAP "welcome.bmp"
	//!define MUI_HEADERIMAGE
	//!define MUI_HEADERIMAGE_BITMAP "header.bmp"
}

func (m *TBuildForm) AssociateFilesClick(sender lcl.IObject) {
	newForm := NewCommonMemoForm(600, 200, `配置应用关联文件`, m)
	newForm.SetDefaultText(strings.Join(m.winAssociateFileArray, "\n"))
	newForm.SetDemoText(`多个换行, 每行使用 | 分割
说明: EXT(后缀名) | FILECLASS(唯一类名) | DESCRIPTION(类型描述)  | ICON(图标路径) | COMMANDTEXT(右键菜单显示文字)
txt | AppTxtFile | My Project File | MyFile.ico | Open with Your App
eng | MyProductEngFile | Custom Config File | YourFile.ico | Open with energy project`)
	newForm.SetOnOK(func(lines []string) {
		m.winAssociateFileArray = lines
	})
	newForm.ShowModal()
}

func (m *TBuildForm) AssociateProtocolsClick(sender lcl.IObject) {
	newForm := NewCommonMemoForm(600, 200, `配置应用关联协议`, m)
	newForm.SetDefaultText(strings.Join(m.winAssociateProtocolArray, "\n"))
	newForm.SetDemoText(`多个换行, 每行使用 | 分割
说明: Scheme(协议头) | DESCRIPTION(协议描述)
myapp | Open My App
fs | fs soft scheme`)
	newForm.SetOnOK(func(lines []string) {
		m.winAssociateProtocolArray = lines
	})
	newForm.ShowModal()
}

func (m *TBuildForm) BannerClick(sender lcl.IObject) {
	// 选择 png 转为 bmp
}

func (m *TBuildForm) LicenseClick(sender lcl.IObject) {
	// 文本保存到临时文件
	licensePath := filepath.Join(bean.ResourcePath(), bean.GProject.Name+"-license.txt")
	licenseText := ""
	if data, err := os.ReadFile(licensePath); err == nil {
		licenseText = string(data)
	}
	newForm := NewCommonMemoForm(600, 400, `设置许可证内容`, m)
	newForm.SetDefaultText(licenseText)
	newForm.SetOnOK(func(lines []string) {
		_ = os.Remove(licensePath)
		if len(lines) > 0 {
			data := strings.Join(lines, "\n")
			_ = os.WriteFile(licensePath, []byte(data), 0644)
			m.license = licensePath
		} else {
			m.license = ""
		}
	})
	newForm.ShowModal()
}
