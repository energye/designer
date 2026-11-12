package editor

import (
	"strings"
	"testing"
)

func TestFilePathToURI_WindowsEscapesSpecialChars(t *testing.T) {
	input := `E:\app\workspace\中文 目录\a#1.go`
	uri := filePathToURI(input, true)
	if strings.Contains(uri, " ") || strings.ContainsAny(uri, "中文") {
		t.Fatalf("URI 未转义空格/中文: %s", uri)
	}
	want := "file:///e:/app/workspace/%E4%B8%AD%E6%96%87%20%E7%9B%AE%E5%BD%95/a%231.go"
	if uri != want {
		t.Fatalf("URI 不正确\n got: %s\nwant: %s", uri, want)
	}
	if got := uriToFilePath(uri, true); got != input {
		t.Fatalf("往返转换丢失信息: %s", got)
	}
}

func TestURIToFilePath_LowercaseDriveIsNormalized(t *testing.T) {
	got := uriToFilePath("file:///e:/app/workspace/a.go", true)
	want := `E:\app\workspace\a.go`
	if got != want {
		t.Fatalf("盘符应归一为大写: got %s want %s", got, want)
	}
	// gopls 回传的小写盘符路径必须能匹配前端打开时使用的 key。
	if canonicalPathKey(got, true) != canonicalPathKey(want, true) {
		t.Fatal("归一化 key 不一致")
	}
}

func TestFilePathToURI_POSIXRoundTrip(t *testing.T) {
	input := "/home/me/项目/a b.go"
	uri := filePathToURI(input, false)
	want := "file:///home/me/%E9%A1%B9%E7%9B%AE/a%20b.go"
	if uri != want {
		t.Fatalf("URI 不正确\n got: %s\nwant: %s", uri, want)
	}
	if got := uriToFilePath(uri, false); got != input {
		t.Fatalf("往返转换丢失信息: %s", got)
	}
}

func TestURIToFilePath_EncodedAndNonFile(t *testing.T) {
	got := uriToFilePath("file:///e:/app/中文%20目录/a%25b.go", true)
	want := `E:\app\中文 目录\a%b.go`
	if got != want {
		t.Fatalf("百分号转义未还原: got %s want %s", got, want)
	}
	if uriToFilePath("", true) != "" {
		t.Fatal("空 URI 应返回空路径")
	}
	if uriToFilePath("not-a-uri", false) != "not-a-uri" {
		t.Fatal("裸路径应原样返回")
	}
}

func TestCanonicalPathKey_IgnoresDriveCaseSeparatorAndURI(t *testing.T) {
	forms := []string{
		`E:\app\中文 目录\a.go`,
		"E:/app/中文 目录/a.go",
		"e:/app/中文 目录/a.go",
		"file:///e:/app/中文 目录/a.go",
		"file:///e:/app/%E4%B8%AD%E6%96%87%20%E7%9B%AE%E5%BD%95/a.go",
	}
	want := "e:/app/中文 目录/a.go"
	for _, form := range forms {
		if got := canonicalPathKey(form, true); got != want {
			t.Fatalf("归一化不一致: %q -> %q, want %q", form, got, want)
		}
	}
	if canonicalPathKey("/home/a.go", false) != canonicalPathKey("file:///home/a.go", false) {
		t.Fatal("POSIX 路径与 file URI 应归一为同一 key")
	}
	if canonicalPathKey("e:/app//a.go", true) != "e:/app/a.go" {
		t.Fatal("重复斜杠未折叠")
	}
}
