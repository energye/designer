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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"strconv"
)

// 项目(应用)配置

// 运行项目(应用)配置窗口
func runConfigApp() {
	// 显示运行项目(应用)配置窗口
	lcl.RunOnMainThreadAsync(func(id uint32) {
		form := NewConfigProjectForm()
		form.ShowModal()
		//form.Show()
	})
}

// saveProjectConfig 保存项目配置信息
// 该函数将表单中的项目配置信息保存到全局项目配置对象中，
// 并异步写入到项目配置文件中。
func (m *TConfigProjectForm) saveProjectConfig() {
	// 项目配置对象
	gProject.AppOption.Title = m.AppTitle()
	gProject.AppOption.Id = m.AppId()
	gProject.AppOption.Desc = m.AppDesc()
	gProject.AppOption.Copyright = m.AppCopyright()
	gProject.AppOption.Version = m.AppVersion()
	gProject.AppOption.Icon = m.appIconData
	gProject.AppOption.Windows.Manifest.CompatibilityOS = m.compatibilityOSBox.ItemIndex()
	gProject.AppOption.Windows.Manifest.DPI = m.dpiBox.ItemIndex()
	gProject.AppOption.Windows.Manifest.RunLevel = m.runLevelBox.ItemIndex()
	gProject.AppOption.Windows.Manifest.UIAccess = m.uiAccessCheckBox.Checked()
	gProject.AppOption.Windows.Manifest.AutoElevate = m.autoElevateBox.Checked()
	gProject.AppOption.Windows.Manifest.DisableTheming = m.disableThemingBox.Checked()
	gProject.AppOption.Windows.Manifest.DisableWindowFiltering = m.disableWindowFilteringBox.Checked()
	gProject.AppOption.Windows.Manifest.HighResolutionScrollingAware = m.highResolutionScrollingAwareBox.Checked()
	gProject.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware = m.ultraHighResolutionScrollingAwareBox.Checked()
	gProject.AppOption.Windows.Manifest.LongPathAware = m.longPathAwareBox.Checked()
	gProject.AppOption.Windows.Manifest.PrinterDriverIsolation = m.printerDriverIsolationBox.Checked()
	gProject.AppOption.Windows.Manifest.GDIScaling = m.gDIScalingBox.Checked()
	gProject.AppOption.Windows.Manifest.SegmentHeap = m.segmentHeapBox.Checked()
	gProject.AppOption.Windows.Manifest.UseCommonControlsV6 = m.useCommonControlsV6Box.Checked()
	go func() {
		// 更新项目配置文件
		if err := WriteEGPConfig(gPath, gProject); err != nil {
			logs.Error("保存-写入项目配置文件失败")
			return
		}
	}()
}

func AppVersionNum(version string) [4]uint16 {
	versionNum := [4]uint16{0, 0, 0, 0}
	for i, v := range tool.Split(version, ".") {
		if i < len(versionNum) {
			vn, _ := strconv.ParseUint(v, 10, 16)
			versionNum[i] = uint16(vn)
		}
	}
	return versionNum
}

// 更新窗口图标
// Windows	16×16	24×24/32×32	.ico	16×16、24×24、32×32
// macOS	16×16	32×32		.icns		16×16、32×32
// Linux	16×16	24×24		.png		16×16、24×24
func updateWindowICON() {
	icon := gProject.AppOption.Icon
	if icon.Data == nil || icon.W <= 0 || icon.H <= 0 {
		logs.Error("更新窗口图标, 图标数据不能为空/大小不正确")
		return
	}

}
