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
	"encoding/json"
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/tool/command"
	"os"
	"path/filepath"
)

// Run 执行构建命令的入口函数
func Run() bool {
	event.ConsoleWriteInfo("CMD-build-run")
	goos := lib.GOOS()
	if goos == "windows" {
		return buildWindows()
	} else if goos == "darwin" {
		return buildDarwin()
	} else if goos == "linux" {
		return buildLinux()
	}
	return false
}

func RunClean() {
	event.ConsoleWriteInfo("CMD-build-clean")
	RunGoCleanCacheCMD()
	Run()
}

// 编译打包信息到执行文件
func xBuildPackVar() (result []string) {
	if bean.GProject == nil {
		return nil
	}
	packMap := make(map[string]string)
	packMap["name"] = bean.GProject.Name
	packMap["id"] = bean.GProject.AppOption.Id
	packMap["version"] = bean.GProject.AppOption.Version
	packMap["arch"] = lib.GOARCH()
	data, _ := json.Marshal(packMap)
	result = append(result, "-X github.com/energye/energy/v3/application/pack.JSON="+string(data))
	return
}

// buildBinFileName 根据构建选项和当前运行环境生成最终的可执行文件名称。
// 如果启用了跨平台构建标记，文件名会自动追加操作系统和架构后缀；
// 若目标平台为 Windows，则确保文件名以 .exe 结尾。
//
// 参数:
//   - buildOption: 包含构建配置信息的选项结构体，用于获取基础文件名。
func buildBinFileName(buildOption bean.TBuildOption) string {
	buildFileName := buildOption.BuildFileName
	goos := lib.GOOS()
	goarch := lib.GOARCH()
	// 在跨平台构建模式下（通常对应 CGO 禁用场景），
	// 为文件名添加操作系统和架构标识，以区分不同平台的产物。
	if IsBuildAllPlatform() {
		// CGO disable
		buildFileName = fmt.Sprintf("%s_%s_%s", buildFileName, goos, goarch)
	}
	// 针对 Windows 平台进行特殊处理，确保生成的二进制文件
	// 具有正确的 .exe 扩展名。
	if goos == "windows" && filepath.Ext(buildFileName) != ".exe" {
		buildFileName += ".exe"
	}
	return buildFileName
}

const buildAllPlatform = "IsBuildAllPlatform"

// IsBuildAllPlatform 检查当前是否处于构建其他平台的环境。
// 通过读取环境变量 "IsBuildOtherPlatform" 来判断，如果值为 "true" 则返回 true。
func IsBuildAllPlatform() bool {
	return os.Getenv(buildAllPlatform) == "true"
}

func RunGoCleanCacheCMD() {
	cmd := command.NewCMD()
	cmd.HideWindow = true
	cmd.Command("go", "clean", "--cache")
}
