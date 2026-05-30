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
	"bytes"
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/winicon"
	"github.com/energye/designer/pkg/winres"
	"github.com/energye/designer/pkg/winres/version"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	m.compatibilityOSText.SetName("ConfigProjectFormCompatibilityOSText")
	m.compatibilityOSText.SetLeft(0)
	m.compatibilityOSText.SetTop(windowsBaseTop)
	m.compatibilityOSText.SetCaption(metadata.GI18n.Dict("ConfigProjectFormCompatibilityOSText.Caption"))
	m.compatibilityOSText.SetWidth(75)
	m.compatibilityOSText.SetAutoSize(false)
	m.compatibilityOSText.SetAlignment(types.TaRightJustify)
	m.compatibilityOSText.SetParent(m.platformTabPageWindows)

	m.compatibilityOSBox = lcl.NewComboBox(m)
	m.compatibilityOSBox.SetName("ConfigProjectFormCompatibilityOSBox")
	m.compatibilityOSBox.SetBounds(m.compatibilityOSText.Left()+80, windowsBaseTop-3, 460, 30)
	m.compatibilityOSBox.SetReadOnly(true)
	m.compatibilityOSBox.SetStyle(types.CsDropDownList)
	m.compatibilityOSBox.AnchorSideTop().SetControl(m.compatibilityOSText)
	m.compatibilityOSBox.AnchorSideTop().SetSide(types.AsrCenter)
	m.compatibilityOSBox.SetParent(m.platformTabPageWindows)

	m.dpiText = lcl.NewLabel(m)
	m.dpiText.SetName("ConfigProjectFormDpiText")
	m.dpiText.SetLeft(0)
	m.dpiText.SetTop(windowsBaseTop + 40)
	m.dpiText.SetCaption(metadata.GI18n.Dict("ConfigProjectFormDpiText.Caption"))
	m.dpiText.SetWidth(75)
	m.dpiText.SetAutoSize(false)
	m.dpiText.SetAlignment(types.TaRightJustify)
	m.dpiText.SetParent(m.platformTabPageWindows)

	m.dpiBox = lcl.NewComboBox(m)
	m.dpiBox.SetName("ConfigProjectFormDpiBox")
	m.dpiBox.SetBounds(m.dpiText.Left()+80, windowsBaseTop+40, 460, 30)
	m.dpiBox.SetReadOnly(true)
	m.dpiBox.SetStyle(types.CsDropDownList)
	m.dpiBox.AnchorSideTop().SetControl(m.dpiText)
	m.dpiBox.AnchorSideTop().SetSide(types.AsrCenter)
	m.dpiBox.SetParent(m.platformTabPageWindows)

	m.runLevelText = lcl.NewLabel(m)
	m.runLevelText.SetName("ConfigProjectFormRunLevelText")
	m.runLevelText.SetLeft(0)
	m.runLevelText.SetTop(windowsBaseTop + 80)
	m.runLevelText.SetWidth(75)
	m.runLevelText.SetAutoSize(false)
	m.runLevelText.SetAlignment(types.TaRightJustify)
	m.runLevelText.SetCaption(metadata.GI18n.Dict("ConfigProjectFormRunLevelText.Caption"))
	m.runLevelText.SetParent(m.platformTabPageWindows)

	m.runLevelBox = lcl.NewComboBox(m)
	m.runLevelBox.SetName("ConfigProjectFormRunLevelBox")
	m.runLevelBox.SetBounds(m.runLevelText.Left()+80, windowsBaseTop+80, 460, 30)
	m.runLevelBox.SetReadOnly(true)
	m.runLevelBox.SetStyle(types.CsDropDownList)
	m.runLevelBox.AnchorSideTop().SetControl(m.runLevelText)
	m.runLevelBox.AnchorSideTop().SetSide(types.AsrCenter)
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
	m.uiAccessCheckBox.SetCaption("uiAccess")
	m.uiAccessCheckBox.SetHint("uiAccess")
	m.uiAccessCheckBox.SetShowHint(true)
	m.uiAccessCheckBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.UIAccess)
	m.uiAccessCheckBox.SetParent(m.platformTabPageWindows)

	m.autoElevateBox = lcl.NewCheckBox(m)
	m.autoElevateBox.SetLeft(160 + 15)
	m.autoElevateBox.SetTop(windowsBaseTop + 125)
	m.autoElevateBox.SetCaption("autoElevate")
	m.autoElevateBox.SetHint("autoElevate")
	m.autoElevateBox.SetShowHint(true)
	m.autoElevateBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.AutoElevate)
	m.autoElevateBox.SetParent(m.platformTabPageWindows)

	m.disableThemingBox = lcl.NewCheckBox(m)
	m.disableThemingBox.SetLeft(160*2 + 10)
	m.disableThemingBox.SetTop(windowsBaseTop + 125)
	m.disableThemingBox.SetCaption("DisableTheming")
	m.disableThemingBox.SetHint("disableTheming")
	m.disableThemingBox.SetShowHint(true)
	m.disableThemingBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.DisableTheming)
	m.disableThemingBox.SetParent(m.platformTabPageWindows)

	m.disableWindowFilteringBox = lcl.NewCheckBox(m)
	m.disableWindowFilteringBox.SetLeft(10)
	m.disableWindowFilteringBox.SetTop(windowsBaseTop + 150)
	m.disableWindowFilteringBox.SetCaption("disableWindowFiltering")
	m.disableWindowFilteringBox.SetHint("disableWindowFiltering (Disable window filtering, takes effect only when DPI virtualization is enabled.)")
	m.disableWindowFilteringBox.SetShowHint(true)
	m.disableWindowFilteringBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.DisableWindowFiltering)
	m.disableWindowFilteringBox.SetParent(m.platformTabPageWindows)

	m.highResolutionScrollingAwareBox = lcl.NewCheckBox(m)
	m.highResolutionScrollingAwareBox.SetLeft(10)
	m.highResolutionScrollingAwareBox.SetTop(windowsBaseTop + 175)
	m.highResolutionScrollingAwareBox.SetCaption("highResolutionScrollingAware")
	m.highResolutionScrollingAwareBox.SetHint("highResolutionScrollingAware (High-resolution scrolling)")
	m.highResolutionScrollingAwareBox.SetShowHint(true)
	m.highResolutionScrollingAwareBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.HighResolutionScrollingAware)
	m.highResolutionScrollingAwareBox.SetParent(m.platformTabPageWindows)

	m.ultraHighResolutionScrollingAwareBox = lcl.NewCheckBox(m)
	m.ultraHighResolutionScrollingAwareBox.SetLeft(10)
	m.ultraHighResolutionScrollingAwareBox.SetTop(windowsBaseTop + 200)
	m.ultraHighResolutionScrollingAwareBox.SetCaption("ultraHighResolutionScrollingAware")
	m.ultraHighResolutionScrollingAwareBox.SetHint("ultraHighResolutionScrollingAware (Ultra-high resolution scrolling Windows 10 2004+ / Windows 11)")
	m.ultraHighResolutionScrollingAwareBox.SetShowHint(true)
	m.ultraHighResolutionScrollingAwareBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware)
	m.ultraHighResolutionScrollingAwareBox.SetParent(m.platformTabPageWindows)

	m.longPathAwareBox = lcl.NewCheckBox(m)
	m.longPathAwareBox.SetLeft(160*2 + 10)
	m.longPathAwareBox.SetTop(windowsBaseTop + 150)
	m.longPathAwareBox.SetCaption("longPathAware")
	m.longPathAwareBox.SetHint("longPathAware (Enable long path support Windows 10 1607 +)")
	m.longPathAwareBox.SetShowHint(true)
	m.longPathAwareBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.LongPathAware)
	m.longPathAwareBox.SetParent(m.platformTabPageWindows)

	m.gDIScalingBox = lcl.NewCheckBox(m)
	m.gDIScalingBox.SetLeft(160*2 + 10)
	m.gDIScalingBox.SetTop(windowsBaseTop + 175)
	m.gDIScalingBox.SetCaption("gdiScaling")
	m.gDIScalingBox.SetHint("gdiScaling (Enable GDI auto-scaling Windows 10 1703+)")
	m.gDIScalingBox.SetShowHint(true)
	m.gDIScalingBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.GDIScaling)
	m.gDIScalingBox.SetParent(m.platformTabPageWindows)

	m.segmentHeapBox = lcl.NewCheckBox(m)
	m.segmentHeapBox.SetLeft(160*2 + 10)
	m.segmentHeapBox.SetTop(windowsBaseTop + 200)
	m.segmentHeapBox.SetCaption("segmentHeap")
	m.segmentHeapBox.SetHint("Enable Segment Heap（Windows 10 2004+）")
	m.segmentHeapBox.SetShowHint(true)
	m.segmentHeapBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.SegmentHeap)
	m.segmentHeapBox.SetParent(m.platformTabPageWindows)

	m.printerDriverIsolationBox = lcl.NewCheckBox(m)
	m.printerDriverIsolationBox.SetLeft(10)
	m.printerDriverIsolationBox.SetTop(windowsBaseTop + 225)
	m.printerDriverIsolationBox.SetCaption("printerDriverIsolation")
	m.printerDriverIsolationBox.SetHint("printerDriverIsolation (Enable print driver isolation, applies only to print driver components)")
	m.printerDriverIsolationBox.SetShowHint(true)
	m.printerDriverIsolationBox.SetChecked(bean.GProject.AppOption.Windows.Manifest.PrinterDriverIsolation)
	m.printerDriverIsolationBox.SetParent(m.platformTabPageWindows)

	m.useCommonControlsV6Box = lcl.NewCheckBox(m)
	m.useCommonControlsV6Box.SetLeft(160*2 + 10)
	m.useCommonControlsV6Box.SetTop(windowsBaseTop + 225)
	m.useCommonControlsV6Box.SetCaption("useCommonControlsV6")
	m.useCommonControlsV6Box.SetHint("useCommonControlsV6 (Enable modern controls with visual styles)")
	m.useCommonControlsV6Box.SetShowHint(true)
	m.useCommonControlsV6Box.SetChecked(bean.GProject.AppOption.Windows.Manifest.UseCommonControlsV6)
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
	m.compatibilityOSBox.SetItemIndex(int32(bean.GProject.AppOption.Windows.Manifest.CompatibilityOS))

	dpiNamePrefix := strings.ToLower(m.dpiBox.Name() + ".Items[")           // > ConfigProjectFormDpiBox.Items[0]
	runLevelNamePrefix := strings.ToLower(m.runLevelBox.Name() + ".Items[") // > ConfigProjectFormRunLevelBox.Items[0]
	dpiList := make(map[winres.DPIAwareness]string)
	runLevelList := make(map[winres.ExecutionLevel]string)
	metadata.GI18n.Iterate(func(name, value string) bool {
		name = strings.ToLower(name)
		if strings.Contains(name, dpiNamePrefix) {
			index := name[len(dpiNamePrefix) : len(dpiNamePrefix)+1]
			i, _ := strconv.Atoi(index)
			dpiList[winres.DPIAwareness(i)] = value
		} else if strings.Contains(name, runLevelNamePrefix) {
			index := name[len(runLevelNamePrefix) : len(runLevelNamePrefix)+1]
			i, _ := strconv.Atoi(index)
			runLevelList[winres.ExecutionLevel(i)] = value
		}
		return false
	})
	// dpi
	{
		dpiBoxItems := m.dpiBox.Items()
		dpiAwareness := []winres.DPIAwareness{winres.DPIAware, winres.DPIUnaware,
			winres.DPIPerMonitor, winres.DPIPerMonitorV2}
		if len(dpiList) == len(dpiAwareness) {
			for _, v := range dpiAwareness {
				dpiBoxItems.Add(dpiList[v])
			}
		} else {
			event.ConsoleWriteError("DPIAwareness i18n Configuration error: Items element invalid for winres.DPIAwareness")
		}
		m.dpiBox.SetItemIndex(int32(bean.GProject.AppOption.Windows.Manifest.DPI))
	}
	// dpi end

	// run level
	{
		runLevelBoxItems := m.runLevelBox.Items()
		executionLevel := []winres.ExecutionLevel{winres.AsInvoker, winres.HighestAvailable, winres.RequireAdministrator}
		if len(runLevelList) == len(executionLevel) {
			for _, v := range executionLevel {
				runLevelBoxItems.Add(runLevelList[v])
			}
		} else {
			event.ConsoleWriteError("ExecutionLevel i18n Configuration error: Items element invalid for winres.ExecutionLevel")
		}
		m.runLevelBox.SetItemIndex(int32(bean.GProject.AppOption.Windows.Manifest.RunLevel))
	}
	// run level end
}

