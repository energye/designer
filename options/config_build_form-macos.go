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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/rtl/version"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"path/filepath"
	"strings"
)

func (m *TBuildForm) initMacOSOptions() {
	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	macOSPackageFmtTitle := lcl.NewLabel(m)
	macOSPackageFmtTitle.SetName("BuildFormMacOSPackageFmtTitle")
	macOSPackageFmtTitle.SetCaption(metadata.GI18n.Dict("BuildFormMacOSPackageFmtTitle.Caption"))
	macOSPackageFmtTitle.SetLeft(10)
	macOSPackageFmtTitle.SetTop(nextTop(5))
	macOSPackageFmtTitle.SetFont(m.titleFontTwo)
	macOSPackageFmtTitle.SetParent(m.platformTabPageMacOS)

	m.macDMGCheckBox = lcl.NewCheckBox(m)
	m.macDMGCheckBox.SetName("BuildFormMacOSDMGCheckBox")
	m.macDMGCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormMacOSDMGCheckBox.Caption"))
	m.macDMGCheckBox.SetLeft(20)
	m.macDMGCheckBox.SetTop(nextTop(25))
	m.macDMGCheckBox.SetFont(m.font)
	m.macDMGCheckBox.SetChecked(bean.GProject.BuildOption.MacDMG)
	m.macDMGCheckBox.SetParent(m.platformTabPageMacOS)

	m.macPKGCheckBox = lcl.NewCheckBox(m)
	m.macPKGCheckBox.SetName("BuildFormMacOSPKGCheckBox")
	m.macPKGCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormMacOSPKGCheckBox.Caption"))
	m.macPKGCheckBox.SetLeft(210)
	m.macPKGCheckBox.SetTop(m.macDMGCheckBox.Top())
	m.macPKGCheckBox.SetFont(m.font)
	m.macPKGCheckBox.SetChecked(bean.GProject.BuildOption.MacPKG)
	m.macPKGCheckBox.SetParent(m.platformTabPageMacOS)

	m.macCommonLibCheckBox = lcl.NewCheckBox(m)
	m.macCommonLibCheckBox.SetName("BuildFormMacOSCommonLibCheckBox")
	m.macCommonLibCheckBox.SetCaption(metadata.GI18n.Dict("BuildFormMacOSCommonLibCheckBox.Caption"))
	m.macCommonLibCheckBox.SetLeft(20)
	m.macCommonLibCheckBox.SetTop(nextTop(30))
	m.macCommonLibCheckBox.SetFont(m.font)
	m.macCommonLibCheckBox.SetChecked(bean.GProject.BuildOption.MacCommonLib)
	m.macCommonLibCheckBox.SetParent(m.platformTabPageMacOS)

	macConfigTitle := lcl.NewLabel(m)
	macConfigTitle.SetName("BuildFormMacOSConfigTitle")
	macConfigTitle.SetCaption(metadata.GI18n.Dict("BuildFormMacOSConfigTitle.Caption"))
	macConfigTitle.SetLeft(10)
	macConfigTitle.SetTop(nextTop(35))
	macConfigTitle.SetFont(m.titleFontTwo)
	macConfigTitle.SetParent(m.platformTabPageMacOS)

	tabColor := colors.ClWhite //colors.TColor(0xF3F4F6)
	btnColor := colors.RGBToColor(234, 239, 249)
	setMacPackConfigTabPageStyle := func(page *wg.TPage) {
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
	macPackConfigBR := types.TRect{Left: 0, Top: nextTop(25)}
	macPackConfigBR.SetWidth(m.platformTabPageMacOS.Width())
	macPackConfigBR.SetHeight(m.platformTabPageMacOS.Height() - macPackConfigBR.Top)
	m.macPackConfigTab = wg.NewTab(m)
	m.macPackConfigTab.Margin = 0
	m.macPackConfigTab.SetBoundsRect(macPackConfigBR)
	m.macPackConfigTab.SetColor(colors.ClWhite)
	m.macPackConfigTab.EnableScrollButton(true)
	m.macPackConfigTab.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
	m.macPackConfigTab.SetParent(m.platformTabPageMacOS)
	m.macPackConfigTab.SetOnChange(func(sender lcl.IObject) {
		page := sender.(*wg.TPage)
		if page == m.macPackConfigTabPageBinSign {
			m.createMacSignCommandList()
		} else if page == m.macPackConfigTabPageAssociateFiles {
			m.createMacAssociateFiles()
		} else if page == m.macPackConfigTabPageAssociateProtocols {
			m.createMacAssociateProtocols()
		} else if page == m.macPackConfigTabPageUniversalLink {
			m.createMacAssociateUniversalLink()
		}
		m.macPackConfigTab.RecalculatePosition()
	})

	{
		m.macPackConfigTabPageBinSign = m.macPackConfigTab.NewPage()
		m.macPackConfigTabPageBinSign.SetName("BuildFormMacOSPackConfigTabPageBinSign")
		m.macPackConfigTabPageBinSign.Button().SetAutoSize(true)
		m.macPackConfigTabPageBinSign.SetCaption(metadata.GI18n.Dict("BuildFormMacOSPackConfigTabPageBinSign.Caption"))
		setMacPackConfigTabPageStyle(m.macPackConfigTabPageBinSign)
		m.macPackConfigTabPageBinSign.SetActive(true)
		m.macSignEnable = NewEnableButton(m.macPackConfigTabPageBinSign)
		m.macSignEnable.DisableText = metadata.GI18n.Dict("BuildFormMacOSSignEnable.DisableText")
		m.macSignEnable.EnableText = metadata.GI18n.Dict("BuildFormMacOSSignEnable.EnableText")
		m.macSignEnable.SetDisable(!bean.GProject.BuildOption.MacSign.Enable)
	}
	{
		m.macPackConfigTabPageAssociateFiles = m.macPackConfigTab.NewPage()
		m.macPackConfigTabPageAssociateFiles.SetName("BuildFormMacOSAssociateFiles")
		m.macPackConfigTabPageAssociateFiles.Button().SetAutoSize(true)
		m.macPackConfigTabPageAssociateFiles.SetCaption(metadata.GI18n.Dict("BuildFormMacOSAssociateFiles.Caption"))
		setMacPackConfigTabPageStyle(m.macPackConfigTabPageAssociateFiles)
	}
	{
		m.macPackConfigTabPageAssociateProtocols = m.macPackConfigTab.NewPage()
		m.macPackConfigTabPageAssociateProtocols.SetName("BuildFormMacOSAssociateProtocols")
		m.macPackConfigTabPageAssociateProtocols.Button().SetAutoSize(true)
		m.macPackConfigTabPageAssociateProtocols.SetCaption(metadata.GI18n.Dict("BuildFormMacOSAssociateProtocols.Caption"))
		setMacPackConfigTabPageStyle(m.macPackConfigTabPageAssociateProtocols)
	}
	{
		// 后续开发
		//m.macPackConfigTabPageUniversalLink = m.macPackConfigTab.NewPage()
		//m.macPackConfigTabPageUniversalLink.SetCaption("通用链接")
		//setMacPackConfigTabPageStyle(m.macPackConfigTabPageUniversalLink)
	}
}

func (m *TBuildForm) createMacAssociateUniversalLink() {
	if m.macPackConfigTabPageUniversalLinkBox == nil {
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.macPackConfigTabPageAssociateProtocols.Width())
		rect.SetHeight(m.macPackConfigTabPageAssociateProtocols.Height())
		m.macPackConfigTabPageUniversalLinkBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormMacOSAssociateUniversalLink.Title"), m.macPackConfigTabPageUniversalLink)
		m.macPackConfigTabPageUniversalLinkBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageUniversalLinkBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacUniversalLink, "\n"))
		m.macPackConfigTabPageUniversalLinkBox.SetDemoText(`Use Example: Multiple lines, Separate fields with |
Field Desc: To Be Written | To Be Written
`)
		m.macPackConfigTabPageUniversalLinkBox.Show()
	}
}

