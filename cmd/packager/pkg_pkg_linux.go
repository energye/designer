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
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/resources/frameworks/lib"
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
