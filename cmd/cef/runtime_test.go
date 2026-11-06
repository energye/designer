package cef

import (
	"net/url"
	"os"
	"path"
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
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-cef109-linux-amd64-gtk3.zip/download"
	got := runtimeArchiveFileName(rawURL, "109", "v3.0.1")
	want := "v3.0.1_libenergy-cef109-linux-amd64-gtk3.zip"
	if got != want {
		t.Fatalf("runtimeArchiveFileName got %q want %q", got, want)
	}
}

func TestRuntimeArchivePathUsesChromiumRoot(t *testing.T) {
	chromiumDir := t.TempDir()
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-cef109-linux-amd64-gtk3.zip/download"
	got := runtimeArchivePath(chromiumDir, rawURL, "109", "v3.0.1")
	want := filepath.Join(chromiumDir, "v3.0.1_libenergy-cef109-linux-amd64-gtk3.zip")
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
	rawURL := "https://sourceforge.net/projects/liblcl/files/v3.0.1/libenergy-cef109-linux-amd64-gtk3.zip"
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
			"github": ["https://github.example/{major}/libenergy-{module-}{os}-{arch}{-ws}.zip"],
			"sourceforge": ["https://sf.example/{version}/libenergy-{module-}{os}-{arch}{-ws}.zip"]
		}
	}`)
	_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", "https://env.example/{major}.zip")
	_ = os.Unsetenv("ENERGY_CEF_RUNTIME_SOURCE")

	urls := RuntimeDownloadURLs("109.1.18", "linux", "amd64", "gtk3")
	got := strings.Join(urls, "\n")
	want := strings.Join([]string{
		"https://env.example/109.zip",
		"https://sf.example/v1.1.1/libenergy-cef109-linux-amd64-gtk3.zip",
		"https://github.example/109/libenergy-cef109-linux-amd64-gtk3.zip",
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

func TestRuntimeFingerprintDetectsAllCustomChanges(t *testing.T) {
	stable := RuntimeOverride{WS: "gtk3", Source: "github"}
	base := NewRuntimeOptions("127.3.5", "linux", "amd64", stable)
	changed := []RuntimeOptions{
		NewRuntimeOptions("127.3.5", "linux", "amd64", RuntimeOverride{WS: "gtk3", Source: "github", Module: RuntimeModuleNone}),
		NewRuntimeOptions("127.3.5", "linux", "amd64", RuntimeOverride{WS: "gtk2", Source: "github"}),
		NewRuntimeOptions("127.3.5", "linux", "amd64", RuntimeOverride{WS: "gtk3", Source: "sourceforge"}),
		NewRuntimeOptions("127.3.5", "linux", "amd64", RuntimeOverride{WS: "gtk3", Source: "github", Release: "v3.0.2"}),
	}
	want := runtimeFingerprint(base)
	for i, options := range changed {
		if got := runtimeFingerprint(options); got == want {
			t.Fatalf("case %d fingerprint should differ: %q", i, got)
		}
	}
	if got := runtimeFingerprint(NewRuntimeOptions("127.3.5", "linux", "amd64", stable)); got != want {
		t.Fatalf("same options should share fingerprint: %q != %q", got, want)
	}
}

func TestRuntimeDefaultSourceIsSourceForge(t *testing.T) {
	if got := RuntimeDefaultSource(); got != "sourceforge" {
		t.Fatalf("RuntimeDefaultSource got %q want sourceforge", got)
	}
}

func TestRuntimeAssetNameFollowsReleaseNaming(t *testing.T) {
	tests := []struct {
		name    string
		options RuntimeOptions
		want    string
	}{
		{
			name:    "linux with gtk3",
			options: RuntimeOptions{Version: "109.1.18", Major: "109", OS: "linux", Arch: "amd64", WS: "gtk3", Module: "cef109", Release: "v3.0.1", Source: "github"},
			want:    "libenergy-cef109-linux-amd64-gtk3.zip",
		},
		{
			name:    "windows without ws",
			options: RuntimeOptions{Version: "127.3.5", Major: "127", OS: "windows", Arch: "386", Module: "cef127", Release: "v3.0.1", Source: "github"},
			want:    "libenergy-cef127-windows-386.zip",
		},
		{
			name:    "darwin arm64",
			options: RuntimeOptions{Version: "147.0.14", Major: "147", OS: "darwin", Arch: "arm64", Module: "cef147", Release: "v3.0.1", Source: "github"},
			want:    "libenergy-cef147-darwin-arm64.zip",
		},
		{
			name:    "no module",
			options: RuntimeOptions{Version: "150.0.21", Major: "150", OS: "linux", Arch: "amd64", WS: "gtk2", Release: "v3.0.1", Source: "github"},
			want:    "libenergy-linux-amd64-gtk2.zip",
		},
		{
			name:    "webview module",
			options: RuntimeOptions{Version: "150.0.21", Major: "150", OS: "linux", Arch: "arm64", WS: "gtk3", Module: "cef150-webview", Release: "v3.0.1", Source: "github"},
			want:    "libenergy-cef150-webview-linux-arm64-gtk3.zip",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.options.AssetName(); got != test.want {
				t.Fatalf("AssetName got %q want %q", got, test.want)
			}
		})
	}
}

func TestRuntimeDownloadURLsFromEmbeddedConfig(t *testing.T) {
	oldSource := config.Config.Chromium.Source
	oldRelease := config.Config.Chromium.RuntimeVersion
	oldURL := os.Getenv("ENERGY_CEF_RUNTIME_URL")
	oldSourceEnv := os.Getenv("ENERGY_CEF_RUNTIME_SOURCE")
	defer func() {
		config.Config.Chromium.Source = oldSource
		config.Config.Chromium.RuntimeVersion = oldRelease
		_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", oldURL)
		_ = os.Setenv("ENERGY_CEF_RUNTIME_SOURCE", oldSourceEnv)
	}()
	config.Config.Chromium.Source = "github"
	config.Config.Chromium.RuntimeVersion = ""
	_ = os.Unsetenv("ENERGY_CEF_RUNTIME_URL")
	_ = os.Unsetenv("ENERGY_CEF_RUNTIME_SOURCE")

	options := NewRuntimeOptions("127.3.5", "linux", "amd64", RuntimeOverride{WS: "gtk3", Source: "github"})
	urls := options.DownloadURLs()
	if len(urls) == 0 {
		t.Fatal("no download URLs resolved from embedded config")
	}
	want := "https://github.com/energye/lib/releases/download/v3.0.1/libenergy-cef127-linux-amd64-gtk3.zip"
	if urls[0] != want {
		t.Fatalf("first URL got %q want %q (all: %v)", urls[0], want, urls)
	}
}

func TestRuntimeOverrideControlsReleaseModuleAndWS(t *testing.T) {
	options := NewRuntimeOptions("109.1.18", "linux", "arm64", RuntimeOverride{
		Release: "v3.0.2",
		Module:  "cef{major}-webview",
		WS:      "gtk2",
		URL:     "https://example.com/{version}/libenergy-{module-}{os}-{arch}{-ws}.zip",
	})
	if options.AssetName() != "libenergy-cef109-webview-linux-arm64-gtk2.zip" {
		t.Fatalf("unexpected asset name: %s", options.AssetName())
	}
	want := "https://example.com/v3.0.2/libenergy-cef109-webview-linux-arm64-gtk2.zip"
	if urls := options.DownloadURLs(); len(urls) == 0 || urls[0] != want {
		t.Fatalf("custom URL got %v want first %q", urls, want)
	}
}

func TestCustomRuntimeURLWinsOverEnv(t *testing.T) {
	oldURL := os.Getenv("ENERGY_CEF_RUNTIME_URL")
	defer func() { _ = os.Setenv("ENERGY_CEF_RUNTIME_URL", oldURL) }()
	_ = os.Setenv("ENERGY_CEF_RUNTIME_URL", "https://env.example/{module}.zip")

	options := NewRuntimeOptions("127.3.5", "linux", "amd64", RuntimeOverride{
		URL: "https://custom.example/{module}.zip",
		WS:  "gtk3",
	})
	urls := options.DownloadURLs()
	if len(urls) < 2 {
		t.Fatalf("expected custom and env urls, got %v", urls)
	}
	if urls[0] != "https://custom.example/cef127.zip" || urls[1] != "https://env.example/cef127.zip" {
		t.Fatalf("unexpected url order: %v", urls[:2])
	}
}

func TestRuntimeModuleNoneProducesPlainAsset(t *testing.T) {
	options := NewRuntimeOptions("127.3.5", "windows", "amd64", RuntimeOverride{Module: RuntimeModuleNone})
	if options.Module != "" {
		t.Fatalf("module should be empty, got %q", options.Module)
	}
	if got := options.AssetName(); got != "libenergy-windows-amd64.zip" {
		t.Fatalf("AssetName got %q", got)
	}
	tmpl := "https://example.com/{version}/libenergy-{module-}{os}-{arch}{-ws}.zip"
	want := "https://example.com/v3.0.1/libenergy-windows-amd64.zip"
	if got := expandRuntimeURLFor(tmpl, options, ""); got != want {
		t.Fatalf("expandRuntimeURLFor got %q want %q", got, want)
	}
}

func TestExpandRuntimeURLForOptionalSegments(t *testing.T) {
	oldChromium := config.Config.Chromium
	defer func() { config.Config.Chromium = oldChromium }()
	config.Config.Chromium.RuntimeModule = ""
	config.Config.Chromium.RuntimeWS = ""

	tmpl := "https://example.com/{version}/libenergy-{module-}{os}-{arch}{-ws}.zip"
	tests := []struct {
		name    string
		options RuntimeOptions
		want    string
	}{
		{
			name:    "module and ws",
			options: NewRuntimeOptions("109.1.18", "linux", "amd64", RuntimeOverride{WS: "gtk3"}),
			want:    "https://example.com/v3.0.1/libenergy-cef109-linux-amd64-gtk3.zip",
		},
		{
			name:    "module without ws",
			options: NewRuntimeOptions("127.3.5", "windows", "386", RuntimeOverride{}),
			want:    "https://example.com/v3.0.1/libenergy-cef127-windows-386.zip",
		},
		{
			name:    "ws without module",
			options: NewRuntimeOptions("147.0.14", "linux", "arm64", RuntimeOverride{WS: "gtk2", Module: RuntimeModuleNone}),
			want:    "https://example.com/v3.0.1/libenergy-linux-arm64-gtk2.zip",
		},
		{
			name:    "neither",
			options: NewRuntimeOptions("150.0.21", "darwin", "arm64", RuntimeOverride{Module: RuntimeModuleNone}),
			want:    "https://example.com/v3.0.1/libenergy-darwin-arm64.zip",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := expandRuntimeURLFor(tmpl, test.options, "v3.0.1"); got != test.want {
				t.Fatalf("got %q want %q", got, test.want)
			}
		})
	}
}

// assetNameFromURL 从展开后的地址中取回 zip 名 (兼容 sourceforge 的 /download 后缀)
func assetNameFromURL(t *testing.T, raw string) string {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	p := u.Path
	if path.Base(p) == "download" {
		p = path.Dir(p)
	}
	return path.Base(p)
}

// TestEmbeddedRuntimeTemplatesMatchAssetName 保证内嵌下载模板与 AssetName() 的命名规则永不脱节
func TestEmbeddedRuntimeTemplatesMatchAssetName(t *testing.T) {
	cases := []RuntimeOptions{
		{Version: "109.1.18", Major: "109", OS: "linux", Arch: "amd64", WS: "gtk3", Module: "cef109", Release: "v3.0.1"},
		{Version: "127.3.5", Major: "127", OS: "windows", Arch: "386", Module: "cef127", Release: "v3.0.1"},
		{Version: "147.0.14", Major: "147", OS: "darwin", Arch: "arm64", Module: "cef147", Release: "v3.0.1"},
		{Version: "150.0.21", Major: "150", OS: "linux", Arch: "arm", WS: "gtk2", Release: "v3.0.1"},
		{Version: "150.0.21", Major: "150", OS: "linux", Arch: "amd64", WS: "gtk3", Module: "cef150-webview", Release: "v3.0.1"},
	}
	count := 0
	for _, cfg := range runtimeSourceConfigs() {
		for source, templates := range cfg.Sources {
			for _, tmpl := range templates {
				for _, options := range cases {
					raw := expandRuntimeURLFor(tmpl, options, "v3.0.1")
					got := assetNameFromURL(t, raw)
					if want := options.AssetName(); got != want {
						t.Fatalf("source %q template %q for %+v: got %q want %q", source, tmpl, options, got, want)
					}
					count++
				}
			}
		}
	}
	if count == 0 {
		t.Fatal("no embedded runtime templates found")
	}
}

func TestExpandRuntimeURLForWSSuffix(t *testing.T) {
	tmpl := "https://example.com/{version}/libenergy-{module-}{os}-{arch}{-ws}.zip"
	options := NewRuntimeOptions("147.0.14", "linux", "amd64", RuntimeOverride{WS: "gtk3"})
	if got := expandRuntimeURLFor(tmpl, options, "v3.0.1"); got != "https://example.com/v3.0.1/libenergy-cef147-linux-amd64-gtk3.zip" {
		t.Fatalf("linux URL got %q", got)
	}
	options = NewRuntimeOptions("147.0.14", "darwin", "arm64", RuntimeOverride{})
	if got := expandRuntimeURLFor(tmpl, options, "v3.0.1"); got != "https://example.com/v3.0.1/libenergy-cef147-darwin-arm64.zip" {
		t.Fatalf("darwin URL got %q", got)
	}
}

func TestRuntimeLibName(t *testing.T) {
	tests := []struct {
		osName, arch, ws, want string
	}{
		{"windows", "amd64", "", "libenergy-amd64.dll"},
		{"darwin", "arm64", "", "libenergy-arm64.dylib"},
		{"linux", "amd64", "gtk3", "libenergy-amd64-gtk3.so"},
		{"linux", "386", "gtk2", "libenergy-386-gtk2.so"},
	}
	for _, test := range tests {
		if got := runtimeLibName(test.osName, test.arch, test.ws); got != test.want {
			t.Fatalf("runtimeLibName(%q,%q,%q) got %q want %q", test.osName, test.arch, test.ws, got, test.want)
		}
	}
}
