package env

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestConfig(t *testing.T) string {
	t.Helper()
	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "C:\\cef",
			"version": "windows_amd64_109.1.18"
		},
		"app": {
			"version": "1.0.0"
		},
		"window": {
			"state": 0,
			"visible": true
		},
		"history_project": [
			"C:/a.egp"
		],
		"env": {
			"myapp": {
				"go_root": [
					"C:\\go"
				]
			},
			"myappcef": {
				"go_root": [
					"D:\\go"
				]
			}
		}
	}`
	if err := os.WriteFile(configFile, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	return configFile
}

func TestReadConfigFuzzy(t *testing.T) {
	configFile := writeTestConfig(t)
	var out bytes.Buffer
	if err := readConfig(configFile, "version", &out); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "app.version=1.0.0") {
		t.Fatalf("missing app version: %q", got)
	}
	if !strings.Contains(got, "chromium.version=windows_amd64_109.1.18") {
		t.Fatalf("missing chromium version: %q", got)
	}
}

func TestReadConfigExactMissing(t *testing.T) {
	configFile := writeTestConfig(t)
	var out bytes.Buffer
	if err := readConfig(configFile, "chromium.missing", &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "\n" {
		t.Fatalf("missing exact path should return empty value, got %q", out.String())
	}
}

func TestReadConfigFuzzyIndexedKey(t *testing.T) {
	configFile := writeTestConfig(t)
	var out bytes.Buffer
	if err := readConfig(configFile, "go_root[0]", &out); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, `env.myapp.go_root[0]=C:\go`) {
		t.Fatalf("missing myapp go root: %q", got)
	}
	if !strings.Contains(got, `env.myappcef.go_root[0]=D:\go`) {
		t.Fatalf("missing myappcef go root: %q", got)
	}
}

func TestWriteConfigExactTypes(t *testing.T) {
	configFile := writeTestConfig(t)
	if err := writeConfig(configFile, "window.state=2", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	if err := writeConfig(configFile, "window.visible=false", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	if err := writeConfig(configFile, "history_project[0]=C:/b.egp", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, root, "window.state", "2")
	assertPath(t, root, "window.visible", false)
	assertPath(t, root, "history_project[0]", "C:/b.egp")
}

func TestWriteConfigWrongTypeFails(t *testing.T) {
	configFile := writeTestConfig(t)
	if err := writeConfig(configFile, "window.state=abc", strings.NewReader(""), ioDiscard{}); err == nil {
		t.Fatal("expected integer write to fail")
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, root, "window.state", "0")
}

func TestWriteConfigFuzzySelect(t *testing.T) {
	configFile := writeTestConfig(t)
	var out bytes.Buffer
	if err := writeConfig(configFile, "version=2.0.0", strings.NewReader("2\n"), &out); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, root, "chromium.version", "2.0.0")
	assertPath(t, root, "app.version", "1.0.0")
	if !strings.Contains(out.String(), "Multiple keys matched:") {
		t.Fatalf("expected selection prompt, got %q", out.String())
	}
}

func TestWriteConfigFuzzyIndexedKeySelect(t *testing.T) {
	configFile := writeTestConfig(t)
	var out bytes.Buffer
	if err := writeConfig(configFile, `go_root[0]=C:\go1`, strings.NewReader("2\n"), &out); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, root, "env.myapp.go_root[0]", `C:\go`)
	assertPath(t, root, "env.myappcef.go_root[0]", `C:\go1`)
	if !strings.Contains(out.String(), "Multiple keys matched:") {
		t.Fatalf("expected selection prompt, got %q", out.String())
	}
}

func TestWriteConfigPreservesKeyOrder(t *testing.T) {
	configFile := writeTestConfig(t)
	if err := writeConfig(configFile, "window.state=3", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	assertBefore(t, content, `"chromium"`, `"app"`)
	assertBefore(t, content, `"app"`, `"window"`)
	assertBefore(t, content, `"window"`, `"history_project"`)
	assertBefore(t, content, `"state"`, `"visible"`)
}

func TestWriteCEFRuntimeSourceCreatesSelectionOnly(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configFile, []byte(`{"window": {}}`), 0644); err != nil {
		t.Fatal(err)
	}

	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	ensureCEFRuntimeSelection(root)
	if err = setPath(root, "cef_runtime.source", "sourceforge"); err != nil {
		t.Fatal(err)
	}

	assertPath(t, root, "cef_runtime.source", "sourceforge")
	if _, ok := getPath(root, "cef_runtime.version"); ok {
		t.Fatal("cef_runtime.version should not be created")
	}
	if _, ok := getPath(root, "cef_runtime.sources"); ok {
		t.Fatal("cef_runtime.sources should not be created")
	}
}

func TestEnsureCEFRuntimeSelectionReplacesNull(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configFile, []byte(`{"cef_runtime": null}`), 0644); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	ensureCEFRuntimeSelection(root)
	if err = setPath(root, "cef_runtime.source", "sourceforge"); err != nil {
		t.Fatal(err)
	}

	assertPath(t, root, "cef_runtime.source", "sourceforge")
}

func assertPath(t *testing.T, root any, path string, want any) {
	t.Helper()
	got, ok := getPath(root, path)
	if !ok {
		t.Fatalf("path not found: %s", path)
	}
	if formatValue(got) != formatValue(want) {
		t.Fatalf("%s got %q want %q", path, formatValue(got), formatValue(want))
	}
}

func assertBefore(t *testing.T, content, left, right string) {
	t.Helper()
	leftIndex := strings.Index(content, left)
	rightIndex := strings.Index(content, right)
	if leftIndex == -1 || rightIndex == -1 {
		t.Fatalf("missing %s or %s in %q", left, right, content)
	}
	if leftIndex > rightIndex {
		t.Fatalf("%s should appear before %s in %q", left, right, content)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestListVersionsCEF(t *testing.T) {
	// Create a temp chromium dir with some version directories
	chromiumDir := t.TempDir()
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_109.1.18"), 0755)
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_127.3.5"), 0755)
	os.MkdirAll(filepath.Join(chromiumDir, "linux_amd64_109.1.18"), 0755)
	// This should be ignored (not a valid version dir name)
	os.MkdirAll(filepath.Join(chromiumDir, "somefile.txt"), 0755)
	os.WriteFile(filepath.Join(chromiumDir, "somefile.txt"), []byte(""), 0644)

	// Write a config.json pointing to this dir with current version set
	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(chromiumDir, `\`, `\\`) + `",
			"version": "windows_amd64_127.3.5"
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	var out bytes.Buffer
	if err := listVersions(configFile, "cef", &out); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "* windows_amd64_127.3.5") {
		t.Fatalf("expected current version marked with *, got:\n%s", got)
	}
	if !strings.Contains(got, "  linux_amd64_109.1.18") {
		t.Fatalf("expected linux version without *, got:\n%s", got)
	}
	if !strings.Contains(got, "  windows_amd64_109.1.18") {
		t.Fatalf("expected old windows version without *, got:\n%s", got)
	}
	if strings.Contains(got, "somefile") {
		t.Fatalf("should not list non-version entries, got:\n%s", got)
	}
}

