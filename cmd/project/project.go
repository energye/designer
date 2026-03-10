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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"os"
	"path/filepath"
	"strings"
)

// LoadProject 加载项目配置文件
//
//	该函数用于读取和解析 .egp 格式的项目配置文件，并将配置信息存储到全局变量中
//	如果项目已经加载过（GPath 不为空且 GProject 不为 nil），则直接返回避免重复加载
//	成功加载后，会将路径和项目名称分别设置到 bean.GPath 和 bean.GProject
//	filePath - 项目配置文件的路径，应为 .egp 后缀的 JSON 格式文件
func LoadProject(filePath string) {
	if bean.GPath != "" && bean.GProject != nil {
		return
	}
	logs.Info("Load project Config:", filePath)
	if filePath == "" {
		logs.Error("Project config .egp path is nil")
		return
	}
	if tool.IsExist(filePath) {
		path, file := filepath.Split(filePath)
		isEgp := strings.ToLower(filepath.Ext(file)) == consts.EGPExt
		if isEgp {
			data, err := os.ReadFile(filePath)
			if err != nil {
				logs.Error("Read project config file error:", err)
				return
			}
			project := &bean.TProject{}
			err = json.Unmarshal(data, project)
			if err != nil {
				logs.Error("Unmarshal project config file error:", err)
				return
			}
			bean.GPath = path
			bean.GProject = project
		} else {
			logs.Error("not .egp project config file:", filePath)
		}
	} else {
		logs.Error("project config file .egp not exist", filePath)
	}
}
