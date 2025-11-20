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

package helperform

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 在设计器中创建项目
// 功能
// 主要: 在指定目录创建一个 energy 应用, 应用有默认模板
// 1. 应用名
// 2. 应用目录
// 3. 所需依赖模块(lcl, cef, webview), 从网络下载, 或设计器内绑定(✔️)
// 4. 模块模式： go.mod (✔️), go.work

/*
# 设计器安装目录（如 Windows：C:\EnergyDesigner，Linux：/opt/EnergyDesigner）
EnergyDesigner/
├── designer.exe  # 设计器主程序
└── frameworks/  # 内置框架根目录
	└── energy/  # 框架模块（必须是标准 Go 模块）
		├── go.mod  # 框架自身的 go.mod（如 module github.com/your-org/energy）
		├── go.sum
		├── core/  # 框架核心代码
		└── v1.2.3/  # （可选）多版本支持，每个版本一个独立模块
			├── go.mod
			└── core/

module my-app  // 项目自身的模块路径

go 1.21  // 项目使用的 Go 版本

// 声明框架依赖（版本需与内置框架的 go.mod 一致）
require github.com/your-org/energy v1.2.3

// 关键：将框架的网络模块路径，替换为设计器内置的本地路径
replace github.com/your-org/energy v1.2.3 => C:/EnergyDesigner/internal/frameworks/energy/v1.2.3

// 若框架无多版本，直接指向根目录：
// replace github.com/your-org/energy => C:/EnergyDesigner/internal/frameworks/energy

// 设计器新建 mod 模式项目时，配置内置框架
func createModProject(projectPath, frameworkBuiltInPath string) error {
    // 1. 创建项目目录
    if err := os.MkdirAll(projectPath, 0755); err != nil {
        return err
    }
    // 2. 初始化 go.mod
    cmd := exec.Command("go", "mod", "init", "my-app")
    cmd.Dir = projectPath
    if err := cmd.Run(); err != nil {
        return err
    }
    // 3. 写入 go.mod（添加 require 和 replace）
    modContent := fmt.Sprintf(`module my-app

go 1.21

require github.com/your-org/energy v1.2.3

replace github.com/your-org/energy v1.2.3 => %s`, frameworkBuiltInPath)
    return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(modContent), 0644)
}
*/

type TCreateProjectForm struct {
	lcl.TEngForm
}

func NewCreateProjectForm() *TCreateProjectForm {
	designerForm := &TCreateProjectForm{}
	lcl.Application.NewForm(designerForm)
	return designerForm
}

func (m *TCreateProjectForm) FormCreate(sender lcl.IObject) {
	logs.Info("TCreateProjectForm FormCreate")
	m.SetWidth(555)
	m.SetHeight(555)
	constr := m.Constraints()
	constr.SetMaxWidth(555)
	constr.SetMaxHeight(555)
	constr.SetMinWidth(555)
	constr.SetMinHeight(555)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.WorkAreaCenter()
}
