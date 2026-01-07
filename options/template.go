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

// main.go 文件代码模板
const runLCLCodeTemplate = `// ==============================================================================
// 📚 应用启动入口文件
// 📌 该文件在创建项目时创建
// ✏️ 可在此文件中添加业务逻辑
// ==============================================================================

package main

import (
	"github.com/energye/energy/v3/lcl"
	"{{.Name}}/app"
	_ "{{.Name}}/resources"
)

func main() {
	// 全局初始化
	lcl.Init(nil, nil)
	// 启动应用程序消息循环
	lcl.Run(app.Forms...)
}
`

// main.go 文件代码模板
const runWVCodeTemplate = `// ==============================================================================
// 📚 应用启动入口文件
// 📌 该文件在创建项目时创建
// ✏️ 可在此文件中添加业务逻辑
// ==============================================================================

package main

import (
	"github.com/energye/energy/v3/application"
	"github.com/energye/energy/v3/wv"
	"{{.Name}}/app"
	_ "{{.Name}}/resources"
)

func main() {
	// 全局初始化
	wvApp := wv.Init(nil, nil)
	wvApp.SetOptions(application.Options{DefaultURL: "about:blank"})
	// 启动应用程序消息循环
	wv.Run(app.Forms...)
}
`

// app.go 文件代码模板
// 用于提供 main.go NewForms 参数使用
// 在项目窗体创建/更新时同步修改
const appCodeTemplate = `// ==============================================================================
// 📚 窗体维护列表
// 🔥 ENERGY GUI 设计器自动生成代码. 不能编辑
// ==============================================================================

package {{.Package}}

import (
	"github.com/energye/lcl/lcl"
	"os"
)

// Forms 应用使用的窗体列表
var Forms = []lcl.IEngForm{
	{{.GoFormNames}}
}

func init() {
	if "{{.GUIRenderFramework}}" == "WV" {
		// linux webkit2 > gtk3
		os.Setenv("--ws", "gtk3")
	}
}
`

// go.mod 模块文件模板
const goModTemplate = `module {{.Name}}

go 1.20

{{.Data}}
`

// resources/resources.go
// 资源代码模板
const resourcesGoTemplate = `// ==============================================================================
// 📚 内嵌资源
// 📌 不存在时自动创建
// ✏️ 可在此文件中添加业务逻辑
// ==============================================================================

package resources

import (
	"embed"
	engLCL "github.com/energye/energy/v3/lcl"
	"github.com/energye/lcl/lcl"
)

//go:embed embed
var icon embed.FS

// Embed 获取内嵌资源
// 函数签名不能修改
func Embed(fileName string) []byte {
	data, _ := icon.ReadFile("embed/" + fileName)
	return data
}

// SetIcon 设置应用程序图标
// 函数签名不能修改
func SetIcon() {
	stream := lcl.NewMemoryStream()
	lcl.StreamHelper.Write(stream, Embed("icon.png"))
	stream.SetPosition(0)
	png := lcl.NewPortableNetworkGraphic()
	png.LoadFromStreamWithStream(stream)
	lcl.Application.Icon().Assign(png)
	png.Free()
	stream.Free()
}

func init() {
	engLCL.SetOnBeforeRun(SetIcon)
}
`
