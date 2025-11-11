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
	"encoding/json"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/tool"
	"os"
	"path/filepath"
	"strings"
)

// 项目文件 xxx.egp 配置文件
// 存在于项目根目录

var (
	// Path 完整项目路径, 打开项目时设置. C:/YouProjectXxx/xxx.egp
	Path    string
	Project TProject
)

// 项目配置文件扩展名
const egp = ".egp"

func init() {
	// TODO 需要通过配置 --test
	if tool.IsWindows() {
		Path = "C:\\app\\workspace\\test" // TODO 测试
	} else if tool.IsLinux() {
		Path = "/home/yanghy/app/projects/workspace/test"
	}
}

// TProject 项目信息 xxx.egp 配置文件
type TProject struct {
	Name         string       `json:"name"`           // 项目名
	Version      string       `json:"version"`        // 项目版本
	Description  string       `json:"description"`    // 项目描述
	Author       string       `json:"author"`         // 项目作者
	Main         string       `json:"main"`           // 主程序入口文件或相对文件目录名
	UIForms      []*TUIForm   `json:"forms"`          // 窗体信息
	ActiveUIForm string       `json:"active_ui_form"` // 当前激活设计的窗体名称
	Lang         string       `json:"lang"`           // 语言
	BuildOption  TBuildOption `json:"build_option"`   // 构建选项
	EnvOption    TEnvOption   `json:"env_option"`     // 环境配置

}

// TUIForm 窗体信息
type TUIForm struct {
	Name       string `json:"name"`      // 窗体名
	UIFile     string `json:"ui_file"`   // UI文件名
	GOFile     string `json:"go_file"`   // UI Go 文件名
	UpdateTime string `json:"date_time"` // 更新时间
	FilePath   string `json:"file_path"` // 文件路径
}

// TBuildOption 构建选项
type TBuildOption struct {
	GoArgument string `json:"go_argument"` // 构建参数
	Output     string `json:"output"`      // 构建输出目录
}

// TEnvOption 环境配置
type TEnvOption struct {
}

// 加载项目
func Load(egpPath string) {
	if tool.IsExist(egpPath) {
		isEgp := strings.ToLower(filepath.Ext(egpPath)) == egp
		if !isEgp {
			logs.Warn("文件目录非 .egp 项目配置文件")
			return
		}
		data, err := os.ReadFile(egpPath)
		if err != nil {
			logs.Error("读取项目配置文件失败:", err)
			return
		}
		err = json.Unmarshal(data, &Project)
		if err != nil {
			logs.Error("解析项目配置文件失败:", err)
			return
		}

	}
}

// 写入项目配置文件
func Write(path string, project *TProject) {
	if project == nil {
		return
	}
	data, err := json.Marshal(project)
	if err != nil {
		logs.Error("解析项目配置文件失败:", err.Error())
		return
	}
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		logs.Error("写入项目配置文件失败:", err.Error())
		return
	}
}
