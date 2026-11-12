package gopls

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/energye/designer/pkg/logs"
)

// 请求超时与前端 pls_bridge.js 的常量成对约束：
// 前端等待时间必须略大于这里的值，否则响应回来时 pending 已被清空。
//
//	DefaultRequestTimeout  30s  <->  PLS_REQUEST_TIMEOUT_MS  35s
//	IndexingRequestTimeout 120s <->  PLS_INDEXING_TIMEOUT_MS 125s
//
// 首次索引（go.work 覆盖大工作区时实测约 45s）期间使用长超时；
// 真机还观测到新增包触发的二次索引会让 CodeAction 排队约 17s，
// 因此默认超时取 30s，不能只靠索引窗口保护。
const (
	DefaultRequestTimeout  = 30 * time.Second
	IndexingRequestTimeout = 120 * time.Second
	indexingGracePeriod    = 120 * time.Second
)

var errRequestTimeout = errors.New("gopls 请求超时")

// lspResponse 区分「服务端正常结果」与「JSON-RPC 错误 / 连接断开」。
// 旧实现把 error 对象当 result 塞进通道，调用方解析失败后静默返回空结果。
type lspResponse struct {
	result json.RawMessage
	err    error
}

type PLSClient struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	reader    *bufio.Reader
	mu        sync.Mutex
	requestID int

	pending   map[int]chan lspResponse
	pendingMu sync.Mutex
	closed    bool
	closeErr  error
	closeOnce sync.Once

	diagnosticsHandler func(uri string, diagnostics []Diagnostic)

	indexingMu       sync.Mutex
	indexing         bool
	indexingDeadline time.Time
	indexingHandler  func(bool)
}

func NewPLSClient(workspaceDir string) (*PLSClient, error) {
	exe, err := resolveGoplsExecutable()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(exe)
	cmd.Dir = workspaceDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("gopls stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("gopls stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("gopls stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("gopls start: %w", err)
	}

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				return
			}
			logs.Debug("gopls stderr:", strings.TrimRight(string(buf[:n]), "\n"))
		}
	}()

	client := &PLSClient{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		reader:  bufio.NewReaderSize(stdout, 10*1024*1024),
		pending: make(map[int]chan lspResponse),
	}

	go client.listenResponses()

	logs.Debug("gopls 启动成功:", exe, "工作目录:", workspaceDir, "PID:", cmd.Process.Pid)
	return client, nil
}

func (c *PLSClient) SetDiagnosticsHandler(handler func(uri string, diagnostics []Diagnostic)) {
	c.mu.Lock()
	c.diagnosticsHandler = handler
	c.mu.Unlock()
}

// SetIndexingHandler 索引窗口开始/结束时回调，用于把状态同步到编辑器页面。
func (c *PLSClient) SetIndexingHandler(handler func(bool)) {
	c.indexingMu.Lock()
	c.indexingHandler = handler
	c.indexingMu.Unlock()
}

// BeginIndexing 在 initialize 完成后开启首次索引窗口。
// 窗口内请求使用长超时；第一个成功的 textDocument 请求会提前结束窗口，
// 兜底期限到达后也会自动结束，避免 gopls 不发进度通知时一直挂着。
func (c *PLSClient) BeginIndexing() {
	c.setIndexing(true)
	// 兜底：即使期间没有任何请求，也要在期限到达后结束窗口，
	// 否则代码编辑页的「正在建立索引」提示会一直挂着。
	go func() {
		<-time.After(indexingGracePeriod + time.Second)
		_ = c.IsIndexing() // 惰性检查会清掉过期状态并回调 handler
	}()
}

func (c *PLSClient) IsIndexing() bool {
	c.indexingMu.Lock()
	expired := c.indexing && !c.indexingDeadline.IsZero() && time.Now().After(c.indexingDeadline)
	c.indexingMu.Unlock()
	if expired {
		c.setIndexing(false)
	}
	c.indexingMu.Lock()
	defer c.indexingMu.Unlock()
	return c.indexing
}

func (c *PLSClient) setIndexing(value bool) {
	c.indexingMu.Lock()
	if c.indexing == value {
		c.indexingMu.Unlock()
		return
	}
	c.indexing = value
	if value {
		c.indexingDeadline = time.Now().Add(indexingGracePeriod)
	} else {
		c.indexingDeadline = time.Time{}
	}
	handler := c.indexingHandler
	c.indexingMu.Unlock()
	if handler != nil {
		handler(value)
	}
}

