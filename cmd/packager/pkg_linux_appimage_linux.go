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
	_ "embed"
	"fmt"
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed gtk-bundle.sh
var gtkBundleScript []byte

// configGTKVersion 根据项目配置获取 GTK 版本和依赖信息
// 返回: gtkVersion, webkitVersion
// gtkVersion: "2" 或 "3"
// webkitVersion: "" (不需要), "4.0", "4.1"
func configGTKVersion() (string, string) {
	proj := bean.GProject
	option := proj.BuildOption

	gtkVersion := "3"
	webkitVersion := ""

	if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_LCL) {
		if !option.UIGtk3 {
			gtkVersion = "2"
		}
	} else if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_WV) {
		// WV 框架需要 WebKit
		if option.BuildCGOEnabled && strings.Contains(option.GoArgs, "webkit2_4_1") {
			webkitVersion = "4.1"
		} else if option.BuildCGOEnabled {
			webkitVersion = "4.0"
		} else {
			// 非 CGO 默认 4.1
			webkitVersion = "4.1"
		}
	}
	// CEF 框架只需要 GTK3，不需要 WebKit

	return gtkVersion, webkitVersion
}

// downloadLinuxdeploy 下载 linuxdeploy 工具
func downloadLinuxdeploy(arch, buildDir string) (string, error) {
	path := filepath.Join(buildDir, fmt.Sprintf("linuxdeploy-%s.AppImage", arch))
	if tool.IsExist(path) {
		return path, nil
	}

	url := fmt.Sprintf("https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-%s.AppImage", arch)
	event.ConsoleWriteInfo("Package - AppImage: downloading linuxdeploy...")

	cmd := exec.Command("curl", "-L", "-o", path, url)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("download linuxdeploy failed: %s\n%s", err, string(output))
	}
	if err := os.Chmod(path, 0755); err != nil {
		return "", err
	}
	return path, nil
}

// downloadAppRun 下载 AppRun
func downloadAppRun(arch, appDir string) error {
	path := filepath.Join(appDir, "AppRun")
	if tool.IsExist(path) {
		return nil
	}

	url := fmt.Sprintf("https://github.com/AppImage/AppImageKit/releases/download/continuous/AppRun-%s", arch)
	event.ConsoleWriteInfo("Package - AppImage: downloading AppRun...")

	cmd := exec.Command("curl", "-L", "-o", path, url)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("download AppRun failed: %s\n%s", err, string(output))
	}
	return os.Chmod(path, 0755)
}

// writeGTKBundleScript 写入 GTK 打包脚本
func writeGTKBundleScript(buildDir string) (string, error) {
	path := filepath.Join(buildDir, "gtk-bundle.sh")
	if tool.IsExist(path) {
		return path, nil
	}
	if err := os.WriteFile(path, gtkBundleScript, 0755); err != nil {
		return "", err
	}
	return path, nil
}

