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

//go:embed linuxdeploy-plugin-gtk.sh
var gtkPluginScript []byte // for: https://github.com/linuxdeploy/linuxdeploy-plugin-gtk

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
	//defer os.RemoveAll(buildDir)

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

	// 写出 GTK 插件脚本
	pluginPath := filepath.Join(buildDir, "linuxdeploy-plugin-gtk.sh")
	if err := os.WriteFile(pluginPath, gtkPluginScript, 0755); err != nil {
		event.ConsoleWriteError("Package - AppImage: write GTK plugin failed:", err.Error())
		return false
	}

	// 自动检测 GTK 版本
	gtkVersion, webkitVersion := detectedVersion()

	event.ConsoleWriteInfo("Package - AppImage: detected GTK version", gtkVersion, "Webkit2Gtk version", webkitVersion)

	// 检查是否需要禁用 strip
	needNoStrip := hasRelrDynSections()
	if needNoStrip {
		event.ConsoleWriteInfo("Package - AppImage: detected .relr.dyn sections, disabling strip")
	}

	// 构建输出文件名
	goarch := lib.GOARCH()
	debArch := debArchName(goarch)
	appImageName := fmt.Sprintf("%s_%s_%s.AppImage", option.PackageName, appOption.Version, debArch)
	if m.AppendPlatform {
		appImageName = fmt.Sprintf("%s_%s_%s_%s.AppImage", option.PackageName, appOption.Version, lib.GOOS(), debArch)
	}
	appImagePath := filepath.Join(output, appImageName)
	_ = os.Remove(appImagePath)

	// 准备环境变量（清除插件模式残留，设置 GTK 版本和输出文件名）

	cmdEnv := append(os.Environ(),
		fmt.Sprintf("DEPLOY_GTK_VERSION=%s", gtkVersion),
		fmt.Sprintf("OUTPUT=%s", appImageName),
	)
	if needNoStrip {
		cmdEnv = append(cmdEnv, "NO_STRIP=1")
	}

	// 执行 linuxdeploy（带 GTK 插件，一次性完成依赖部署和打包）
	event.ConsoleWriteInfo("Package - AppImage: building AppImage with linuxdeploy (GTK plugin)...")
	cmd := exec.Command(linuxdeployPath, "--appimage-extract-and-run", "--appdir", appDir, "--output", "appimage", "--plugin", "gtk")
	cmd.Dir = buildDir
	cmd.Env = cmdEnv

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		event.ConsoleWriteError("Package - AppImage: linuxdeploy failed:", err.Error())
		event.ConsoleWriteError("Output:", string(outputBytes))
		return false
	}

	// 移动生成的 AppImage 到输出目录
	srcAppImage := filepath.Join(buildDir, appImageName)
	if !tool.IsExist(srcAppImage) {
		event.ConsoleWriteError("Package - AppImage: generated AppImage not found:", srcAppImage)
		return false
	}
	if err := os.Rename(srcAppImage, appImagePath); err != nil {
		event.ConsoleWriteError("Package - AppImage: move failed:", err.Error())
		return false
	}

	_ = os.RemoveAll(appDir)
	event.ConsoleWriteInfo("Package - AppImage:", appImageName, "created successfully")
	return true
}

// detectedVersion 根据项目配置获取 GTK 版本和依赖信息
// 返回: gtkVersion, webkitVersion
// gtkVersion: "2" 或 "3"
// webkitVersion: "" (不需要), "4.0", "4.1"
func detectedVersion() (string, string) {
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

// hasRelrDynSections 检查系统库是否包含 .relr.dyn 段（现代工具链）
func hasRelrDynSections() bool {
	testLibs := []string{
		"/usr/lib/libgtk-4.so.1",
		"/usr/lib64/libgtk-4.so.1",
		"/usr/lib/x86_64-linux-gnu/libgtk-4.so.1",
		"/usr/lib/aarch64-linux-gnu/libgtk-4.so.1",
		"/usr/lib/libgtk-3.so.0",
		"/usr/lib64/libgtk-3.so.0",
		"/usr/lib/x86_64-linux-gnu/libgtk-3.so.0",
		"/usr/lib/aarch64-linux-gnu/libgtk-3.so.0",
	}
	for _, lib := range testLibs {
		if _, err := os.Stat(lib); err == nil {
			cmd := exec.Command("readelf", "-S", lib)
			output, err := cmd.Output()
			if err == nil && strings.Contains(string(output), ".relr.dyn") {
				return true
			}
		}
	}
	return false
}
