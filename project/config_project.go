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
	"bytes"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project/bean"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"image/png"
	"os"
	"path/filepath"
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
	{
		// windows manifest
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
	}
	{
		// macOS plist.info
		gProject.AppOption.MacOS.Plist.CFBundleExecutable = m.AppBundleExecutable() // 从构建配置里获取
		gProject.AppOption.MacOS.Plist.CFBundleName = m.AppBundleName()
		gProject.AppOption.MacOS.Plist.CFBundleDisplayName = gProject.AppOption.Title
		gProject.AppOption.MacOS.Plist.CFBundleLocalizations = m.AppBundleLocalizations()
		gProject.AppOption.MacOS.Plist.CFBundleIdentifier = gProject.AppOption.Id
		gProject.AppOption.MacOS.Plist.CFBundleVersion = gProject.AppOption.Version
		gProject.AppOption.MacOS.Plist.CFBundleShortVersionString = gProject.AppOption.Version
		gProject.AppOption.MacOS.Plist.CFBundleGetInfoString = gProject.AppOption.Desc
		gProject.AppOption.MacOS.Plist.CFBundleIconFile = gProject.Name + ".icns"
		gProject.AppOption.MacOS.Plist.NSHumanReadableCopyright = gProject.AppOption.Copyright
		gProject.AppOption.MacOS.Plist.LSUIElement = m.AppLSUIElement()
		gProject.AppOption.MacOS.Plist.LSMinimumSystemVersion = m.AppLSMinimumSystemVersion()
	}
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
// TODO ?? 下面是标准尺寸
// Windows	16×16	24×24/32×32	.ico	16×16、24×24、32×32
// macOS	16×16	32×32		.icns		16×16、32×32
// Linux	16×16	24×24		.png		16×16、24×24
func updateWindowICON() {
	icon := gProject.AppOption.Icon
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
