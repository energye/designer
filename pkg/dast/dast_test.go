package dast

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	lclTypes "github.com/energye/lcl/types"
	"go/ast"
	"os"
	"path/filepath"
	"testing"
)

var (
	wd, _            = os.Getwd()
	testTestFilePath = filepath.Join(wd, "dast_bean_test.go")
)

type TCursor = int16

const (
	CrHigh = TCursor(0)

	CrDefault = TCursor(0)
	CrNone    = TCursor(-1)
	CrArrow   = TCursor(-2)
)

type TTest2 = int16

const (
	TestHigh TTest2 = iota
	TestDefault
	TestNone
	TestArrow
)

type ITestInterface interface {
	Method1(pam1 lcl.IObject, param2 string, param3 TTestStruct, param4 lclTypes.TRect) lcl.IPanel
	Method2(pam1 lcl.IObject, param2 string, param3 TTestStruct, param4 lclTypes.TRect) lcl.IPanel
}

type TTestStruct struct {
	Name string
	Age  int
}

func TestCreateMethod(t *testing.T) {
	code, _ := CreateMethod(testTestFilePath, func(file *ast.File) {
		newMethod := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{ast.NewIdent("m")},
					Type:  &ast.StarExpr{X: ast.NewIdent("NewIdent")},
				}},
			},
			Name: ast.NewIdent("NewMethod"),
			Type: &ast.FuncType{},
			Body: &ast.BlockStmt{},
		}
		file.Decls = append(file.Decls, newMethod)
		file.Imports = append(file.Imports, &ast.ImportSpec{})
	})
	t.Log(string(code))
	//ast.Ident
	//ast.SelectorExpr

	//code = DeleteMethod(testTestFilePath, "TTestStruct", "NewTestStruct")
	//t.Log(string(code))
}

func TestUpdateRecvMethodByTypeName(t *testing.T) {
	newCode, isUpdate, err := UpdateMethodRecv(testTestFilePath, "TTestStruct", "NewTestStruct")
	t.Log(isUpdate, err)
	t.Log(string(newCode))
}

func TestFindRecvMethod(t *testing.T) {
	var funcs []TFuncInfo
	FindRecvMethod(testTestFilePath, "TTestEvent", func(funcDecl *ast.FuncDecl) {
		name := funcDecl.Name.Name
		params := funcDecl.Type.Params
		results := funcDecl.Type.Results
		t.Log(name, params, results)
		funcInfo := TFuncInfo{
			Name:    name,
			Params:  ParseFields(params),
			Results: ParseFields(results),
		}
		funcs = append(funcs, funcInfo)
	})
	t.Log(funcs)
}

func TestFixFieldList(t *testing.T) {

	GLCLFuncTypeAliases := GetAllFuncTypeAliases(testTestFilePath)
	GLCLFuncTypeAliases.Mod = "lcl"
	GLCLFuncTypeAliases.Imports.Add(GLCLFuncTypeAliases.Mod, consts.DmLCL)

	funcType := GLCLFuncTypeAliases.Funcs.Get("TAlignPositionEvent")

	code, err := CreateMethod(testTestFilePath, func(file *ast.File) {
		currImports := PackageImportToHashMap(file.Imports) // 当前包
		allImports := tool.NewHashMap[string, string]()     // 所有包
		addImports := tool.NewHashMap[string, string]()     // 添加包
		GLCLFuncTypeAliases.Imports.Iterate(func(key string, value string) bool {
			allImports.Add(key, value)
			return false
		})
		currImports.Iterate(func(key string, value string) bool {
			allImports.Add(key, value)
			return false
		})
		// 可能添加的新导入包
		params := FixFieldList(allImports, currImports, addImports, GLCLFuncTypeAliases.Mod, funcType.Params)
		results := FixFieldList(allImports, currImports, addImports, GLCLFuncTypeAliases.Mod, funcType.Results)
		newMethod := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{ast.NewIdent("m")},
					Type:  &ast.StarExpr{X: ast.NewIdent("NewFormName")},
				}},
			},
			Name: ast.NewIdent("NewEventName"),
			Type: &ast.FuncType{
				Params:  params,
				Results: results,
			},
			Body: &ast.BlockStmt{},
		}
		file.Decls = append(file.Decls, newMethod)

		addImports.Iterate(func(alias string, pkgPath string) bool {
			file.Imports = append(file.Imports, CreateImport(alias, pkgPath))
			return false
		})
	})
	t.Log(string(code), err)
}
