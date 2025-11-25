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

package project

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/winres"
	"github.com/energye/designer/project/bean"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
)

// initWindowsOptions 初始化 Windows 平台配置选项界面控件
//
// 该方法用于在 Windows 配置标签页中创建并布局各类输入控件（如文本框、下拉框、复选框等），
// 以便用户可以设置应用程序的元数据、兼容性、DPI 设置及各种高级 Manifest 属性。
// 控件包括应用名称、描述、版本、操作系统兼容性、执行权限级别以及多个与系统行为相关的开关选项。
func (m *TConfigProjectForm) initWindowsOptions() {
	logs.Debug("TConfigProjectForm initWindowsOptions")

	// windows manifest
	windowsBaseTop := int32(10)

	m.compatibilityOSText = lcl.NewLabel(m)
	m.compatibilityOSText.SetLeft(10)
	m.compatibilityOSText.SetTop(windowsBaseTop)
	m.compatibilityOSText.SetCaption("　　兼容")
	m.compatibilityOSText.SetParent(m.platformTabPageWindows)
	m.compatibilityOSBox = lcl.NewComboBox(m)
	m.compatibilityOSBox.SetBounds(m.compatibilityOSText.Left()+60, windowsBaseTop-5, 475, 30)
	m.compatibilityOSBox.SetFont(m.font)
	m.compatibilityOSBox.SetReadOnly(true)
	m.compatibilityOSBox.SetStyle(types.CsDropDownList)
	m.compatibilityOSBox.SetParent(m.platformTabPageWindows)

	m.dpiText = lcl.NewLabel(m)
	m.dpiText.SetLeft(10)
	m.dpiText.SetTop(windowsBaseTop + 40)
	m.dpiText.SetCaption("　　 DPI")
	m.dpiText.SetParent(m.platformTabPageWindows)
	m.dpiBox = lcl.NewComboBox(m)
	m.dpiBox.SetBounds(m.dpiText.Left()+60, windowsBaseTop+35, 475, 30)
	m.dpiBox.SetFont(m.font)
	m.dpiBox.SetReadOnly(true)
	m.dpiBox.SetStyle(types.CsDropDownList)
	m.dpiBox.SetParent(m.platformTabPageWindows)

	m.runLevelText = lcl.NewLabel(m)
	m.runLevelText.SetLeft(10)
	m.runLevelText.SetTop(windowsBaseTop + 80)
	m.runLevelText.SetCaption(" 执行等级")
	m.runLevelText.SetParent(m.platformTabPageWindows)
	m.runLevelBox = lcl.NewComboBox(m)
	m.runLevelBox.SetBounds(m.runLevelText.Left()+60, windowsBaseTop+75, 475, 30)
	m.runLevelBox.SetFont(m.font)
	m.runLevelBox.SetReadOnly(true)
	m.runLevelBox.SetStyle(types.CsDropDownList)
	m.runLevelBox.SetParent(m.platformTabPageWindows)

	bg := wg.NewButton(m)
	bg.SetDisabledColor(colors.RGBToColor(204, 232, 255), colors.RGBToColor(204, 232, 255))
	bg.SetBounds(5, windowsBaseTop+120, m.Width()-10, 130)
	bg.SetRadius(8)
	bg.SetDisable(true)
	bg.SetParent(m.platformTabPageWindows)

	m.uiAccessCheckBox = lcl.NewCheckBox(m)
	m.uiAccessCheckBox.SetLeft(10)
	m.uiAccessCheckBox.SetTop(windowsBaseTop + 125)
	m.uiAccessCheckBox.SetFont(m.font)
	m.uiAccessCheckBox.SetCaption("uiAccess (用户界面访问)")
	m.uiAccessCheckBox.SetHint("uiAccess (用户界面访问)")
	m.uiAccessCheckBox.SetShowHint(true)
	m.uiAccessCheckBox.SetParent(m.platformTabPageWindows)

	m.autoElevateBox = lcl.NewCheckBox(m)
	m.autoElevateBox.SetLeft(160 + 15)
	m.autoElevateBox.SetTop(windowsBaseTop + 125)
	m.autoElevateBox.SetFont(m.font)
	m.autoElevateBox.SetCaption("autoElevate (自动提权)")
	m.autoElevateBox.SetHint("autoElevate (自动提权)")
	m.autoElevateBox.SetShowHint(true)
	m.autoElevateBox.SetParent(m.platformTabPageWindows)

	m.disableThemingBox = lcl.NewCheckBox(m)
	m.disableThemingBox.SetLeft(160*2 + 10)
	m.disableThemingBox.SetTop(windowsBaseTop + 125)
	m.disableThemingBox.SetFont(m.font)
	m.disableThemingBox.SetCaption("DisableTheming (禁用主题)")
	m.disableThemingBox.SetHint("disableTheming (禁用主题)")
	m.disableThemingBox.SetShowHint(true)
	m.disableThemingBox.SetParent(m.platformTabPageWindows)

	m.disableWindowFilteringBox = lcl.NewCheckBox(m)
	m.disableWindowFilteringBox.SetLeft(10)
	m.disableWindowFilteringBox.SetTop(windowsBaseTop + 150)
	m.disableWindowFilteringBox.SetFont(m.font)
	m.disableWindowFilteringBox.SetCaption("disableWindowFiltering (禁用窗口过滤)")
	m.disableWindowFilteringBox.SetHint("disableWindowFiltering (禁用窗口过滤仅在 DPI 虚拟化启用时生效)")
	m.disableWindowFilteringBox.SetShowHint(true)
	m.disableWindowFilteringBox.SetParent(m.platformTabPageWindows)

	m.highResolutionScrollingAwareBox = lcl.NewCheckBox(m)
	m.highResolutionScrollingAwareBox.SetLeft(10)
	m.highResolutionScrollingAwareBox.SetTop(windowsBaseTop + 175)
	m.highResolutionScrollingAwareBox.SetFont(m.font)
	m.highResolutionScrollingAwareBox.SetCaption("highResolutionScrollingAware (高分辨率滚动)")
	m.highResolutionScrollingAwareBox.SetHint("highResolutionScrollingAware (高分辨率滚动)")
	m.highResolutionScrollingAwareBox.SetShowHint(true)
	m.highResolutionScrollingAwareBox.SetParent(m.platformTabPageWindows)

	m.ultraHighResolutionScrollingAwareBox = lcl.NewCheckBox(m)
	m.ultraHighResolutionScrollingAwareBox.SetLeft(10)
	m.ultraHighResolutionScrollingAwareBox.SetTop(windowsBaseTop + 200)
	m.ultraHighResolutionScrollingAwareBox.SetFont(m.font)
	m.ultraHighResolutionScrollingAwareBox.SetCaption("ultraHighResolutionScrollingAware (超高分辨率滚动)")
	m.ultraHighResolutionScrollingAwareBox.SetHint("ultraHighResolutionScrollingAware (超高分辨率滚动Windows 10 2004+ / Windows 11)")
	m.ultraHighResolutionScrollingAwareBox.SetShowHint(true)
	m.ultraHighResolutionScrollingAwareBox.SetParent(m.platformTabPageWindows)

	m.longPathAwareBox = lcl.NewCheckBox(m)
	m.longPathAwareBox.SetLeft(160*2 + 10)
	m.longPathAwareBox.SetTop(windowsBaseTop + 150)
	m.longPathAwareBox.SetFont(m.font)
	m.longPathAwareBox.SetCaption("longPathAware (启用长路径支持)")
	m.longPathAwareBox.SetHint("longPathAware (启用长路径支持 Windows 10 1607 +)")
	m.longPathAwareBox.SetShowHint(true)
	m.longPathAwareBox.SetParent(m.platformTabPageWindows)

	m.gDIScalingBox = lcl.NewCheckBox(m)
	m.gDIScalingBox.SetLeft(160*2 + 10)
	m.gDIScalingBox.SetTop(windowsBaseTop + 175)
	m.gDIScalingBox.SetFont(m.font)
	m.gDIScalingBox.SetCaption("gdiScaling (GDI 自动缩放)")
	m.gDIScalingBox.SetHint("gdiScaling (启用 GDI 自动缩放 Windows 10 1703+)")
	m.gDIScalingBox.SetShowHint(true)
	m.gDIScalingBox.SetParent(m.platformTabPageWindows)

	m.segmentHeapBox = lcl.NewCheckBox(m)
	m.segmentHeapBox.SetLeft(160*2 + 10)
	m.segmentHeapBox.SetTop(windowsBaseTop + 200)
	m.segmentHeapBox.SetFont(m.font)
	m.segmentHeapBox.SetCaption("segmentHeap (分段堆)")
	m.segmentHeapBox.SetHint("启用 Segment Heap（Windows 10 2004+）")
	m.segmentHeapBox.SetShowHint(true)
	m.segmentHeapBox.SetEnabled(true)
	m.segmentHeapBox.SetParent(m.platformTabPageWindows)

	m.printerDriverIsolationBox = lcl.NewCheckBox(m)
	m.printerDriverIsolationBox.SetLeft(10)
	m.printerDriverIsolationBox.SetTop(windowsBaseTop + 225)
	m.printerDriverIsolationBox.SetFont(m.font)
	m.printerDriverIsolationBox.SetCaption("printerDriverIsolation (打印驱动隔离)")
	m.printerDriverIsolationBox.SetHint("printerDriverIsolation (启用打印驱动隔离, 仅适用于打印驱动组件)")
	m.printerDriverIsolationBox.SetShowHint(true)
	m.printerDriverIsolationBox.SetParent(m.platformTabPageWindows)

	m.useCommonControlsV6Box = lcl.NewCheckBox(m)
	m.useCommonControlsV6Box.SetLeft(160*2 + 10)
	m.useCommonControlsV6Box.SetTop(windowsBaseTop + 225)
	m.useCommonControlsV6Box.SetFont(m.font)
	m.useCommonControlsV6Box.SetCaption("useCommonControlsV6 (视觉样式)")
	m.useCommonControlsV6Box.SetHint("useCommonControlsV6 (启用视觉样式的现代控件)")
	m.useCommonControlsV6Box.SetShowHint(true)
	m.useCommonControlsV6Box.SetParent(m.platformTabPageWindows)

	m.manifestDataInit()
}

