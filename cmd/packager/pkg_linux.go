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
	"os/exec"
)

func (m *Package) platformPackage() {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build - project GProject is nil")
		return
	}
	event.ConsoleWriteInfo("CMD-package-run", "GOOS:", lib.GOOS(), "GOARCH:", lib.GOARCH())

	if !build.Run() {
		return
	}
	event.ConsoleWriteInfo("CMD-package-run")
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
		if !m.dpkg() {
			return false
		}
	}
	if option.LinuxRPM {
		if !m.rpmbuild() {
			return false
		}
	}
	if option.LinuxAppImage {
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
		env.Put("--ws", "gtk3")
	}
}
