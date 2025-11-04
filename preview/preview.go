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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project"
	lclTool "github.com/energye/lcl/tool"
	"github.com/energye/lcl/tool/command"
)

var cmd *command.CMD

// 构建项目
func build(output string) {
	buildCmd := command.NewCMD()
	buildCmd.Dir = project.Path
	buildCmd.MessageCallback = func(bytes []byte, err error) {
		info := string(bytes)
		logs.Info(info)
	}
	// TODO 需要通过配置, 构建参数
	buildCmd.Command("go", "build", "-o", output)
}

// 执行应用程序的预览功能
// 根据项目配置预览当前项目
func runPreview(state chan<- any) {
	if cmd != nil {
		return
	}
	state <- consts.PsStarting
	var output string
	// TODO 需要通过配置 --test
	if lclTool.IsWindows() {
		output = "./build/main.exe"
	} else {
		output = "./build/main"
	}
	// 构建项目二进制
	build(output)
	// 运行命令
	cmd = command.NewCMD()
	cmd.Dir = project.Path
	cmd.MessageCallback = func(bytes []byte, err error) {
		info := string(bytes)
		logs.Info(info)
		if tool.Equal(info, "exit") {
			// 退出
			//state <- 0
		}
	}
	// 开始运行
	state <- consts.PsStarted
	cmd.Command(output)
	state <- consts.PsStop
	close(state)
	logs.Debug("run preview end")
	cmd = nil
}

func stopPreview() {
	// 停止运行
	if cmd != nil {
		logs.Debug("停止预览, 进程ID:", cmd.Cmd.Process.Pid)
		err := cmd.Cmd.Process.Kill()
		logs.Debug("停止预览, 进程ID:", cmd.Cmd.Process.Pid, "结果:", err)
	}
	cmd = nil
}
