package gopls

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type discardWriteCloser struct{}

func (discardWriteCloser) Write(p []byte) (int, error) { return len(p), nil }
func (discardWriteCloser) Close() error                { return nil }

// fakeServer 用 io.Pipe 模拟 gopls 进程的 stdin/stdout，
// 不启动真实进程即可驱动 listenResponses / sendRequest。
type fakeServer struct {
	client       *PLSClient
	toClient     *io.PipeWriter
	clientReader *bufio.Reader
}

func newFakeServer(t *testing.T) *fakeServer {
	t.Helper()
	outR, outW := io.Pipe()
	inR, inW := io.Pipe()
	c := &PLSClient{
		stdin:   inW,
		stdout:  outR,
		reader:  bufio.NewReader(outR),
		pending: make(map[int]chan lspResponse),
	}
	go c.listenResponses()
	t.Cleanup(func() {
		_ = outW.Close()
		_ = inW.Close()
		_ = outR.Close()
		_ = inR.Close()
	})
	return &fakeServer{client: c, toClient: outW, clientReader: bufio.NewReader(inR)}
}

func (s *fakeServer) writeFrame(t *testing.T, v interface{}) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprintf(s.toClient, "Content-Length: %d\r\n\r\n%s", len(data), data); err != nil {
		t.Fatal(err)
	}
}

func (s *fakeServer) readFrame(t *testing.T) map[string]interface{} {
	t.Helper()
	length := 0
	for {
		line, err := s.clientReader.ReadString('\n')
		if err != nil {
			t.Fatal("读客户端消息失败:", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			length, err = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:")))
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	if length <= 0 {
		t.Fatal("非法 Content-Length")
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(s.clientReader, buf); err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

type requestResult struct {
	data []byte
	err  error
}

func startRequest(s *fakeServer, method string, timeout time.Duration) chan requestResult {
	ch := make(chan requestResult, 1)
	go func() {
		d, e := s.client.sendRequestTimeout(method, map[string]interface{}{}, timeout)
		ch <- requestResult{data: d, err: e}
	}()
	return ch
}

func requestID(t *testing.T, frame map[string]interface{}) int {
	t.Helper()
	id, ok := frame["id"].(float64)
	if !ok {
		t.Fatalf("请求缺少 id: %v", frame)
	}
	return int(id)
}

func TestResponseLoop_NormalResult(t *testing.T) {
	s := newFakeServer(t)
	ch := startRequest(s, "initialize", 2*time.Second)
	id := requestID(t, s.readFrame(t))
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": id,
		"result": map[string]interface{}{"ready": true},
	})
	select {
	case r := <-ch:
		if r.err != nil {
			t.Fatal(r.err)
		}
		if !strings.Contains(string(r.data), `"ready":true`) {
			t.Fatalf("结果不正确: %s", r.data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到响应")
	}
}

func TestResponseLoop_JSONRPCErrorIsNotResult(t *testing.T) {
	s := newFakeServer(t)
	ch := startRequest(s, "textDocument/completion", 2*time.Second)
	id := requestID(t, s.readFrame(t))
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": id,
		"error": map[string]interface{}{"code": -32601, "message": "boom"},
	})
	select {
	case r := <-ch:
		if r.err == nil {
			t.Fatal("JSON-RPC error 应作为错误返回")
		}
		if r.data != nil {
			t.Fatalf("错误响应不应带 result: %s", r.data)
		}
		if !strings.Contains(r.err.Error(), "boom") {
			t.Fatalf("错误信息丢失: %v", r.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到错误响应")
	}
}

func TestResponseLoop_ServerRequestDoesNotMatchPending(t *testing.T) {
	s := newFakeServer(t)
	ch := startRequest(s, "textDocument/completion", 3*time.Second)
	id := requestID(t, s.readFrame(t))

	// 服务端主动请求，故意复用同一个 id —— 旧实现会把它当响应投递。
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": id, "method": "workspace/configuration",
		"params": map[string]interface{}{"items": []interface{}{}},
	})
	reply := s.readFrame(t)
	if _, ok := reply["result"]; !ok {
		t.Fatalf("应回复服务端请求，实际: %v", reply)
	}

	select {
	case <-ch:
		t.Fatal("服务端请求被错误地当成响应投递给 pending")
	case <-time.After(100 * time.Millisecond):
	}

	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": id, "result": []interface{}{},
	})
	select {
	case r := <-ch:
		if r.err != nil {
			t.Fatal(r.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("真正的响应未送达")
	}
}

func TestResponseLoop_EOFWakesAllPending(t *testing.T) {
	s := newFakeServer(t)
	const n = 3
	results := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			_, e := s.client.sendRequestTimeout("x", map[string]interface{}{}, 5*time.Second)
			results <- e
		}()
	}
	for i := 0; i < n; i++ {
		s.readFrame(t)
	}
	_ = s.toClient.Close() // 模拟 gopls 进程退出

	start := time.Now()
	for i := 0; i < n; i++ {
		select {
		case e := <-results:
			if e == nil {
				t.Fatal("连接断开时应返回错误")
			}
		case <-time.After(2 * time.Second):
			t.Fatal("pending 未被唤醒")
		}
	}
	if time.Since(start) > time.Second {
		t.Fatal("唤醒不应等待请求超时")
	}
}

func TestConnectionClosed_RequestFailsFast(t *testing.T) {
	s := newFakeServer(t)
	_ = s.toClient.Close()
	for i := 0; i < 100; i++ {
		s.client.pendingMu.Lock()
		closed := s.client.closed
		s.client.pendingMu.Unlock()
		if closed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	start := time.Now()
	if _, err := s.client.sendRequest("x", map[string]interface{}{}); err == nil {
		t.Fatal("连接关闭后请求应立即失败")
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Fatal("关闭后请求不应挂到 10s 超时")
	}
}

func TestResponseLoop_PublishDiagnostics(t *testing.T) {
	s := newFakeServer(t)
	got := make(chan string, 1)
	s.client.SetDiagnosticsHandler(func(uri string, diagnostics []Diagnostic) {
		got <- uri
	})
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         "file:///a.go",
			"diagnostics": []interface{}{},
		},
	})
	select {
	case uri := <-got:
		if uri != "file:///a.go" {
			t.Fatalf("诊断 URI 不正确: %s", uri)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("诊断处理器未被调用")
	}
}

