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
	"path/filepath"
	"strings"
)

const signtool = "signtool.exe"

func packager() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - GProject is nil")
		return false
	}
	event.ConsoleWriteInfo("Package - project check config options")
	option := proj.BuildOption
	if !option.WinSign.Enable {
		event.ConsoleWriteInfo("Package - Not Enabled cert")
	}
	if option.WinExe {
		packageNSIS()
	} else if option.WinMsi {

	}
	return false
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

// 签名 windows 二进制文件
func signWindowsBinary(binaryFilePath string) {
	event.ConsoleWriteInfo("Signing binary file:", binaryFilePath)
	if bean.GProject.BuildOption.WinSign.Enable {
		if !checkToolCMD(signtool) {
			event.ConsoleWriteError("signtool.exe is not installed, is not found in PATH. Please install Windows SDK")
			return
		}
		event.ConsoleWriteInfo("Sign Command Configuration:", strings.Join(bean.GProject.BuildOption.WinSign.Cert, " "))
		if len(bean.GProject.BuildOption.WinSign.Cert) > 0 {
			certLine := bean.GProject.BuildOption.WinSign.Cert[0]
			dataLine := strings.Split(certLine, "=") // file=xxxx  or  auto=xxx
			if len(dataLine) == 2 {
				name := strings.TrimSpace(dataLine[0])
				cmd := strings.TrimSpace(dataLine[1])
				cmdArray := strings.Split(cmd, " ")
				if strings.Contains(cmdArray[0], "signtool") {
					// 删除 signtool
					cmdArray = cmdArray[1:]
				}
				success := false
				// 检查配置是否正确
				if success = name == "auto"; success {

				} else if success = name == "file"; success {
					for i := 0; i < len(cmdArray); i++ {
						v := strings.TrimSpace(cmdArray[i])
						if v == "/f" && i < len(cmdArray) {
							i++             // next
							v = cmdArray[i] // 证书文件名
							// 处理证书相对目录
							if v[0] == '@' {
								v = v[1:]
								// 相对目录, 从项目的 resources 目录找证书
								cmdArray[i] = filepath.Join(bean.ResourcePath(), v)
							}
							break
						}
					}
				}
				if success {
					args := append(cmdArray, binaryFilePath)
					err := RunCMD("", signtool, args...)
					if err != nil {
						event.ConsoleWriteError(err.Error())
						return
					}
				} else {
					event.ConsoleWriteError("Signature config must be 'auto' or 'file'. See example.")
					return
				}
			} else {
				event.ConsoleWriteError("Incorrect certificate signature configuration. Please check the example.")
				return
			}
		} else {
			// 未配置任何签名参数命令
			event.ConsoleWriteError("No signature parameters are set. Please check the configuration example.")
			return
		}
		event.ConsoleWriteInfo("Binary file signed successfully.", binaryFilePath)
	} else {
		event.ConsoleWriteInfo("Signing not enabled.")
	}
}
