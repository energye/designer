package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCEFManifestSupportsOldFormat(t *testing.T) {
	chromium := TChromium{Dir: t.TempDir()}
	data := []byte(`{
		"linux_amd64_127.3.5": [
			{"name": "libcef.so", "size": 123}
		]
	}`)
	if err := os.WriteFile(filepath.Join(chromium.Dir, ".versions"), data, 0644); err != nil {
		t.Fatal(err)
	}

	entry := chromium.LoadCEFManifest()["linux_amd64_127.3.5"]
	if len(entry.Files) != 1 || entry.Files[0].Name != "libcef.so" || entry.Files[0].Size != 123 {
		t.Fatalf("old manifest files not loaded: %#v", entry.Files)
	}
	if entry.Runtime != nil {
		t.Fatalf("old manifest should not have runtime info: %#v", entry.Runtime)
	}
}

func TestSaveCEFRuntimeManifestPreservesCEFFiles(t *testing.T) {
	chromium := TChromium{Dir: t.TempDir()}
	oav := "linux_amd64_127.3.5"
	if err := chromium.SaveCEFManifest(oav, []CEFFileInfo{{Name: "libcef.so", Size: 123}}); err != nil {
		t.Fatal(err)
	}
	if err := chromium.SaveCEFRuntimeManifest(oav, "v3.0.1", CEFFileInfo{Name: "libenergy-amd64-gtk3.so", Size: 456}); err != nil {
		t.Fatal(err)
	}

	entry := chromium.LoadCEFManifest()[oav]
	if len(entry.Files) != 1 || entry.Files[0].Name != "libcef.so" {
		t.Fatalf("CEF files were not preserved: %#v", entry.Files)
	}
	if entry.Runtime == nil || entry.Runtime.Release != "v3.0.1" || entry.Runtime.File.Size != 456 {
		t.Fatalf("runtime manifest not saved: %#v", entry.Runtime)
	}
}
