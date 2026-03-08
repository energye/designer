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
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

// 项目(应用)配置

// 运行项目(应用)配置窗口
func runAppConfig() {
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
	bean.GProject.AppOption.Title = m.AppTitle()
	bean.GProject.AppOption.Id = m.AppId()
	bean.GProject.AppOption.Desc = m.AppDesc()
	bean.GProject.AppOption.Copyright = m.AppCopyright()
	bean.GProject.AppOption.Version = m.AppVersion()
	bean.GProject.AppOption.Icon = m.appIconData
	{
		// windows manifest
		bean.GProject.AppOption.Windows.Manifest.CompatibilityOS = m.compatibilityOSBox.ItemIndex()
		bean.GProject.AppOption.Windows.Manifest.DPI = m.dpiBox.ItemIndex()
		bean.GProject.AppOption.Windows.Manifest.RunLevel = m.runLevelBox.ItemIndex()
		bean.GProject.AppOption.Windows.Manifest.UIAccess = m.uiAccessCheckBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.AutoElevate = m.autoElevateBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.DisableTheming = m.disableThemingBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.DisableWindowFiltering = m.disableWindowFilteringBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.HighResolutionScrollingAware = m.highResolutionScrollingAwareBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware = m.ultraHighResolutionScrollingAwareBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.LongPathAware = m.longPathAwareBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.PrinterDriverIsolation = m.printerDriverIsolationBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.GDIScaling = m.gDIScalingBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.SegmentHeap = m.segmentHeapBox.Checked()
		bean.GProject.AppOption.Windows.Manifest.UseCommonControlsV6 = m.useCommonControlsV6Box.Checked()
	}
	{
		// macOS plist.info
		bean.GProject.AppOption.MacOS.PList.CFBundleExecutable = m.AppBundleExecutable() // 从构建配置里获取
		bean.GProject.AppOption.MacOS.PList.CFBundleName = m.AppBundleName()
		bean.GProject.AppOption.MacOS.PList.CFBundleDisplayName = bean.GProject.AppOption.Title
		bean.GProject.AppOption.MacOS.PList.CFBundleLocalizations = m.AppBundleLocalizations()
		bean.GProject.AppOption.MacOS.PList.CFBundleIdentifier = bean.GProject.AppOption.Id
		bean.GProject.AppOption.MacOS.PList.CFBundleVersion = bean.GProject.AppOption.Version
		bean.GProject.AppOption.MacOS.PList.CFBundleShortVersionString = bean.GProject.AppOption.Version
		bean.GProject.AppOption.MacOS.PList.CFBundleGetInfoString = bean.GProject.AppOption.Desc
		bean.GProject.AppOption.MacOS.PList.CFBundleIconFile = bean.GProject.Name + ".icns"
		bean.GProject.AppOption.MacOS.PList.NSHumanReadableCopyright = bean.GProject.AppOption.Copyright
		bean.GProject.AppOption.MacOS.PList.LSUIElementIndex = m.LSUIElementBox.ItemIndex()
		bean.GProject.AppOption.MacOS.PList.LSUIElement = m.AppLSUIElement()
		bean.GProject.AppOption.MacOS.PList.LSMinimumSystemVersionIndex = m.LSMinimumSystemVersionBox.ItemIndex()
		bean.GProject.AppOption.MacOS.PList.LSMinimumSystemVersion = m.AppLSMinimumSystemVersion()
	}
	go func() {
		// 更新项目配置文件
		if err := WriteEGPConfig(bean.GPath, bean.GProject); err != nil {
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
// TODO ?? 下面是标准尺寸
// Windows	16×16	24×24/32×32	.ico	16×16、24×24、32×32
// macOS	16×16	32×32		.icns		16×16、32×32
// Linux	16×16	24×24		.png		16×16、24×24
func updateWindowICON() {
	icon := bean.GProject.AppOption.Icon
	if icon.Data == nil {
		// 使用默认图标
		icon = bean.TAppIcon{
			Data: resources.Images("icons/window-icon_256x256.png"),
			W:    256,
			H:    256,
		}
	}
	if icon.Data == nil || icon.W <= 0 || icon.H <= 0 {
		logs.Error("updateWindowICON, 图标数据不能为空/大小不正确")
		return
	}
	iconData := icon.Data // png
	// 这里没使用标准尺寸, 统一最大尺寸为: 256x256
	if icon.W > 256 || icon.H > 256 {
		iconData = tool.Scale(iconData, 256, 256)
	}
	pngBuf := new(bytes.Buffer)
	pngBuf.Write(iconData)
	pngImage, err := png.Decode(pngBuf)
	if err != nil {
		logs.Error("updateWindowICON, 图标转为 png 对象失败:", err.Error())
		return
	}
	icoBuf := new(bytes.Buffer)
	err = tool.Encode(icoBuf, pngImage)
	if err != nil {
		logs.Error("updateWindowICON, png 转为 ico 对象失败:", err.Error())
		return
	}
	// 保存图标
	// icon.ico
	embedPath := EmbedPath()
	iconIcoFilePath := filepath.Join(embedPath, "icon.ico")
	err = os.WriteFile(iconIcoFilePath, icoBuf.Bytes(), 0666)
	if err != nil {
		logs.Error("updateWindowICON, 写 windows icon.ico 失败:", err.Error())
		return
	}
	// icon.png
	iconPngFilePath := filepath.Join(embedPath, "icon.png")
	err = os.WriteFile(iconPngFilePath, iconData, 0666)
	if err != nil {
		logs.Error("updateWindowICON, 写 windows icon.ico 失败:", err.Error())
		return
	}
	// macOS
	// icon.icns
}
