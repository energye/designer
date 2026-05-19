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

func buildDarwin(env Envs) bool {
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
	isArm64 := lib.GOARCH() == "arm64"
	if isAmd64 {
		if !option.ArchAmd64 {
			event.ConsoleWriteWarn("Build - amd64 architecture not enabled for Project Settings > Build Configurations")
			return false
		}
		env.Put("MACOSX_DEPLOYMENT_TARGET", "10.15")
		env.Put("CGO_CFLAGS", "-mmacosx-version-min=10.15")
		env.Put("CGO_LDFLAGS", "-mmacosx-version-min=10.15")
	}
	if isArm64 {
		if !option.ArchArm64 {
			event.ConsoleWriteWarn("Build - arm64 architecture not enabled for Project Settings > Build Configurations")
			return false
		}
		env.Put("MACOSX_DEPLOYMENT_TARGET", "11.0")
		env.Put("CGO_CFLAGS", "-mmacosx-version-min=11.0")
		env.Put("CGO_LDFLAGS", "-mmacosx-version-min=11.0")
	}
	if !option.UICocoa {
		event.ConsoleWriteWarn("Build - UI Cocoa is not enabled for the project.Project Settings > Build Configurations")
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
	// macOS 去除 -H windowsgui
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
		if option.BuildModeRelease && tool.IsDarwin {
			event.ConsoleWriteInfo("strip", output)
			cmd.Command("strip", output)
		}
	}
	/*
	 * 根据配置决定是否构建 macOS 通用二进制文件 (Universal Binary)。
	 * 如果启用了 MacCommonLib 选项，则分别编译 amd64 和 arm64 架构的临时二进制文件，
	 * 随后使用 lipo 工具将它们合并为一个支持多架构的最终产物，并清理中间文件。
	 * 否则，执行标准的单架构构建流程。
	 */
	if option.MacCommonLib && tool.IsDarwin {
		// build amd64
		amd64OutputFilename := filepath.Join(output, "temp_amd64_"+option.BuildFileName)
		runGoBuild([]string{"GOARCH=amd64"}, amd64OutputFilename)
		// build arm64
		arm64OutputFilename := filepath.Join(output, "temp_arm64_"+option.BuildFileName)
		runGoBuild([]string{"GOARCH=arm64"}, arm64OutputFilename)
		defer func() {
			_ = os.Remove(amd64OutputFilename)
			_ = os.Remove(arm64OutputFilename)
		}()
		// merge universal
		outputFilename := filepath.Join(output, buildBinFileName(env, option))
		mergeUniversal(amd64OutputFilename, arm64OutputFilename, outputFilename)
		verifyUniversal(outputFilename)
	} else {
		outputFilename := filepath.Join(output, buildBinFileName(env, option))
		runGoBuild(env.ToArray(), outputFilename)
		verifyUniversal(outputFilename)
	}

	event.ConsoleWriteInfo("Build Successfully")
	return true
}

func mergeUniversal(amd64File, arm64File, output string) {
	if !tool.IsDarwin {
		return
	}
	event.ConsoleWriteInfo("Build - merge universal")
	cmd := command.NewCMD()
	cmd.Dir = bean.GPath
	cmd.HideWindow = true
	cmd.Command("lipo", "-create", amd64File, arm64File, "-output", output)
}

func verifyUniversal(filePath string) {
	if !tool.IsDarwin {
		return
	}
	event.ConsoleWriteInfo("Build - verify universal")
	cmd := command.NewCMD()
	cmd.Dir = bean.GPath
	cmd.HideWindow = true
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.Command("lipo", "-info", filePath)
}