func (m *TBuildForm) createMacAssociateProtocols() {
	if m.macPackConfigTabPageAssociateProtocolsBox == nil {
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.macPackConfigTabPageAssociateProtocols.Width())
		rect.SetHeight(m.macPackConfigTabPageAssociateProtocols.Height())
		m.macPackConfigTabPageAssociateProtocolsBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormMacOSAssociateProtocols.Title"), m.macPackConfigTabPageAssociateProtocols)
		m.macPackConfigTabPageAssociateProtocolsBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageAssociateProtocolsBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacAssociateProtocolList, "\n"))
		m.macPackConfigTabPageAssociateProtocolsBox.SetDemoText(`Use Example: Multiple lines, Separate fields with |
Field Desc: Protocol Name | Protocol Description
myapp | Open My App
fs | fs soft scheme`)
		m.macPackConfigTabPageAssociateProtocolsBox.Show()
	}
}

func (m *TBuildForm) createMacAssociateFiles() {
	if m.macPackConfigTabPageAssociateFilesBox == nil {
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.macPackConfigTabPageAssociateFiles.Width())
		rect.SetHeight(m.macPackConfigTabPageAssociateFiles.Height())
		m.macPackConfigTabPageAssociateFilesBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormMacOSAssociatFiles.Title"), m.macPackConfigTabPageAssociateFiles)
		m.macPackConfigTabPageAssociateFilesBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageAssociateFilesBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacAssociateFileList, "\n"))
		m.macPackConfigTabPageAssociateFilesBox.SetDemoText(`Use Example: Multiple lines, Separate fields with |
Field Description: Extension | Name | Role(Editor/Viewer) | Priority(Owner/Default) | Icon(png/.icns) | MIME(Optional)
txt | AppTxtFile | Editor | Owner | MyIcon.icns | application/x-gproj
eng | MyProductEngFile | Editor | Owner | MyIcon.png | application/x-gproj`)
		m.macPackConfigTabPageAssociateFilesBox.Show()
	}
}