func (c *PLSClient) Initialize(rootURI string) error {
	logs.Debug("gopls Initialize rootURI:", rootURI)
	params := map[string]interface{}{
		"processId": nil,
		"rootUri":   rootURI,
		"capabilities": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"completion": map[string]interface{}{
					"completionItem": map[string]interface{}{
						"snippetSupport":    true,
						"deprecatedSupport": true,
						// 不声明 resolveSupport：gopls 未实现 completionItem/resolve，
						// 声明后会对每个候选项发出注定失败的请求（真机已验证）。
						"insertTextModeSupport": map[string]interface{}{
							"valueSet": []int{1, 2},
						},
						"labelDetailsSupport": true,
					},
					"insertTextMode": 2,
					"contextSupport": true,
				},
				"signatureHelp": map[string]interface{}{
					"signatureInformation": map[string]interface{}{
						"documentationFormat": []string{"markdown", "plaintext"},
						"parameterInformation": map[string]interface{}{
							"labelOffsetSupport": true,
						},
					},
					"contextSupport": true,
				},
				"hover": map[string]interface{}{
					"contentFormat": []string{"markdown", "plaintext"},
				},
				"codeAction": map[string]interface{}{
					"codeActionLiteralSupport": map[string]interface{}{
						"codeActionKind": map[string]interface{}{
							"valueSet": []string{
								"quickfix",
								"source", "source.organizeImports",
							},
						},
					},
					"isPreferredSupport": true,
				},
				"publishDiagnostics": map[string]interface{}{
					"relatedInformation": true,
					"versionSupport":     false,
					"tagSupport": map[string]interface{}{
						"valueSet": []int{1, 2},
					},
					"codeDescriptionSupport": true,
				},
				"definition": map[string]interface{}{
					"linkSupport": true,
				},
				"synchronization": map[string]interface{}{
					"didSave": true,
				},
			},
		},
	}

	_, err := c.sendRequest("initialize", params)
	if err != nil {
		logs.Error("gopls Initialize 失败:", err)
		return err
	}

	c.sendNotification("initialized", map[string]interface{}{})
	logs.Debug("gopls Initialize 完成")
	return nil
}

func (c *PLSClient) Completion(fileURI string, line, column int, triggerKind int, triggerChar string) ([]CompletionItem, error) {
	if ShouldIgnoreFile(fileURI, nil) {
		return nil, errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls Completion 请求: uri=", fileURI, "line=", line, "column=", column, "triggerKind=", triggerKind, "triggerChar=", triggerChar)
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"position": map[string]int{
			"line":      line,
			"character": column,
		},
		"context": map[string]interface{}{
			"triggerKind":      triggerKind,
			"triggerCharacter": triggerChar,
		},
	}

	resp, err := c.sendRequest("textDocument/completion", params)
	if err != nil {
		logs.Error("gopls Completion 请求失败:", err)
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("gopls 返回空响应")
	}

	var completionList struct {
		Items        []CompletionItem `json:"items"`
		IsIncomplete bool             `json:"isIncomplete"`
	}
	if err := json.Unmarshal(resp, &completionList); err != nil {
		logs.Error("gopls Completion JSON解析失败:", err)
		return nil, nil
	}

	logs.Debug("gopls Completion 返回", len(completionList.Items), "个建议项, isIncomplete=", completionList.IsIncomplete)
	// 不再对 isIncomplete 的候选项逐个调用 completionItem/resolve：
	// 当前 gopls 未实现该方法，只会产生一批 JSON-RPC 错误和额外延迟，
	// 候选项本身已包含 label/detail/documentation/additionalTextEdits。
	return completionList.Items, nil
}

func (c *PLSClient) SignatureHelp(fileURI string, line, column int) (*SignatureHelpResult, error) {
	if ShouldIgnoreFile(fileURI, nil) {
		return nil, errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls SignatureHelp 请求: uri=", fileURI, "line=", line, "column=", column)
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"position": map[string]int{
			"line":      line,
			"character": column,
		},
	}

	resp, err := c.sendRequest("textDocument/signatureHelp", params)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var result SignatureHelpResult
	if err := json.Unmarshal(resp, &result); err != nil {
		logs.Error("gopls SignatureHelp JSON解析失败:", err)
		return nil, nil
	}

	logs.Debug("gopls SignatureHelp 返回", len(result.Signatures), "个签名, activeSignature=", result.ActiveSignature, "activeParameter=", result.ActiveParameter)
	return &result, nil
}