// saveOrUpdateWindowsManifest 用于生成并保存 Windows 平台所需的资源文件（如图标、版本信息等）。
// 该函数会根据项目配置生成不同尺寸的 ICO 图标，并将应用程序的元数据（如版本号、版权信息等）
// 写入到 .syso 资源对象文件中，供编译时链接使用。
//
//	无显式参数，依赖全局变量：
//	  - gProject: 当前项目的配置信息
//	  - gPath: 项目根路径
//	 调用时机: 创建项目/修改项目配置时
func saveOrUpdateWindowsManifest() {
	// 图标
	iconData := bean.GProject.AppOption.Icon
	if iconData.Data == nil {
		// 使用默认图标
		iconData.Data = resources.Images("icons/window-icon_256x256.png")
		iconData.W = 256
		iconData.H = 256
	}

	var err error
	// 图标转为 ico 集合: [256, 128, 64, 48, 32, 16]
	icoSetBuf := tool.Buffer{}
	err = winicon.GenerateIcon(bytes.NewBuffer(iconData.Data), &icoSetBuf, []int{256, 128, 64, 48, 32, 16})
	if err != nil {
		event.ConsoleWriteError("windows Application Configuration - Save Configuration GenerateIcon: ", err.Error())
		return
	}

	rs := &winres.ResourceSet{}

	ico, err := winres.LoadICO(bytes.NewReader(icoSetBuf.Bytes()))
	if err != nil {
		event.ConsoleWriteError("windows Application Configuration - Save Configuration LoadICO: ", err.Error())
		return
	}
	err = rs.SetIcon(winres.RT_ICON, ico)
	if err != nil {
		event.ConsoleWriteError("windows Application Configuration - Save Configuration SetIcon: ", err.Error())
		return
	}
	rs.SetManifest(NewManifest())

	// 文件版本信息
	v := version.Info{}
	v.ProductVersion = AppVersionNum(bean.GProject.AppOption.Version)
	v.FileVersion = AppVersionNum(bean.GProject.AppOption.Version)
	v.Flags.SpecialBuild = true
	v.Timestamp = time.Now()
	// langID: 2052(中文) 1033(英语)
	v.Set(2052, version.CompanyName, bean.GProject.AppOption.Id)
	v.Set(2052, version.ProductName, bean.GProject.AppOption.Title)
	v.Set(2052, version.LegalCopyright, bean.GProject.AppOption.Copyright)
	v.Set(2052, version.FileDescription, bean.GProject.AppOption.Desc)
	v.Set(2052, version.ProductVersion, bean.GProject.AppOption.Version)
	v.Set(2052, version.FileVersion, bean.GProject.AppOption.Version)
	v.Set(2052, version.Comments, bean.GProject.AppOption.Desc)
	rs.SetVersionInfo(v)

	// 保存到 resource 目录
	resourcesPath := bean.ResourceMetadataPath()
	for _, arch := range []winres.Arch{winres.ArchAMD64 /*winres.ArchARM64,*/, winres.ArchI386} {
		sysoOutBuf := tool.Buffer{}
		err = rs.WriteObject(&sysoOutBuf, arch)
		if err != nil {
			event.ConsoleWriteError("windows Application Configuration - Save Configuration WriteObject: ", err.Error())
			return
		}
		sysoOutFile := fmt.Sprintf("%s-windows_%v.syso", bean.GProject.Name, arch)
		// 保存到项目的 resources 目录
		err = os.WriteFile(filepath.Join(resourcesPath, sysoOutFile), sysoOutBuf.Bytes(), 0644)
		if err != nil {
			event.ConsoleWriteError("windows Application Configuration - Save Configuration WriteFile: ", err.Error())
		}
	}
}

