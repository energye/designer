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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/uigen"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	packageName    = "forms"                                          // 包名
	lcl            = `"github.com/energye/lcl/lcl"`                   // lcl 组件库
	cef            = `"github.com/energye/cef"`                       // cef 组件库
	wvWindows      = `"github.com/energye/wv/windows"`                // webview2 组件库
	wvLinux        = `"github.com/energye/wv/linux"`                  // gtk webkit2 组件库
	wvDarwin       = `"github.com/energye/wv/darwin"`                 // cocoa darwin 组件库
	lclTypes       = `lclTypes "github.com/energye/lcl/types"`        // lcl 类型
	cefTypes       = `cefTypes "github.com/energye/cef/types"`        // cef 类型
	wvTypesWindows = `cefTypes "github.com/energye/wv/types/windows"` // webview2 windows 类型
	wvTypesLinux   = `cefTypes "github.com/energye/wv/types/linux"`   // gtk webkit2 linux 类型
	wvTypesDarwin  = `cefTypes "github.com/energye/wv/types/darwin"`  // cocoa webkit2 darwin 类型
)

// go 代码生成 自动时时生成
// 依赖 uigen UI 布局文件
// 生成触发条件: UI 布局文件修改后

// 根据UI文件生成Go代码
func runGenerateCode(uiFilePath string) error {
	// 读取并解析UI文件
	data, err := os.ReadFile(uiFilePath)
	if err != nil {
		return fmt.Errorf("读取UI文件失败: %w", err)
	}

	var uiComponent uigen.TUIComponent
	if err := json.Unmarshal(data, &uiComponent); err != nil {
		return fmt.Errorf("解析UI文件失败: %w", err)
	}

	// 生成自动代码文件 (form[n].ui.go)
	if err := generateAutoCode(uiFilePath, &uiComponent); err != nil {
		return fmt.Errorf("生成自动代码失败: %w", err)
	}

	// 生成用户代码文件 (form[n].go) - 仅当文件不存在时
	if err := generateUserCode(uiFilePath, &uiComponent); err != nil {
		return fmt.Errorf("生成用户代码失败: %w", err)
	}

	return nil
}

// 生成自动代码文件
func generateAutoCode(uiFilePath string, component *uigen.TUIComponent) error {
	baseName := strings.TrimSuffix(filepath.Base(uiFilePath), filepath.Ext(uiFilePath))
	// 构建模板数据
	data := buildAutoTemplateData(component)
	data.BaseInfo = &TBaseInfo{
		DesignerVersion: config.Config.Version, DateTime: time.Now().Format("2006-01-02 15:04:05"),
		UIFile: baseName + ".ui", UserFile: baseName + ".go",
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
		//return fmt.Errorf("格式化代码失败: %w", err)
		logs.Error("格式化代码失败:", err.Error())
		formatted = []byte(buf.String())
	}

	// 写入文件
	autoFileName := baseName + ".ui.go"
	autoFilePath := filepath.Join(filepath.Dir(uiFilePath), autoFileName)

	if err := os.WriteFile(autoFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("写入自动代码文件失败: %w", err)
	}

	return nil
}

// generateUserCode 生成用户代码文件
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
