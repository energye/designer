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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project"
	"github.com/energye/designer/uigen/bean"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// 构建模板数据

// go 代码生成 自动时时生成
// 依赖 uigen UI 布局文件
// 生成触发条件: UI 布局文件修改后

// 生成自动代码文件
func generateAutoCode(formTab *designer.FormTab, component *bean.TUIComponent) error {
	// 构建模板数据
	data := buildAutoTemplateData(component)
	data.BaseInfo = &TBaseInfo{
		DesignerVersion: config.Config.Version, DateTime: time.Now().Format("2006-01-02 15:04:05"),
		UIFile:   formTab.UIFile(),
		UserFile: formTab.GOUserFile(),
	}
	data.Imports.Add(lcl)
	data.IncludePackage()

	// 解析模板
	tmpl, err := template.New("auto").Parse(autoCodeTemplate)
	if err != nil {
		return fmt.Errorf("解析自动代码模板失败: %w", err)
	}

	// 生成代码
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行自动代码模板失败: %w", err)
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		logs.Error("AutoCode 格式化代码失败:", err.Error())
		formatted = []byte(buf.String())
	}

	// 写入文件
	goUIFilePath := filepath.Join(project.Path(), project.Project().Package, formTab.GOFile())
	if err := os.WriteFile(goUIFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("写入自动代码文件失败: %w", err)
	}

	return nil
}

// 构建自动代码模板数据
func buildAutoTemplateData(component *bean.TUIComponent) *TFormData {
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
