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
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"path/filepath"
	"regexp"
	"strings"
)

var runCmd *command.CMD

// 构建项目
func runBuild(output string) (err error) {
	option := bean.GProject.BuildOption
	// 编译参数
	buildArgs := option.GoArgs
	buildArgs = strings.ReplaceAll(buildArgs, "'", "\"")
	reTags := regexp.MustCompile(`-tags\s+([^\s-]+)`)
	reLdflags := regexp.MustCompile(`-ldflags\s+"([^"]+)"`)
	// 提取 tags
	tagMatches := reTags.FindStringSubmatch(buildArgs)
	customTags := ""
	customLdflags := ""
	if len(tagMatches) > 1 {
		customTags = tagMatches[1]
	}
	// 提取 ldflags
	ldMatches := reLdflags.FindStringSubmatch(buildArgs)
	if len(ldMatches) > 1 {
		customLdflags = ldMatches[1]
	}
	// 合并 tags
	tags := tool.MergeTags("dev", customTags)
	// 合并 ldflags
	ldflags := tool.MergeLdflags("", customLdflags)
	// 其它参数
	otherArgs := tool.ExtractOtherBuildArgs(option.GoArgs)
	if !tool.IsWindows {
		// macOS, linux 去除 -H windowsgui
		tempNewLdflags := []string{}
		for _, v := range ldflags {
			if v == "-H" || v == "windowsgui" {
				continue
			}
			tempNewLdflags = append(tempNewLdflags, v)
		}
		ldflags = tempNewLdflags
	}

	buildCmd := command.NewCMD()
	buildCmd.IsPrint = false
	buildCmd.HideWindow = true
	buildCmd.Dir = bean.GPath
	buildCmd.Console = func(data string, level command.Level) {
		logs.Info("Level", level.String(), data)
		event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: data}}) //正常消息
		if level == command.LError && err == nil {
			err = errors.New(data)
		}
	}
	args := []string{"build", "-v"}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	if len(ldflags) > 0 {
		args = append(args, "-ldflags", strings.Join(ldflags, " "))
	}
	args = append(args, "-o", output)
	if len(otherArgs) > 0 {
		args = append(args, otherArgs...)
	}
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: "go " + strings.Join(args, " ")}})
	buildCmd.Command("go", args...)
	return
}

// 执行应用程序的预览功能
// 根据项目配置预览当前项目
func runPreview(state chan<- any) {
	if runCmd != nil {
		return
	}
	event.ConsoleWriteClear() //清空控制台消息
	state <- consts.PsStarting
	option := bean.GProject.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	output = filepath.Join(output, option.BuildFileName)

	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: "构建程序: " + output}})
	// 构建项目二进制
	if err := runBuild(output); err != nil {
		msg := fmt.Sprintf("构建程序失败: %v", err.Error())
		logs.Error(msg)
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