// Hover 返回 textDocument/hover 的结果。LSP 的 contents 有多种形态
// （MarkupContent、MarkedString、字符串/对象数组、null），统一解析成
// []HoverContent 后再交给前端拼 Markdown。
func (c *PLSClient) Hover(fileURI string, line, column int) (*Hover, error) {
	if ShouldIgnoreFile(fileURI, nil) {
		return nil, errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls Hover 请求: uri=", fileURI, "line=", line, "column=", column)
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"position": map[string]int{
			"line":      line,
			"character": column,
		},
	}

	resp, err := c.sendRequest("textDocument/hover", params)
	if err != nil {
		return nil, err
	}
	if resp == nil || string(resp) == "null" {
		return nil, nil
	}

	var raw struct {
		Contents json.RawMessage `json:"contents"`
		Range    *Range          `json:"range"`
	}
	if err := json.Unmarshal(resp, &raw); err != nil {
		logs.Error("gopls Hover JSON解析失败:", err)
		return nil, nil
	}
	result := &Hover{Contents: parseHoverContents(raw.Contents), Range: raw.Range}
	if len(result.Contents) == 0 {
		return nil, nil
	}
	logs.Debug("gopls Hover 返回", len(result.Contents), "段内容")
	return result, nil
}

func parseHoverContents(raw json.RawMessage) []HoverContent {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err == nil {
		result := make([]HoverContent, 0, len(items))
		for _, item := range items {
			if c, ok := parseHoverContent(item); ok {
				result = append(result, c)
			}
		}
		return result
	}
	// 单个 MarkupContent / MarkedString / 纯字符串。
	if c, ok := parseHoverContent(raw); ok {
		return []HoverContent{c}
	}
	return nil
}

