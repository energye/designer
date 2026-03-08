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

func LoadProject(filePath string) {
	if bean.GPath != "" && bean.GProject != nil {
		return
	}
	logs.Info("加载项目配置文件:", filePath)
	if filePath == "" {
		logs.Error("项目配置文件 .egp 路径为空")
		return
	}
	if tool.IsExist(filePath) {
		path, file := filepath.Split(filePath)
		isEgp := strings.ToLower(filepath.Ext(file)) == consts.EGPExt
		if isEgp {
			data, err := os.ReadFile(filePath)
			if err != nil {
				logs.Error("读取项目配置文件失败:", err)
				return
			}
			project := &bean.TProject{}
			err = json.Unmarshal(data, project)
			if err != nil {
				logs.Error("解析项目配置文件失败:", err)
				return
			}
			bean.GPath = path
			bean.GProject = project
		} else {
			logs.Error("非 .egp 项目配置文件:", filePath)
		}
	} else {
		logs.Error("项目配置文件 .egp 不存在")
	}
}
