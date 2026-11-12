package webview

import "testing"

func TestClassifyURL(t *testing.T) {
	const origin = "http://127.0.0.1:22022"
	cases := []struct {
		name string
		url  string
		want urlAction
	}{
		{"monaco index", origin + "/index.html", urlActionInternal},
		{"monaco asset", origin + "/vs/editor/editor.main.js?v=1", urlActionInternal},
		{"uppercase scheme", "HTTP://127.0.0.1:22022/index.html", urlActionInternal},
		{"about blank", "about:blank", urlActionInternal},
		{"webview error page", "chrome-error://chromewebdata/", urlActionInternal},
		{"external https", "https://go.dev/doc", urlActionExternal},
		{"same host other scheme", "https://127.0.0.1:22022/index.html", urlActionExternal},
		{"external http", "http://example.com", urlActionExternal},
		{"same host other port", "http://127.0.0.1:8080/index.html", urlActionExternal},
		{"localhost", "http://localhost:22022/index.html", urlActionExternal},
		{"mailto", "mailto:dev@example.com", urlActionExternal},
		{"empty host", "http://", urlActionBlock},
		{"command", "command:editor.action.formatDocument", urlActionBlock},
		{"file", "file:///C:/Windows/win.ini", urlActionBlock},
		{"javascript", "javascript:alert(1)", urlActionBlock},
		{"data", "data:text/html,<b>x</b>", urlActionBlock},
		{"relative path", "index.html", urlActionBlock},
		{"empty", "", urlActionBlock},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := classifyURL(c.url, origin); got != c.want {
				t.Fatalf("classifyURL(%q) = %v, want %v", c.url, got, c.want)
			}
		})
	}
}

func TestDecideNavigation(t *testing.T) {
	const origin = "http://127.0.0.1:22022"
	cases := []struct {
		name             string
		url              string
		cancel, external bool
	}{
		{"internal", origin + "/index.html", false, false},
		{"external", "https://go.dev/doc", true, true},
		{"mailto", "mailto:dev@example.com", true, true},
		{"blocked", "command:editor.action.formatDocument", true, false},
		{"empty", "", true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cancel, external := decideNavigation(c.url, origin)
			if cancel != c.cancel || external != c.external {
				t.Fatalf("decideNavigation(%q) = (%v,%v), want (%v,%v)", c.url, cancel, external, c.cancel, c.external)
			}
		})
	}
}

func TestDecideNewWindow(t *testing.T) {
	const origin = "http://127.0.0.1:22022"
	if decideNewWindow(origin+"/index.html", origin) {
		t.Fatal("同源新窗口不应交给系统浏览器")
	}
	if !decideNewWindow("https://pkg.go.dev/net/http", origin) {
		t.Fatal("外部新窗口应交给系统浏览器")
	}
	if !decideNewWindow("mailto:dev@example.com", origin) {
		t.Fatal("mailto 新窗口应交给系统默认程序")
	}
	if decideNewWindow("file:///C:/x.txt", origin) {
		t.Fatal("file 新窗口不应交给系统")
	}
}

func TestMonacoAssetsOrigin(t *testing.T) {
	if got := monacoAssetsOrigin(0); got != "http://127.0.0.1:22022" {
		t.Fatalf("monacoAssetsOrigin(0) = %q", got)
	}
	if got := monacoAssetsOrigin(12345); got != "http://127.0.0.1:12345" {
		t.Fatalf("monacoAssetsOrigin(12345) = %q", got)
	}
}
