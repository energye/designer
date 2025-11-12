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
	"github.com/energye/designer/pkg/logs"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// 创建项目目录结构
/*
	[app]	// 应用主目录, 生成代码存放目录 (xxx.go xxx.ui.go xxx.ui)
	[resources]	// 资源存放目录, 图标等静态资源文件
		| icon
			| icon.md
		| icon.go
		| windows_[386|amd64].syso ?? 根据设计器功能动态生成, 只适用于 windows
	go.mod
	main.go
*/
func createProjectDir() {
	if gProject == nil || gPath == "" {
		return
	}
	appRoot := gPath
	// 代码存放目录
	appCodePath := filepath.Join(appRoot, consts.AppPackageName)
	// 资源存放目录
	resourcesPath := filepath.Join(appRoot, "resources")
	resourcesIconPath := filepath.Join(resourcesPath, "icon")
	paths := []string{appCodePath, resourcesPath, resourcesIconPath}
	for _, path := range paths {
		if err := os.Mkdir(path, fs.ModePerm); err != nil {
			logs.Error("创建项目目录失败:", err.Error())
		}
	}
	// 文件创建
	files := []struct {
		path string
		name string
		data string
	}{
		{appCodePath, consts.FormListFileName, buildTemplateData(appCodeTemplate)},
		{resourcesPath, "resources.go", buildTemplateData(resourcesGoTemplate)},
		{resourcesIconPath, "icon.md", ""},
		{appRoot, "go.mod", buildTemplateData(goModTemplate)},
		{appRoot, "main.go", buildTemplateData(runCodeTemplate)},
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(file.path, file.name), []byte(file.data), 0666); err != nil {
			logs.Error("创建项目文件失败:", err.Error())
		}
	}
}

// 构建填充模板数据
func buildTemplateData(templateData string) string {
	// 解析模板
	tmpl, err := template.New("project").Parse(templateData)
	if err != nil {
		logs.Error("解析自动代码模板失败:", err.Error())
		return ""
	}

	// 生成代码
	var buf strings.Builder
	if err := tmpl.Execute(&buf, gProject); err != nil {
		logs.Error("执行自动代码模板失败:", err.Error())
		return ""
	}

	return buf.String()
}
