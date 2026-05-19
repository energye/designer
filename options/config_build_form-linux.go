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
)

func (m *TBuildForm) initLinuxOptions() {
	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	linuxPackageFmtTitle := lcl.NewLabel(m)
	linuxPackageFmtTitle.SetCaption("打包格式")
	linuxPackageFmtTitle.SetLeft(10)
	linuxPackageFmtTitle.SetTop(nextTop(5))
	linuxPackageFmtTitle.SetFont(m.titleFontTwo)
	linuxPackageFmtTitle.SetParent(m.platformTabPageLinux)

	m.linuxDEBCheckBox = lcl.NewCheckBox(m)
	m.linuxDEBCheckBox.SetCaption("deb")
	m.linuxDEBCheckBox.SetLeft(20)
	m.linuxDEBCheckBox.SetTop(nextTop(30))
	m.linuxDEBCheckBox.SetFont(m.font)
	m.linuxDEBCheckBox.SetChecked(bean.GProject.BuildOption.LinuxDEB)
	m.linuxDEBCheckBox.SetParent(m.platformTabPageLinux)

	m.linuxRPMCheckBox = lcl.NewCheckBox(m)
	m.linuxRPMCheckBox.SetCaption("rpm")
	m.linuxRPMCheckBox.SetLeft(210)
	m.linuxRPMCheckBox.SetTop(m.linuxDEBCheckBox.Top())
	m.linuxRPMCheckBox.SetFont(m.font)
	m.linuxRPMCheckBox.SetChecked(bean.GProject.BuildOption.LinuxRPM)
	m.linuxRPMCheckBox.SetParent(m.platformTabPageLinux)

	m.linuxAppImageCheckBox = lcl.NewCheckBox(m)
	m.linuxAppImageCheckBox.SetCaption("AppImage")
	m.linuxAppImageCheckBox.SetLeft(390)
	m.linuxAppImageCheckBox.SetTop(m.linuxRPMCheckBox.Top())
	m.linuxAppImageCheckBox.SetFont(m.font)
	m.linuxAppImageCheckBox.SetChecked(bean.GProject.BuildOption.LinuxAppImage)
	m.linuxAppImageCheckBox.SetParent(m.platformTabPageLinux)

	configTitle := lcl.NewLabel(m)
	configTitle.SetCaption("配置选项")
	configTitle.SetLeft(10)
	configTitle.SetTop(nextTop(35))
	configTitle.SetFont(m.titleFontTwo)
	configTitle.SetParent(m.platformTabPageLinux)

	anchors := types.NewSet(types.AkLeft, types.AkRight, types.AkTop)
	m.dependsEdit = lcl.NewEdit(m)
	m.dependsEdit.SetBounds(20, nextTop(35), 515, 30)
	m.dependsEdit.SetFont(m.font)
	m.dependsEdit.SetTextHint("CSV list of dependencies, e.g: libc6 (>= 2.14)")
	m.dependsEdit.SetShowHint(true)
	m.dependsEdit.SetText(bean.GProject.BuildOption.Depends)
	m.dependsEdit.SetAnchors(anchors)
	m.dependsEdit.SetParent(m.platformTabPageLinux)

	m.categoriesEdit = lcl.NewEdit(m)
	m.categoriesEdit.SetBounds(20, nextTop(35), 515, 30)
	m.categoriesEdit.SetFont(m.font)
	m.categoriesEdit.SetTextHint("Categories, e.g: Development;Utility;")
	m.categoriesEdit.SetText(bean.GProject.AppOption.Linux.Categories)
	m.categoriesEdit.SetAnchors(anchors)
	m.categoriesEdit.SetParent(m.platformTabPageLinux)

	m.homepageEdit = lcl.NewEdit(m)
	m.homepageEdit.SetBounds(20, nextTop(35), 515, 30)
	m.homepageEdit.SetFont(m.font)
	m.homepageEdit.SetTextHint("Homepage, e.g: https://example.com")
	m.homepageEdit.SetText(bean.GProject.AppOption.Linux.Homepage)
	m.homepageEdit.SetAnchors(anchors)
	m.homepageEdit.SetParent(m.platformTabPageLinux)

	m.maintainerEdit = lcl.NewEdit(m)
	m.maintainerEdit.SetBounds(20, nextTop(35), 515, 30)
	m.maintainerEdit.SetFont(m.font)
	m.maintainerEdit.SetTextHint("Maintainer, e.g: Name <email@example.com>")
	m.maintainerEdit.SetText(bean.GProject.AppOption.Linux.Maintainer)
	m.maintainerEdit.SetAnchors(anchors)
	m.maintainerEdit.SetParent(m.platformTabPageLinux)

	m.licenseEdit = lcl.NewEdit(m)
	m.licenseEdit.SetBounds(20, nextTop(35), 515, 30)
	m.licenseEdit.SetFont(m.font)
	m.licenseEdit.SetTextHint("License, e.g: MIT, GPL-3.0")
	m.licenseEdit.SetText(bean.GProject.AppOption.Linux.License)
	m.licenseEdit.SetAnchors(anchors)
	m.licenseEdit.SetParent(m.platformTabPageLinux)
}