func parseHoverContent(raw json.RawMessage) (HoverContent, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return HoverContent{Value: s}, s != ""
	}
	var obj struct {
		Kind     string `json:"kind"`
		Value    string `json:"value"`
		Language string `json:"language"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil && obj.Value != "" {
		return HoverContent{Value: obj.Value, Language: obj.Language}, true
	}
	return HoverContent{}, false
}

func (c *PLSClient) CodeAction(fileURI string, startLine, startChar, endLine, endChar int, kinds []string, diagnostics []Diagnostic) ([]CodeAction, error) {
	if ShouldIgnoreFile(fileURI, nil) {
		return nil, errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls CodeAction 请求: uri=", fileURI, "kinds=", kinds, "diagnostics=", len(diagnostics))

	diagInterfaces := make([]interface{}, len(diagnostics))
	for i, d := range diagnostics {
		diagInterfaces[i] = d
	}

	// If no diagnostics provided but we have kinds like source.organizeImports,
	// send empty diagnostics - gopls will still process source actions
	if len(diagInterfaces) == 0 {
		diagInterfaces = []interface{}{}
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"range": map[string]interface{}{
			"start": map[string]int{"line": startLine, "character": startChar},
			"end":   map[string]int{"line": endLine, "character": endChar},
		},
		"context": map[string]interface{}{
			"diagnostics": diagInterfaces,
			"only":        kinds,
			"triggerKind": 2, // Invoked (not automatic)
		},
	}

	resp, err := c.sendRequest("textDocument/codeAction", params)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	// gopls may return either []CodeAction or []Command
	// Try CodeAction first
	var actions []CodeAction
	if err := json.Unmarshal(resp, &actions); err != nil {
		logs.Error("gopls CodeAction JSON解析失败:", err, "raw:", string(resp[:goplsMin(len(resp), 300)]))
		return nil, nil
	}
	logs.Debug("gopls CodeAction 返回", len(actions), "个操作")
	return actions, nil
}

func (c *PLSClient) DidOpen(fileURI, languageID, content string, version int) error {
	if ShouldIgnoreFile(fileURI, nil) {
		return errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls DidOpen: uri=", fileURI, "lang=", languageID, "version=", version, "contentLen=", len(content))
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":        fileURI,
			"languageId": languageID,
			"version":    version,
			"text":       content,
		},
	}

	return c.sendNotification("textDocument/didOpen", params)
}

func (c *PLSClient) DidChange(fileURI string, version int, content string) error {
	if ShouldIgnoreFile(fileURI, nil) {
		return errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls DidChange: uri=", fileURI, "version=", version, "contentLen=", len(content))
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":     fileURI,
			"version": version,
		},
		"contentChanges": []map[string]interface{}{
			{
				"text": content,
			},
		},
	}

	if err := c.sendNotification("textDocument/didChange", params); err != nil {
		logs.Error("gopls DidChange 发送失败:", err)
		return err
	}
	return nil
}

func (c *PLSClient) DidSave(fileURI string, text string) error {
	if ShouldIgnoreFile(fileURI, nil) {
		return errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls DidSave: uri=", fileURI, "textLen=", len(text))
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri": fileURI,
		},
		"text": text,
	}

	return c.sendNotification("textDocument/didSave", params)
}

func (c *PLSClient) Formatting(fileURI string) ([]TextEdit, error) {
	if ShouldIgnoreFile(fileURI, nil) {
		return nil, errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls Formatting 请求: uri=", fileURI)
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"options": map[string]interface{}{
			"tabSize":                4,
			"insertSpaces":           true,
			"trimTrailingWhitespace": true,
			"insertFinalNewline":     true,
		},
	}

	resp, err := c.sendRequest("textDocument/formatting", params)
	if err != nil {
		logs.Error("gopls Formatting 请求失败:", err)
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var edits []TextEdit
	if err := json.Unmarshal(resp, &edits); err != nil {
		logs.Error("gopls Formatting JSON解析失败:", err)
		return nil, nil
	}

	logs.Debug("gopls Formatting 返回", len(edits), "个编辑")
	return edits, nil
}

func (c *PLSClient) DidClose(fileURI string) error {
	if ShouldIgnoreFile(fileURI, nil) {
		return errors.New("ignored file directories: " + fileURI)
	}
	logs.Debug("gopls DidClose: uri=", fileURI)
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
	}

	return c.sendNotification("textDocument/didClose", params)
}

// requestTimeout 与前端 plsRequestTimeout() 成对约束：索引窗口内放宽到 120s。
func (c *PLSClient) requestTimeout(method string) time.Duration {
	// initialize 本身也可能因为工作区较大而变慢。
	if c.IsIndexing() || method == "initialize" {
		return IndexingRequestTimeout
	}
	return DefaultRequestTimeout
}

// shouldRetryAfterIndexing 首个请求在索引期间失败、且索引已结束时重试一次。
func shouldRetryAfterIndexing(method string, err error, indexingAtStart, indexingNow bool) bool {
	return indexingAtStart && !indexingNow && isTextDocumentRequest(method) && errors.Is(err, errRequestTimeout)
}

func (c *PLSClient) sendRequest(method string, params interface{}) ([]byte, error) {
	indexing := c.IsIndexing()
	resp, err := c.sendRequestTimeout(method, params, c.requestTimeout(method))
	if err == nil {
		// 首个成功的文本请求说明 gopls 已经能处理文档，索引窗口可以提前结束。
		if indexing && isTextDocumentRequest(method) {
			c.setIndexing(false)
		}
		return resp, nil
	}
	// 统一记录失败原因，避免补全/诊断静默变成「没有结果」。
	logs.Error("gopls 请求失败:", method, err)

	// 首次索引期间发出的只读请求若失败，而索引窗口已经结束，重试一次：
	// 覆盖「首个 CodeAction 在索引完成前发出、被服务端拒绝」的情况。
	if shouldRetryAfterIndexing(method, err, indexing, c.IsIndexing()) {
		resp, err = c.sendRequestTimeout(method, params, DefaultRequestTimeout)
		if err != nil {
			logs.Error("gopls 索引结束后重试失败:", method, err)
		}
	}
	return resp, err
}

func isTextDocumentRequest(method string) bool {
	switch method {
	case "textDocument/completion", "textDocument/signatureHelp", "textDocument/codeAction",
		"textDocument/formatting", "textDocument/definition", "textDocument/hover":
		return true
	}
	return false
}

func (c *PLSClient) sendRequestTimeout(method string, params interface{}, timeout time.Duration) ([]byte, error) {
	c.mu.Lock()
	c.requestID++
	id := c.requestID
	c.mu.Unlock()

	respChan := make(chan lspResponse, 1)
	c.pendingMu.Lock()
	if c.closed {
		err := c.closeErr
		c.pendingMu.Unlock()
		return nil, err
	}
	c.pending[id] = respChan
	c.pendingMu.Unlock()

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}
	if err := c.writeMessage(request); err != nil {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, err
	}

	select {
	case resp := <-respChan:
		if resp.err != nil {
			return nil, resp.err
		}
		return resp.result, nil
	case <-time.After(timeout):
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, fmt.Errorf("%w: %s", errRequestTimeout, method)
	}
}

// writeMessage 统一序列化并加锁写入 stdin（请求、通知、服务端请求回复共用）。
func (c *PLSClient) writeMessage(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stdin == nil {
		return io.ErrClosedPipe
	}
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return err
	}
	_, err = c.stdin.Write(data)
	return err
}

func (c *PLSClient) sendNotification(method string, params interface{}) error {
	return c.writeMessage(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	})
}

func (c *PLSClient) listenResponses() {
	// 无论正常退出还是读到 EOF，都要唤醒所有等待中的请求，
	// 否则它们会一直挂到 10s 超时。
	defer c.failPending(fmt.Errorf("gopls 进程已退出"))
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			// 进程退出（EOF）或主动关闭管道都按正常结束处理。
			if err != io.EOF && !errors.Is(err, io.ErrClosedPipe) {
				logs.Error("gopls listenResponses: 读取header失败:", err)
			}
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "Content-Length:") {
			continue
		}

		lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
		length, err := strconv.Atoi(lengthStr)
		if err != nil || length <= 0 {
			continue
		}

		for {
			headerLine, err := c.reader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.TrimSpace(headerLine) == "" {
				break
			}
		}

		buf := make([]byte, length)
		if _, err := io.ReadFull(c.reader, buf); err != nil {
			return
		}

		var resp struct {
			ID     int             `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
			Result json.RawMessage `json:"result"`
			Error  *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(buf, &resp); err != nil {
			logs.Warn("gopls 消息解析失败:", err)
			continue
		}

		// 服务端主动请求同时带 id 和 method，它不是响应。
		// 必须回一个回复让 gopls 继续，且绝不能投递到 pending 造成错配。
		if resp.ID > 0 && resp.Method != "" {
			c.replyToServerRequest(resp.ID, resp.Method)
			continue
		}

		if resp.ID > 0 {
			c.pendingMu.Lock()
			ch, ok := c.pending[resp.ID]
			if ok {
				delete(c.pending, resp.ID)
			}
			c.pendingMu.Unlock()
			if ok {
				var r lspResponse
				if resp.Error != nil {
					r.err = fmt.Errorf("gopls 返回错误(%d): %s", resp.Error.Code, resp.Error.Message)
				} else {
					r.result = resp.Result
				}
				select {
				case ch <- r:
				default:
				}
			}
			continue
		}

		if resp.Method == "textDocument/publishDiagnostics" && len(resp.Params) > 0 {
			c.handleDiagnostics(resp.Params)
		}
	}
}

