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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/pkg/logs"
	"os"
	"path/filepath"
	"time"
)

// 项目配置更新, 在窗体创建/更新时运行
// 1. 更新 TUIForm 配置
// 2. 更新 app/app.go 代码文件
func runUpdate(formTab *designer.FormTab) {
	logs.Debug("运行项目更新")
	uiForms := gProject.UIForms // 窗体列表
	isExist := false            // 是否存在
	for i := range uiForms {
		// 更新
		// TODO 删除和其它待增加
		if uiForms[i].Name == formTab.Name {
			uiForms[i].UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			isExist = true
			break
		}
	}
	if !isExist {
		// 不存在添加
		uiForms = append(uiForms, TUIForm{
			Name:       formTab.Name,
			UIFile:     formTab.UIFile(),
			GOFile:     formTab.GOFile(),
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	gProject.UIForms = uiForms
	// 1. 更新 TUIForm 配置
	if err := Write(gPath, gProject); err != nil {
		logs.Error("项目更新, 写入项目配置失败:", err.Error())
		return
	}
	// 2. 更新 app/app.go 代码文件
	appRoot := gPath
	appCodePath := filepath.Join(appRoot, consts.AppPackageName, consts.FormListFileName)
	code := buildTemplateData(appCodeTemplate)
	if err := os.WriteFile(appCodePath, []byte(code), 0666); err != nil {
		logs.Error("创建项目文件失败:", err.Error())
	}
}
