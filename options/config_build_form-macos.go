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
	m.macDMGCheckBox.SetTop(nextTop(30))
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

	m.certCheckBox = lcl.NewCheckBox(m)
	m.certCheckBox.SetCaption("签名")
	m.certCheckBox.SetLeft(20)
	m.certCheckBox.SetTop(nextTop(30))
	m.certCheckBox.SetFont(m.font)
	m.certCheckBox.SetChecked(bean.GProject.BuildOption.Cert)
	m.certCheckBox.SetParent(m.platformTabPageMacOS)
	m.certCheckBox.SetOnChange(func(sender lcl.IObject) {
		m.certListBtn.SetVisible(m.certCheckBox.Checked())
	})

	// 文件签名配置按钮
	m.certListBtn = wg.NewButton(m)
	m.certListBtn.SetVisible(m.certCheckBox.Checked())
	m.certListBtn.SetText("二进制签名")
	m.certListBtn.Font().SetColor(colors.ClWhite)
	m.certListBtn.SetRadius(0)
	certListBtnRect := types.TRect{Left: 75, Top: m.certCheckBox.Top()}
	certListBtnRect.SetWidth(120)
	certListBtnRect.SetHeight(20)
	m.certListBtn.SetBoundsRect(certListBtnRect)
	m.certListBtn.SetColor(colors.RGBToColor(59, 130, 246))
	m.certListBtn.SetParent(m.platformTabPageMacOS)
	m.certListBtn.SetOnClick(m.macCertCommandList)

	{
		// macOS 文件签名配置
		m.macCertArray = bean.GProject.BuildOption.MacCertList
		m.macCommonLibCheckBox = lcl.NewCheckBox(m)
		m.macCommonLibCheckBox.SetCaption("‌通用二进制文件(Universal Binary)")
		m.macCommonLibCheckBox.SetLeft(210)
		m.macCommonLibCheckBox.SetTop(m.certCheckBox.Top())
		m.macCommonLibCheckBox.SetFont(m.font)
		if version.OSVersion.Major <= 10 {
			// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
			bean.GProject.BuildOption.MacCommonLib = false
			m.macCommonLibCheckBox.SetEnabled(false)
		}
		m.macCommonLibCheckBox.SetChecked(bean.GProject.BuildOption.MacCommonLib)
		m.macCommonLibCheckBox.SetParent(m.platformTabPageMacOS)
	}
}

func (m *TBuildForm) macCertCommandList(sender lcl.IObject) {
	newCertCommandListForm := lcl.NewEngForm(nil)
	newCertCommandListForm.SetBorderStyleToFormBorderStyle(types.BsNone)
	newCertCommandListForm.SetBounds(0, 0, 250, 200)
	newCertCommandListForm.WorkAreaCenter()
	memo := lcl.NewMemo(newCertCommandListForm)
	//memo.SetBounds(25, 0, 250, 145)
	memo.SetBounds(0, 0, 250, 170)
	memo.SetShowHint(true)
	memo.SetHint(`签名文件命令列表, 多个换行. 按深度顺序添加
codesign -f -s "Developer ID Application: 你的名字 (团队ID)" "$APP_NAME/Contents/Frameworks/your.dylib"
codesign -f -s "Developer ID Application: 你的名字 (团队ID)" --options runtime "$APP_NAME"
...
`)
	lines := memo.Lines()
	if len(m.macCertArray) > 0 {
		for _, line := range m.macCertArray {
			lines.Add(line)
		}
	}
	memo.Font().SetSize(8)
	memo.SetParent(newCertCommandListForm)
	saveBtn := lcl.NewButton(newCertCommandListForm)
	saveBtn.SetBounds(140, 173, 50, 25)
	saveBtn.SetCaption("确定")
	saveBtn.SetParent(newCertCommandListForm)
	saveBtn.SetOnClick(func(sender lcl.IObject) {
		var macCertList []string
		lines = memo.Lines()
		for i := int32(0); i < lines.Count(); i++ {
			line := strings.TrimSpace(lines.Strings(i))
			if line == "" {
				continue
			}
			macCertList = append(macCertList, line)
		}
		m.macCertArray = macCertList
		newCertCommandListForm.Close()
	})
	closeBtn := lcl.NewButton(newCertCommandListForm)
	closeBtn.SetBounds(195, 173, 50, 25)
	closeBtn.SetCaption("取消")
	closeBtn.SetParent(newCertCommandListForm)
	closeBtn.SetOnClick(func(sender lcl.IObject) {
		newCertCommandListForm.Close()
	})
	newCertCommandListForm.ShowModal()
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
