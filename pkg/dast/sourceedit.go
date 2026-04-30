package dast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// UpdateFormStructName updates a designer form user struct and its directly
// related UI embedding and package variable.
//
// Example:
//
//	type TForm4 struct {
//		TForm4UI
//	}
//
//	var Form4 TForm4
//
// UpdateFormStructName(src, "TForm4", "TMainForm") rewrites it to:
//
//	type TMainForm struct {
//		TMainFormUI
//	}
//
//	var MainForm TMainForm
func UpdateFormStructName(srcCode []byte, oldStructName, newStructName string) ([]byte, error) {
	fset, file, err := parseGoFile(srcCode)
	if err != nil {
		return nil, err
	}

	oldUIName := oldStructName + "UI"
	newUIName := newStructName + "UI"

	structChanged := false
	uiEmbedChanged := false

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch genDecl.Tok {
		case token.TYPE:
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != oldStructName {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				typeSpec.Name.Name = newStructName
				structChanged = true

				for _, field := range structType.Fields.List {
					if len(field.Names) != 0 {
						continue
					}
					ident, ok := field.Type.(*ast.Ident)
					if ok && ident.Name == oldUIName {
						ident.Name = newUIName
						uiEmbedChanged = true
					}
				}
			}
		}
	}

	if !structChanged {
		return nil, fmt.Errorf("struct %q not found", oldStructName)
	}
	if !uiEmbedChanged {
		return nil, fmt.Errorf("embedded UI type %q in struct %q not found", oldUIName, oldStructName)
	}
	return formatGoFile(fset, file)
}

func UpdateStructName(srcCode []byte, oldStructName, newStructName string) ([]byte, error) {
	fset, file, err := parseGoFile(srcCode)
	if err != nil {
		return nil, err
	}

	changed := false
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != oldStructName {
				continue
			}
			if _, ok := typeSpec.Type.(*ast.StructType); !ok {
				continue
			}
			typeSpec.Name.Name = newStructName
			changed = true
		}
	}

	if !changed {
		return nil, fmt.Errorf("struct %q not found", oldStructName)
	}
	return formatGoFile(fset, file)
}

func UpdateStructFieldType(srcCode []byte, structName, oldFieldType, newFieldType string) ([]byte, error) {
	fset, file, err := parseGoFile(srcCode)
	if err != nil {
		return nil, err
	}

	changed := false
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				fieldType, err := formatExpr(fset, field.Type)
				if err != nil {
					return nil, err
				}
				if fieldType != oldFieldType {
					continue
				}
				newType, err := parser.ParseExpr(newFieldType)
				if err != nil {
					return nil, fmt.Errorf("parse new field type %q: %w", newFieldType, err)
				}
				field.Type = newType
				changed = true
			}
		}
	}

	if !changed {
		return nil, fmt.Errorf("field type %q in struct %q not found", oldFieldType, structName)
	}
	return formatGoFile(fset, file)
}

func UpdateVarName(srcCode []byte, oldName, newName string) ([]byte, error) {
	fset, file, err := parseGoFile(srcCode)
	if err != nil {
		return nil, err
	}

	changed := false
	ast.Inspect(file, func(node ast.Node) bool {
		valueSpec, ok := node.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for _, name := range valueSpec.Names {
			if name.Name == oldName {
				name.Name = newName
				changed = true
			}
		}
		return false
	})

	if !changed {
		return nil, fmt.Errorf("var %q not found", oldName)
	}
	return formatGoFile(fset, file)
}

func UpdateVarType(srcCode []byte, oldType, newType string) ([]byte, error) {
	fset, file, err := parseGoFile(srcCode)
	if err != nil {
		return nil, err
	}
	newTypeExpr, err := parser.ParseExpr(newType)
	if err != nil {
		return nil, fmt.Errorf("parse new var type %q: %w", newType, err)
	}

	changed := false
	ast.Inspect(file, func(node ast.Node) bool {
		valueSpec, ok := node.(*ast.ValueSpec)
		if !ok || valueSpec.Type == nil {
			return true
		}

		varType, err := formatExpr(fset, valueSpec.Type)
		if err != nil || varType != oldType {
			return true
		}

		valueSpec.Type = newTypeExpr
		changed = true
		return false
	})

	if !changed {
		return nil, fmt.Errorf("var type %q not found", oldType)
	}
	return formatGoFile(fset, file)
}

func parseGoFile(srcCode []byte) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", srcCode, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return fset, file, nil
}

func formatGoFile(fset *token.FileSet, file *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatExpr(fset *token.FileSet, expr ast.Expr) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, expr); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func formVarName(structName string) string {
	return strings.TrimPrefix(structName, "T")
}