func TestClose_Idempotent(t *testing.T) {
	c := &PLSClient{
		stdin:   discardWriteCloser{},
		pending: make(map[int]chan lspResponse),
	}
	start := time.Now()
	c.Close()
	c.Close()
	if time.Since(start) > 2*time.Second {
		t.Fatal("Close 不应卡住")
	}
	c.pendingMu.Lock()
	closed := c.closed
	c.pendingMu.Unlock()
	if !closed {
		t.Fatal("Close 应标记连接已关闭")
	}
	if _, err := c.sendRequest("x", map[string]interface{}{}); err == nil {
		t.Fatal("关闭后请求应失败")
	}
}

func TestResolveGoplsExecutable_ExplicitEnv(t *testing.T) {
	p := filepath.Join(t.TempDir(), "my-gopls")
	if err := os.WriteFile(p, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(GoplsPathEnv, p)
	got, err := resolveGoplsExecutable()
	if err != nil {
		t.Fatal(err)
	}
	if got != p {
		t.Fatalf("应返回配置路径 %s，实际 %s", p, got)
	}
}

func TestRequestTimeout_IndexingWindow(t *testing.T) {
	if DefaultRequestTimeout != 30*time.Second || IndexingRequestTimeout != 120*time.Second {
		t.Fatal("超时常量与前端 pls_bridge.js 的注释约束不一致")
	}
	c := &PLSClient{pending: make(map[int]chan lspResponse)}
	if c.requestTimeout("textDocument/codeAction") != DefaultRequestTimeout {
		t.Fatal("正常状态应使用 10s 超时")
	}
	if c.requestTimeout("initialize") != IndexingRequestTimeout {
		t.Fatal("initialize 应使用长超时")
	}
	c.BeginIndexing()
	if c.requestTimeout("textDocument/codeAction") != IndexingRequestTimeout {
		t.Fatal("首次索引期间首个 CodeAction 不应被 10s 超时丢弃")
	}
	c.setIndexing(false)
	if c.requestTimeout("textDocument/codeAction") != DefaultRequestTimeout {
		t.Fatal("索引结束后应恢复 10s 超时")
	}
}

func TestIndexingState_CallbackAndGraceDeadline(t *testing.T) {
	c := &PLSClient{pending: make(map[int]chan lspResponse)}
	var mu sync.Mutex
	var events []bool
	c.SetIndexingHandler(func(v bool) {
		mu.Lock()
		events = append(events, v)
		mu.Unlock()
	})
	c.BeginIndexing()
	c.BeginIndexing() // 重复调用不应重复通知
	if !c.IsIndexing() {
		t.Fatal("BeginIndexing 后应处于索引窗口")
	}
	c.setIndexing(false)
	if c.IsIndexing() {
		t.Fatal("索引结束后应为 false")
	}

	// 兜底期限到达后自动结束，避免 gopls 不发进度通知时一直挂着长超时。
	c.BeginIndexing()
	c.indexingMu.Lock()
	c.indexingDeadline = time.Now().Add(-time.Second)
	c.indexingMu.Unlock()
	if c.IsIndexing() {
		t.Fatal("超过兜底期限应自动结束索引窗口")
	}
	mu.Lock()
	defer mu.Unlock()
	if len(events) != 4 || !events[0] || events[1] || !events[2] || events[3] {
		t.Fatalf("索引状态变化通知不正确: %v", events)
	}
}

func TestShouldRetryAfterIndexing(t *testing.T) {
	timeoutErr := fmt.Errorf("%w: textDocument/codeAction", errRequestTimeout)
	if !shouldRetryAfterIndexing("textDocument/codeAction", timeoutErr, true, false) {
		t.Fatal("索引期间超时且索引已结束时应重试一次")
	}
	if shouldRetryAfterIndexing("textDocument/codeAction", timeoutErr, true, true) {
		t.Fatal("仍在索引期间不应立即重试")
	}
	if shouldRetryAfterIndexing("textDocument/codeAction", timeoutErr, false, false) {
		t.Fatal("非索引期间的请求不应重试")
	}
	if shouldRetryAfterIndexing("shutdown", timeoutErr, true, false) {
		t.Fatal("shutdown/initialize 等有副作用的请求不应重试")
	}
	if shouldRetryAfterIndexing("textDocument/codeAction", errors.New("boom"), true, false) {
		t.Fatal("非超时的 JSON-RPC 错误不应重试")
	}
}

type hoverResult struct {
	hover *Hover
	err   error
}

func startHover(s *fakeServer) chan hoverResult {
	ch := make(chan hoverResult, 1)
	go func() {
		h, e := s.client.Hover("file:///a.go", 2, 3)
		ch <- hoverResult{hover: h, err: e}
	}()
	return ch
}

func TestHover_MarkupContentAndRange(t *testing.T) {
	s := newFakeServer(t)
	ch := startHover(s)
	frame := s.readFrame(t)
	if frame["method"] != "textDocument/hover" {
		t.Fatalf("请求方法不正确: %v", frame["method"])
	}
	params := frame["params"].(map[string]interface{})
	pos := params["position"].(map[string]interface{})
	if pos["line"] != float64(2) || pos["character"] != float64(3) {
		t.Fatalf("位置参数不正确: %v", pos)
	}
	if uri := params["textDocument"].(map[string]interface{})["uri"]; uri != "file:///a.go" {
		t.Fatalf("URI 不正确: %v", uri)
	}

	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": requestID(t, frame),
		"result": map[string]interface{}{
			"contents": map[string]interface{}{"kind": "markdown", "value": "```go\nfunc F()\n```"},
			"range": map[string]interface{}{
				"start": map[string]interface{}{"line": 1, "character": 2},
				"end":   map[string]interface{}{"line": 1, "character": 9},
			},
		},
	})
	select {
	case r := <-ch:
		if r.err != nil {
			t.Fatal(r.err)
		}
		if r.hover == nil || len(r.hover.Contents) != 1 {
			t.Fatalf("hover 内容解析失败: %+v", r.hover)
		}
		if r.hover.Contents[0].Value != "```go\nfunc F()\n```" {
			t.Fatalf("Markdown 内容丢失: %q", r.hover.Contents[0].Value)
		}
		if r.hover.Range == nil || r.hover.Range.Start.Line != 1 || r.hover.Range.End.Character != 9 {
			t.Fatalf("range 解析失败: %+v", r.hover.Range)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到 hover 结果")
	}
}

