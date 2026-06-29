package cef

import (
	"os"
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

func TestRuntimeVersionedLibName(t *testing.T) {
	tests := map[string]string{
		"libenergy-amd64-gtk3.so": "libenergy-amd64-gtk3-109.so",
		"libenergy-amd64.dll":     "libenergy-amd64-109.dll",
		"libenergy-arm64.dylib":   "libenergy-arm64-109.dylib",
	}
	for input, want := range tests {
		if got := runtimeVersionedLibName(input, "109"); got != want {
			t.Fatalf("runtimeVersionedLibName(%q) got %q want %q", input, got, want)
		}
	}
}

func TestRuntimeArchiveFileNameSourceForgeDownload(t *testing.T) {
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-linux-amd64-gtk3-109.zip/download"
	got := runtimeArchiveFileName(rawURL, "109", 0)
	want := "0_libenergy-linux-amd64-gtk3-109.zip"
	if got != want {
		t.Fatalf("runtimeArchiveFileName got %q want %q", got, want)
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
	oldRuntime := config.Config.CEFRuntime
	oldDesignerRuntime := config.DesignerConfig.CEFRuntime
	oldURL := os.Getenv("ENERGY_CEF_RUNTIME_URL")
	oldSource := os.Getenv("ENERGY_CEF_RUNTIME_SOURCE")
	defer func() {
		config.Config.CEFRuntime = oldRuntime
		config.DesignerConfig.CEFRuntime = oldDesignerRuntime
		_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", oldURL)
		_ = os.Setenv("ENERGY_CEF_RUNTIME_SOURCE", oldSource)
	}()

	config.Config.CEFRuntime = []byte(`{
		"version": "v1.1.1",
		"source": "sourceforge",
		"sources": {
			"github": ["https://github.example/{major}/{os}-{arch}-{ws}.zip"],
			"sourceforge": ["https://sf.example/{version}/{os}-{arch}-{ws}.zip"]
		}
	}`)
	config.DesignerConfig.CEFRuntime = nil
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
