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

package codegen

import (
	"encoding/json"
	"fmt"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/project"
	"github.com/energye/designer/uigen/bean"
	"os"
	"path/filepath"
)

// 根据UI文件生成Go代码
func runGenerateCode(formTab *designer.FormTab) error {
	// ui 布局文件名
	goUIFileName := formTab.GOFile()

	// 验证UI布局文件名
	var uiForm *project.TUIForm
	for _, form := range project.Project().UIForms {
		if form.Id == formTab.Id {
			uiForm = &form
			break
		}
	}
	if uiForm != nil && uiForm.GOFile != goUIFileName {
		// 修改 xxx.ui.go 布局文件名
		oldGoUIFilePath := filepath.Join(project.Path(), project.Project().Package, uiForm.GOFile)
		newGoUIFilePath := filepath.Join(project.Path(), project.Project().Package, goUIFileName)
		if err := os.Rename(oldGoUIFilePath, newGoUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return err
		}
		uiForm.GOFile = formTab.GOFile()
		// 修改 xxx.go 用户代码文件名
		oldGoUIUserFilePath := filepath.Join(project.Path(), project.Project().Package, uiForm.GOUserFile)
		newGoUIUserFilePath := filepath.Join(project.Path(), project.Project().Package, formTab.GOUserFile())
		if err := os.Rename(oldGoUIUserFilePath, newGoUIUserFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return err
		}
		uiForm.GOUserFile = formTab.GOUserFile()
	}
	uiFilePath := filepath.Join(project.Path(), project.Project().Package, formTab.UIFile())
	// 读取并解析UI文件
	data, err := os.ReadFile(uiFilePath)
	if err != nil {
		return fmt.Errorf("读取UI文件失败: %w", err)
	}
	// 解析UI文件
	var uiComponent bean.TUIComponent
	if err := json.Unmarshal(data, &uiComponent); err != nil {
		return fmt.Errorf("解析UI文件失败: %w", err)
	}
	// 生成自动代码文件 (xxx.ui.go)
	if err := generateAutoCode(formTab, &uiComponent); err != nil {
		return fmt.Errorf("生成自动代码失败: %w", err)
	}
	// 生成用户代码文件 (xxx.go) - 仅当文件不存在时
	if err := generateUserCode(formTab, &uiComponent); err != nil {
		return fmt.Errorf("生成用户代码失败: %w", err)
	}

	return nil
}
