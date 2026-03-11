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

//go:build darwin

package build

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func build() {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build - project GProject is nil")
		return
	}
	event.ConsoleWriteInfo("Build - project check config options")
	option := proj.BuildOption
	if !option.PlatformMacOS {
		event.ConsoleWriteWarn("Build - Project has not enabled Project Settings > Build Configurations")
		return
	}
	isAmd64 := runtime.GOARCH == "amd64"
	isArm64 := runtime.GOARCH == "arm64"
	if isAmd64 {
		if !option.ArchX86_64 {
			event.ConsoleWriteWarn("Build - amd64 architecture not enabled for Project Settings > Build Configurations")
			return
		}
		os.Setenv("MACOSX_DEPLOYMENT_TARGET", "10.15")
		os.Setenv("CGO_CFLAGS", "-mmacosx-version-min=10.15")
		os.Setenv("CGO_LDFLAGS", "-mmacosx-version-min=10.15")
	}
	if isArm64 {
		if !option.ArchAarch64 {
			event.ConsoleWriteWarn("Build - arm64 architecture not enabled for Project Settings > Build Configurations")
			return
		}
		os.Setenv("MACOSX_DEPLOYMENT_TARGET", "11.0")
		os.Setenv("CGO_CFLAGS", "-mmacosx-version-min=11.0")
		os.Setenv("CGO_LDFLAGS", "-mmacosx-version-min=11.0")
	}
	if !option.UICocoa {
		event.ConsoleWriteWarn("Build - UI Cocoa is not enabled for the project.Project Settings > Build Configurations")
		return
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
	// 合并 tags prod
	tags := tool.MergeTags("prod", customTags)
	// 合并 ldflags
	ldflags := tool.MergeLdflags("", customLdflags)
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
		cmd.BeforeRun = func(cmd *exec.Cmd) {
			cmd.Env = append(os.Environ(), env...)
		}
		cmd.Command("go", tempArgs...)
		event.ConsoleWriteInfo("strip", output)
		cmd.Command("strip", output)
		event.ConsoleWriteInfo("Build - output", output)
	}
	if option.MacCommonLib {
		// build amd64
		amd64OutputFilename := filepath.Join(output, "temp_amd64_"+option.BuildFileName)
		runGoBuild([]string{"GOOS=darwin", "GOARCH=amd64", "CGO_ENABLED=1"}, amd64OutputFilename)
		// build arm64
		arm64OutputFilename := filepath.Join(output, "temp_arm64_"+option.BuildFileName)
		runGoBuild([]string{"GOOS=darwin", "GOARCH=arm64", "CGO_ENABLED=1"}, arm64OutputFilename)
		defer func() {
			_ = os.Remove(amd64OutputFilename)
			_ = os.Remove(arm64OutputFilename)
		}()
		// merge universal
		outputFilename := filepath.Join(output, option.BuildFileName)
		mergeUniversal(amd64OutputFilename, arm64OutputFilename, outputFilename)
		verifyUniversal(outputFilename)
	} else {
		outputFilename := filepath.Join(output, option.BuildFileName)
		runGoBuild(nil, outputFilename)
		verifyUniversal(outputFilename)
	}

	event.ConsoleWriteInfo("Build Successfully")
}

func mergeUniversal(amd64File, arm64File, output string) {
	event.ConsoleWriteInfo("Build - merge universal")
	cmd := command.NewCMD()
	cmd.Dir = bean.GPath
	cmd.Command("lipo", "-create", amd64File, arm64File, "-output", output)
}

func verifyUniversal(filePath string) {
	event.ConsoleWriteInfo("Build - verify universal")
	cmd := command.NewCMD()
	cmd.Dir = bean.GPath
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.Command("lipo", "-info", filePath)
}
