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
	"github.com/energye/designer/cmd/run"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
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
	state <- consts.PsStarting
	// build app
	if !build.Run(nil) {
		state <- consts.PsStop
		return
	}
	// macOS bundle > xxx.app
	if !packager.AppBundle(nil) {
		state <- consts.PsStop
		return
	}
	// macOS > xxx.app, windows > xxx.exe, linux > xxx
	output := run.AppExecutable()
	event.ConsoleWriteInfo("运行预览:", output)

	if tool.IsDarwin {
		// /Users/xxx/xxx/build/myapp.app/Contents/MacOS/myapp
		// 预览 dev 签名
		event.ConsoleWriteInfo("run preview app-codesign")
		xAppDir := filepath.Join(output, "../../../") // /Users/xxx/xxx/build/myapp.app
		//  codesign --force --deep --sign - xxx.app
		RunCMD("", "codesign", "--force", "--deep", "--sign", "-", xAppDir)
		event.ConsoleWriteInfo("run preview app-codesign end")
	}

	// 开始运行
	state <- consts.PsStarted // 运行命令
	runCmd = command.NewCMD()
	runCmd.HideWindow = true
	run.Run(runCmd)
	state <- consts.PsStop
	close(state)
	logs.Debug("run preview end")
	runCmd = nil
	event.ConsoleWriteInfo("结束预览")
}

// 停止预览
func stopPreview() {
	// 停止运行
	if runCmd != nil && runCmd.Cmd != nil && runCmd.Cmd.Process != nil {
		logs.Debug("停止预览, 进程ID:", runCmd.Cmd.Process.Pid)
		err := runCmd.Cmd.Process.Kill()
		logs.Debug("停止预览, 进程ID:", runCmd.Cmd.Process.Pid, "结果:", err)
	}
	runCmd = nil
}

func RunCMD(dir, name string, args ...string) {
	cmd := command.NewCMD()
	cmd.HideWindow = true
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.Command(name, args...)
}
