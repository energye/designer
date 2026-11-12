package assetserve

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed assets
var assets embed.FS

func newTestServer() *httpServer {
	server := New()
	server.AssetsFSName = "assets"
	server.Assets = &assets
	return server
}

func TestServeEmbeddedAsset(t *testing.T) {
	server := newTestServer()

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets-test.md", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("embedded asset body is empty")
	}
	if !strings.Contains(rec.Body.String(), "assets") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestServeMissingAsset(t *testing.T) {
	server := newTestServer()

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/not-found.md", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestServeWithoutAssets(t *testing.T) {
	server := New()

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "resource directory is not configured") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestServeHeaderKeyRejectsMismatch(t *testing.T) {
	server := newTestServer()
	AssetsServerHeaderKeyName = "X-Test-Key"
	AssetsServerHeaderKeyValue = "secret"
	t.Cleanup(func() {
		AssetsServerHeaderKeyValue = ""
	})

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets-test.md", nil))
	if rec.Body.Len() != 0 {
		t.Fatalf("mismatched header should not write body, got %q", rec.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/assets-test.md", nil)
	req.Header.Set("X-Test-Key", "secret")
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Body.Len() == 0 {
		t.Fatalf("matched header failed: status=%d body=%q", rec.Code, rec.Body.String())
	}
}