func TestListVersionsUnsupportedModule(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.json")
	os.WriteFile(configFile, []byte("{}"), 0644)
	var out bytes.Buffer
	err := listVersions(configFile, "unsupported", &out)
	if err == nil {
		t.Fatal("expected error for unsupported module")
	}
	if !strings.Contains(err.Error(), "unsupported module") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListVersionsEmptyModule(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.json")
	os.WriteFile(configFile, []byte("{}"), 0644)
	var out bytes.Buffer
	err := listVersions(configFile, "", &out)
	if err == nil {
		t.Fatal("expected error for empty module")
	}
}

func TestIsValidVersionDir(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"windows_amd64_127.3.5", true},
		{"linux_arm64_109.1.18", true},
		{"win_x64_1", true},
		{"no_underscore", false},
		{"a_b", false},
		{"os_arch_noDigit", false},
	}
	for _, tt := range tests {
		if got := isValidVersionDir(tt.name); got != tt.want {
			t.Errorf("isValidVersionDir(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestWriteConfigFuzzyChromiumVersionSingleMatch(t *testing.T) {
	chromiumDir := t.TempDir()
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_127.3.5"), 0755)
	os.MkdirAll(filepath.Join(chromiumDir, "linux_amd64_109.1.18"), 0755)

	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(chromiumDir, `\`, `\\`) + `",
			"version": "linux_amd64_109.1.18"
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	// "127" should auto-resolve to "windows_amd64_127.3.5" (only one match)
	if err := writeConfig(configFile, "version=127", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := getPath(root, "chromium.version")
	if !ok {
		t.Fatal("chromium.version not found")
	}
	if formatValue(got) != "windows_amd64_127.3.5" {
		t.Fatalf("expected windows_amd64_127.3.5, got %q", formatValue(got))
	}
}

func TestWriteConfigFuzzyChromiumVersionMultipleMatch(t *testing.T) {
	chromiumDir := t.TempDir()
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_127.3.5"), 0755)
	os.MkdirAll(filepath.Join(chromiumDir, "windows_386_127.3.5"), 0755)
	os.MkdirAll(filepath.Join(chromiumDir, "linux_amd64_109.1.18"), 0755)

	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(chromiumDir, `\`, `\\`) + `",
			"version": "linux_amd64_109.1.18"
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	var out bytes.Buffer
	// "127" matches two dirs; user selects the first one
	if err := writeConfig(configFile, "version=127", strings.NewReader("1\n"), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Multiple versions matched") {
		t.Fatalf("expected selection prompt, got %q", out.String())
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := getPath(root, "chromium.version")
	if !ok {
		t.Fatal("chromium.version not found")
	}
	if formatValue(got) != "windows_386_127.3.5" {
		t.Fatalf("expected windows_386_127.3.5 (first sorted), got %q", formatValue(got))
	}
}

func TestWriteConfigFuzzyChromiumVersionExactPassthrough(t *testing.T) {
	chromiumDir := t.TempDir()
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_127.3.5"), 0755)

	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(chromiumDir, `\`, `\\`) + `",
			"version": ""
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	// Exact match should pass through without fuzzy logic
	if err := writeConfig(configFile, "version=windows_amd64_127.3.5", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := getPath(root, "chromium.version")
	if !ok {
		t.Fatal("chromium.version not found")
	}
	if formatValue(got) != "windows_amd64_127.3.5" {
		t.Fatalf("expected exact passthrough, got %q", formatValue(got))
	}
}

func TestWriteConfigFuzzyChromiumVersionPrefixMatch(t *testing.T) {
	chromiumDir := t.TempDir()
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_127.3.5"), 0755)
	os.MkdirAll(filepath.Join(chromiumDir, "windows_386_127.3.5"), 0755)

	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(chromiumDir, `\`, `\\`) + `",
			"version": ""
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	// "windows_amd64_127" should match "windows_amd64_127.3.5" as prefix
	if err := writeConfig(configFile, "version=windows_amd64_127", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := getPath(root, "chromium.version")
	if !ok {
		t.Fatal("chromium.version not found")
	}
	if formatValue(got) != "windows_amd64_127.3.5" {
		t.Fatalf("expected windows_amd64_127.3.5, got %q", formatValue(got))
	}
}

func TestWriteConfigFuzzyChromiumVersionNoMatch(t *testing.T) {
	chromiumDir := t.TempDir()
	os.MkdirAll(filepath.Join(chromiumDir, "windows_amd64_127.3.5"), 0755)

	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(chromiumDir, `\`, `\\`) + `",
			"version": ""
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	// No match -> write is skipped entirely, config unchanged
	if err := writeConfig(configFile, "version=999", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := getPath(root, "chromium.version")
	if !ok {
		t.Fatal("chromium.version not found")
	}
	if formatValue(got) != "" {
		t.Fatalf("expected value unchanged (empty), got %q", formatValue(got))
	}
}

func TestWriteConfigFuzzyChromiumVersionNoInstalledDir(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.json")
	data := `{
		"chromium": {
			"dir": "` + strings.ReplaceAll(filepath.Join(t.TempDir(), "nonexistent"), `\`, `\\`) + `",
			"version": ""
		}
	}`
	os.WriteFile(configFile, []byte(data), 0644)

	// Directory doesn't exist -> no installed versions to match, set as-is
	if err := writeConfig(configFile, "version=127", strings.NewReader(""), ioDiscard{}); err != nil {
		t.Fatal(err)
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := getPath(root, "chromium.version")
	if !ok {
		t.Fatal("chromium.version not found")
	}
	if formatValue(got) != "127" {
		t.Fatalf("expected raw value 127, got %q", formatValue(got))
	}
}

func TestIsFuzzyVersionMatch(t *testing.T) {
	tests := []struct {
		dir   string
		value string
		want  bool
	}{
		{"windows_amd64_127.3.5", "127", true},
		{"windows_amd64_127.3.5", "127.3", true},
		{"windows_amd64_127.3.5", "windows_amd64_127", true},
		{"windows_amd64_127.3.5", "windows_amd64_127.3.5", true},
		{"windows_amd64_127.3.5", "109", false},
		{"windows_amd64_127.3.5", "linux", false},
		{"windows_386_127.3.5", "127", true},
		{"windows_386_127.3.5", "windows_386_127", true},
	}
	for _, tt := range tests {
		if got := isFuzzyVersionMatch(tt.dir, tt.value); got != tt.want {
			t.Errorf("isFuzzyVersionMatch(%q, %q) = %v, want %v", tt.dir, tt.value, got, tt.want)
		}
	}
}
