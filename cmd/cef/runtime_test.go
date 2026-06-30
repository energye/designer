package cef

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/energye/designer/pkg/config"
)

func TestMajorVersion(t *testing.T) {
	tests := map[string]string{
		"109.1.18": "109",
		"127":      "127",
		"":         "",
	}
	for input, want := range tests {
		if got := MajorVersion(input); got != want {
			t.Fatalf("MajorVersion(%q) got %q want %q", input, got, want)
		}
	}
}

func TestRuntimeArchiveFileNameSourceForgeDownload(t *testing.T) {
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-linux-amd64-gtk3-109.zip/download"
	got := runtimeArchiveFileName(rawURL, "109", "v3.0.1")
	want := "v3.0.1_libenergy-linux-amd64-gtk3-109.zip"
	if got != want {
		t.Fatalf("runtimeArchiveFileName got %q want %q", got, want)
	}
}

func TestRuntimeArchivePathUsesChromiumRoot(t *testing.T) {
	chromiumDir := t.TempDir()
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-linux-amd64-gtk3-109.zip/download"
	got := runtimeArchivePath(chromiumDir, rawURL, "109", "v3.0.1")
	want := filepath.Join(chromiumDir, "v3.0.1_libenergy-linux-amd64-gtk3-109.zip")
	if got != want {
		t.Fatalf("runtimeArchivePath got %q want %q", got, want)
	}
}

func TestRuntimeArchivePrefixFallback(t *testing.T) {
	if got := runtimeArchivePrefix("", "109"); got != "109" {
		t.Fatalf("runtimeArchivePrefix got %q want 109", got)
	}
}

func TestRuntimeTempLibPathKeepsLibraryExtension(t *testing.T) {
	tests := map[string]string{
		"/runtime/libenergy-amd64.dll":     "/runtime/libenergy-amd64.download.dll",
		"/runtime/libenergy-amd64-gtk3.so": "/runtime/libenergy-amd64-gtk3.download.so",
		"/runtime/libenergy-arm64.dylib":   "/runtime/libenergy-arm64.download.dylib",
	}
	for input, want := range tests {
		if got := runtimeTempLibPath(input); got != want {
			t.Fatalf("runtimeTempLibPath(%q) got %q want %q", input, got, want)
		}
	}
}

func TestNormalizeSourceForgeDownloadURL(t *testing.T) {
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-linux-amd64-gtk3-109.zip"
	got := normalizeDownloadURL(rawURL)
	want := rawURL + "/download"
	if got != want {
		t.Fatalf("normalizeDownloadURL got %q want %q", got, want)
	}
}

func TestRuntimeDownloadURLs(t *testing.T) {
	oldSourceConfig := config.Config.Chromium.Source
	oldDesignerRuntime := config.DesignerConfig.CEFRuntime
	oldURL := os.Getenv("ENERGY_CEF_RUNTIME_URL")
	oldSource := os.Getenv("ENERGY_CEF_RUNTIME_SOURCE")
	defer func() {
		config.Config.Chromium.Source = oldSourceConfig
		config.DesignerConfig.CEFRuntime = oldDesignerRuntime
		_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", oldURL)
		_ = os.Setenv("ENERGY_CEF_RUNTIME_SOURCE", oldSource)
	}()

	config.Config.Chromium.Source = "sourceforge"
	config.DesignerConfig.CEFRuntime = []byte(`{
		"version": "v1.1.1",
		"source": "github",
		"sources": {
			"github": ["https://github.example/{major}/{os}-{arch}-{ws}.zip"],
			"sourceforge": ["https://sf.example/{version}/{os}-{arch}-{ws}.zip"]
		}
	}`)
	_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", "https://env.example/{major}.zip")
	_ = os.Unsetenv("ENERGY_CEF_RUNTIME_SOURCE")

	urls := RuntimeDownloadURLs("109.1.18", "linux", "amd64", "gtk3")
	got := strings.Join(urls, "\n")
	want := strings.Join([]string{
		"https://env.example/109.zip",
		"https://sf.example/v1.1.1/linux-amd64-gtk3.zip",
		"https://github.example/109/linux-amd64-gtk3.zip",
	}, "\n")
	if got != want {
		t.Fatalf("RuntimeDownloadURLs got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRuntimeDownloadURLsIgnoreUserRuntimeTemplates(t *testing.T) {
	oldSourceConfig := config.Config.Chromium.Source
	oldDesignerRuntime := config.DesignerConfig.CEFRuntime
	oldURL := os.Getenv("ENERGY_CEF_RUNTIME_URL")
	oldSource := os.Getenv("ENERGY_CEF_RUNTIME_SOURCE")
	defer func() {
		config.Config.Chromium.Source = oldSourceConfig
		config.DesignerConfig.CEFRuntime = oldDesignerRuntime
		_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", oldURL)
		_ = os.Setenv("ENERGY_CEF_RUNTIME_SOURCE", oldSource)
	}()

	config.Config.Chromium.Source = "sourceforge"
	config.DesignerConfig.CEFRuntime = []byte(`{
		"version": "v1.1.1",
		"source": "github",
		"sources": {
			"github": ["https://github.example/{version}/{major}.zip"],
			"sourceforge": ["https://sf.example/{version}/{major}.zip"]
		}
	}`)
	_ = os.Unsetenv("ENERGY_CEF_RUNTIME_URL")
	_ = os.Unsetenv("ENERGY_CEF_RUNTIME_SOURCE")

	urls := RuntimeDownloadURLs("109.1.18", "linux", "amd64", "gtk3")
	got := strings.Join(urls, "\n")
	want := strings.Join([]string{
		"https://sf.example/v1.1.1/109.zip",
		"https://github.example/v1.1.1/109.zip",
	}, "\n")
	if got != want {
		t.Fatalf("RuntimeDownloadURLs got:\n%s\nwant:\n%s", got, want)
	}
}

func TestExpandRuntimeURLUsesReleaseVersionAndOptionalWS(t *testing.T) {
	tmpl := "https://sourceforge.example/{version}/libenergy-{os}-{arch}-{ws}-{major}.zip"
	got := expandRuntimeURL(tmpl, "109.1.18", "109", "windows", "amd64", "", "v1.1.1")
	want := "https://sourceforge.example/v1.1.1/libenergy-windows-amd64-109.zip"
	if got != want {
		t.Fatalf("expandRuntimeURL without ws got %q want %q", got, want)
	}

	got = expandRuntimeURL(tmpl, "109.1.18", "109", "linux", "amd64", "gtk3", "v1.1.1")
	want = "https://sourceforge.example/v1.1.1/libenergy-linux-amd64-gtk3-109.zip"
	if got != want {
		t.Fatalf("expandRuntimeURL with ws got %q want %q", got, want)
	}
}
