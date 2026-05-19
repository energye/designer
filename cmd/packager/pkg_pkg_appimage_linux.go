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

//go:build linux

package packager

import (
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
	"path/filepath"
)

// appImage 构建 AppImage 包
func (m *Package) appImage() bool {
	event.ConsoleWriteInfo("Package - AppImage")
	if !checkToolCMD("appimagetool") {
		event.ConsoleWriteError("Package - AppImage: appimagetool command not found. Download from: https://github.com/AppImage/AppImageKit/releases")
		return false
	}

	proj := bean.GProject
	option := proj.BuildOption
	appOption := proj.AppOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}

	// 创建 AppDir 目录结构
	appDir := filepath.Join(output, option.PackageName+".AppDir")
	dirs := []string{
		filepath.Join(appDir, "usr", "bin"),
		filepath.Join(appDir, "usr", "lib"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			event.ConsoleWriteError("Package - AppImage: mkdir failed:", err.Error())
			return false
		}
	}

	// 复制可执行文件
	srcBin := filepath.Join(output, option.BuildFileName)
	dstBin := filepath.Join(appDir, "usr", "bin", option.PackageName)
	if err := tool.CopyFile(srcBin, dstBin); err != nil {
		event.ConsoleWriteError("Package - AppImage: copy binary failed:", err.Error())
		return false
	}
	if err := os.Chmod(dstBin, 0755); err != nil {
		event.ConsoleWriteError("Package - AppImage: chmod binary failed:", err.Error())
		return false
	}

	// 生成 AppRun 脚本
	appRunTemplate := app.Packager("linux/AppRun")
	if appRunTemplate == nil {
		event.ConsoleWriteError("Package - AppImage: AppRun template not found")
		return false
	}
	appRunData, err := tool.RenderTemplate(string(appRunTemplate), map[string]string{
		"PackageName": option.PackageName,
	})
	if err != nil {
		event.ConsoleWriteError("Package - AppImage: render AppRun failed:", err.Error())
		return false
	}
	appRunPath := filepath.Join(appDir, "AppRun")
	if err := os.WriteFile(appRunPath, appRunData, 0755); err != nil {
		event.ConsoleWriteError("Package - AppImage: write AppRun failed:", err.Error())
		return false
	}
	if err := os.Chmod(appRunPath, 0755); err != nil {
		event.ConsoleWriteError("Package - AppImage: chmod AppRun failed:", err.Error())
		return false
	}

	// 生成 .desktop 文件 (放在 AppDir 根目录)
	desktopData := renderDesktopFile(option.PackageName, option.PackageName, option.PackageName,
		appOption.Desc, option.PackageName, appOption.Linux.Categories)
	if desktopData == nil {
		event.ConsoleWriteError("Package - AppImage: render desktop failed")
		return false
	}
	desktopPath := filepath.Join(appDir, option.PackageName+".desktop")
	if err := os.WriteFile(desktopPath, desktopData, 0644); err != nil {
		event.ConsoleWriteError("Package - AppImage: write desktop failed:", err.Error())
		return false
	}

	// 复制图标 (放在 AppDir 根目录)
	srcIcon := filepath.Join(bean.ResourceEmbedPath(), "icon.png")
	if tool.IsExist(srcIcon) {
		dstIcon := filepath.Join(appDir, option.PackageName+".png")
		if err := tool.CopyFile(srcIcon, dstIcon); err != nil {
			event.ConsoleWriteError("Package - AppImage: copy icon failed:", err.Error())
			return false
		}
	}

	// 复制运行时库 libenergy.so
	libName := lib.GetDLLName()
	srcLib := filepath.Join(config.Config.FrameworkRuntimePath(), libName)
	if tool.IsExist(srcLib) {
		dstLib := filepath.Join(appDir, "usr", "lib", libName)
		if err := tool.CopyFile(srcLib, dstLib); err != nil {
			event.ConsoleWriteError("Package - AppImage: copy libenergy failed:", err.Error())
			return false
		}
	} else {
		event.ConsoleWriteWarn("Package - AppImage: runtime library not found:", srcLib)
	}

	// 构建 AppImage
	goarch := lib.GOARCH()
	debArch := debArchName(goarch)
	appImageName := fmt.Sprintf("%s_%s_%s.AppImage", option.PackageName, appOption.Version, debArch)
	if m.AppendPlatform {
		appImageName = fmt.Sprintf("%s_%s_%s_%s.AppImage", option.PackageName, appOption.Version, lib.GOOS(), debArch)
	}
	appImagePath := filepath.Join(output, appImageName)
	_ = os.Remove(appImagePath)

	cmd := RunCMD(output, "appimagetool", "--no-appstream", appDir, appImagePath)
	if cmd != nil {
		event.ConsoleWriteError("Package - AppImage: build failed:", cmd.Error())
		return false
	}

	// 清理临时目录
	_ = os.RemoveAll(appDir)

	event.ConsoleWriteInfo("Package - AppImage:", appImageName, "created successfully")
	return true
}
