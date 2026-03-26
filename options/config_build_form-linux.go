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

func (m *TBuildForm) initLinuxOptions() {
	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	linuxPackageTitle := lcl.NewLabel(m)
	linuxPackageTitle.SetFont(m.titleFont)
	linuxPackageTitle.SetCaption("Linux 打包配置")
	linuxPackageTitle.SetTop(nextTop(5))
	linuxPackageTitle.SetLeft(10)
	linuxPackageTitle.SetParent(m.platformTabPageLinux)

	linuxPackageFmtTitle := lcl.NewLabel(m)
	linuxPackageFmtTitle.SetCaption("打包格式")
	linuxPackageFmtTitle.SetLeft(10)
	linuxPackageFmtTitle.SetTop(nextTop(30))
	linuxPackageFmtTitle.SetFont(m.titleFontTwo)
	linuxPackageFmtTitle.SetParent(m.platformTabPageLinux)

	m.linuxDEBCheckBox = lcl.NewCheckBox(m)
	m.linuxDEBCheckBox.SetCaption("DEB 包")
	m.linuxDEBCheckBox.SetLeft(20)
	m.linuxDEBCheckBox.SetTop(nextTop(30))
	m.linuxDEBCheckBox.SetFont(m.font)
	m.linuxDEBCheckBox.SetChecked(bean.GProject.BuildOption.LinuxDEB)
	m.linuxDEBCheckBox.SetParent(m.platformTabPageLinux)

	dependsTitle := lcl.NewLabel(m)
	dependsTitle.SetCaption("依赖项")
	dependsTitle.SetLeft(10)
	dependsTitle.SetTop(nextTop(35))
	dependsTitle.SetFont(m.titleFontTwo)
	dependsTitle.SetParent(m.platformTabPageLinux)

	m.dependsEdit = lcl.NewEdit(m)
	m.dependsEdit.SetBounds(20, nextTop(30), 515, 30)
	m.dependsEdit.SetFont(m.font)
	m.dependsEdit.SetTextHint("用逗号分隔的依赖项列表, 如: libc6 (>= 2.14)")
	m.dependsEdit.SetText(bean.GProject.BuildOption.Depends)
	m.dependsEdit.SetParent(m.platformTabPageLinux)
}
