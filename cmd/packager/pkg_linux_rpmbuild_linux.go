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
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// rpmArchName 将 GOARCH 转换为 RPM 架构名
func rpmArchName(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i686"
	case "arm64":
		return "aarch64"
	case "arm":
		return "armv7hl"
	case "loong64":
		return "loongarch64"
	default:
		return goarch
	}
}

// rpmbuild 构建 RPM 包
func (m *Package) rpmbuild() bool {
	event.ConsoleWriteInfo("Package - RPM")
	if !checkToolCMD("rpmbuild") {
		event.ConsoleWriteError("Package - RPM: rpmbuild command not found. Install with: sudo dnf install rpm-build")
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

	goarch := lib.GOARCH()
	rpmArch := rpmArchName(goarch)

	// 创建 rpmbuild 目录结构
	rpmBuildDir := filepath.Join(output, "rpmbuild")
	specsDir := filepath.Join(rpmBuildDir, "SPECS")
	stageDir := filepath.Join(rpmBuildDir, "stage", "usr")

	dirs := []string{
		filepath.Join(stageDir, "bin"),
		filepath.Join(stageDir, "share", "applications"),
		filepath.Join(stageDir, "share", "icons"),
		filepath.Join(stageDir, "lib"),
		specsDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			event.ConsoleWriteError("Package - RPM: mkdir failed:", err.Error())
			return false
		}
	}

	// 复制可执行文件
	srcBin := filepath.Join(output, option.BuildFileName)
	dstBin := filepath.Join(stageDir, "bin", option.PackageName)
	if err := tool.CopyFile(srcBin, dstBin); err != nil {
		event.ConsoleWriteError("Package - RPM: copy binary failed:", err.Error())
		return false
	}
	if err := os.Chmod(dstBin, 0755); err != nil {
		event.ConsoleWriteError("Package - RPM: chmod binary failed:", err.Error())
		return false
	}

	// 生成 .desktop 文件
	desktopData := renderDesktopFile(option.PackageName, option.PackageName, option.PackageName,
		appOption.Desc, option.PackageName, appOption.Linux.Categories)
	if desktopData == nil {
		event.ConsoleWriteError("Package - RPM: render desktop failed")
		return false
	}
	dstDesktop := filepath.Join(stageDir, "share", "applications", option.PackageName+".desktop")
	if err := os.WriteFile(dstDesktop, desktopData, 0644); err != nil {
		event.ConsoleWriteError("Package - RPM: write desktop failed:", err.Error())
		return false
	}

	// 复制图标
	srcIcon := filepath.Join(bean.ResourceEmbedPath(), "icon.png")
	if tool.IsExist(srcIcon) {
		dstIcon := filepath.Join(stageDir, "share", "icons", option.PackageName+".png")
		if err := tool.CopyFile(srcIcon, dstIcon); err != nil {
			event.ConsoleWriteError("Package - RPM: copy icon failed:", err.Error())
			return false
		}
	}

	// 复制运行时库 libenergy.so
	libName := lib.GetDLLName()
	srcLib := filepath.Join(config.Config.FrameworkRuntimePath(), libName)
	if !tool.IsExist(srcLib) {
		event.ConsoleWriteError("Package - RPM: runtime library not found:", srcLib)
		return false
	}
	dstLib := filepath.Join(stageDir, "lib", libName)
	if err := os.MkdirAll(filepath.Dir(dstLib), 0755); err != nil {
		event.ConsoleWriteError("Package - RPM: mkdir lib failed:", err.Error())
		return false
	}
	if err := tool.CopyFile(srcLib, dstLib); err != nil {
		event.ConsoleWriteError("Package - RPM: copy libenergy failed:", err.Error())
		return false
	}

	// 渲染 spec 文件
	specTemplate := app.Packager("linux/app.spec")
	if specTemplate == nil {
		event.ConsoleWriteError("Package - RPM: spec template not found")
		return false
	}

	_, rpmAutoDeps := linuxAutoDeps()
	_, rpmAutoDeps = linuxUserOverrideWebKit("", rpmAutoDeps, appOption.Linux.Depends)
	depends := rpmAutoDeps
	if userDeps := parseDependsToRequires(appOption.Linux.Depends); userDeps != "" {
		if depends != "" {
			depends += "\n"
		}
		depends += userDeps
	}

	specData := map[string]interface{}{
		"PackageName": option.PackageName,
		"App":         appOption,
		"Linux":       appOption.Linux,
		"Depends":     depends,
		"LibName":     libName,
	}
	specContent, err := tool.RenderTemplate(string(specTemplate), specData)
	if err != nil {
		event.ConsoleWriteError("Package - RPM: render spec failed:", err.Error())
		return false
	}
	specFile := filepath.Join(specsDir, option.PackageName+".spec")
	if err := os.WriteFile(specFile, specContent, 0644); err != nil {
		event.ConsoleWriteError("Package - RPM: write spec failed:", err.Error())
		return false
	}

	// 构建 RPM 包
	stagePath := filepath.Join(rpmBuildDir, "stage")
	args := []string{
		"-bb",
		"--define", fmt.Sprintf("_topdir %s", rpmBuildDir),
		"--define", fmt.Sprintf("_appdir %s", stagePath),
		"--target", rpmArch,
	}
	// 交叉编译时禁用 strip 等后处理，避免宿主 strip 无法识别目标架构二进制
	if runtime.GOARCH != goarch {
		args = append(args, "--define", "__os_install_post %{nil}")
	}
	args = append(args, specFile)
	cmd := RunCMD(output, "rpmbuild", args...)
	if cmd != nil {
		event.ConsoleWriteError("Package - RPM: build failed:", cmd.Error())
		return false
	}

	// 从 RPMS 目录复制到输出目录
	rpmsDir := filepath.Join(rpmBuildDir, "RPMS", rpmArch)
	var rpmFile string
	_ = filepath.Walk(rpmsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".rpm" {
			rpmFile = path
		}
		return nil
	})
	if rpmFile == "" {
		event.ConsoleWriteError("Package - RPM: output file not found")
		return false
	}

	debArch := debArchName(goarch)
	rpmFileName := fmt.Sprintf("%s_%s_%s.rpm", option.PackageName, appOption.Version, debArch)
	if m.AppendPlatform {
		rpmFileName = fmt.Sprintf("%s_%s_%s_%s.rpm", option.PackageName, appOption.Version, lib.GOOS(), debArch)
	}
	dstRpm := filepath.Join(output, rpmFileName)
	_ = os.Remove(dstRpm)
	if err := tool.CopyFile(rpmFile, dstRpm); err != nil {
		event.ConsoleWriteError("Package - RPM: copy output failed:", err.Error())
		return false
	}

	// 清理临时目录
	_ = os.RemoveAll(rpmBuildDir)

	event.ConsoleWriteInfo("Package - RPM:", rpmFileName, "created successfully")
	return true
}

// renderDesktopFile 渲染 .desktop 文件
func renderDesktopFile(name, exec, icon, comment, wmClass, categories string) []byte {
	desktopTemplate := app.Packager("linux/app.desktop")
	if desktopTemplate == nil {
		return nil
	}
	if categories == "" {
		categories = "Utility;"
	}
	categories = strings.TrimRight(categories, ";") + ";"
	data := map[string]string{
		"Name":       name,
		"Exec":       exec,
		"Icon":       icon,
		"Comments":   comment,
		"WMClass":    wmClass,
		"Categories": categories,
	}
	rendered, _ := tool.RenderTemplate(string(desktopTemplate), data)
	return rendered
}