// NewManifest 创建一个新的 Windows 应用程序清单配置
// 该函数初始化并返回一个 winres.AppManifest 结构体，其中填充了来自全局项目配置的应用程序元数据
// 和 Windows 特定的清单设置。
func NewManifest() winres.AppManifest {
	return winres.AppManifest{
		Identity: winres.AssemblyIdentity{
			Name:    bean.GProject.AppOption.Id,
			Version: AppVersionNum(bean.GProject.AppOption.Version),
		},
		Description:                       bean.GProject.AppOption.Desc,
		UIAccess:                          bean.GProject.AppOption.Windows.Manifest.UIAccess,
		AutoElevate:                       bean.GProject.AppOption.Windows.Manifest.AutoElevate,
		DisableTheming:                    bean.GProject.AppOption.Windows.Manifest.DisableTheming,
		DisableWindowFiltering:            bean.GProject.AppOption.Windows.Manifest.DisableWindowFiltering,
		HighResolutionScrollingAware:      bean.GProject.AppOption.Windows.Manifest.HighResolutionScrollingAware,
		UltraHighResolutionScrollingAware: bean.GProject.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware,
		LongPathAware:                     bean.GProject.AppOption.Windows.Manifest.LongPathAware,
		PrinterDriverIsolation:            bean.GProject.AppOption.Windows.Manifest.PrinterDriverIsolation,
		GDIScaling:                        bean.GProject.AppOption.Windows.Manifest.GDIScaling,
		SegmentHeap:                       bean.GProject.AppOption.Windows.Manifest.SegmentHeap,
		UseCommonControlsV6:               bean.GProject.AppOption.Windows.Manifest.UseCommonControlsV6,
		ExecutionLevel:                    winres.ExecutionLevel(bean.GProject.AppOption.Windows.Manifest.RunLevel),
		Compatibility:                     winres.SupportedOS(bean.GProject.AppOption.Windows.Manifest.CompatibilityOS),
		DPIAwareness:                      winres.DPIAwareness(bean.GProject.AppOption.Windows.Manifest.DPI),
	}
}
