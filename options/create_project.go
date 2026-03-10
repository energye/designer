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

package options

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
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
		| embed
			| embed.md
		| resources.go
		| windows_[386|amd64].syso ?? 根据设计器功能动态生成, 只适用于 windows
	go.mod
	main.go
*/
func createProjectDir() {
	if bean.GProject == nil || bean.GPath == "" {
		return
	}
	appRoot := bean.GPath
	// 代码存放目录
	appCodePath := filepath.Join(appRoot, consts.AppPackageName)
	// 资源存放目录
	resourcesPath := ResourcePath()
	resourcesEmbedPath := filepath.Join(resourcesPath, "embed")
	paths := []string{appCodePath, resourcesPath, resourcesEmbedPath}
	for _, path := range paths {
		if err := os.Mkdir(path, fs.ModePerm); err != nil {
			logs.Error("创建项目目录失败:", err.Error())
		}
	}
	// 模板数据
	data := *bean.GProject
	// 本地模式
	localModule := tool.Buffer{}
	localModule.WriteString("replace (", "\n")
	localModule.WriteString("  github.com/energye/lcl", " => ", config.Config.FrameworkDirForLCL(), "\n")
	localModule.WriteString("  github.com/energye/wv", " => ", config.Config.FrameworkDirForWV(), "\n")
	localModule.WriteString("  github.com/energye/cef", " => ", config.Config.FrameworkDirForCEF(), "\n")
	localModule.WriteString("  github.com/energye/energy/v3", " => ", config.Config.FrameworkDirForENERGY(), "\n")
	localModule.WriteString(")", "\n")
	data.Data = localModule.String()

	// 文件创建
	type TFile struct {
		path string // 文件目录
		name string // 文件名
		data string // 文件内容
	}
	files := []TFile{
		{appCodePath, consts.FormListFileName, buildTemplateData(appCodeTemplate, &data)},
		{resourcesPath, "resources.go", buildTemplateData(resourcesGoTemplate, &data)},
		{resourcesEmbedPath, "embed.md", ""},
		{appRoot, "go.mod", buildTemplateData(goModTemplate, &data)},
		//{appRoot, "main.go", buildTemplateData(runCodeTemplate, &data)},
	}
	switch data.GUIRenderFramework {
	case bean.GUIRenderFramework_LCL:
		files = append(files, TFile{appRoot, "main.go", buildTemplateData(runLCLCodeTemplate, &data)})
	case bean.GUIRenderFramework_WV:
		files = append(files, TFile{appRoot, "main.go", buildTemplateData(runWVCodeTemplate, &data)})
	case bean.GUIRenderFramework_CEF:
		files = append(files, TFile{appRoot, "main.go", buildTemplateData(runLCLCodeTemplate, &data)})
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(file.path, file.name), []byte(file.data), 0666); err != nil {
			event.ConsoleWriteError("创建项目文件失败:", err.Error())
		}
	}
}

// 构建填充模板数据
func buildTemplateData(templateData string, data any) string {
	// 解析模板
	tmpl, err := template.New("project").Parse(templateData)
	if err != nil {
		logs.Error("解析自动代码模板失败:", err.Error())
		return ""
	}

	// 生成代码
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		logs.Error("执行自动代码模板失败:", err.Error())
		return ""
	}

	return buf.String()
}
