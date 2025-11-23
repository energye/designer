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

package designer

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/tool/exec"
)

// 初始化依赖模块的一些数据信息
func initDependencyModule() {
	go initModuleTypeInfo()
}

// 初始化模块类型信息
func initModuleTypeInfo() {
	// go list -m -f '{{.Dir}}' github.com/energye/widget
	// go list -f '{{.Path}}  {{.Dir}}' -m all
	cmd := command.NewCMD()
	cmd.Dir = exec.Dir
	cmd.HideWindow = true
	cmd.IsPrint = false
	cmd.Console = func(data string, level command.Level) {
		logs.Info("initModuleTypeInfo:", level, data)
	}
	cmd.Command("go", "list", "-m", "-f", "{{.Dir}}", consts.DmLCL)

}
