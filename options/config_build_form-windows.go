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
	m.winMsiCheckBox.SetCaption("MSI 安装包(MakeAppx)")
	m.winMsiCheckBox.SetLeft(20)
	m.winMsiCheckBox.SetTop(nextTop(25))
	m.winMsiCheckBox.SetFont(m.font)
	m.winMsiCheckBox.SetChecked(bean.GProject.BuildOption.WinMsi)
	m.winMsiCheckBox.SetParent(m.platformTabPageWindows)
	m.winExeCheckBox = lcl.NewCheckBox(m)
	m.winExeCheckBox.SetCaption("EXE 安装包(makensis)")
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
}
