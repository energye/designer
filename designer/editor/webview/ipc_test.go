package webview

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/gopls"
)

// fakeContext 只实现 ipc.IContext 的三个方法，避免单测依赖 LCL/WebView。
type fakeContext struct {
	data    any
	results []any
}

func (c *fakeContext) Data() any          { return c.data }
func (c *fakeContext) Result(data ...any) { c.results = append(c.results, data...) }
func (c *fakeContext) BrowserId() uint32  { return 0 }

func TestParseIPCParams(t *testing.T) {
	ctx := &fakeContext{data: []any{map[string]any{
		"requestID": 7, "file": "a.go", "line": 3, "column": 4,
	}}}
	var params CompletionParams
	if !parseIPCParams(ctx, &params, "[]") {
		t.Fatal("合法参数应解析成功")
	}
	if params.RequestID != 7 || params.File != "a.go" || params.Line != 3 || params.Column != 4 {
		t.Fatalf("参数解析错误: %+v", params)
	}
	if len(ctx.results) != 0 {
		t.Fatal("解析成功不应写入默认结果")
	}

	bad := &fakeContext{data: []any{}}
	var out CompletionParams
	if parseIPCParams(bad, &out, "[]") {
		t.Fatal("空参数应解析失败")
	}
	if len(bad.results) != 1 || bad.results[0] != "[]" {
		t.Fatalf("解析失败应回写默认结果: %v", bad.results)
	}
}

func TestParseIPCString(t *testing.T) {
	ok := &fakeContext{data: []any{"a.go"}}
	if s, good := parseIPCString(ok, "ok"); !good || s != "a.go" {
		t.Fatalf("字符串参数解析错误: %q %v", s, good)
	}
	bad := &fakeContext{data: []any{42}}
	if _, good := parseIPCString(bad, "error"); good {
		t.Fatal("非字符串参数应解析失败")
	}
	if len(bad.results) != 1 || bad.results[0] != "error" {
		t.Fatalf("解析失败应回写默认结果: %v", bad.results)
	}
}

func TestDidOpenNotification_EscapesURI(t *testing.T) {
	file := filepath.Join(t.TempDir(), "中文 目录", "a#1.go")
	n := didOpenNotification(DidOpenParams{
		File: file, LanguageID: "go", Content: "package main", Version: 9,
	})
	if n.Kind != editor.PLSNotifyOpen || n.LanguageID != "go" || n.Content != "package main" || n.Version != 9 {
		t.Fatalf("didOpen 参数丢失: %+v", n)
	}
	if strings.Contains(n.URI, " ") || strings.ContainsAny(n.URI, "中文") || strings.Contains(n.URI, "#") {
		t.Fatalf("URI 未正确转义: %s", n.URI)
	}
	if !strings.HasPrefix(n.URI, "file:///") || !strings.Contains(n.URI, "%") {
		t.Fatalf("URI 格式不正确: %s", n.URI)
	}
	if got := editor.URIToFilePath(n.URI); got == "" {
		t.Fatal("URI 应能还原成本地路径")
	}
}

func TestDidChangeCloseSaveNotifications(t *testing.T) {
	const file = "a.go"
	change := didChangeNotification(DidChangeParams{File: file, Content: "new", Version: 12})
	if change.Kind != editor.PLSNotifyChange || change.Content != "new" || change.Version != 12 {
		t.Fatalf("didChange 参数丢失: %+v", change)
	}
	closeN := didCloseNotification(file)
	if closeN.Kind != editor.PLSNotifyClose || closeN.URI != change.URI {
		t.Fatal("didClose 与 didChange 必须指向同一 URI")
	}
	saveN := didSaveNotification(FileData{File: file, Content: "saved"})
	if saveN.Kind != editor.PLSNotifySave || saveN.Content != "saved" || saveN.URI != change.URI {
		t.Fatalf("didSave 参数丢失: %+v", saveN)
	}
}