// manifestDataInit 初始化清单配置相关的下拉框数据
// 该函数用于初始化兼容性操作系统、DPI感知模式和运行级别三个下拉框的选项内容，
// 并设置默认选中项为第一个选项
func (m *TConfigProjectForm) manifestDataInit() {
	compatibilityOSBoxItems := m.compatibilityOSBox.Items()
	for _, item := range bean.CompatibilityOSList.Values() {
		compatibilityOSBoxItems.Add(item)
	}
	m.compatibilityOSBox.SetItemIndex(int32(winres.WinVistaAndAbove))

	dpiBoxItems := m.dpiBox.Items()
	for _, item := range bean.DPIList.Values() {
		dpiBoxItems.Add(item)
	}
	m.dpiBox.SetItemIndex(int32(winres.DPIAware))

	runLevelBoxItems := m.runLevelBox.Items()
	for _, item := range bean.RunLevelList.Values() {
		runLevelBoxItems.Add(item)
	}
	m.runLevelBox.SetItemIndex(int32(winres.AsInvoker))
}

// saveWindows 保存 Windows 平台的配置信息
// 该函数负责将 Windows 平台相关的配置数据进行持久化存储，
// 包括应用程序的元数据、兼容性设置、DPI 配置以及各种 Manifest 属性等。
// 最后生成一个完整的 Manifest 文件，并将其保存到项目目录中。
func (m *TConfigProjectForm) saveWindows() {
	// TODO: 实现 Windows 配置的保存逻辑
	// 需要保存的配置包括：
	// 1. 应用程序基本信息：名称、描述、版本
	// 2. 兼容性设置：目标操作系统
	// 3. DPI 设置：DPI 感知模式
	// 4. 执行权限：运行级别
	// 5. 各种功能开关状态：uiAccess、autoElevate 等复选框状态

}
