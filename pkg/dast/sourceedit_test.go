package dast

import (
	"strings"
	"testing"
)

func TestUpdateStructName(t *testing.T) {
	src := []byte(`package app

type A struct {
	Field int
}
`)

	got, err := UpdateStructName(src, "A", "B")
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, got, "type B struct")
	assertNotContains(t, got, "type A struct")
}

func TestUpdateFormStructName(t *testing.T) {
	src := []byte(`package app

type TForm4 struct {
	TForm4UI
	UserField string
}
`)

	got, err := UpdateFormStructName(src, "TForm4", "TMainForm")
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, got, "type TMainForm struct")
	assertContains(t, got, "TMainFormUI")
	assertContains(t, got, "UserField string")
	assertNotContains(t, got, "type TForm4 struct")
	assertNotContains(t, got, "TForm4UI")
}

func TestUpdateStructFieldType(t *testing.T) {
	src := []byte(`package app

type A struct {
	Field int
}
`)

	got, err := UpdateStructFieldType(src, "A", "int", "string")
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, got, "Field string")
	assertNotContains(t, got, "Field int")
}

func TestUpdateVarName(t *testing.T) {
	src := []byte(`package app

var a int
`)

	got, err := UpdateVarName(src, "a", "b")
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, got, "var b int")
	assertNotContains(t, got, "var a int")
}

func TestUpdateVarType(t *testing.T) {
	src := []byte(`package app

var a int
`)

	got, err := UpdateVarType(src, "int", "string")
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, got, "var a string")
	assertNotContains(t, got, "var a int")
}

func assertContains(t *testing.T, got []byte, want string) {
	t.Helper()
	if !strings.Contains(string(got), want) {
		t.Fatalf("expected source to contain %q:\n%s", want, got)
	}
}

func assertNotContains(t *testing.T, got []byte, unwanted string) {
	t.Helper()
	if strings.Contains(string(got), unwanted) {
		t.Fatalf("expected source not to contain %q:\n%s", unwanted, got)
	}
}
