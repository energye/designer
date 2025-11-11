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

package project

import (
	"github.com/energye/designer/designer"
)

// 项目文件 xxx.egp 配置文件, 存在项目根目录
// 功能说明
// 1. 创建项目
// 2. 加载项目
// 3. 更新项目信息
// 4. 更新项目的窗体信息
// 5. 更新其它配置和选项
// 6. UI 布局文件加载, 恢复到设计器
// 6.1 UI 布局同步更新: 新增/修改/删除

var (
	// 全局 Path 完整项目路径, 打开项目时设置. C:/YouProjectXxx/xxx.egp > C:/YouProjectXxx
	gPath string
	// 全局项目配置, 在创建或加载项目时设置
	gProject *TProject
)

// TProject 项目信息 xxx.egp 配置文件
type TProject struct {
	Name         string       `json:"name"`           // 项目名
	Version      string       `json:"version"`        // 项目版本
	Description  string       `json:"description"`    // 项目描述
	Author       string       `json:"author"`         // 项目作者
	Package      string       `json:"package"`        // 项目(应用)包名
	Main         string       `json:"main"`           // 主程序入口文件或相对文件目录名
	UIForms      []TUIForm    `json:"forms"`          // 窗体信息
	ActiveUIForm string       `json:"active_ui_form"` // 当前激活设计的窗体名称
	Lang         string       `json:"lang"`           // 语言 zh_CN
	BuildOption  TBuildOption `json:"build_option"`   // 构建配置
	EnvOption    TEnvOption   `json:"env_option"`     // 环境配置
}

// 当前项目的设计器配置
type TDesignerConfig struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
}

// TUIForm 窗体信息
type TUIForm struct {
	Name       string `json:"name"`      // 窗体名
	UIFile     string `json:"ui_file"`   // UI文件名
	GOFile     string `json:"go_file"`   // UI Go 文件名
	UpdateTime string `json:"date_time"` // 更新时间
	FilePath   string `json:"file_path"` // 文件路径
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

// SetGlobalProject 设置全局项目路径和项目对象
// path: 项目路径
// project: 项目对象指针
func SetGlobalProject(path string, project *TProject) {
	gPath = path
	gProject = project
	if path == "" || project == nil {
		designer.SetEnableFuncComponent(false)
	} else {
		designer.SetEnableFuncComponent(true)
	}
}

func Path() string {
	return gPath
}

func Project() *TProject {
	return gProject
}
