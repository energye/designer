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
// 📌 该文件不存在时自动创建
// ✏️ 可在此文件中添加业务逻辑
// ==============================================================================

package main

import (
	"github.com/energye/lcl/lcl" 
	{{.WindowsSyso}}
)

func main() {
	lcl.Init(nil, nil)
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms({{.Forms}})
	lcl.Application.Run()
}
`

// go.mod 模块文件模板
const goModTemplate = `module {{.Module}}

go 1.20

require (
	github.com/energye/lcl/lcl
)
`

// resources/resources。go
// 资源代码模板
const resourcesGoTemplate = `// ==============================================================================
// 📚 项目资源文件
// 📌 该文件不存在时自动创建
// ✏️ 可在此文件中添加业务逻辑
// ==============================================================================

package resources

import "embed"

//go:embed icon
var icon embed.FS

// Icon 获取图片数据
func Icon(fileName string) []byte {
	data, _ := icon.ReadFile("icon/" + fileName)
	return data
}
`
