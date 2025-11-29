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
	"crypto/md5"
	"encoding/hex"
	"github.com/energye/designer/pkg/tool"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

var astMap *tool.HashMap[string, *ast.File]

func init() {
	astMap = tool.NewHashMap[string, *ast.File]()
}

func MustFile(filename string, src any) *ast.File {
	hash := md5.Sum([]byte(filename))
	key := hex.EncodeToString(hash[:])
	if astFile := astMap.Get(key); astFile != nil {
		return astFile
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	//node, err := parser.ParseFile(fset, filename, src, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	astMap.Add(key, node)
	return node
}

// FindFunction 在Go源文件中查找函数声明
func FindFunction(filename string, functionName string) *ast.FuncDecl {
	node := MustFile(filename, nil) // 使用缓存版本
	if node == nil {
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

// GetAllFuncTypeAliases 在Go源文件中获取所有函数类型别名
func GetAllFuncTypeAliases(filename string) *tool.HashMap[string, *ast.FuncType] {
	node := MustFile(filename, nil)
	if node == nil {
		return nil
	}
	result := tool.NewHashMap[string, *ast.FuncType]()
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			funcType, ok := typeSpec.Type.(*ast.FuncType)
			if !ok {
				continue
			}
			result.Add(typeSpec.Name.Name, funcType)
		}
	}
	return result
}

// FindConst 查找常量声明
func FindConst(filename string, constName string) *ast.ValueSpec {
	node := MustFile(filename, nil)
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

// UpdateMethodRecv 更新指定go代码文件的所有方法接收者
func UpdateMethodRecv(filename string, oldTypeName, newTypename string) ([]byte, bool, error) {
	isUpdate := false
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == oldTypeName {
						ident.Name = newTypename
						isUpdate = true
					}
				}
			}
		}
	}
	if isUpdate {
		var buf bytes.Buffer
		err := format.Node(&buf, fset, node)
		return buf.Bytes(), true, err
	}
	return nil, false, nil
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

// GetConstValue 获取常量值, 在指定 go 源码文件获取常量值
func GetConstValue(filename string, name string) any {
	value := FindConst(filename, name)
	if value != nil && len(value.Names) > 0 {
		ident := value.Names[0]
		return ident.Obj.Data
	}
	return nil
}