// replyToServerRequest 回复 gopls 发起的 server→client 请求。
func (c *PLSClient) replyToServerRequest(id int, method string) {
	logs.Debug("gopls 服务端请求:", method)
	if method == "workspace/configuration" {
		_ = c.writeMessage(map[string]interface{}{
			"jsonrpc": "2.0", "id": id, "result": []interface{}{},
		})
		return
	}
	if method == "window/workDoneProgress/create" {
		// LSP 规定 create 回复 null；回 MethodNotFound 会让部分 gopls 版本放弃进度上报。
		_ = c.writeMessage(map[string]interface{}{
			"jsonrpc": "2.0", "id": id, "result": nil,
		})
		return
	}
	_ = c.writeMessage(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    -32601,
			"message": "method not supported: " + method,
		},
	})
}

// failPending 在连接断开时唤醒所有等待中的请求。
func (c *PLSClient) failPending(err error) {
	c.pendingMu.Lock()
	if c.closed {
		c.pendingMu.Unlock()
		return
	}
	c.closed = true
	c.closeErr = err
	for id, ch := range c.pending {
		delete(c.pending, id)
		select {
		case ch <- lspResponse{err: err}:
		default:
		}
	}
	c.pendingMu.Unlock()
}

// 诊断信息 触发时机
//
//	打开一个文件 (textDocument/didOpen)
//	修改了文件内容 (textDocument/didChange)
//	保存了文件 (textDocument/didSave)
func (c *PLSClient) handleDiagnostics(params json.RawMessage) {
	var notification struct {
		URI         string       `json:"uri"`
		Diagnostics []Diagnostic `json:"diagnostics"`
	}

	if err := json.Unmarshal(params, &notification); err != nil {
		return
	}

	c.mu.Lock()
	handler := c.diagnosticsHandler
	c.mu.Unlock()

	if handler != nil {
		handler(notification.URI, notification.Diagnostics)
	}
}

