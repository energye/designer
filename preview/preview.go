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
	"errors"
	"fmt"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project"
	lclTool "github.com/energye/lcl/tool"
	"github.com/energye/lcl/tool/command"
	"strings"
)

var runCmd *command.CMD

// 构建项目
func build(output string) (err error) {
	buildCmd := command.NewCMD()
	buildCmd.IsPrint = false
	buildCmd.Dir = project.Path
	buildCmd.Console = func(data string, level command.Level) {
		logs.Info("Level", level.String(), data)
		event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 0, Data: data}}) //正常消息
		if level == command.LError && err == nil {
			err = errors.New(data)
		}
	}
	// TODO 需要通过配置, 构建参数
	args := []string{"build", "-v", "-o", output}
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 0, Data: "go " + strings.Join(args, " ")}})
	buildCmd.Command("go", args...)
	return
}

// 执行应用程序的预览功能
// 根据项目配置预览当前项目
func runPreview(state chan<- any) {
	if runCmd != nil {
		return
	}
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 1}}) //清空控制台消息

	state <- consts.PsStarting
	var output string
	// TODO 需要通过配置 --test
	if lclTool.IsWindows() {
		output = "./build/main.exe"
	} else {
		output = "./build/main"
	}
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 0, Data: "构建程序: " + output}})
	// 构建项目二进制
	if err := build(output); err != nil {
		msg := fmt.Sprintf("构建程序失败: %v", err.Error())
		logs.Error(msg)
		state <- consts.PsStop
		return
	}
	// 运行命令
	runCmd = command.NewCMD()
	runCmd.IsPrint = false
	runCmd.Dir = project.Path
	runCmd.Console = func(data string, level command.Level) {
		logs.Info("[", level.String(), "]", data)
		event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 0, Data: data}}) //正常消息
		if tool.Equal(data, "exit") {
			// 退出
			//state <- 0
		}
	}
	// 开始运行
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 0, Data: "运行预览: " + output}})
	state <- consts.PsStarted
	runCmd.Command(output)
	state <- consts.PsStop
	close(state)
	logs.Debug("run preview end")
	runCmd = nil
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: 0, Data: "结束预览"}}) //运行结束消息
}

func stopPreview() {
	// 停止运行
	if runCmd != nil {
		logs.Debug("停止预览, 进程ID:", runCmd.Cmd.Process.Pid)
		err := runCmd.Cmd.Process.Kill()
		logs.Debug("停止预览, 进程ID:", runCmd.Cmd.Process.Pid, "结果:", err)
	}
	runCmd = nil
}
