package dast

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"
)

var (
	wd, _            = os.Getwd()
	testTestFilePath = filepath.Join(wd, "dast_test.go")
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
	code := CreateMethod(testTestFilePath, "TTestStruct", "NewTestStruct", fnc.Params.List, fnc.Results.List)
	t.Log(string(code))
	code = DeleteMethod(testTestFilePath, "TTestStruct", "NewTestStruct")
	t.Log(string(code))
}
