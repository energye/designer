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

// main.go 文件代码模板
const runCodeTemplate = `// ==============================================================================
// 📚 应用启动入口文件
// 📌 该文件在创建项目时创建
// ✏️ 可在此文件中添加业务逻辑
// ==============================================================================

package main

import (
	"github.com/energye/lcl/lcl" 
	"{{.Name}}/app"
	. "{{.Name}}/resources"
)

func main() {
	// 使用默认参数初始化LCL库
	lcl.Init(nil, nil)
	// 初始化应用程序实例
	lcl.Application.Initialize()
	// 配置应用程序设置，使主窗体在Windows任务栏上显示
	lcl.Application.SetMainFormOnTaskBar(true)
	// 启用自动缩放功能以支持高DPI显示器
	lcl.Application.SetScaled(true)
	// 设置应用程序图标
	SetIcon(lcl.Application)
	// 创建所有窗体
	lcl.Application.NewForms(app.Forms...)
	// 启动应用程序消息循环
	lcl.Application.Run()
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

import "github.com/energye/lcl/lcl"

// Forms 应用使用的窗体列表
var Forms = []lcl.IEngForm{
	{{.GoFormNames}}
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
//
//	app - 应用程序接口对象, 用于设置图标
func SetIcon(app lcl.IApp) {
	stream := lcl.NewMemoryStream()
	lcl.StreamHelper.Write(stream, Embed("icon.png"))
	stream.SetPosition(0)
	png := lcl.NewPortableNetworkGraphic()
	png.LoadFromStreamWithStream(stream)
	app.Icon().Assign(png)
	png.Free()
	stream.Free()
}
`
