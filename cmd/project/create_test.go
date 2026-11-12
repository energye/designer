package project

import (
	projecttpl "github.com/energye/designer/internal/templates/project"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveCreateDirUsesSelectedDirWhenNameMatches(t *testing.T) {
	root := t.TempDir()
	selectedDir := filepath.Join(root, "myapp")
	if err := os.MkdirAll(selectedDir, 0755); err != nil {
		t.Fatal(err)
	}

	got, err := resolveCreateDir("myapp", selectedDir)
	if err != nil {
		t.Fatal(err)
	}

	want, err := filepath.Abs(selectedDir)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("resolveCreateDir() = %q, want %q", got, want)
	}
}

func TestResolveCreateDirCreatesProjectDirWhenNameDiffers(t *testing.T) {
	root := t.TempDir()

	got, err := resolveCreateDir("myapp", root)
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(root, "myapp")
	if got != want {
		t.Fatalf("resolveCreateDir() = %q, want %q", got, want)
	}
	if info, err := os.Stat(want); err != nil {
		t.Fatal(err)
	} else if !info.IsDir() {
		t.Fatalf("%q is not a directory", want)
	}
}

func TestRunCEFCodeTemplate(t *testing.T) {
	code := buildTemplateData(projecttpl.RunCEFCodeTemplate, map[string]string{"Name": "cefcodegen"})
	for _, want := range []string{
		`"github.com/energye/energy/v3/cef"`,
		"cefApp := cef.Init()",
		"forms := make([]cef.IAppWindow, 0, len(app.Forms))",
		"cef.Run(forms...)",
		`"cefcodegen/app"`,
	} {
		if !strings.Contains(code, want) {
			t.Fatalf("CEF main template missing %q:\n%s", want, code)
		}
	}
}

func TestResolveCreateDirErrorsWhenProjectPathIsFile(t *testing.T) {
	root := t.TempDir()
	projectPath := filepath.Join(root, "myapp")
	if err := os.WriteFile(projectPath, []byte("not a dir"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := resolveCreateDir("myapp", root); err == nil {
		t.Fatal("resolveCreateDir() error = nil, want error")
	}
}
