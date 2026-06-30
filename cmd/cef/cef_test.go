package cef

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/energye/designer/pkg/config"
)

func TestInstalledVersionsIncludesDirectoryWithoutManifest(t *testing.T) {
	oldDir := config.Config.Chromium.Dir
	oldVersion := config.Config.Chromium.Version
	t.Cleanup(func() {
		config.Config.Chromium.Dir = oldDir
		config.Config.Chromium.Version = oldVersion
	})

	root := t.TempDir()
	oav := OSArchVersion("linux", "amd64", "127.0.0")
	config.Config.Chromium.Dir = root
	config.Config.Chromium.Version = oav

	versionDir := filepath.Join(root, oav)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(versionDir, "libcef.so"), []byte("cef"), 0644); err != nil {
		t.Fatal(err)
	}

	if !IsInstalled(oav) {
		t.Fatal("expected CEF directory without manifest to be recognized as installed")
	}
	got := InstalledVersions("linux", "amd64")
	if len(got) != 1 || got[0] != oav {
		t.Fatalf("InstalledVersions got %v want [%s]", got, oav)
	}
}

func TestIsInstalledRejectsDirectoryWithoutCEFMarker(t *testing.T) {
	oldDir := config.Config.Chromium.Dir
	t.Cleanup(func() {
		config.Config.Chromium.Dir = oldDir
	})

	root := t.TempDir()
	oav := OSArchVersion("linux", "amd64", "127.0.0")
	config.Config.Chromium.Dir = root
	if err := os.MkdirAll(filepath.Join(root, oav), 0755); err != nil {
		t.Fatal(err)
	}

	if IsInstalled(oav) {
		t.Fatal("empty CEF directory should not be recognized as installed")
	}
}
