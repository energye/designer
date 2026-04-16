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
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"os"
	"path/filepath"
	"strings"
)

func (m *TConfigProjectForm) initMacOSOptions() {
	logs.Debug("TConfigProjectForm initMacOSOptions")
	baseTop := int32(10)
	baseLeft := int32(110)
	nextTop := func(v int32) int32 {
		baseTop += v
		return baseTop
	}

	m.CFBundleNameEdit = lcl.NewLabeledEdit(m)
	m.CFBundleNameEdit.EditLabel().SetCaption("短标题")
	m.CFBundleNameEdit.SetBounds(baseLeft, nextTop(0), 430, 30)
	m.CFBundleNameEdit.SetFont(m.font)
	m.CFBundleNameEdit.SetShowHint(true)
	m.CFBundleNameEdit.SetTextHint("应用的短显示名称, 默认: 构建二进制文件名")
	m.CFBundleNameEdit.SetHint("应用的短显示名称, 默认: 构建二进制文件名")
	m.CFBundleNameEdit.SetText(bean.GProject.AppOption.MacOS.PList.CFBundleName)
	m.CFBundleNameEdit.SetLabelPosition(types.LpLeft)
	m.CFBundleNameEdit.SetParent(m.platformTabPageMacOS)

	m.CFBundleLocalizationsEdit = lcl.NewLabeledEdit(m)
	m.CFBundleLocalizationsEdit.EditLabel().SetCaption("本地化语言列表")
	m.CFBundleLocalizationsEdit.SetBounds(baseLeft, nextTop(35), 230, 30)
	m.CFBundleLocalizationsEdit.SetFont(m.font)
	m.CFBundleLocalizationsEdit.SetShowHint(true)
	m.CFBundleLocalizationsEdit.SetTextHint("本地化语言列表, 豆号分隔 zh_CN,en")
	m.CFBundleLocalizationsEdit.SetHint("本地化语言列表, 豆号分隔 zh_CN,en")
	m.CFBundleLocalizationsEdit.SetText(strings.Join(bean.GProject.AppOption.MacOS.PList.CFBundleLocalizations, ","))
	m.CFBundleLocalizationsEdit.SetLabelPosition(types.LpLeft)
	m.CFBundleLocalizationsEdit.SetParent(m.platformTabPageMacOS)

	m.LSUIElementText = lcl.NewLabel(m)
	m.LSUIElementText.SetLeft(15)
	m.LSUIElementText.SetTop(nextTop(35))
	m.LSUIElementText.SetCaption("　控制运行模式")
	m.LSUIElementText.SetParent(m.platformTabPageMacOS)
	m.LSUIElementBox = lcl.NewComboBox(m)
	m.LSUIElementBox.SetBounds(baseLeft, m.LSUIElementText.Top(), 430, 30)
	m.LSUIElementBox.SetFont(m.font)
	m.LSUIElementBox.SetReadOnly(true)
	m.LSUIElementBox.SetStyle(types.CsDropDownList)
	m.LSUIElementBox.SetShowHint(true)
	m.LSUIElementBox.SetHint("控制运行模式, 默认: 常规前台应用")
	m.LSUIElementBox.SetParent(m.platformTabPageMacOS)

	m.LSMinimumSystemVersionText = lcl.NewLabel(m)
	m.LSMinimumSystemVersionText.SetLeft(15)
	m.LSMinimumSystemVersionText.SetTop(nextTop(35))
	m.LSMinimumSystemVersionText.SetCaption("　最低系统版本")
	m.LSMinimumSystemVersionText.SetParent(m.platformTabPageMacOS)
	m.LSMinimumSystemVersionBox = lcl.NewComboBox(m)
	m.LSMinimumSystemVersionBox.SetBounds(baseLeft, m.LSMinimumSystemVersionText.Top(), 430, 30)
	m.LSMinimumSystemVersionBox.SetFont(m.font)
	m.LSMinimumSystemVersionBox.SetReadOnly(true)
	m.LSMinimumSystemVersionBox.SetStyle(types.CsDropDownList)
	m.LSMinimumSystemVersionBox.SetShowHint(true)
	m.LSMinimumSystemVersionBox.SetHint("支持最低系统版本")
	m.LSMinimumSystemVersionBox.SetParent(m.platformTabPageMacOS)

	m.pListDataInit()
}

func (m *TConfigProjectForm) pListDataInit() {
	LSUIElementBoxItems := m.LSUIElementBox.Items()
	bean.LSUIElementList.Iterate(func(key bean.MacOSUIElementList, value string) bool {
		LSUIElementBoxItems.Add(value)
		return false
	})
	m.LSUIElementBox.SetItemIndex(bean.GProject.AppOption.MacOS.PList.LSUIElementIndex)

	LSMinimumSystemVersionItems := m.LSMinimumSystemVersionBox.Items()
	bean.LSMinimumSystemVersionList.Iterate(func(key bean.LSMinimumSystemVersion, value string) bool {
		LSMinimumSystemVersionItems.Add(value)
		return false
	})
	m.LSMinimumSystemVersionBox.SetItemIndex(bean.GProject.AppOption.MacOS.PList.LSMinimumSystemVersionIndex)
}

func createAppLocalizations() {
	// Contents/Resources/xxx.lproj
	resourcesWindowsMetadataPath := bean.ResourceMetadataPath()
	for _, local := range bean.GProject.AppOption.MacOS.PList.CFBundleLocalizations {
		resourcesLocal := filepath.Join(resourcesWindowsMetadataPath, local+".lproj")
		if tool.IsExist(resourcesLocal) {
			continue
		}
		if err := os.MkdirAll(resourcesLocal, 0755); err != nil {
			event.ConsoleWriteError("Unable to create localizations:", err.Error())
			continue
		}
		localizations := `/* localizations */
CFBundleDisplayName = "{{CFBundleDisplayName}}";
CFBundleName = "{{CFBundleName}}";
`
		localizations = strings.Replace(localizations, "{{CFBundleDisplayName}}", bean.GProject.AppOption.MacOS.PList.CFBundleDisplayName, 1)
		localizations = strings.Replace(localizations, "{{CFBundleName}}", bean.GProject.AppOption.MacOS.PList.CFBundleName, 1)
		_ = os.WriteFile(filepath.Join(resourcesLocal, "InfoPlist.strings"), []byte(localizations), 0644)
	}
}
