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
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/rtl/version"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
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
	macOSPackageFmtTitle.SetCaption("打包格式")
	macOSPackageFmtTitle.SetLeft(10)
	macOSPackageFmtTitle.SetTop(nextTop(5))
	macOSPackageFmtTitle.SetFont(m.titleFontTwo)
	macOSPackageFmtTitle.SetParent(m.platformTabPageMacOS)

	m.macDMGCheckBox = lcl.NewCheckBox(m)
	m.macDMGCheckBox.SetCaption("DMG 镜像")
	m.macDMGCheckBox.SetLeft(20)
	m.macDMGCheckBox.SetTop(nextTop(25))
	m.macDMGCheckBox.SetFont(m.font)
	m.macDMGCheckBox.SetChecked(bean.GProject.BuildOption.MacDMG)
	m.macDMGCheckBox.SetParent(m.platformTabPageMacOS)

	m.macPKGCheckBox = lcl.NewCheckBox(m)
	m.macPKGCheckBox.SetCaption("PKG 安装包")
	m.macPKGCheckBox.SetLeft(210)
	m.macPKGCheckBox.SetTop(m.macDMGCheckBox.Top())
	m.macPKGCheckBox.SetFont(m.font)
	m.macPKGCheckBox.SetChecked(bean.GProject.BuildOption.MacPKG)
	m.macPKGCheckBox.SetParent(m.platformTabPageMacOS)

	m.macCommonLibCheckBox = lcl.NewCheckBox(m)
	m.macCommonLibCheckBox.SetCaption("‌通用二进制(Universal Binary)")
	m.macCommonLibCheckBox.SetLeft(20)
	m.macCommonLibCheckBox.SetTop(nextTop(30))
	m.macCommonLibCheckBox.SetFont(m.font)
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
		bean.GProject.BuildOption.MacCommonLib = false
		m.macCommonLibCheckBox.SetEnabled(false)
	}
	m.macCommonLibCheckBox.SetChecked(bean.GProject.BuildOption.MacCommonLib)
	m.macCommonLibCheckBox.SetParent(m.platformTabPageMacOS)

	macConfigTitle := lcl.NewLabel(m)
	macConfigTitle.SetCaption("配置选项")
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
	m.macPackConfigTab.EnableScrollButton(false)
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
		}
	})

	{
		m.macPackConfigTabPageBinSign = m.macPackConfigTab.NewPage()
		m.macPackConfigTabPageBinSign.SetCaption("签 名")
		setMacPackConfigTabPageStyle(m.macPackConfigTabPageBinSign)
		m.macPackConfigTabPageBinSign.SetActive(true)
		m.macSignEnable = NewEnableButton(m.macPackConfigTabPageBinSign)
		m.macSignEnable.DisableText = "已启用"
		m.macSignEnable.EnableText = "已禁用"
		m.macSignEnable.SetDisable(!bean.GProject.BuildOption.MacSign.Enable)
	}
	{
		m.macPackConfigTabPageAssociateFiles = m.macPackConfigTab.NewPage()
		m.macPackConfigTabPageAssociateFiles.SetCaption("关联文件")
		setMacPackConfigTabPageStyle(m.macPackConfigTabPageAssociateFiles)
	}
	{
		m.macPackConfigTabPageAssociateProtocols = m.macPackConfigTab.NewPage()
		m.macPackConfigTabPageAssociateProtocols.SetCaption("关联协议")
		setMacPackConfigTabPageStyle(m.macPackConfigTabPageAssociateProtocols)
	}
}

func (m *TBuildForm) createMacAssociateProtocols() {
	if m.macPackConfigTabPageAssociateProtocolsBox == nil {
		rect := types.TRect{Top: 0, Left: 0}
		rect.SetWidth(m.macPackConfigTabPageAssociateProtocols.Width())
		rect.SetHeight(m.macPackConfigTabPageAssociateProtocols.Height())
		m.macPackConfigTabPageAssociateProtocolsBox = NewCommonMemoBox(rect, "配置应用关联协议", m.macPackConfigTabPageAssociateProtocols)
		m.macPackConfigTabPageAssociateProtocolsBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageAssociateProtocolsBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacAssociateProtocolList, "\n"))
		m.macPackConfigTabPageAssociateProtocolsBox.SetDemoText(`使用说明: 多个换行, 字段使用 | 分割
字段说明: 协议名称 | 协议描述
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
		m.macPackConfigTabPageAssociateFilesBox = NewCommonMemoBox(rect, "配置应用关联文件", m.macPackConfigTabPageAssociateFiles)
		m.macPackConfigTabPageAssociateFilesBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageAssociateFilesBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacAssociateFileList, "\n"))
		m.macPackConfigTabPageAssociateFilesBox.SetDemoText(`使用说明: 多个换行, 字段使用 | 分割
字段说明: 扩展名 | 名称 | 角色(Editor/Viewer) | 优先级(Owner/Default) | 图标(png/.icns)｜ MIME(允许空)
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
		m.macPackConfigTabPageBinSignMemoBox = NewCommonMemoBox(rect, "配置 MacOS codesign 签名命令列表", m.macPackConfigTabPageBinSign)
		m.macPackConfigTabPageBinSignMemoBox.SetAnchors(types.NewSet(types.AkTop, types.AkBottom, types.AkLeft, types.AkRight))
		m.macPackConfigTabPageBinSignMemoBox.SetDefaultText(strings.Join(bean.GProject.BuildOption.MacSign.Cert, "\n"))
		m.macPackConfigTabPageBinSignMemoBox.SetDemoText(`签名文件命令列表, 多个换行. 按深度顺序添加
codesign -f -s "Developer ID Application: 你的名字 (团队ID)" "$APP_NAME/Contents/Frameworks/your.dylib"
codesign -f -s "Developer ID Application: 你的名字 (团队ID)" --options runtime "$APP_NAME"`)
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
