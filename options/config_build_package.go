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
	"github.com/energye/designer/cmd/packager"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
)

// configBuildPackage 执行构建打包流程
// 只给当前系统和架构构建打包
func configBuildPackage() {
	logs.Debug("构建配置-打包")
	event.ConsoleWriteClear()
	if envs, ok := packagePlatformENVs[lib.GOOS()]; ok {
		pack := packager.Default()
		pack.AppendPlatform = true
		pack.AppendArch = true
		for _, arch := range envs {
			_ = os.Setenv("GOARCH", arch)
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
