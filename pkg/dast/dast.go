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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/tool"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

var astMap *tool.HashMap[string, *ast.File]

func init() {
	astMap = tool.NewHashMap[string, *ast.File]()
	// 初始化内置类型
	initBuiltinTypesMap()
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

type TFuncTypeAlias struct {
	Mod     consts.Mod
	Imports *tool.HashMap[string, string]
	Funcs   *tool.HashMap[string, *ast.FuncType]
}

func PackageImportToHashMap(importSpec []*ast.ImportSpec) *tool.HashMap[string, string] {
	imports := tool.NewHashMap[string, string]()
	for _, spec := range importSpec {
		name := ""
		importPath := ""
		if spec.Name != nil {
			name = spec.Name.Name
		}
		importPath = strings.Trim(spec.Path.Value, "\"")
		if name == "" {
			name = filepath.Base(importPath)
		}
		imports.Add(name, importPath)
	}
	return imports
}

// GetAllFuncTypeAliasesByCode 在Go源文件中获取所有函数类型别名
func GetAllFuncTypeAliasesByCode(filename string, code []byte) *TFuncTypeAlias {
	node := MustFile(filename, code)
	if node == nil {
		return nil
	}
	funcs := &TFuncTypeAlias{Funcs: tool.NewHashMap[string, *ast.FuncType]()}
	funcs.Imports = PackageImportToHashMap(node.Imports)
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
			lowerType := strings.ToLower(typeSpec.Name.Name)
			funcs.Funcs.Add(lowerType, funcType)
		}
	}
	return funcs
}

// GetAllFuncTypeAliases 在Go源文件中获取所有函数类型别名
func GetAllFuncTypeAliases(filename string) *TFuncTypeAlias {
	node := MustFile(filename, nil)
	if node == nil {
		return nil
	}
	funcs := &TFuncTypeAlias{Funcs: tool.NewHashMap[string, *ast.FuncType]()}
	funcs.Imports = PackageImportToHashMap(node.Imports)
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
			lowerType := strings.ToLower(typeSpec.Name.Name)
			funcs.Funcs.Add(lowerType, funcType)
		}
	}
	return funcs
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

// FindRecvMethod 更新指定go代码文件的所有方法接收者
func FindRecvMethod(filename string, recvName string, callback func(funcDecl *ast.FuncDecl)) {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, filename, nil, parser.SkipObjectResolution)
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if recvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := recvType.X.(*ast.Ident); ok && ident.Name == recvName {
						callback(funcDecl)
					}
				}
			}
		}
	}

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

