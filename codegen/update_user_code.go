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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/designer/dependmod"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
)

// 更新用户代码文件
// 此处功能为:
//   1. 事件绑定/删除/修改
//	 2. 自引用

// 执行更新自引用
// 被执行时说明当前窗体名称被修改. 窗体 name 修改过程 => file name (xxx.ui.go xxx.go xxx.ui) => go self code
// ast 用户文件
//
//	功能: 修改窗体绑定的自引用结构名. 例如 窗体结构 TForm1 > TForm2,  原: func (m *Form1) 新名: func (m *TForm2)
//	仅修改 xxx.go 用户文件, 其它文件需手动修改
func doUpdateSelf(uiGenData designer.TUIGenerationData) {
	formTab := uiGenData.Component.FormTab
	goUserFilePath := filepath.Join(bean.CodePath(), formTab.GOUserFile())
	newFormName := "T" + formTab.FormRoot.Name()
	oldFormName := "T" + uiGenData.Component.FormTab.OldFormName
	logs.Debug("NewFormName:", newFormName, "OldFormName:", oldFormName, "FilePath:", goUserFilePath)
	// 查找 OldFormName 的类型函数名
	// 修改 OldFormName 类型函数为 newFormName 的函数
	// 写入文件
	newCode, isUpdate, err := dast.UpdateMethodRecv(goUserFilePath, oldFormName, newFormName)
	if err != nil {
		logs.Error("执行更新自引用-UpdateMethodRecv:", err.Error())
		return
	}
	if !isUpdate {
		return
	}
	st, err := os.Stat(goUserFilePath)
	if err != nil {
		logs.Error("执行更新自引用-Stat:", err.Error())
		return
	}
	if err = os.WriteFile(goUserFilePath, newCode, st.Mode()); err != nil {
		logs.Error("执行更新自引用-WriteFile:", err.Error())
		return
	}
}

// 执行更新绑定事件
// 被执行时说明组件事件被修改. 组件属性-事件 => file code (xxx.ui.go xxx.go xxx.ui) => go event code
// ast 用户文件
//
//	xxx.ui.go => 绑定事件
//	xxx.go => 事件定义
//	xxx.ui => 事件记录
func doUpdateEvent(uiGenData designer.TUIGenerationData) {
	if uiGenData.NodeData == nil {
		logs.Error("UserCode-doUpdateEvent 执行更新绑定事件-事件节点数据为 nil")
		return
	}
	formTab := uiGenData.Component.FormTab
	goUserFilePath := filepath.Join(bean.CodePath(), formTab.GOUserFile())
	formName := "T" + formTab.FormRoot.Name() // T + [FormName]
	nodeData := uiGenData.NodeData.EditNodeData
	bindEventFuncName := uiGenData.NodeData.EditStringValue()

	switch nodeData.EventState {
	case consts.EsUpdate:
		// A > B 忽略
		logs.Debug("UserCode-doUpdateEvent 绑定事件忽略 Update")
	case consts.EsDelete:
		// 忽略
		logs.Debug("UserCode-doUpdateEvent 绑定事件忽略 Delete")
	case consts.EsAdd:
		logs.Debug("UserCode-doUpdateEvent 绑定事件添加 Add")
		// > A
		// 模块类型别名集合 TODO 多模块时需要动态控制 lcl, cef, wv
		funcTypeAliases := dependmod.GetFuncTypeAliases(uiGenData.Component.GetMod())
		// 生成绑定事件
		lowerType := strings.ToLower(nodeData.Metadata.Type)
		funcType := funcTypeAliases.Funcs.Get(lowerType)
		if funcType == nil {
			logs.Error("UserCode-doUpdateEvent 获取函数类型别名返回 nil", lowerType)
		} else {
			s, err := os.Stat(goUserFilePath)
			if err != nil {
				logs.Error("UserCode-doUpdateEvent 获取代码文件信息失败.", err.Error())
				return
			}
			// 创建新方法
			newCode, err := dast.CreateMethod(goUserFilePath, func(file *ast.File) {
				currImports := dast.PackageImportToHashMap(file.Imports) // 当前包
				allImports := tool.NewHashMap[string, string]()          // 所有包
				addImports := tool.NewHashMap[string, string]()          // 添加包
				// 合并包
				funcTypeAliases.Imports.Iterate(func(key string, value string) bool {
					allImports.Add(key, value)
					return false
				})
				// 合并包
				currImports.Iterate(func(key string, value string) bool {
					allImports.Add(key, value)
					return false
				})
				// 修改参数和返回值
				params := dast.FixFieldList(allImports, currImports, addImports, funcTypeAliases.Mod, funcType.Params)
				results := dast.FixFieldList(allImports, currImports, addImports, funcTypeAliases.Mod, funcType.Results)
				newMethod := &ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{{
							Names: []*ast.Ident{ast.NewIdent("m")},
							Type:  &ast.StarExpr{X: ast.NewIdent(formName)},
						}},
					},
					Name: ast.NewIdent(bindEventFuncName),
					Type: &ast.FuncType{
						Params:  params,
						Results: results,
					},
					Body: &ast.BlockStmt{},
				}
				file.Decls = append(file.Decls, newMethod)
				// 添加导入包
				addImports.Iterate(func(alias string, pkgPath string) bool {
					file.Imports = append(file.Imports, dast.CreateImport(alias, pkgPath))
					return false
				})
			})
			if err != nil {
				logs.Error("UserCode-doUpdateEvent 创建方法失败.", bindEventFuncName, err.Error())
				return
			}
			err = os.WriteFile(goUserFilePath, newCode, s.Mode())
			if err != nil {
				logs.Error("UserCode-doUpdateEvent 创建方法失败.", bindEventFuncName, goUserFilePath, err.Error())
				return
			}
			// 成功, 立即更新一次文件改变监听
			detectFileChange()
		}
	}

}
