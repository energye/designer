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

package dast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"strings"
)

// 函数信息
type TFuncInfo struct {
	Name    string
	Params  []TFieldInfo
	Results []TFieldInfo
}

// 字段信息
type TFieldInfo struct {
	Name string
	Type string
}

// 将 ast.Expr（语法树中的类型节点）转换为字符串
func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	// 使用 go/format 格式化语法树节点为可读字符串
	if err := format.Node(&buf, token.NewFileSet(), expr); err != nil {
		log.Printf("格式化类型失败: %v", err)
		return ""
	}
	return buf.String()
}

// ParseFields 解析 ast.FieldList
func ParseFields(fieldList *ast.FieldList) []TFieldInfo {
	var items []TFieldInfo
	if fieldList == nil { // 无参数/返回值时直接返回空切片
		return items
	}
	for _, field := range fieldList.List {
		// 解析类型字符串
		typStr := exprToString(field.Type)
		_, isPtr := field.Type.(*ast.StarExpr)
		if typStrIdx := strings.LastIndex(typStr, "."); typStrIdx != -1 {
			// 去除包声明
			typStr = typStr[typStrIdx+1:]
			if isPtr {
				typStr = "*" + typStr
			}
		}
		if len(field.Names) == 0 {
			// 无名称场景（如返回值 `(int, string)` 或参数 `func(int)`）
			items = append(items, TFieldInfo{
				Name: "",
				Type: typStr,
			})
		} else {
			// 多名称场景（如 `a, b int` 会生成两个 TParam）
			for _, name := range field.Names {
				items = append(items, TFieldInfo{
					Name: name.Name,
					Type: typStr,
				})
			}
		}
	}
	return items
}