func TestHover_MarkedStringArray(t *testing.T) {
	s := newFakeServer(t)
	ch := startHover(s)
	frame := s.readFrame(t)
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": requestID(t, frame),
		"result": map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{"language": "go", "value": "func Selftest(a int) int"},
				"文档说明",
			},
		},
	})
	select {
	case r := <-ch:
		if r.err != nil || r.hover == nil {
			t.Fatalf("hover 解析失败: %v %+v", r.err, r.hover)
		}
		if len(r.hover.Contents) != 2 {
			t.Fatalf("应有 2 段内容，实际 %d", len(r.hover.Contents))
		}
		if r.hover.Contents[0].Language != "go" || r.hover.Contents[0].Value != "func Selftest(a int) int" {
			t.Fatalf("MarkedString 解析失败: %+v", r.hover.Contents[0])
		}
		if r.hover.Contents[1].Value != "文档说明" || r.hover.Contents[1].Language != "" {
			t.Fatalf("纯字符串解析失败: %+v", r.hover.Contents[1])
		}
		if r.hover.Range != nil {
			t.Fatal("无 range 时不应填充")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到 hover 结果")
	}
}

func TestHover_NullAndEmptyContents(t *testing.T) {
	s := newFakeServer(t)
	ch := startHover(s)
	frame := s.readFrame(t)
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": requestID(t, frame), "result": nil,
	})
	select {
	case r := <-ch:
		if r.err != nil || r.hover != nil {
			t.Fatalf("null 结果应返回 nil hover: %v %+v", r.err, r.hover)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到 hover 结果")
	}
}

