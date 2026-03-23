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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/tool"
	"go/ast"
)

// 全局存放各依赖模块的函数类型别名
// key: 函数类型别名
// value: 函数类型别名对应的AST节点
var (
	GLCLFuncTypeAliases       *dast.TFuncTypeAlias
	GCEFFuncTypeAliases       *dast.TFuncTypeAlias
	GWVWindowsFuncTypeAliases *dast.TFuncTypeAlias
	GWVDarwinFuncTypeAliases  *dast.TFuncTypeAlias
	GWVLinuxFuncTypeAliases   *dast.TFuncTypeAlias
	// 全部
	gAllFuncTypeAliases = tool.NewHashMap[string, *ast.FuncType]()
)

func GetFuncTypeAliases(mod consts.Mod) *dast.TFuncTypeAlias {
	switch mod {
	case consts.ModLCL:
		return GLCLFuncTypeAliases
	case consts.ModCEF:
		return GCEFFuncTypeAliases
	case consts.ModWVWindows:
		return GWVWindowsFuncTypeAliases
	case consts.ModWVDarwin:
		return GWVDarwinFuncTypeAliases
	case consts.ModWVLinux:
		return GWVLinuxFuncTypeAliases
	}
	return nil
}

// InitDependencyModule 初始化模块类型信息
func InitDependencyModule(success func()) {
	go func() {
		// initModuleTypeInfoFormEmbed() 不再使用
		//go initModuleTypeInfoFormEmbed()

		// 根据 designer/resources/config.json 配置依赖模块下载模块
		downloadMod()
		// 从模块缓存 初始化模块类型信息
		initModuleTypeInfoFormModCache()

		// 完成回调
		if success != nil {
			success()
		}
	}()
}