// Close 按 LSP 正常关停：shutdown → exit → 关 stdin → 等待进程，
// 超时则强制结束；并唤醒所有仍在等待的请求。可重复调用。
func (c *PLSClient) Close() {
	c.closeOnce.Do(func() {
		_, _ = c.sendRequestTimeout("shutdown", map[string]interface{}{}, 500*time.Millisecond)
		_ = c.sendNotification("exit", map[string]interface{}{})
		if c.stdin != nil {
			_ = c.stdin.Close()
		}
		if c.cmd != nil && c.cmd.Process != nil {
			done := make(chan struct{})
			go func() {
				_ = c.cmd.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				_ = c.cmd.Process.Kill()
				<-done
			}
		}
		c.failPending(fmt.Errorf("gopls 已关闭"))
	})
}

func goplsMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- PLS Type Definitions ---

// HoverContent 归一化后的 hover 内容段。Language 非空表示这段是带语言的
// 代码（LSP MarkedString），前端会包成 Markdown fenced code。
type HoverContent struct {
	Value    string `json:"value"`
	Language string `json:"language,omitempty"`
}

type Hover struct {
	Contents []HoverContent `json:"contents"`
	Range    *Range         `json:"range,omitempty"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

func (c *PLSClient) Definition(fileURI string, line, column int) ([]Location, error) {
	if ShouldIgnoreFile(fileURI, nil) {
		return nil, errors.New("ignored file directories: " + fileURI)
	}
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"position": map[string]int{
			"line":      line,
			"character": column,
		},
	}

	resp, err := c.sendRequest("textDocument/definition", params)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	// Response can be a single Location, an array of Locations, or an array of LocationLink
	var locations []Location

	// Try array first
	var arr []Location
	if err := json.Unmarshal(resp, &arr); err == nil && len(arr) > 0 {
		return arr, nil
	}

	// Try single Location
	var loc Location
	if err := json.Unmarshal(resp, &loc); err == nil && loc.URI != "" {
		locations = append(locations, loc)
		return locations, nil
	}

	return locations, nil
}

type CompletionItem struct {
	Label               string      `json:"label"`
	Kind                int         `json:"kind"`
	Detail              interface{} `json:"detail,omitempty"`
	Documentation       interface{} `json:"documentation,omitempty"`
	SortText            string      `json:"sortText,omitempty"`
	FilterText          string      `json:"filterText,omitempty"`
	InsertText          string      `json:"insertText,omitempty"`
	InsertTextFormat    int         `json:"insertTextFormat,omitempty"`
	Preselect           bool        `json:"preselect,omitempty"`
	AdditionalTextEdits []TextEdit  `json:"additionalTextEdits,omitempty"`
	TextEdit            interface{} `json:"textEdit,omitempty"`
	Deprecated          bool        `json:"deprecated,omitempty"`
}

func (item *CompletionItem) GetDocumentation() string {
	return ExtractString(item.Documentation)
}

func (item *CompletionItem) GetDetail() string {
	return ExtractString(item.Detail)
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type SignatureHelpResult struct {
	Signatures      []SignatureInformation `json:"signatures"`
	ActiveSignature int                    `json:"activeSignature"`
	ActiveParameter int                    `json:"activeParameter"`
}

type SignatureInformation struct {
	Label         string                 `json:"label"`
	Documentation interface{}            `json:"documentation,omitempty"`
	Parameters    []ParameterInformation `json:"parameters,omitempty"`
}

func (s *SignatureInformation) GetDocumentation() string {
	return ExtractString(s.Documentation)
}

type ParameterInformation struct {
	Label         interface{} `json:"label"`
	Documentation interface{} `json:"documentation,omitempty"`
}

func (p *ParameterInformation) GetLabel() string {
	switch v := p.Label.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) == 2 {
			return fmt.Sprintf("[%v,%v]", v[0], v[1])
		}
	}
	data, _ := json.Marshal(p.Label)
	return string(data)
}

type CodeAction struct {
	Title       string         `json:"title"`
	Kind        string         `json:"kind,omitempty"`
	Edit        *WorkspaceEdit `json:"edit,omitempty"`
	IsPreferred bool           `json:"isPreferred,omitempty"`
}

type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}

type Diagnostic struct {
	Range struct {
		Start struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"start"`
		End struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"end"`
	} `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
}

func ExtractString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case map[string]interface{}:
		if s, ok := val["value"].(string); ok {
			return s
		}
	}
	data, _ := json.Marshal(v)
	return string(data)
}
