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
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"os"
	"path/filepath"
	"strings"
)

func (m *TConfigProjectForm) initMacOSOptions() {
	logs.Debug("TConfigProjectForm initMacOSOptions")
	// macOS plist.info

	//box := lcl.NewScrollBox(m)
	//box.SetAlign(types.AlClient)
	//box.SetBorderStyleToBorderStyle(types.BsNone)
	//box.VertScrollBar().SetVisible(true)
	//box.VertScrollBar().SetTracking(true)
	//box.HorzScrollBar().SetVisible(false)
	//box.SetParent(m.platformTabPageMacOS)

	/**
	   	CFBundleExecutable 可执行文件名称 = 在编译选项配置
	   	CFBundleDisplayName 应用的短显示名称 = 应用标题
	    CFBundleIdentifier 应用的唯一标识 = 标识
		CFBundleVersion 内部版本号 = 版本
		CFBundleShortVersionString 对外版本号 = 版本
		CFBundleGetInfoString 应用备注信息 = 描述
		CFBundleIconFile 应用图标 = 图标需生成 xxx.icns
		NSHumanReadableCopyright 应用版权信息 = 版权

	   	CFBundleName 应用的短显示名称 = 新增
		CFBundleLocalizations 应用本地化语言列表 = 新增, 文本框 格式 豆号分隔 zh-CN,en-US
		LSUIElement 控制运行模式 = 新增
			- true：无 Dock 图标 / 菜单栏，仅状态栏运行（后台应用）；
			- false：常规前台应用。 设为true后需通过状态栏触发界面，不可与 LSBackgroundOnly 同时为true。
		LSMinimumSystemVersion 系统最低版本 = 新增 10.15 或 11.0

	*/
	baseTop := int32(10)

	m.CFBundleNameText = lcl.NewLabel(m)
	m.CFBundleNameText.SetLeft(10)
	m.CFBundleNameText.SetTop(baseTop)
	m.CFBundleNameText.SetCaption("　　　　短标题")
	m.CFBundleNameText.SetParent(m.platformTabPageMacOS)
	m.CFBundleNameEdit = lcl.NewEdit(m)
	m.CFBundleNameEdit.SetBounds(m.CFBundleNameText.Left()+90, baseTop-5, 440, 30)
	m.CFBundleNameEdit.SetFont(m.font)
	m.CFBundleNameEdit.SetShowHint(true)
	m.CFBundleNameEdit.SetTextHint("应用的短显示名称, 默认: 构建二进制文件名")
	m.CFBundleNameEdit.SetHint("应用的短显示名称, 默认: 构建二进制文件名")
	m.CFBundleNameEdit.SetText(gProject.AppOption.MacOS.PList.CFBundleName)
	m.CFBundleNameEdit.SetParent(m.platformTabPageMacOS)

	m.CFBundleLocalizationsText = lcl.NewLabel(m)
	m.CFBundleLocalizationsText.SetLeft(10)
	m.CFBundleLocalizationsText.SetTop(baseTop + 40)
	m.CFBundleLocalizationsText.SetCaption("本地化语言列表")
	m.CFBundleLocalizationsText.SetParent(m.platformTabPageMacOS)
	m.CFBundleLocalizationsEdit = lcl.NewEdit(m)
	m.CFBundleLocalizationsEdit.SetBounds(m.CFBundleLocalizationsText.Left()+90, baseTop+35, 440, 30)
	m.CFBundleLocalizationsEdit.SetFont(m.font)
	m.CFBundleLocalizationsEdit.SetShowHint(true)
	m.CFBundleLocalizationsEdit.SetTextHint("本地化语言列表, 豆号分隔 zh_CN,en_US, 默认: zh_CN")
	m.CFBundleLocalizationsEdit.SetHint("本地化语言列表, 豆号分隔 zh_CN,en_US, 默认: zh_CN")
	m.CFBundleLocalizationsEdit.SetText(strings.Join(gProject.AppOption.MacOS.PList.CFBundleLocalizations, ","))
	m.CFBundleLocalizationsEdit.SetParent(m.platformTabPageMacOS)

	m.LSUIElementText = lcl.NewLabel(m)
	m.LSUIElementText.SetLeft(10)
	m.LSUIElementText.SetTop(baseTop + 80)
	m.LSUIElementText.SetCaption("　控制运行模式")
	m.LSUIElementText.SetParent(m.platformTabPageMacOS)
	m.LSUIElementBox = lcl.NewComboBoxEx(m)
	m.LSUIElementBox.SetBounds(m.LSUIElementText.Left()+90, baseTop+75, 440, 30)
	m.LSUIElementBox.SetFont(m.font)
	m.LSUIElementBox.SetReadOnly(true)
	m.LSUIElementBox.SetStyle(types.CsDropDownList)
	m.LSUIElementBox.SetShowHint(true)
	m.LSUIElementBox.SetHint("控制运行模式, 默认: 常规前台应用")
	m.LSUIElementBox.SetParent(m.platformTabPageMacOS)

	m.LSMinimumSystemVersionText = lcl.NewLabel(m)
	m.LSMinimumSystemVersionText.SetLeft(10)
	m.LSMinimumSystemVersionText.SetTop(baseTop + 120)
	m.LSMinimumSystemVersionText.SetCaption("　最低系统版本")
	m.LSMinimumSystemVersionText.SetParent(m.platformTabPageMacOS)
	m.LSMinimumSystemVersionBox = lcl.NewComboBoxEx(m)
	m.LSMinimumSystemVersionBox.SetBounds(m.LSMinimumSystemVersionText.Left()+90, baseTop+115, 440, 30)
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
	m.LSUIElementBox.SetItemIndex(gProject.AppOption.MacOS.PList.LSUIElementIndex)

	LSMinimumSystemVersionItems := m.LSMinimumSystemVersionBox.Items()
	bean.LSMinimumSystemVersionList.Iterate(func(key bean.LSMinimumSystemVersion, value string) bool {
		LSMinimumSystemVersionItems.Add(value)
		return false
	})
	m.LSMinimumSystemVersionBox.SetItemIndex(gProject.AppOption.MacOS.PList.LSMinimumSystemVersionIndex)
}

// 保存或更新 macOS 配置并生成程序信息
func saveOrUpdateMacOSPList() {
	pListInfoTemplate := app.Packager("darwin/Info.plist")
	if pListInfoTemplate == nil {
		logs.Error("macOS 应用配置-保存配置 info.plist 模板获取失败, 模板内容为 nil")
		return
	}
	pListInfo, err := tool.RenderTemplate(string(pListInfoTemplate), gProject.AppOption.MacOS)
	if err != nil {
		logs.Error("macOS 应用配置-保存配置 info.plist 内容渲染失败:", err.Error())
		return
	}
	// 保存到 resources/Info.plist
	resourcesPath := ResourcePath()
	pListOutFile := "Info.plist"
	err = os.WriteFile(filepath.Join(resourcesPath, pListOutFile), pListInfo, 0666)
	if err != nil {
		logs.Error("macOS 应用配置-保存配置-WriteFile: ", err.Error())
	}
}