func TestHover_JSONRPCError(t *testing.T) {
	s := newFakeServer(t)
	ch := startHover(s)
	frame := s.readFrame(t)
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": requestID(t, frame),
		"error": map[string]interface{}{"code": -32601, "message": "no hover"},
	})
	select {
	case r := <-ch:
		if r.err == nil || r.hover != nil {
			t.Fatal("JSON-RPC error 应作为错误返回")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到 hover 错误")
	}
}

func TestParseHoverContents_EmptyArray(t *testing.T) {
	if got := parseHoverContents(json.RawMessage(`[]`)); len(got) != 0 {
		t.Fatalf("空数组应返回空内容，实际 %v", got)
	}
	if got := parseHoverContents(json.RawMessage(`null`)); got != nil {
		t.Fatalf("null 应返回 nil，实际 %v", got)
	}
	if got := parseHoverContents(json.RawMessage(`"plain text"`)); len(got) != 1 || got[0].Value != "plain text" {
		t.Fatalf("单个字符串解析失败: %v", got)
	}
}

func TestResponseLoop_WorkDoneProgressCreateReply(t *testing.T) {
	s := newFakeServer(t)
	s.writeFrame(t, map[string]interface{}{
		"jsonrpc": "2.0", "id": 7, "method": "window/workDoneProgress/create",
		"params": map[string]interface{}{"token": "gopls/workspace"},
	})
	reply := s.readFrame(t)
	result, ok := reply["result"]
	if !ok {
		t.Fatalf("create 请求必须回复 result，实际: %v", reply)
	}
	if result != nil {
		t.Fatalf("create 的 result 应为 null，实际: %v", result)
	}
	if _, ok := reply["error"]; ok {
		t.Fatalf("create 不应回 MethodNotFound: %v", reply)
	}
}