// CreateMethod 在指定的Go源文件中创建一个新的方法定义
//
//	filename - 要修改的源文件路径
//	typeName - 接收者类型名称
//	methodName - 新方法的名称
//	params - 方法参数列表的AST节点
//	returns - 方法返回值列表的AST节点
//
//
//	[]byte - 格式化后的源码字节流
//	error - 格式化过程中可能产生的错误
func CreateMethod(filename string, callback func(file *ast.File)) ([]byte, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	callback(node)
	var buf bytes.Buffer
	err = format.Node(&buf, fset, node)
	return buf.Bytes(), err
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

var (
	// 内置类型
	builtinTypesMap = tool.NewHashMap[string, struct{}]()
)

// 初始化内置类型
func initBuiltinTypesMap() {
	builtinTypesMap.Add("bool", struct{}{})
	builtinTypesMap.Add("int", struct{}{})
	builtinTypesMap.Add("int8", struct{}{})
	builtinTypesMap.Add("int16", struct{}{})
	builtinTypesMap.Add("int32", struct{}{})
	builtinTypesMap.Add("int64", struct{}{})
	builtinTypesMap.Add("uint", struct{}{})
	builtinTypesMap.Add("uint8", struct{}{})
	builtinTypesMap.Add("uint16", struct{}{})
	builtinTypesMap.Add("uint32", struct{}{})
	builtinTypesMap.Add("uint64", struct{}{})
	builtinTypesMap.Add("uintptr", struct{}{})
	builtinTypesMap.Add("byte", struct{}{})
	builtinTypesMap.Add("rune", struct{}{})
	builtinTypesMap.Add("float32", struct{}{})
	builtinTypesMap.Add("float64", struct{}{})
	builtinTypesMap.Add("complex64", struct{}{})
	builtinTypesMap.Add("complex128", struct{}{})
	builtinTypesMap.Add("string", struct{}{})
	builtinTypesMap.Add("interface", struct{}{})
	builtinTypesMap.Add("any", struct{}{})
	builtinTypesMap.Add("chan", struct{}{})
	builtinTypesMap.Add("map", struct{}{})
	builtinTypesMap.Add("func", struct{}{})
}

// IsBuiltinTypeName 简单的判断类型名是否为内置类型
func IsBuiltinTypeName(typeName string) bool {
	if x := strings.LastIndex(typeName, "*"); x != -1 {
		typeName = typeName[x+1:]
	}
	return builtinTypesMap.ContainsKey(typeName)
}

// FixFieldList 用于修复字段列表中的类型引用问题，确保非内置类型的字段正确地加上包前缀。
// 同时处理已存在或缺失的导入项，并更新字段类型表达式。
//
//   - allImports: 所有可用的导入映射，键为包名，值为包路径。
//   - currImports: 当前文件中已经存在的导入项。
//   - addImports: 需要新增的导入项集合。
//   - mod: 当前模块名称（通常作为包前缀使用）。
//   - fieldList: 原始的字段列表。
//
// 返回值：
//   - newFieldList: 修复后的字段列表，其中类型引用已被修正。
func FixFieldList(allImports, currImports, addImports *tool.HashMap[string, string], mod consts.Mod, fieldList *ast.FieldList) (newFieldList *ast.FieldList) {
	if fieldList == nil {
		return
	}
	var processType func(fieldType ast.Expr) ast.Expr
	processType = func(fieldType ast.Expr) ast.Expr {
		switch t := fieldType.(type) {
		case *ast.Ident:
			// 判断参数类型 t.Name 是否 Go 内置类型, 如果不是则使用当前模块的包
			if !IsBuiltinTypeName(t.Name) {
				if !currImports.ContainsKey(mod) {
					packageImport := allImports.Get(mod)
					addImports.Add("", packageImport)
				}
				return &ast.SelectorExpr{
					X:   ast.NewIdent(mod),
					Sel: ast.NewIdent(t.Name),
				}
			} else {
				return t
			}
		case *ast.StarExpr:
			return &ast.StarExpr{
				X:    processType(t.X),
				Star: t.Star,
			}
		case *ast.SelectorExpr:
			if identX, ok := t.X.(*ast.Ident); ok {
				if !currImports.ContainsKey(identX.Name) {
					// 不存在，准备添加，但需要判断一下，除了有别名以外，导入的包路径是否相同，如果包路径相同，将type改为原有的
					packageImport := allImports.Get(identX.Name)
					// 导入包路径是相同的, 修改参数类型导入名称
					if pkg, ok := currImports.ContainsValue(packageImport); ok {
						return &ast.SelectorExpr{
							X:   ast.NewIdent(pkg),
							Sel: ast.NewIdent(t.Sel.Name),
						}
					} else {
						addImports.Add("", packageImport)
					}
				}
			}
			return t
		}
		return fieldType
	}

	newFieldList = &ast.FieldList{
		Opening: fieldList.Opening,
		List:    nil,
		Closing: fieldList.Closing,
	}
	for _, field := range fieldList.List {
		newField := &ast.Field{
			Doc:     field.Doc,
			Names:   field.Names,
			Type:    nil,
			Tag:     field.Tag,
			Comment: field.Comment,
		}
		newField.Type = processType(field.Type)
		newFieldList.List = append(newFieldList.List, newField)
	}
	return
}

// CreateImport 创建一个导入项
func CreateImport(alias, pkgPath string) *ast.ImportSpec {
	if pkgPath == "" {
		return nil
	}
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(pkgPath), // 自动添加双引号，符合 Go 语法
		},
	}
	if alias != "" {
		importSpec.Name = ast.NewIdent(alias)
	}
	return importSpec
}

// IsFuncInfoMatch 比较两个函数信息是否匹配
// 主要比较函数的参数类型和返回值类型
// source: 源函数信息
// target: 目标函数信息
// 返回值: 如果两个函数的参数和返回值类型都相同则返回true，否则返回false
func IsFuncInfoMatch(source, target *TFuncInfo) bool {
	if source == nil || target == nil {
		return false
	}
	if len(source.Params) != len(target.Params) || len(source.Results) != len(target.Results) {
		return false
	}
	for i, param := range source.Params {
		if param.Type != target.Params[i].Type {
			return false
		}
	}
	for i, result := range source.Results {
		if result.Type != target.Results[i].Type {
			return false
		}
	}
	return true
}
