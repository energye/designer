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
	"fmt"
	"github.com/energye/designer/designer"
	codegentpl "github.com/energye/designer/internal/templates/codegen"
	projBean "github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/uigen/bean"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// 生成用户代码文件
// 生成条件: 文件未创建, 绑定事件, self 修改
func generateUserCode(formTab *designer.FormTab, component *bean.TUIComponent) error {
	goUIUserFilePath := filepath.Join(projBean.CodePath(), formTab.GOUserFile())
	// 检查文件是否已存在
	// 如果文件已存在，不覆盖
	if _, err := os.Stat(goUIUserFilePath); err == nil {
		return nil // 文件已存在，直接返回
	}
	// 创建用户代码文件
	// 构建模板数据
	data := buildUserTemplateData(component)
	data.BaseInfo = &TBaseInfo{
		DesignerVersion: config.DesignerConfig.Version, DateTime: time.Now().Format("2006-01-02 15:04:05"),
		UIFile:   formTab.UIFile(),
		UserFile: formTab.GOUserFile(),
	}
	data.Imports.Add(lcl)
	data.IncludePackage()

	// 解析模板
	tmpl, err := template.New("user").Parse(codegentpl.UserCodeTemplate)
	if err != nil {
		return fmt.Errorf("解析用户代码模板失败: %w", err)
	}

	// 生成代码
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行用户代码模板失败: %w", err)
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		logs.Error("UserCode 格式化代码失败:", err.Error(), "id:", formTab.Id, "name:", formTab.FormRoot.Name())
		// 错误的代码也写到文件, 有助于调试
		formatted = []byte(buf.String())
	}

	// 写入文件
	if err := os.WriteFile(goUIUserFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("写入用户代码文件失败: %w", err)
	}

	return nil
}

// 构建用户代码模板数据
func buildUserTemplateData(component *bean.TUIComponent) *TFormData {
	formData := &TFormData{BaseInfo: &TBaseInfo{}, PackageName: projBean.GProject.Package, Imports: tool.NewHashSet[string]()}
	formData.Form = &TComponentData{
		Name:       component.Name,
		ClassName:  component.ClassName,
		Mod:        component.Mod,
		Type:       component.Type,
		Properties: uiPropertiesToTemplateProperties(component.Properties),
		Children:   make([]*TComponentData, 0),
	}
	formData.Form.Children = formData.Form.buildComponents(component)
	return formData
}
