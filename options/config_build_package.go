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

package options

import (
	"context"
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/cmd/packager"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"runtime"
)

// 执行构建打包流程
//
// 只给当前系统和架构构建打包
func configBuildPackage(ctx context.Context) {
	logs.Debug("构建配置-打包")
	if envs, ok := packagePlatformENVs[runtime.GOOS]; ok {
		pack := packager.Default()
		pack.AppendPlatform = true
		pack.AppendArch = true
		pack.Context = ctx

		defer env.Clear()

		// 当项目配置为禁用 CGO 且启用了跨平台构建时，执行当前系统其他架构编译
		cgoEnabled := bean.GProject.BuildOption.BuildCGOEnabled
		//buildOtherPlatform := bean.GProject.BuildOption.BuildOtherPlatform
		for _, arch := range envs {
			env.Put("GOARCH", arch)
			env.Put("GOOS", runtime.GOOS)
			if cgoEnabled {
				env.Put("CGO_ENABLED", "1")
			} else {
				env.Put("CGO_ENABLED", "0")
			}
			if runtime.GOOS == "linux" && cgoEnabled && runtime.GOARCH != arch {
				// linux 启用 CGO 当前架构和目标架构不同 忽略其他架构
				event.ConsoleWriteWarn("Build and Package - Linux CGO=ENABLED or Platform=false skip-arch:", arch, "current-arch:", runtime.GOARCH)
				continue
			}
			packager.Run(pack)
		}
		logs.Debug("构建配置-打包-完成")
	} else {
		logs.Debug("构建配置-打包失败, 不支持的系统")
	}
}

var packagePlatformENVs = map[string][]string{
	"windows": {"amd64", "386"},
	"darwin":  {"amd64", "arm64"},
	"linux":   {"amd64", "386", "arm64", "arm"},
}
