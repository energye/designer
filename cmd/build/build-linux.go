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

package build

import (
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/tool/command"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func buildLinux() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build - project GProject is nil")
		return false
	}
	event.ConsoleWriteInfo("Build - project check config options")
	option := proj.BuildOption
	if !option.PlatformMacOS {
		event.ConsoleWriteWarn("Build - Project has not enabled Project Settings > Build Configurations")
		return false
	}
	isAmd64 := lib.GOARCH() == "amd64"
	is386 := lib.GOARCH() == "386"
	isArm64 := lib.GOARCH() == "arm64"
	isArm := lib.GOARCH() == "arm"
	if isAmd64 {
		if !option.ArchAmd64 {
			event.ConsoleWriteWarn("Build - amd64 architecture not enabled for Project Settings > Build Configurations")
			return false
		}
	}
	if is386 {
		if !option.Arch386 {
			event.ConsoleWriteWarn("Build - 386 architecture not enabled for Project Settings > Build Configurations")
			return false
		}
	}
	if isArm64 {
		if !option.ArchArm64 {
			event.ConsoleWriteWarn("Build - arm64 architecture not enabled for Project Settings > Build Configurations")
			return false
		}
	}
	if isArm {
		if !option.ArchArm {
			event.ConsoleWriteWarn("Build - arm architecture not enabled for Project Settings > Build Configurations")
			return false
		}
	}
	if !option.UIGtk3 && !option.UIGtk2 {
		event.ConsoleWriteWarn("Build - UI Gtk2 or Gtk3 is not enabled for the project.Project Settings > Build Configurations")
		return false
	}
	event.ConsoleWriteInfo("Build - start build", proj.Name)

	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
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
	buildMode := "dev"
	if option.BuildModeRelease {
		buildMode = "prod"
	}
	// 合并 tags prod
	tags := tool.MergeTags(buildMode, customTags)
	// 合并 ldflags
	ldflags := tool.MergeLdflags("", customLdflags)
	// -X
	for _, pack := range xBuildPackVar() {
		ldflags = tool.MergeLdflags(pack, strings.Join(ldflags, " "))
	}
	// 其它参数
	otherArgs := tool.ExtractOtherBuildArgs(option.GoArgs)
	// linux 去除 -H windowsgui
	tempNewLdflags := []string{}
	for _, v := range ldflags {
		if v == "-H" || v == "windowsgui" {
			continue
		}
		tempNewLdflags = append(tempNewLdflags, v)
	}
	ldflags = tempNewLdflags
	args := []string{"build", "-v"}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	// 编译模式
	if option.BuildModeDebug {
		// debug
		if len(ldflags) > 0 {
			args = append(args, "-ldflags", strings.Join(ldflags, " "))
		}
	} else {
		// release
		ldflags = tool.MergeLdflags("-s -w", strings.Join(ldflags, " "))
		args = append(args, "-ldflags", strings.Join(ldflags, " "))
		args = append(args, "-trimpath")
	}

	runGoBuild := func(env []string, output string) {
		tempArgs := args[:]
		tempArgs = append(tempArgs, "-o", output)
		if len(otherArgs) > 0 {
			tempArgs = append(tempArgs, otherArgs...)
		}
		event.ConsoleWriteInfo("Build - output", output)
		event.ConsoleWriteInfo("go", strings.Join(tempArgs, " "))
		// 编译命令
		cmd := command.NewCMD()
		cmd.Dir = bean.GPath
		cmd.HideWindow = true
		cmd.BeforeRun = func(cmd *exec.Cmd) {
			cmd.Env = append(os.Environ(), env...)
		}
		cmd.Console = func(data string, level command.Level) {
			event.ConsoleWriteInfo(data)
		}
		cmd.Command("go", tempArgs...)
		if option.BuildModeRelease && tool.IsLinux {
			event.ConsoleWriteInfo("strip", output)
			cmd.Command("strip", output)
		}
	}

	outputFilename := filepath.Join(output, buildBinFileName(option))
	runGoBuild(env.ToEnviron(), outputFilename)

	event.ConsoleWriteInfo("Build Successfully")
	return true
}
