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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

func (m *TConfigProjectForm) initWindowsOptions() {
	logs.Debug("TConfigProjectForm initWindowsOptions")

	// windows manifest
	windowsBaseTop := int32(10)
	m.appNameText = lcl.NewLabel(m)
	m.appNameText.SetLeft(10)
	m.appNameText.SetTop(windowsBaseTop)
	m.appNameText.SetCaption("名称")
	m.appNameText.SetParent(m.platformTabPageWindows)
	m.appNameEdit = lcl.NewEdit(m)
	m.appNameEdit.SetBounds(45, windowsBaseTop-5, 225, 30)
	m.appNameEdit.SetFont(m.font)
	m.appNameEdit.SetTextHint("company.product.app")
	m.appNameEdit.SetParent(m.platformTabPageWindows)

	m.appDescText = lcl.NewLabel(m)
	m.appDescText.SetLeft(m.appNameEdit.Left() + m.appNameEdit.Width() + 15)
	m.appDescText.SetTop(windowsBaseTop)
	m.appDescText.SetCaption("描述")
	m.appDescText.SetParent(m.platformTabPageWindows)
	m.appDescEdit = lcl.NewEdit(m)
	m.appDescEdit.SetBounds(m.appDescText.Left()+35, windowsBaseTop-5, 225, 30)
	m.appDescEdit.SetFont(m.font)
	m.appDescEdit.SetTextHint("your application description.")
	m.appDescEdit.SetParent(m.platformTabPageWindows)

	m.appVersionText = lcl.NewLabel(m)
	m.appVersionText.SetLeft(10)
	m.appVersionText.SetTop(windowsBaseTop + 45)
	m.appVersionText.SetCaption("版本")
	m.appVersionText.SetParent(m.platformTabPageWindows)
	m.appVersionEdit = lcl.NewEdit(m)
	m.appVersionEdit.SetBounds(m.appVersionText.Left()+35, windowsBaseTop+40, 225, 30)
	m.appVersionEdit.SetFont(m.font)
	m.appVersionEdit.SetTextHint("1.2.3.4")
	m.appVersionEdit.SetParent(m.platformTabPageWindows)

	m.compatibilityOSText = lcl.NewLabel(m)
	m.compatibilityOSText.SetLeft(m.appVersionEdit.Left() + m.appVersionEdit.Width() + 15)
	m.compatibilityOSText.SetTop(windowsBaseTop + 45)
	m.compatibilityOSText.SetCaption("兼容")
	m.compatibilityOSText.SetParent(m.platformTabPageWindows)
	m.compatibilityOSBox = lcl.NewComboBox(m)
	m.compatibilityOSBox.SetBounds(m.compatibilityOSText.Left()+35, windowsBaseTop+40, 225, 30)
	m.compatibilityOSBox.SetFont(m.font)
	m.compatibilityOSBox.SetReadOnly(true)
	m.compatibilityOSBox.SetStyle(types.CsDropDownList)
	m.compatibilityOSBox.SetParent(m.platformTabPageWindows)

	m.dpiText = lcl.NewLabel(m)
	m.dpiText.SetLeft(10)
	m.dpiText.SetTop(windowsBaseTop + 90)
	m.dpiText.SetCaption(" DPI")
	m.dpiText.SetParent(m.platformTabPageWindows)
	m.dpiBox = lcl.NewComboBox(m)
	m.dpiBox.SetBounds(m.dpiText.Left()+35, windowsBaseTop+85, 205, 30)
	m.dpiBox.SetFont(m.font)
	m.dpiBox.SetReadOnly(true)
	m.dpiBox.SetStyle(types.CsDropDownList)
	m.dpiBox.SetParent(m.platformTabPageWindows)

	m.runLevelText = lcl.NewLabel(m)
	m.runLevelText.SetLeft(m.dpiBox.Left() + m.dpiBox.Width() + 10)
	m.runLevelText.SetTop(windowsBaseTop + 90)
	m.runLevelText.SetCaption(" 执行等级")
	m.runLevelText.SetParent(m.platformTabPageWindows)
	m.runLevelBox = lcl.NewComboBox(m)
	m.runLevelBox.SetBounds(m.runLevelText.Left()+60, windowsBaseTop+85, 225, 30)
	m.runLevelBox.SetFont(m.font)
	m.runLevelBox.SetReadOnly(true)
	m.runLevelBox.SetStyle(types.CsDropDownList)
	m.runLevelBox.SetParent(m.platformTabPageWindows)

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

func (m *TConfigProjectForm) manifestDataInit() {
	compatibilityOSBoxItems := m.compatibilityOSBox.Items()
	compatibilityOSBoxItems.Add("Windows Vista")
	compatibilityOSBoxItems.Add("Windows 7")
	compatibilityOSBoxItems.Add("Windows 8")
	compatibilityOSBoxItems.Add("Windows 8.1")
	compatibilityOSBoxItems.Add("Windows 10")
	compatibilityOSBoxItems.Add("Windows 11")
	m.compatibilityOSBox.SetItemIndex(0)

	dpiBoxItems := m.dpiBox.Items()
	dpiBoxItems.Add("System (true)")
	dpiBoxItems.Add("UnAware (false)")
	dpiBoxItems.Add("PerMonitor (true/PM)")
	dpiBoxItems.Add("PerMonitorV2 (true/PM-V2)")
	m.dpiBox.SetItemIndex(0)

	runLevelBoxItems := m.runLevelBox.Items()
	runLevelBoxItems.Add("AsInvoker (当前用户)")
	runLevelBoxItems.Add("HighestAvailable (最高可用权限)")
	runLevelBoxItems.Add("RequireAdministrator (要求管理员)")
	m.runLevelBox.SetItemIndex(0)
}
