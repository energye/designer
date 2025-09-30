package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// TypeInfo 存储类型定义信息
type TypeInfo struct {
	Name       string // 类型名（如TAlign）
	Underlying string // 底层类型（如int32）
}

// ConstInfo 存储常量信息
type ConstInfo struct {
	Name  string // 常量名（如AlTop）
	Type  string // 常量类型（如TAlign，可能为空）
	Value string // 常量值（如"0"、"iota"）
}

func main() {
	// 接收命令行参数：包路径（如./mypackage）
	pkgPath := "C:\\app\\workspace\\gen\\gout\\lcl\\go\\types"
	if len(os.Args) > 1 {
		pkgPath = os.Args[1]
	}

	// 扫描包并获取类型和常量信息
	types, consts, err := scanPackage(pkgPath)
	if err != nil {
		fmt.Printf("扫描失败: %v\n", err)
		os.Exit(1)
	}

	// 打印结果
	fmt.Println("===== 类型定义 =====")
	for _, t := range types {
		fmt.Printf("类型名: %-20s 底层类型: %s\n", t.Name, t.Underlying)
	}

	fmt.Println("\n===== 常量定义 =====")
	for _, c := range consts {
		if c.Type != "" {
			fmt.Printf("常量名: %-20s 类型: %-15s 值: %s\n", c.Name, c.Type, c.Value)
		} else {
			fmt.Printf("常量名: %-20s 值: %s\n", c.Name, c.Value)
		}
	}
}

// scanPackage 扫描指定包，返回类型和常量信息
func scanPackage(pkgPath string) ([]TypeInfo, []ConstInfo, error) {
	var types []TypeInfo
	var consts []ConstInfo

	// 1. 获取包的源文件列表
	pkg, err := build.ImportDir(pkgPath, build.ImportComment)
	if err != nil {
		return nil, nil, fmt.Errorf("获取包信息失败: %v", err)
	}

	// 2. 解析每个源文件的AST
	fset := token.NewFileSet() // 用于记录文件名和位置信息
	for _, file := range pkg.GoFiles {
		filePath := filepath.Join(pkgPath, file)
		astFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, fmt.Errorf("解析文件 %s 失败: %v", filePath, err)
		}

		// 3. 遍历AST提取类型和常量
		for _, decl := range astFile.Decls {
			// 处理类型声明（如 type TAlign int32）
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				switch genDecl.Tok {
				case token.TYPE:
					// 提取类型定义
					extractTypes(genDecl, &types)
				case token.CONST:
					// 提取常量定义
					extractConsts(genDecl, &consts)
				}
			}
		}
	}

	return types, consts, nil
}

// extractTypes 从类型声明中提取类型信息
func extractTypes(genDecl *ast.GenDecl, types *[]TypeInfo) {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		// 获取底层类型的字符串表示
		underlying := exprToString(typeSpec.Type)

		*types = append(*types, TypeInfo{
			Name:       typeSpec.Name.Name,
			Underlying: underlying,
		})
	}
}

// extractConsts 从常量声明中提取常量信息
func extractConsts(genDecl *ast.GenDecl, consts *[]ConstInfo) {
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// 常量类型（可能为空，如 const AlNone = 0）
		var typ string
		if valueSpec.Type != nil {
			typ = exprToString(valueSpec.Type)
		}

		// 处理批量声明的常量（如 const a, b = 1, 2）
		for i, name := range valueSpec.Names {
			var value string
			if i < len(valueSpec.Values) && valueSpec.Values[i] != nil {
				value = exprToString(valueSpec.Values[i])
			}
			*consts = append(*consts, ConstInfo{
				Name:  name.Name,
				Type:  typ,
				Value: value,
			})
		}
	}
}

// exprToString 将ast.Expr转换为简洁的字符串表示
func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.BasicLit:
		return e.Value
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + exprToString(e.Elt)
		}
		return "[" + exprToString(e.Len) + "]" + exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + exprToString(e.Key) + "]" + exprToString(e.Value)
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	case *ast.FuncType:
		return "func"
	case *ast.ChanType:
		return "chan " + exprToString(e.Value)
	case *ast.BinaryExpr:
		return exprToString(e.X) + " " + e.Op.String() + " " + exprToString(e.Y)
	case *ast.UnaryExpr:
		return e.Op.String() + exprToString(e.X)
	case *ast.CallExpr:
		return exprToString(e.Fun) + "()"
	case *ast.ParenExpr:
		return "(" + exprToString(e.X) + ")"
	case *ast.CompositeLit:
		return exprToString(e.Type) + "{}"
	case *ast.IndexExpr:
		return exprToString(e.X) + "[" + exprToString(e.Index) + "]"
	case *ast.SliceExpr:
		base := exprToString(e.X) + "["
		if e.Low != nil {
			base += exprToString(e.Low)
		}
		base += ":"
		if e.High != nil {
			base += exprToString(e.High)
		}
		if e.Slice3 {
			base += ":"
			if e.Max != nil {
				base += exprToString(e.Max)
			}
		}
		return base + "]"
	default:
		// 对于未知类型，使用简化表示
		return fmt.Sprintf("%T", expr)
	}
}
