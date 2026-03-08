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
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/tool/command"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func build() {
	logs.Info("构建项目, 检查配置选项")
	proj := bean.GProject
	option := proj.BuildOption
	if !option.PlatformMacOS {
		logs.Warn("项目未启用 MacOS, 项目配置 > 构建配置")
		return
	}
	isAmd64 := runtime.GOARCH == "amd64"
	isArm64 := runtime.GOARCH == "arm64"
	if isAmd64 {
		if !option.ArchX86_64 {
			logs.Warn("项目未启用架构 amd64, 项目配置 > 构建配置")
			return
		}
		os.Setenv("MACOSX_DEPLOYMENT_TARGET", "10.15")
		os.Setenv("CGO_CFLAGS", "-mmacosx-version-min=10.15")
		os.Setenv("CGO_LDFLAGS", "-mmacosx-version-min=10.15")
	}
	if isArm64 {
		if !option.ArchAarch64 {
			logs.Warn("项目未启用架构 arm64, 项目配置 > 构建配置")
			return
		}
		os.Setenv("MACOSX_DEPLOYMENT_TARGET", "11.0")
		os.Setenv("CGO_CFLAGS", "-mmacosx-version-min=11.0")
		os.Setenv("CGO_LDFLAGS", "-mmacosx-version-min=11.0")
	}
	if !option.UICocoa {
		logs.Warn("项目未启用 UI Cocoa, 项目配置 > 构建配置")
		return
	}
	logs.Info("构建项目, 开始构建项目", proj.Name)

	// go build
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	outputFilename := filepath.Join(output, option.BuildFileName)
	logs.Info("Building", outputFilename)
	cmd := command.NewCMD()
	cmd.Dir = bean.GPath
	args := []string{"build"}
	// go build -ldflags="-H windowsgui -s -w" -trimpath -o build/designer
	{
		buildArgs := option.GoArgs
		buildArgs = strings.ReplaceAll(buildArgs, "'", "\"")
		reTags := regexp.MustCompile(`-tags\s+([^\s-]+)`)
		reLdflags := regexp.MustCompile(`-ldflags\s+"([^"]+)"`)
		// 提取tags
		tagMatches := reTags.FindStringSubmatch(buildArgs)
		customTags := ""
		customLdflags := ""
		if len(tagMatches) > 1 {
			customTags = tagMatches[1]
		}
		ldMatches := reLdflags.FindStringSubmatch(buildArgs)
		if len(ldMatches) > 1 {
			customLdflags = ldMatches[1]
		}
		newTags := mergeTags(customTags, " a,b,d a proc")
		newLdflags := mergeLdflags(customLdflags, "-H windowsgui -s -w -X main.Version=v1.0.0")
		println("newtags:", strings.Join(newTags, " "))
		println("newldflags:", strings.Join(newLdflags, " "))
		// macOS 去除 -H windowsgui

	}

	if option.BuildModeDebug {

	} else {
		// -ldflags="-s -w" -trimpath
		args = append(args, "-trimpath")
	}
	args = append(args, "-o", outputFilename)
	cmd.Command("go", args...)
	cmd.Command("strip", outputFilename)
	logs.Info("Build Successfully")
}
