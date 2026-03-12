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

package preview

import (
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/cmd/packager"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"path/filepath"
)

var runCmd *command.CMD

// 执行应用程序的预览功能
// 根据项目配置预览当前项目
func runPreview(state chan<- any) {
	if runCmd != nil {
		return
	}
	event.ConsoleWriteClear() //清空控制台消息
	state <- consts.PsStarting

	// macOS > xxx.app, windows > xxx.exe, linux > xxx
	output := appExecutable()

	// build app
	if !build.Run() {
		state <- consts.PsStop
		return
	}
	// macOS bundle > xxx.app
	if !packager.AppBundle() {
		state <- consts.PsStop
		return
	}
	// 运行命令
	runCmd = command.NewCMD()
	runCmd.IsPrint = false
	runCmd.HideWindow = true
	runCmd.Dir = bean.GPath
	runCmd.Console = func(data string, level command.Level) {
		logs.Info("[", level.String(), "]", data)
		event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: data}}) //正常消息
		if tool.Equal(data, "exit") {
			// 退出
			//state <- 0
		}
	}
	// 开始运行
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: "运行预览: " + output}})
	state <- consts.PsStarted
	runCmd.Command(output)
	state <- consts.PsStop
	close(state)
	logs.Debug("run preview end")
	runCmd = nil
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: "结束预览"}}) //运行结束消息
}

// 停止预览
func stopPreview() {
	// 停止运行
	if runCmd != nil {
		logs.Debug("停止预览, 进程ID:", runCmd.Cmd.Process.Pid)
		err := runCmd.Cmd.Process.Kill()
		logs.Debug("停止预览, 进程ID:", runCmd.Cmd.Process.Pid, "结果:", err)
	}
	runCmd = nil
}

// appExecutable 获取应用程序可执行文件的完整路径
//
//	根据项目构建选项和操作系统类型，计算并返回可执行文件的绝对路径
//	在 macOS 上会构建 .app 包的可执行文件路径，在 Windows/Linux 上直接返回构建文件路径
func appExecutable() string {
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	if tool.IsWindows || tool.IsLinux {
		buildFile := filepath.Join(output, option.BuildFileName)
		return buildFile
	} else if tool.IsDarwin {
		packageName := option.PackageName + ".app"
		appRoot := filepath.Join(output, packageName)
		macOS := filepath.Join(appRoot, "Contents", "MacOS")
		buildFile := filepath.Join(macOS, proj.AppOption.MacOS.PList.CFBundleExecutable)
		return buildFile
	}
	return ""
}
