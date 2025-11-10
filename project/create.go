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
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/types"
	"os"
	"path/filepath"
	"strings"
)

// 项目创建, 在指定目录创建新项目
// 检查目录是否为空

// 创建项目
func runCreate(dir string) {
	logs.Debug("运行创建项目 目录:", dir)
	if !tool.IsExist(dir) {
		logs.Error("目录不存在:", dir)
		return
	}
	de, err := os.ReadDir(dir)
	if err != nil {
		logs.Error("读取目录失败:", err.Error())
		return
	}
	var (
		isNotEmpty bool   // 当前目录是否为空
		existEgp   string // 目录的egp文件
		isCreate   bool   // 是否创建
	)
	for _, entry := range de {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.LastIndex(name, egp) != -1 {
			existEgp = entry.Name()
			break
		}
		isNotEmpty = true // 非空目录
	}
	// 已存在项目 egp 文件, 提示覆盖
	if existEgp != "" {
		msg := fmt.Sprintf("当前目录已存在项目配置 %s\n是否覆盖？", existEgp)
		logs.Warn("当前目录已存在项目配置")
		isCreate = api.MessageDlg(msg, types.MtCustom, types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
	}
	if isCreate {
		// 覆盖并创建项目
		logs.Warn("当前目录已存在项目配置", existEgp, "创建并覆盖")
	} else if isNotEmpty {
		// 目录非空并且没有项目配置文件 egp, 提示是否在当前目录创建项目
		logs.Warn("当前目录非空")
		isCreate = api.MessageDlg("当前目录非空是否创建？", types.MtCustom, types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
	}
	if isCreate {
		// 设置项目目录
		Path = dir
		logs.Info("开始创建项目")
		_, name := filepath.Split(dir)
		newProject := new(TProject)
		newProject.Name = name

	} else {
		logs.Info("取消创建项目")
	}
}
