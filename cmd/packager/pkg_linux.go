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
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/api"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func (m *Package) platformPackage() {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build - project GProject is nil")
		return
	}
	event.ConsoleWriteInfo("CMD-package-run", "GOOS:", lib.GOOS(), "GOARCH:", lib.GOARCH())
	if m.packager() {
		event.ConsoleWriteInfo("Package Successfully")
	}
}

func (m *Package) packager() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - GProject is nil")
		return false
	}
	event.ConsoleWriteInfo("Package - project check config options")
	option := proj.BuildOption
	if option.LinuxDEB {
		if !build.Run(m.Context) {
			return false
		}
		if !m.dpkg() {
			return false
		}
	}
	if option.LinuxRPM {
		if !build.Run(m.Context) {
			return false
		}
		if !m.rpmbuild() {
			return false
		}
	}
	if option.LinuxAppImage {
		// AppImage, must enable CGO
		oldCGOEnable := option.BuildCGOEnabled
		defer func() {
			(&proj.BuildOption).BuildCGOEnabled = oldCGOEnable
		}()
		(&proj.BuildOption).BuildCGOEnabled = true
		if !build.Run(m.Context) {
			return false
		}
		if !m.appImage() {
			return false
		}
	}
	return true
}

func (m *Package) createAppBundle() bool {
	return true
}

// checkToolCMD 检查命令工具是否可用
func checkToolCMD(name string) bool {
	_, err := exec.LookPath(name)
	if err != nil {
		return false
	}
	return true
}

// choiceLibEnergySOWS 根据GUI渲染框架和GTK版本选择对应的libenergy.so运行时库
//
// 该函数根据不同的GUI渲染框架(LCL/WV/CEF)和GTK版本(GTK2/GTK3)配置,
// 设置相应的窗口系统类型,以确保使用正确的libenergy.so动态链接库。
//
// 选择规则:
//   - LCL + GTK3启用    -> GTK3
//   - LCL + GTK2        -> GTK2
//   - LCL + WV          -> GTK3
//   - LCL + CEF         -> GTK3
//   - 默认情况          -> GTK2
func choiceLibEnergySOWS() {
	// 处理使用的 libenergy.so 运行时库, lib.GetDLLName()
	// linux 分为 gtk2/gtk3
	// 根据渲染框架选择不同的 libenergy.so
	// 根据构建选项 UI (GTK2, GTK3) 配置
	// GUIRenderFramework：分为 LCL, LCL + WV, LCL + CEF 3种
	// 规则：
	//	GUIRenderFramework			UI				libenergy.so
	//	LCL							GTK2 or GTK3	GTK3
	//	LCL							GTK2			GTK2
	//	LCL + WV 					GTK3			GTK3
	//	LCL + CEF					GTK3			GTK3
	proj := bean.GProject
	option := proj.BuildOption
	wsGtk := api.WtGTK2 // 默认 GTK2
	// 渲染框架
	if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_LCL) {
		// LCL 启用了 GTK3
		if option.UIGtk3 {
			wsGtk = api.WtGTK3
		}
	} else if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_WV) {
		// Webkit2Gtk 使用 GTK3
		wsGtk = api.WtGTK3
	} else if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_CEF) {
		// CEF 使用 GTK3
		wsGtk = api.WtGTK3
	}
	if wsGtk == api.WtGTK3 {
		env.Put("ENERGY_WS", "gtk3")
	} else {
		env.Put("ENERGY_WS", "gtk2")
	}
}

// linuxAutoDeps 根据 GUIRenderFramework 和 UI 配置生成自动依赖
// 返回 DEB 和 RPM 格式的依赖字符串
// RPM 使用文件依赖(.so)，兼容所有发行版
func linuxAutoDeps() (debDeps, rpmDeps string) {
	proj := bean.GProject
	option := proj.BuildOption
	gtk3 := false
	webkit := false

	if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_LCL) {
		if option.UIGtk3 {
			gtk3 = true
		}
	} else if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_WV) {
		gtk3 = true
		webkit = true
	} else if tool.Equal(proj.GUIRenderFramework, bean.GUIRenderFramework_CEF) {
		gtk3 = true
	}

	// 64位架构需要 (64bit) 后缀
	soSuffix := ""
	goarch := lib.GOARCH()
	if goarch == "amd64" || goarch == "arm64" || goarch == "loong64" {
		soSuffix = "()(64bit)"
	}
	var (
		deb []string // ,
		rpm []string // /n
	)
	if gtk3 {
		deb = append(deb, "libgtk-3-0 (>= 3.24.24)", "libglib2.0-0 (>= 2.66.0)")
		rpm = append(rpm, "Requires: libgtk-3.so.0"+soSuffix+" >= 3.24.24", "Requires: libglib-2.0.so.0"+soSuffix+" >= 2.66.0")

		if webkit {
			// deb := " (>= 2.40.0)"
			// rpm := " >= 2.40.0"
			if option.BuildCGOEnabled && strings.Contains(option.GoArgs, "webkit2_4_1") {
				// CGO + webkit2_4_1 build tag -> 4.1
				deb = append(deb, "libwebkit2gtk-4.1-0")
				rpm = append(rpm, "Requires: libwebkit2gtk-4.1.so.0"+soSuffix)
			} else if option.BuildCGOEnabled {
				// CGO 默认 -> 4.0
				deb = append(deb, "libwebkit2gtk-4.0-37")
				rpm = append(rpm, "Requires: libwebkit2gtk-4.1.so.0"+soSuffix)
			} else {
				// 非 CGO dlopen 自动降级
				deb = append(deb, "libwebkit2gtk-4.1-0 | libwebkit2gtk-4.0-37")
				rpm = append(rpm, "Requires: (libwebkit2gtk-4.1.so.0"+soSuffix+" or libwebkit2gtk-4.0.so.37"+soSuffix+")")
			}
		}
	} else {
		deb = append(deb, "libgtk2.0-0 (>= 2.24.0)")
		rpm = append(rpm, "Requires: libgtk-2.0.so.0"+soSuffix+" >= 2.24.0")
	}
	deb = append(deb, "libharfbuzz-gobject0", "libgl1")
	rpm = append(rpm, "Requires: libharfbuzz-gobject.so.0"+soSuffix, "Requires: libGL.so.1"+soSuffix)

	debDeps = strings.Join(deb, ",")
	rpmDeps = strings.Join(rpm, "\n")
	return
}

