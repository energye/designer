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
	Name        string   `json:"name"`        // 项目名称
	Version     string   `json:"version"`     // 项目版本
	Description string   `json:"description"` // 项目描述
	Author      string   `json:"author"`      // 项目作者
	Main        string   `json:"main"`        // 主程序入口文件或相对文件目录名
	Forms       []*TForm `json:"forms"`       // 窗体信息
}

// TForm 窗体信息
type TForm struct {
	Name       string `json:"name"`      // 窗体名称
	UpdateTime string `json:"date_time"` // 更新时间
	Active     bool   `json:"active"`    // 是否激活设计, 同一时间只允许一个窗体激活设计
}

func Create() {

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
