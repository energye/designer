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

package bean

import "github.com/energye/designer/pkg/tool"

// TProject 项目信息 xxx.egp 配置文件
type TProject struct {
	Name         string       `json:"name"`           // 项目名
	EGPName      string       `json:"egp_name"`       // 项目配置文件名
	Version      string       `json:"version"`        // 项目版本
	Description  string       `json:"description"`    // 项目描述
	Author       string       `json:"author"`         // 项目作者
	Package      string       `json:"package"`        // 项目(应用)包名
	Main         string       `json:"main"`           // 主程序入口文件或相对文件目录名
	UIForms      []TUIForm    `json:"ui_forms"`       // 窗体信息
	ActiveUIForm int          `json:"active_ui_form"` // 当前激活设计的窗体Id
	Lang         string       `json:"lang"`           // 语言 zh_CN
	BuildOption  TBuildOption `json:"build_option"`   // 构建配置
	EnvOption    TEnvOption   `json:"env_option"`     // 环境配置
}

// TUIForm 窗体信息
type TUIForm struct {
	Id         int    `json:"id"`           // 设计窗体Id
	Name       string `json:"name"`         // 窗体名
	UIFile     string `json:"ui_file"`      // UI文件名
	GOFile     string `json:"go_file"`      // UI Go 文件名
	GOUserFile string `json:"go_user_file"` // UI Go 用户文件名
	UpdateTime string `json:"date_time"`    // 更新时间
}

// TBuildOption 构建配置
type TBuildOption struct {
	GoArgument string `json:"go_argument"` // 构建参数
	Output     string `json:"output"`      // 构建输出目录
}

// TEnvOption 环境配置
type TEnvOption struct {
	GoRoot string `json:"go_root"` // Go 安装目录
}

// 模板调用 返回当前项目的所有窗体名称
func (m *TProject) GoFormNames() string {
	buf := tool.Buffer{}
	for _, form := range m.UIForms {
		buf.WriteString("&", form.Name, ", ")
	}
	return buf.String()
}