// linuxUserOverrideWebKit 检测用户是否指定了 webkit2gtk 4.0，调整依赖顺序
// 非 CGO 模式下，如果用户在 Depends 里指定了 4.0，将优先级调整为用户的版本
func linuxUserOverrideWebKit(debDeps, rpmDeps, userDeps string) (string, string) {
	if !strings.Contains(userDeps, "webkit2gtk-4.0") && !strings.Contains(userDeps, "webkit2gtk4.0") {
		return debDeps, rpmDeps
	}
	// DEB: swap OR order
	debDeps = strings.Replace(debDeps, "libwebkit2gtk-4.1-0 | libwebkit2gtk-4.0-37", "libwebkit2gtk-4.0-37 | libwebkit2gtk-4.1-0", 1)
	// RPM (.so 文件依赖): 交换 4.1 <-> 4.0
	// 用临时占位符避免二次替换冲突
	rpmDeps = strings.ReplaceAll(rpmDeps, "libwebkit2gtk-4.1.so.0", "libwebkit2gtk-__TMP_A__")
	rpmDeps = strings.ReplaceAll(rpmDeps, "libwebkit2gtk-4.0.so.37", "libwebkit2gtk-4.1.so.0")
	rpmDeps = strings.ReplaceAll(rpmDeps, "libwebkit2gtk-__TMP_A__", "libwebkit2gtk-4.0.so.37")
	return debDeps, rpmDeps
}

type AppImageDepsInstall struct {
	Name    string // 库名称
	LibDir  string // 库文件所在目录，如 /usr/lib/x86_64-linux-gnu
	LibPath string // 库文件路径
}

func findDeps(deps []*AppImageDepsInstall) {
	// 再次尝试从常规目录获取
	for _, dep := range deps {
		if dep.LibDir == "" {
			// 为空时尝试从常规目录获取
			findByCommonPaths(dep)
		}
	}
}

func findLdconfig(name string) []string {
	cmd := exec.Command("ldconfig", "-p")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	result := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		for _, f := range fields {
			if strings.HasPrefix(f, "/") && strings.Contains(f, name) {
				result = append(result, f)
			}
		}
	}
	return result
}

func findLdd(libSOPath, name string) []string {
	cmd := exec.Command("ldd", libSOPath)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	result := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		for _, f := range fields {
			if strings.HasPrefix(f, "/") && strings.Contains(f, name) {
				result = append(result, f)
			}
		}
	}
	return result
}

func findByCommonPaths(dep *AppImageDepsInstall) bool {
	possiblePaths := []string{
		"/usr/lib",
		"/usr/lib64",
	}
	// 尝试从 uname -m 构造搜索目录
	if arch := platformArch(); arch != "" {
		possiblePaths = append(possiblePaths, "/usr/lib/"+arch+"-linux-gnu")
	}
	for _, libDir := range possiblePaths {
		if libPath := filepath.Join(libDir, dep.Name); tool.IsExist(libPath) {
			dep.LibDir = libDir
			dep.LibPath = libPath
			return true
		} else if libPath := filepath.Join(libDir, dep.Name+".so"); tool.IsExist(libPath) {
			dep.LibDir = libDir
			dep.LibPath = libPath
			return true
		} else if libPath := filepath.Join(libDir, "lib"+dep.Name+".so"); tool.IsExist(libPath) {
			dep.LibDir = libDir
			dep.LibPath = libPath
			return true
		}
	}
	return false
}

func findFiles(deps []*AppImageDepsInstall) error {
	err := filepath.Walk("/usr/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, dep := range deps {
			if strings.Contains(path, dep.Name) {
				dep.LibPath = path
				dep.LibDir = filepath.Dir(path)
				break
			}
		}
		return nil
	})
	return err
}

func platformArch() string {
	archs := map[string]string{
		"arm64":   "aarch64",
		"aarch64": "aarch64",
		"amd64":   "x86_64",
		"x86_64":  "x86_64",
		"386":     "i686",
		"i386":    "i686",
		"i686":    "i686",
		"armhf":   "armhf",
		"arm":     "armhf",
		"arm32":   "armhf",
	}
	arch := archs[runtime.GOARCH]
	return arch
}
