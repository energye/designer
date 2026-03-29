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

//go:build windows

package packager

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"os/exec"
)

// signtool.exe

func packager() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - GProject is nil")
		return false
	}
	event.ConsoleWriteInfo("Package - project check config options")
	option := proj.BuildOption
	if option.Cert {
		if !cert() {
			return false
		}
	} else {
		event.ConsoleWriteInfo("Package - Not Enabled cert")
	}
	if option.WinExe {
		packageNSIS()
	} else if option.WinMsi {

	}
	return false
}

func cert() bool {
	return true
}

func createAppBundle() bool {
	// empty impl
	return true
}

func checkToolCMD(name string) bool {
	//_, err := exec.LookPath(name)
	//if err != nil {
	//	return false
	//}
	cmd := exec.Command("where", name)
	if tool.IsWindows {
		cmd.SysProcAttr = command.HideWindow(true)
	}
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
