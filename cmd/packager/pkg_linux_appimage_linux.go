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
	"github.com/energye/lcl/tool/command"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// https://github.com/linuxdeploy/linuxdeploy-plugin-gtk
//
//	fix: im-ibus.so load fail.
//	 "cat >> "$HOOKFILE" <<EOF
//	 export LD_LIBRARY_PATH="\$APPDIR//usr/lib:\$LD_LIBRARY_PATH"
//	 EOF"
//
//go:embed linuxdeploy-plugin-gtk.sh
var gtkPluginScript []byte

// downloadAppImageBuildTool 下载并准备 AppImage 构建工具
//
// 该函数负责下载和配置构建 AppImage 包所需的三个关键组件：
// 1. linuxdeploy AppImage 工具（用于打包）
// 2. AppRun 运行时文件（AppImage 的启动入口）
// 3. GTK 插件脚本（用于支持 GTK 应用）
//
// 如果文件已存在则跳过下载，避免重复操作。
//
// 返回值:
//   - linuxdeploy: linuxdeploy AppImage 工具的完整路径
//   - appRun: AppRun 运行时文件的完整路径
//   - err: 错误信息，如果任何步骤失败则返回错误
func downloadAppImageBuildTool() (buildDir, linuxdeploy, appRun string, err error) {
	buildDir = filepath.Join(config.Config.FrameworkDir, "appimage-build")
	_ = os.MkdirAll(buildDir, 0755)

	arch := platformArch()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		path := filepath.Join(buildDir, fmt.Sprintf("linuxdeploy-%s.AppImage", arch))
		if !tool.IsExist(path) {
			url := fmt.Sprintf("https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-%s.AppImage", arch)
			event.ConsoleWriteInfo("Package - AppImage downloading linuxdeploy...")
			err = tool.Download(url, path)
			if err != nil {
				event.ConsoleWriteInfo("Package - AppImage downloading linuxdeploy failed:", err.Error())
				goto done
			}
			err = os.Chmod(path, 0755)
			if err != nil {
				event.ConsoleWriteInfo("Package - AppImage downloading linuxdeploy . chmod execution permission failed:", err.Error())
				goto done
			}
		}
		linuxdeploy = path
	done:
		wg.Done()
	}()
	go func() {
		appRunName := fmt.Sprintf("AppRun-%s", arch)
		path := filepath.Join(buildDir, fmt.Sprintf("AppRun-%s", arch))
		if !tool.IsExist(path) {
			url := fmt.Sprintf("https://github.com/AppImage/AppImageKit/releases/download/continuous/%s", appRunName)
			event.ConsoleWriteInfo("Package - AppImage downloading", appRunName+"...")
			err = tool.Download(url, path)
			if err != nil {
				event.ConsoleWriteInfo("Package - AppImage downloading", appRunName, " failed:", err.Error())
				goto done
			}
			err = os.Chmod(path, 0755)
			if err != nil {
				event.ConsoleWriteInfo("Package - AppImage downloading", appRunName, ". chmod execution permission failed:", err.Error())
				goto done
			}
		}
		appRun = path
	done:
		wg.Done()
	}()
	pluginPath := filepath.Join(buildDir, "linuxdeploy-plugin-gtk.sh")
	err = os.WriteFile(pluginPath, gtkPluginScript, 0755)
	wg.Wait()
	return
}

// appImageWebkit2GtkDeps 获取 AppImage 打包所需的 WebKit2GTK 依赖库列表
//
// 该函数用于收集构建 AppImage 包时需要的 WebKit2GTK 相关动态库依赖，
// 包括 webkit2gtk-4.0 和 webkit2gtk-4.1 两个版本。
// 通过 findDeps 函数自动查找这些库的实际安装路径。
//
// 返回值:
//   - []*AppImageDepsInstall: 包含依赖库信息的切片，每个元素包含库名称、库文件目录和完整路径
func appImageWebkit2GtkDeps(version string) []*AppImageDepsInstall {
	// /usr/lib/x86_64-linux-gnu/webkit2gtk-4.0/
	// /usr/lib64/webkit2gtk-4.0/
	// injected-bundle/libwebkit2gtkinjectedbundle.so
	// WebKitGPUProcess  WebKitNetworkProcess  WebKitWebProcess
	webkit2gtk4_x := fmt.Sprintf("webkit2gtk-%s", version)
	deps := []*AppImageDepsInstall{{Name: webkit2gtk4_x}}
	findDeps(deps)
	dep := deps[0]
	if dep.LibDir == "" {
		return nil
	}

	webkit2gtk := filepath.Join(dep.LibDir, webkit2gtk4_x)

	deps = []*AppImageDepsInstall{}
	deps = append(deps, &AppImageDepsInstall{
		Name:    "WebKitGPUProcess",
		LibDir:  webkit2gtk,
		LibPath: filepath.Join(webkit2gtk, "WebKitGPUProcess"),
	})
	deps = append(deps, &AppImageDepsInstall{
		Name:    "WebKitNetworkProcess",
		LibDir:  webkit2gtk,
		LibPath: filepath.Join(webkit2gtk, "WebKitNetworkProcess"),
	})
	deps = append(deps, &AppImageDepsInstall{
		Name:    "WebKitWebProcess",
		LibDir:  webkit2gtk,
		LibPath: filepath.Join(webkit2gtk, "WebKitWebProcess"),
	})
	deps = append(deps, &AppImageDepsInstall{
		Name:    "libwebkit2gtkinjectedbundle.so",
		LibDir:  filepath.Join(webkit2gtk, "injected-bundle"),
		LibPath: filepath.Join(webkit2gtk, "injected-bundle/libwebkit2gtkinjectedbundle.so"),
	})
	return deps
}

