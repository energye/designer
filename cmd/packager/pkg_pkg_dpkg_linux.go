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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
	"path/filepath"
	"strings"
)

// debArchName 将 GOARCH 转换为 DEB 架构名
func debArchName(goarch string) string {
	switch goarch {
	case "amd64":
		return "amd64"
	case "386":
		return "i386"
	case "arm64":
		return "arm64"
	case "arm":
		return "armhf"
	case "loong64":
		return "loong64"
	default:
		return goarch
	}
}

// dpkg 构建 DEB 包
func (m *Package) dpkg() bool {
	event.ConsoleWriteInfo("Package - DEB")
	if !checkToolCMD("dpkg-deb") {
		event.ConsoleWriteError("Package - DEB: dpkg-deb command not found")
		return false
	}

	proj := bean.GProject
	option := proj.BuildOption
	appOption := proj.AppOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}

	goarch := lib.GOARCH()
	debArch := debArchName(goarch)

	// 包名: name_version_arch
	debDirName := fmt.Sprintf("%s_%s_%s", option.PackageName, appOption.Version, debArch)
	debRoot := filepath.Join(output, debDirName)

	// 创建目录结构
	dirs := []string{
		filepath.Join(debRoot, "DEBIAN"),
		filepath.Join(debRoot, "usr", "bin"),
		filepath.Join(debRoot, "usr", "share", "applications"),
		filepath.Join(debRoot, "usr", "share", "icons"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			event.ConsoleWriteError("Package - DEB: mkdir failed:", err.Error())
			return false
		}
	}

	// 复制可执行文件
	srcBin := filepath.Join(output, option.BuildFileName)
	dstBin := filepath.Join(debRoot, "usr", "bin", option.PackageName)
	if err := tool.CopyFile(srcBin, dstBin); err != nil {
		event.ConsoleWriteError("Package - DEB: copy binary failed:", err.Error())
		return false
	}
	if err := os.Chmod(dstBin, 0755); err != nil {
		event.ConsoleWriteError("Package - DEB: chmod binary failed:", err.Error())
		return false
	}

	// 复制 .desktop 文件
	srcDesktop := filepath.Join(bean.ResourceMetadataPath(), option.PackageName+".desktop")
	if !tool.IsExist(srcDesktop) {
		// 使用模板生成
		desktopTemplate := app.Packager("linux/app.desktop")
		if desktopTemplate == nil {
			event.ConsoleWriteError("Package - DEB: desktop template not found")
			return false
		}
		categories := appOption.Linux.Categories
		if categories == "" {
			categories = "Utility;"
		}
		categories = strings.TrimRight(categories, ";") + ";"
		data := map[string]string{
			"Name":       option.PackageName,
			"Exec":       option.PackageName,
			"Icon":       option.PackageName,
			"Comments":   appOption.Desc,
			"WMClass":    option.PackageName,
			"Categories": categories,
		}
		rendered, err := tool.RenderTemplate(string(desktopTemplate), data)
		if err != nil {
			event.ConsoleWriteError("Package - DEB: render desktop failed:", err.Error())
			return false
		}
		srcDesktop = filepath.Join(output, option.PackageName+".desktop")
		if err := os.WriteFile(srcDesktop, rendered, 0644); err != nil {
			event.ConsoleWriteError("Package - DEB: write desktop failed:", err.Error())
			return false
		}
	}
	dstDesktop := filepath.Join(debRoot, "usr", "share", "applications", option.PackageName+".desktop")
	if err := tool.CopyFile(srcDesktop, dstDesktop); err != nil {
		event.ConsoleWriteError("Package - DEB: copy desktop failed:", err.Error())
		return false
	}

	// 复制图标
	srcIcon := filepath.Join(bean.ResourceEmbedPath(), "icon.png")
	if tool.IsExist(srcIcon) {
		dstIcon := filepath.Join(debRoot, "usr", "share", "icons", option.PackageName+".png")
		if err := tool.CopyFile(srcIcon, dstIcon); err != nil {
			event.ConsoleWriteError("Package - DEB: copy icon failed:", err.Error())
			return false
		}
	} else {
		event.ConsoleWriteWarn("Package - DEB: icon not found, skipping")
	}

	// 渲染 control 文件
	controlTemplate := app.Packager("linux/control")
	if controlTemplate == nil {
		event.ConsoleWriteError("Package - DEB: control template not found")
		return false
	}

	// 计算安装大小 (KB)
	installedSize := calcDirSize(debRoot) / 1024

	depends := option.Depends
	if depends == "" {
		depends = "libc6"
	}

	maintainer := appOption.Linux.Maintainer
	if maintainer == "" {
		maintainer = "Unknown <unknown@example.com>"
	}

	homepage := appOption.Linux.Homepage
	if homepage == "" {
		homepage = ""
	}

	section := "utils"

	controlData := map[string]interface{}{
		"PackageName":   option.PackageName,
		"Version":       appOption.Version,
		"Section":       section,
		"Arch":          debArch,
		"InstalledSize": installedSize,
		"Depends":       depends,
		"Maintainer":    maintainer,
		"Homepage":      homepage,
		"Description":   appOption.Desc,
	}

	controlData2, err := tool.RenderTemplate(string(controlTemplate), controlData)
	if err != nil {
		event.ConsoleWriteError("Package - DEB: render control failed:", err.Error())
		return false
	}
	controlPath := filepath.Join(debRoot, "DEBIAN", "control")
	if err := os.WriteFile(controlPath, controlData2, 0644); err != nil {
		event.ConsoleWriteError("Package - DEB: write control failed:", err.Error())
		return false
	}

	// 构建 DEB 包
	debFileName := fmt.Sprintf("%s_%s_%s.deb", option.PackageName, appOption.Version, debArch)
	if m.AppendPlatform {
		debFileName = fmt.Sprintf("%s_%s_%s_%s.deb", option.PackageName, appOption.Version, lib.GOOS(), debArch)
	}

	debFilePath := filepath.Join(output, debFileName)
	_ = os.Remove(debFilePath)

	cmd := RunCMD(output, "dpkg-deb", "--build", debRoot, debFilePath)
	if cmd != nil {
		event.ConsoleWriteError("Package - DEB: build failed:", cmd.Error())
		return false
	}

	// 清理临时目录
	_ = os.RemoveAll(debRoot)

	event.ConsoleWriteInfo("Package - DEB:", debFileName, "created successfully")
	return true
}

// calcDirSize 计算目录大小 (bytes)
func calcDirSize(dir string) int64 {
	var size int64
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// parseDependsToRequires 将 DEB 格式的 Depends 转为 RPM Requires 行
func parseDependsToRequires(depends string) string {
	if depends == "" {
		return ""
	}
	var lines []string
	for _, dep := range strings.Split(depends, ",") {
		dep = strings.TrimSpace(dep)
		if dep == "" {
			continue
		}
		// 处理 | (或关系), 取第一个
		if strings.Contains(dep, "|") {
			dep = strings.TrimSpace(strings.Split(dep, "|")[0])
		}
		// 处理版本约束: "libc6 (>= 2.17)" -> "libc6 >= 2.17"
		dep = strings.ReplaceAll(dep, "(", "")
		dep = strings.ReplaceAll(dep, ")", "")
		dep = strings.TrimSpace(dep)
		lines = append(lines, "Requires: "+dep)
	}
	return strings.Join(lines, "\n")
}