// appImage 构建 AppImage 包
func (m *Package) appImage() bool {
	event.ConsoleWriteInfo("Package - AppImage")

	// AppImage 不支持跨架构构建
	targetArch := lib.GOARCH()
	if runtime.GOARCH != targetArch {
		event.ConsoleWriteWarn("Package - AppImage: skipped, cross-architecture build is not supported. Current:", runtime.GOARCH, ", Target:", targetArch)
		return true
	}

	// 检查必要工具
	if !checkToolCMD("curl") {
		event.ConsoleWriteError("Package - AppImage: curl not found")
		return false
	}

	proj := bean.GProject
	option := proj.BuildOption
	appOption := proj.AppOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}

	// 处理 libenergy.so 运行时库
	defer env.Delete("ENERGY_WS")
	choiceLibEnergySOWS()

	// 架构
	arch := "x86_64"
	if runtime.GOARCH == "arm64" {
		arch = "aarch64"
	}

	// 创建目录
	appDir := filepath.Join(output, option.PackageName+".AppDir")
	buildDir := filepath.Join(output, ".appimage-build")
	defer os.RemoveAll(buildDir)

	for _, dir := range []string{
		filepath.Join(appDir, "usr", "bin"),
		filepath.Join(appDir, "usr", "lib"),
		buildDir,
	} {
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

	// 复制图标
	srcIcon := filepath.Join(bean.ResourceEmbedPath(), "icon.png")
	if tool.IsExist(srcIcon) {
		dstIcon := filepath.Join(appDir, ".DirIcon")
		if err := tool.CopyFile(srcIcon, dstIcon); err != nil {
			event.ConsoleWriteError("Package - AppImage: copy icon failed:", err.Error())
			return false
		}
		iconLink := filepath.Join(appDir, option.PackageName+".png")
		_ = os.Remove(iconLink)
		_ = os.Symlink(".DirIcon", iconLink)
	}

	// 生成 .desktop 文件
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

	// 复制运行时库 libenergy.so
	srcLibName := lib.GetDLLName()
	dstLibName := packageLibName()
	srcLib := filepath.Join(config.Config.FrameworkRuntimePath(), srcLibName)
	if !tool.IsExist(srcLib) {
		event.ConsoleWriteError("Package - AppImage: runtime library not found:", srcLib)
		return false
	}
	dstLib := filepath.Join(appDir, "usr", "lib", dstLibName)
	if err := tool.CopyFile(srcLib, dstLib); err != nil {
		event.ConsoleWriteError("Package - AppImage: copy libenergy failed:", err.Error())
		return false
	}

	// 下载 linuxdeploy
	linuxdeployPath, err := downloadLinuxdeploy(arch, buildDir)
	if err != nil {
		event.ConsoleWriteError("Package - AppImage:", err.Error())
		return false
	}

	// 下载 AppRun
	if err := downloadAppRun(arch, appDir); err != nil {
		event.ConsoleWriteError("Package - AppImage:", err.Error())
		return false
	}

	// 获取 GTK 版本和依赖信息
	gtkVersion, webkitVersion := configGTKVersion()
	event.ConsoleWriteInfo("Package - AppImage: GTK version", gtkVersion, "WebKit version", webkitVersion)

	// 写入 GTK 打包脚本
	gtkScriptPath, err := writeGTKBundleScript(buildDir)
	if err != nil {
		event.ConsoleWriteError("Package - AppImage: write GTK script failed:", err.Error())
		return false
	}

	// 运行 GTK 打包脚本
	event.ConsoleWriteInfo("Package - AppImage: running GTK bundle script...")
	gtkCmd := exec.Command(gtkScriptPath, appDir, gtkVersion, webkitVersion)
	gtkCmd.Env = append(os.Environ(), "LINUXDEPLOY="+linuxdeployPath)
	gtkOutput, err := gtkCmd.CombinedOutput()
	if err != nil {
		event.ConsoleWriteError("Package - AppImage: GTK bundle failed:", err.Error())
		event.ConsoleWriteError("Output:", string(gtkOutput))
		return false
	}
	event.ConsoleWriteInfo("GTK bundle output:", string(gtkOutput))

	// 运行 linuxdeploy 生成 AppImage
	goarch := lib.GOARCH()
	debArch := debArchName(goarch)
	appImageName := fmt.Sprintf("%s_%s_%s.AppImage", option.PackageName, appOption.Version, debArch)
	if m.AppendPlatform {
		appImageName = fmt.Sprintf("%s_%s_%s_%s.AppImage", option.PackageName, appOption.Version, lib.GOOS(), debArch)
	}
	appImagePath := filepath.Join(output, appImageName)
	_ = os.Remove(appImagePath)

	event.ConsoleWriteInfo("Package - AppImage: building AppImage...")
	cmd := exec.Command(linuxdeployPath, "--appimage-extract-and-run", "--appdir", appDir, "--output", "appimage")
	cmd.Dir = buildDir
	cmd.Env = append(os.Environ(), "OUTPUT="+appImageName)
	linuxdeployOutput, err := cmd.CombinedOutput()
	if err != nil {
		event.ConsoleWriteError("Package - AppImage: build failed:", err.Error())
		event.ConsoleWriteError("Output:", string(linuxdeployOutput))
		return false
	}

	// 移动到输出目录
	srcAppImage := filepath.Join(buildDir, appImageName)
	if tool.IsExist(srcAppImage) {
		if err := os.Rename(srcAppImage, appImagePath); err != nil {
			event.ConsoleWriteError("Package - AppImage: move failed:", err.Error())
			return false
		}
	}

	_ = os.RemoveAll(appDir)
	event.ConsoleWriteInfo("Package - AppImage:", appImageName, "created successfully")
	return true
}
