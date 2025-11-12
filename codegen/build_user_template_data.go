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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project"
	"github.com/energye/designer/uigen"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// 生成用户代码文件
func generateUserCode(uiFilePath string, component *uigen.TUIComponent) error {
	// 检查文件是否已存在
	baseName := strings.TrimSuffix(filepath.Base(uiFilePath), filepath.Ext(uiFilePath))
	userFileName := baseName + ".go"
	userFilePath := filepath.Join(filepath.Dir(uiFilePath), userFileName)

	// 如果文件已存在，不覆盖
	if _, err := os.Stat(userFilePath); err == nil {
		return nil // 文件已存在，直接返回
	}

	// 构建模板数据
	data := buildUserTemplateData(component)
	data.BaseInfo = &TBaseInfo{
		DesignerVersion: config.Config.Version, DateTime: time.Now().Format("2006-01-02 15:04:05"),
		UIFile: baseName + ".ui", UserFile: baseName + ".go",
	}
	data.Imports.Add(lcl)
	data.IncludePackage()

	// 解析模板
	tmpl, err := template.New("user").Parse(userCodeTemplate)
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
		return fmt.Errorf("格式化代码失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(userFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("写入用户代码文件失败: %w", err)
	}

	return nil
}

// 构建用户代码模板数据
func buildUserTemplateData(component *uigen.TUIComponent) *TFormData {
	formData := &TFormData{BaseInfo: &TBaseInfo{}, PackageName: project.Project().Package, Imports: tool.NewHashSet()}
	formData.Form = &TComponentData{
		Name:       component.Name,
		ClassName:  component.ClassName,
		Type:       component.Type,
		Properties: uiPropertiesToTemplateProperties(component.Properties),
		Children:   make([]*TComponentData, 0),
	}
	formData.Form.Children = formData.Form.buildComponents(component)
	return formData
}
