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
	"github.com/energye/designer/pkg/tool"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
)

var astMap *tool.HashMap[*ast.File]

func init() {
	astMap = tool.NewHashMap[*ast.File]()
}

func mustFile(filename string, src any) *ast.File {
	_, file := filepath.Split(filename)
	if astFile := astMap.Get(file); astFile != nil {
		return astFile
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil
	}
	astMap.Add(file, node)
	return node
}

// FindFunction 在Go源文件中查找函数声明
func FindFunction(filename string, functionName string) *ast.FuncDecl {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil
	}
	// 遍历文件中的所有声明
	for _, decl := range node.Decls {
		// 检查是否为函数声明
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 检查函数名是否匹配
			if funcDecl.Name.Name == functionName {
				return funcDecl
			}
		}
	}
	return nil
}

// FindConst 查找常量声明
func FindConst(filename string, constName string) *ast.ValueSpec {
	node := mustFile(filename, nil)
	if node == nil {
		return nil
	}
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if name.Name == constName {
							return valueSpec
						}
					}
				}
			}
		}
	}
	return nil
}

// FindType 在Go源文件中查找类型声明
func FindType(filename string, typeName string) *ast.TypeSpec {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil
	}
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == typeName {
					return typeSpec
				}
			}
		}
	}
	return nil
}

// DeleteMethod 从Go源文件中删除方法
func DeleteMethod(filename string, typeName string, methodName string) []byte {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	newDecls := []ast.Decl{}
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 检查是否是要删除的方法
			if funcDecl.Name.Name == methodName {
				// 检查是否有正确的接收者
				if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
					if recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == typeName {
							continue // 跳过此声明（删除它）
						}
					}
				}
			}
		}
		newDecls = append(newDecls, decl)
	}
	node.Decls = newDecls
	var buf bytes.Buffer
	format.Node(&buf, fset, node)
	return buf.Bytes()
}

// 创建方法
func CreateMethod(filename string, typeName string, methodName string, params []*ast.Field, returns []*ast.Field) []byte {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == typeName {
					// 创建新方法
					method := &ast.FuncDecl{
						Recv: &ast.FieldList{
							List: []*ast.Field{{
								Names: []*ast.Ident{ast.NewIdent("self")},
								Type:  &ast.StarExpr{X: ast.NewIdent(typeName)},
							}},
						},
						Name: ast.NewIdent(methodName),
						Type: &ast.FuncType{
							Params:  &ast.FieldList{List: params},
							Results: &ast.FieldList{List: returns},
						},
						Body: &ast.BlockStmt{},
					}
					node.Decls = append(node.Decls, method)
				}
			}
		}
	}
	var buf bytes.Buffer
	format.Node(&buf, fset, node)
	return buf.Bytes()
}

// 获取常量值
func GetConstValue(filename string, name string) any {
	value := FindConst(filename, name)
	if value != nil && len(value.Names) > 0 {
		ident := value.Names[0]
		return ident.Obj.Data
	}
	return nil
}
