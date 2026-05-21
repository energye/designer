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
	"bufio"
	"fmt"
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// findSOPath 通过 ldconfig -p 查找 .so 文件的系统路径
func findSOPath(soName string) string {
	out, err := exec.Command("ldconfig", "-p").Output()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		// 格式: "	libfoo.so.0 (libc6,x86-64) => /usr/lib/x86_64-linux-gnu/libfoo.so.0"
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, soName+" ") {
			continue
		}
		idx := strings.Index(trimmed, "=> ")
		if idx < 0 {
			continue
		}
		return strings.TrimSpace(trimmed[idx+3:])
	}
	return ""
}

// copySODeps 将依赖 .so 文件及其版本符号链接复制到目标目录
func copySODeps(soName, dstDir string) bool {
	srcPath := findSOPath(soName)
	if srcPath == "" {
		event.ConsoleWriteWarn("Package - AppImage: library not found:", soName, "(skipping)")
		return true // 非致命，跳过
	}

	srcDir := filepath.Dir(srcPath)
	base := filepath.Base(srcPath)

	// 复制主文件（真实文件）
	dstFile := filepath.Join(dstDir, base)
	if err := tool.CopyFile(srcPath, dstFile); err != nil {
		event.ConsoleWriteError("Package - AppImage: copy", soName, "failed:", err.Error())
		return false
	}

	// 查找并复制同目录下同前缀的符号链接
	// 例如 libfoo.so.0 -> libfoo.so.0.3600.0，以及 libfoo.so -> libfoo.so.0
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return true // 已复制主文件，链接失败非致命
	}
	prefix := soName // e.g. "libgtk-3.so.0"
	for _, entry := range entries {
		name := entry.Name()
		if name == base {
			continue
		}
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		// 检查是否为符号链接
		linkPath := filepath.Join(srcDir, name)
		info, err := os.Lstat(linkPath)
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		// 读取链接目标
		target, err := os.Readlink(linkPath)
		if err != nil {
			continue
		}
		// 创建符号链接
		dstLink := filepath.Join(dstDir, name)
		_ = os.Remove(dstLink)
		if err := os.Symlink(target, dstLink); err != nil {
			event.ConsoleWriteWarn("Package - AppImage: create symlink", name, "failed:", err.Error())
		}
	}
	return true
}

// appImage 构建 AppImage 包
func (m *Package) appImage() bool {
	event.ConsoleWriteInfo("Package - AppImage")

	// AppImage 内嵌的 runtime 是宿主架构的，不支持跨架构构建
	targetArch := lib.GOARCH()
	if runtime.GOARCH != targetArch {
		event.ConsoleWriteWarn("Package - AppImage: skipped, cross-architecture build is not supported. Current:", runtime.GOARCH, ", Target:", targetArch)
		return true
	}

	if !checkToolCMD("appimagetool") {
		event.ConsoleWriteError(`Package - AppImage: appimagetool command not found. Download from: https://github.com/AppImage/AppImageKit/releases
	chmod +x appimagetool-x86_64.AppImage
  	sudo mv appimagetool-x86_64.AppImage /usr/local/bin/appimagetool
  	sudo apt install libfuse2`)
		return false
	}

	proj := bean.GProject
	option := proj.BuildOption
	appOption := proj.AppOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}

	// 处理使用的 libenergy.so 运行时库, lib.GetDLLName()
	defer env.Delete("--ws") // 删除环境变量 --ws
	choiceLibEnergySOWS()

	// 创建 AppDir 目录结构
	appDir := filepath.Join(output, option.PackageName+".AppDir")

	// 清理临时目录
	defer os.RemoveAll(appDir)

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

	// 复制依赖库 (GTK, WebKit 等)
	libDir := filepath.Join(appDir, "usr", "lib")
	for _, soName := range linuxAppImageDeps() {
		if !copySODeps(soName, libDir) {
			return false
		}
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

	event.ConsoleWriteInfo("Package - AppImage:", appImageName, "created successfully")
	return true
}