func (m *TBuildForm) createMacSignCommandList() {
	if m.macPackConfigTabPageBinSignMemoBox == nil {
		rect := types.TRect{Top: 40, Left: 0}
		rect.SetWidth(m.macPackConfigTabPageBinSign.Width())
		rect.SetHeight(m.macPackConfigTabPageBinSign.Height())
		m.macPackConfigTabPageBinSignMemoBox = NewCommonMemoBox(rect, metadata.GI18n.Dict("BuildFormMacOSSignCommandList.Title"), m.macPackConfigTabPageBinSign)
		m.macPackConfigTabPageBinSignMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageBinSignMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacSign.Cert, "\n"))
		m.macPackConfigTabPageBinSignMemoBox.SetDemoText(`Signature Commands | Multiple lines | Add by depth order
codesign -f -s "Developer ID Application: Name (Team ID)" "$APP_NAME/Contents/Frameworks/your.dylib"
codesign -f -s "Developer ID Application: Name (Team ID)" --options runtime "$APP_NAME"`)
		m.macPackConfigTabPageBinSignMemoBox.Show()
	}
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
		outputLibPath := config.Config.FrameworkRuntimePath()
		tempArm64LibName := libArm64.OutputFilename
		tempAmd64LibName := libAmd64.OutputFilename
		arm64LibFilePath := filepath.Join(outputLibPath, tempArm64LibName)
		amd64LibFilePath := filepath.Join(outputLibPath, tempAmd64LibName)
		universalLibFilePath := filepath.Join(outputLibPath, libname.DarwinUniversalBinaryName)
		event.ConsoleWriteInfo("Merge macOS UniversalBinary, arm64LibFilePath:", arm64LibFilePath)
		event.ConsoleWriteInfo("Merge macOS UniversalBinary, amd64LibFilePath:", amd64LibFilePath)
		_ = os.Remove(universalLibFilePath)
		cmd := command.NewCMD()
		cmd.HideWindow = true
		cmd.Command("lipo", "-create", amd64LibFilePath, arm64LibFilePath, "-output", universalLibFilePath)
		event.ConsoleWriteInfo("Merge macOS UniversalBinary, universalLibFilePath:", universalLibFilePath)
	}
	return nil
}