// appImageDeps 获取 AppImage 打包所需的基础依赖库列表
//
// 该函数用于收集构建 AppImage 包时需要的基础动态库依赖，
// 通过 findDeps 函数自动查找这些库的实际安装路径。
//
// 返回值:
//   - []*AppImageDepsInstall: 包含依赖库信息的切片，每个元素包含库名称、库文件目录和完整路径
func appImageDeps() []*AppImageDepsInstall {
	deps := []*AppImageDepsInstall{
		{Name: "libGLESv2.so.2"},
		{Name: "libharfbuzz-gobject.so.0"},
		{Name: "libGL.so.1"},
	}
	findDeps(deps)
	return deps
}

// appImage 构建 AppImage 包
func (m *Package) appImage() bool {
	event.ConsoleWriteInfo("Package - AppImage")

	// AppImage 不支持跨架构构建
	targetArch := lib.GOARCH()
	if runtime.GOARCH != targetArch {
		event.ConsoleWriteWarn("Package - AppImage skipped, cross-architecture build is not supported. Current:", runtime.GOARCH, ", Target:", targetArch)
		return true
	}

	buildDir, linuxdeployPath, appRunPath, err := downloadAppImageBuildTool()
	if err != nil {
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
	choiceLibEnergySOWS()
	defer env.Delete("ENERGY_WS")

	// 创建 app dir 目录
	appDir := filepath.Join(output, option.PackageName+".AppDir")
	for _, dir := range []string{
		filepath.Join(appDir, "usr", "bin"),
		filepath.Join(appDir, "usr", "lib"),
	} {
		if err = os.MkdirAll(dir, 0755); err != nil {
			event.ConsoleWriteError("Package - AppImage mkdir failed:", err.Error())
			return false
		}
	}

	defer os.RemoveAll(appDir)

	// 复制 AppRun
	if err = tool.CopyFile(appRunPath, filepath.Join(appDir, "AppRun")); err != nil {
		event.ConsoleWriteError("Package - AppImage copy AppRun failed:", err.Error())
		return false
	}

	// 复制可执行文件
	srcBin := filepath.Join(output, option.BuildFileName)
	dstBin := filepath.Join(appDir, "usr", "bin", option.BuildFileName)
	_ = os.Remove(dstBin)
	if err = tool.CopyFile(srcBin, dstBin); err != nil {
		event.ConsoleWriteError("Package - AppImage copy binary failed:", err.Error())
		return false
	}
	if err = os.Chmod(dstBin, 0755); err != nil {
		event.ConsoleWriteError("Package - AppImage chmod binary failed:", err.Error())
		return false
	}

	// 复制图标
	srcIcon := filepath.Join(bean.ResourceEmbedPath(), "icon.png")
	if tool.IsExist(srcIcon) {
		dstIcon := filepath.Join(appDir, "icon.png")
		if err = tool.CopyFile(srcIcon, dstIcon); err != nil {
			event.ConsoleWriteError("Package - AppImage copy icon failed:", err.Error())
			return false
		}
	}

	// 生成 .desktop
	desktopData := renderDesktopFile(proj.Name, option.BuildFileName, "icon",
		appOption.Desc, option.PackageName, appOption.Linux.Categories)
	if desktopData == nil {
		event.ConsoleWriteError("Package - AppImage render desktop failed")
		return false
	}
	// 复制 .desktop
	desktopPath := filepath.Join(appDir, option.PackageName+".desktop")
	if err = os.WriteFile(desktopPath, desktopData, 0644); err != nil {
		event.ConsoleWriteError("Package - AppImage write desktop failed:", err.Error())
		return false
	}

	// 复制运行时库 libenergy.so
	srcLibName := lib.GetDLLName()
	dstLibName := packageLibName()
	srcLib := filepath.Join(config.Config.FrameworkRuntimePath(), srcLibName)
	if !tool.IsExist(srcLib) {
		event.ConsoleWriteError("Package - AppImage runtime library not found:", srcLib)
		return false
	}
	dstLib := filepath.Join(appDir, "usr", "lib", dstLibName)
	if err = tool.CopyFile(srcLib, dstLib); err != nil {
		event.ConsoleWriteError("Package - AppImage copy libenergy failed:", err.Error())
		return false
	}

	// 检测版本
	gtkVersion, webkitVersion := detectedVersion()

	// 使用 webkit2gtk
	if webkitVersion != "" {
		// 复制 webkit2gtk 运行时库
		wvDeps := appImageWebkit2GtkDeps(webkitVersion)
		// webkit2gtk 有 4.0 和 4.1, 根据 webkitVersion 选择
		if wvDeps == nil {
			event.ConsoleWriteError("Package - AppImage copy webkit2gtk runtime not found", webkitVersion)
			return false
		}
		for _, wvDep := range wvDeps {
			if err = tool.CopyFile(wvDep.LibPath, filepath.Join(appDir, wvDep.LibPath)); err != nil {
				event.ConsoleWriteError("Package - AppImage copy webkit2gtk runtime failed:", err.Error())
				return false
			}
		}
	}
	// 复制其他依赖库
	othDeps := appImageDeps()
	for _, dep := range othDeps {
		if err = tool.CopyFile(dep.LibPath, filepath.Join(appDir, dep.LibPath)); err != nil {
			event.ConsoleWriteError("Package - AppImage copy dep runtime failed:", err.Error())
			return false
		}
	}

	event.ConsoleWriteInfo("Package - AppImage detected GTK-v:", gtkVersion, "webkit2gtk-v:", webkitVersion)

	// 构建输出文件名
	goarch := lib.GOARCH()
	debArch := debArchName(goarch)
	appImageName := fmt.Sprintf("%s_%s_%s.AppImage", option.PackageName, appOption.Version, debArch)
	if m.AppendPlatform {
		appImageName = fmt.Sprintf("%s_%s_%s_%s.AppImage", option.PackageName, appOption.Version, lib.GOOS(), debArch)
	}
	appImagePath := filepath.Join(output, appImageName)
	_ = os.Remove(appImagePath)

	// 执行 linuxdeploy（带 GTK 插件，一次性完成依赖部署和打包）
	event.ConsoleWriteInfo("Package - AppImage building AppImage with linuxdeploy...")

	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.HideWindow = true
	cmd.Dir = buildDir
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.BeforeRun = func(cmd *exec.Cmd) {
		cmdEnv := append(os.Environ(),
			fmt.Sprintf("DEPLOY_GTK_VERSION=%s", gtkVersion),
			fmt.Sprintf("OUTPUT=%s", appImagePath),
		)
		cmd.Env = cmdEnv
	}
	cmd.Command(linuxdeployPath, "--appimage-extract-and-run",
		"--appdir", appDir,
		"--executable", filepath.Join(appDir, "usr", "bin", option.BuildFileName),
		"--desktop-file", filepath.Join(appDir, option.PackageName+".desktop"),
		"--icon-file", filepath.Join(appDir, "icon.png"),
		"--plugin", "gtk", "--output", "appimage")

	event.ConsoleWriteInfo("Package - AppImage", appImageName, "end")
	return true
}

// detectedVersion 根据项目配置获取 GTK 版本和依赖信息
// 返回: gtkVersion, webkitVersion
// gtkVersion: "2" 或 "3"
// webkitVersion: "4.0", "4.1"
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
		// webkit2gtk
		if option.BuildCGOEnabled && strings.Contains(option.GoArgs, "webkit2_4_1") {
			webkitVersion = "4.1"
		} else if option.BuildCGOEnabled {
			webkitVersion = "4.0"
		} else {
			// 非 CGO 从系统检测以安装的 webkit2gtk, 优先 4.1 > 4.0
			deps := []*AppImageDepsInstall{
				{Name: "webkit2gtk-4.1"},
				{Name: "webkit2gtk-4.0"},
			}
			findDeps(deps)
			for _, dep := range deps {
				if strings.Contains(dep.LibPath, "webkit2gtk-4.1") {
					webkitVersion = "4.1"
					break // 优先 4.1
				} else if strings.Contains(dep.LibPath, "webkit2gtk-4.0") {
					webkitVersion = "4.0"
				}
			}
		}
	}
	// CEF 框架只需要 GTK3，不需要 WebKit
	return gtkVersion, webkitVersion
}