func TestWorkspaceEditToJS_ConvertsURIKeys(t *testing.T) {
	file := filepath.Join(t.TempDir(), "中文 目录", "a.go")
	uri := editor.FilePathToURI(file)
	edit := &gopls.WorkspaceEdit{
		Changes: map[string][]gopls.TextEdit{
			uri: {{NewText: "import \"fmt\"", Range: gopls.Range{
				Start: gopls.Position{Line: 1, Character: 2},
				End:   gopls.Position{Line: 1, Character: 5},
			}}},
		},
	}
	js := workspaceEditToJS(edit)
	if js == nil || len(js.Changes) != 1 {
		t.Fatalf("编辑项转换失败: %+v", js)
	}
	edits, ok := js.Changes[editor.URIToFilePath(uri)]
	if !ok || len(edits) != 1 {
		t.Fatalf("URI key 未还原成本地路径: %+v", js.Changes)
	}
	if edits[0].NewText != "import \"fmt\"" || edits[0].Range.Start.Line != 1 || edits[0].Range.End.Character != 5 {
		t.Fatalf("编辑内容丢失: %+v", edits[0])
	}
	if workspaceEditToJS(nil) != nil || workspaceEditToJS(&gopls.WorkspaceEdit{}) != nil {
		t.Fatal("空 WorkspaceEdit 应返回 nil")
	}
}

func TestJSTextEdits_MapsAllFields(t *testing.T) {
	got := jsTextEdits([]gopls.TextEdit{
		{NewText: "a", Range: gopls.Range{Start: gopls.Position{Line: 0, Character: 1}, End: gopls.Position{Line: 2, Character: 3}}},
		{NewText: "b", Range: gopls.Range{Start: gopls.Position{Line: 4, Character: 5}, End: gopls.Position{Line: 6, Character: 7}}},
	})
	if len(got) != 2 || got[0].NewText != "a" || got[1].NewText != "b" {
		t.Fatalf("文本映射错误: %+v", got)
	}
	if got[1].Range.Start.Line != 4 || got[1].Range.End.Character != 7 {
		t.Fatalf("坐标映射错误: %+v", got[1])
	}
}

func TestHoverToJS_MarkupAndRange(t *testing.T) {
	h := &gopls.Hover{
		Contents: []gopls.HoverContent{
			{Value: "```go\nfunc F()\n```"},
			{Value: "func F()", Language: "go"},
		},
		Range: &gopls.Range{
			Start: gopls.Position{Line: 1, Character: 2},
			End:   gopls.Position{Line: 1, Character: 9},
		},
	}
	js := hoverToJS(h)
	if js == nil || len(js.Contents) != 2 {
		t.Fatalf("hover 内容转换失败: %+v", js)
	}
	if js.Contents[0].Value != "```go\nfunc F()\n```" || js.Contents[0].Language != "" {
		t.Fatalf("Markdown 内容丢失: %+v", js.Contents[0])
	}
	if js.Contents[1].Language != "go" || js.Contents[1].Value != "func F()" {
		t.Fatalf("带语言的代码段丢失: %+v", js.Contents[1])
	}
	if js.Range == nil || js.Range.Start.Line != 1 || js.Range.End.Character != 9 {
		t.Fatalf("range 转换失败: %+v", js.Range)
	}
}

func TestHoverToJS_NilOrEmpty(t *testing.T) {
	if hoverToJS(nil) != nil {
		t.Fatal("nil hover 应返回 nil")
	}
	if hoverToJS(&gopls.Hover{}) != nil {
		t.Fatal("空内容 hover 应返回 nil")
	}
	js := hoverToJS(&gopls.Hover{Contents: []gopls.HoverContent{{Value: "doc"}}})
	if js.Range != nil {
		t.Fatal("无 range 时不应填充")
	}
}

func TestParseIPCParams_HoverParams(t *testing.T) {
	ctx := &fakeContext{data: []any{map[string]any{
		"requestID": 11, "file": "a.go", "line": 7, "column": 8,
	}}}
	var params HoverParams
	if !parseIPCParams(ctx, &params, "null") {
		t.Fatal("hover 参数应解析成功")
	}
	if params.RequestID != 11 || params.File != "a.go" || params.Line != 7 || params.Column != 8 {
		t.Fatalf("hover 参数解析错误: %+v", params)
	}
}

func TestEmitPLSResponse_IsInjectable(t *testing.T) {
	orig := emitPLSResponse
	t.Cleanup(func() { emitPLSResponse = orig })
	var gotReqID int
	var gotEvent, gotPayload string
	emitPLSResponse = func(reqID int, event, payload string) {
		gotReqID, gotEvent, gotPayload = reqID, event, payload
	}
	emitPLSResponse(3, "gopls-completion-response", "[]")
	if gotReqID != 3 || gotEvent != "gopls-completion-response" || gotPayload != "[]" {
		t.Fatalf("响应注入失效: %d %s %s", gotReqID, gotEvent, gotPayload)
	}
}
