package dast

import (
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

type TTestStruct struct {
	Name string
	Age  int
}

type TTestEvent func(pam1 int32, param2 string) string

func TestCreateMethod(t *testing.T) {
	ts := FindType(testTestFilePath, "TTestEvent")
	if ts == nil {
		t.Fatal("ts is nil")
	}
	fnc, ok := ts.Type.(*ast.FuncType)
	if !ok {
		t.Fatal("ts.Type is not *ast.FuncType")
	}
	code, _ := CreateMethod(testTestFilePath, "TTestStruct", "NewTestStruct", fnc.Params, fnc.Results)
	t.Log(string(code))
	code = DeleteMethod(testTestFilePath, "TTestStruct", "NewTestStruct")
	t.Log(string(code))
}

func TestUpdateRecvMethodByTypeName(t *testing.T) {
	newCode, isUpdate, err := UpdateMethodRecv(testTestFilePath, "TTestStruct", "NewTestStruct")
	t.Log(isUpdate, err)
	t.Log(string(newCode))
}

func TestFindRecvMethod(t *testing.T) {
	var funcs []TFuncInfo
	FindRecvMethod(testTestFilePath, "TTestStruct", func(funcDecl *ast.FuncDecl) {
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
