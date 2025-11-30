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

package dependmod

import (
	"github.com/energye/designer/pkg/tool"
	"go/ast"
)

// 全局存放各依赖模块的函数类型别名
// key: 函数类型别名
// value: 函数类型别名对应的AST节点
var (
	gLCLFuncTypeAliases       = tool.NewHashMap[string, *ast.FuncType]()
	gCEFFuncTypeAliases       = tool.NewHashMap[string, *ast.FuncType]()
	gWVWindowsFuncTypeAliases = tool.NewHashMap[string, *ast.FuncType]()
	gWVDarwinFuncTypeAliases  = tool.NewHashMap[string, *ast.FuncType]()
	gWVLinuxFuncTypeAliases   = tool.NewHashMap[string, *ast.FuncType]()
	// 全部
	gAllFuncTypeAliases = tool.NewHashMap[string, *ast.FuncType]()
)

// GetLCLFuncTypeAlias 根据名称获取LCL函数类型别名的AST节点
//
//	name - 要查找的函数类型别名名称
func GetLCLFuncTypeAlias(name string) *ast.FuncType {
	if gLCLFuncTypeAliases == nil {
		return nil
	}
	return gLCLFuncTypeAliases.Get(name)
}
