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
	"github.com/energye/designer/pkg/tool"
	"os"
	"path/filepath"
	"strings"
)

func runLoad(filePath string) {
	logs.Debug("运行加载项目/UI 文件目录:", filePath)
	Load(filePath)
}

// 加载项目
func Load(egpPath string) {
	if tool.IsExist(egpPath) {
		path, file := filepath.Split(egpPath)
		isEgp := strings.ToLower(filepath.Ext(file)) == egp
		if !isEgp {
			logs.Warn("文件目录非 .egp 项目配置文件")
			SetGlobalProject("", nil)
			return
		}
		data, err := os.ReadFile(egpPath)
		if err != nil {
			logs.Error("读取项目配置文件失败:", err)
			SetGlobalProject("", nil)
			return
		}
		loadProject := &TProject{}
		err = json.Unmarshal(data, loadProject)
		if err != nil {
			logs.Error("解析项目配置文件失败:", err)
			SetGlobalProject("", nil)
			return
		}
		SetGlobalProject(path, loadProject)
	}
}
