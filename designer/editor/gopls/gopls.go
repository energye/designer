package gopls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/energye/designer/pkg/logs"
)

// go install golang.org/x/tools/gopls@latest

type PLSClient struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	reader    *bufio.Reader
	mu        sync.Mutex
	requestID int

	pending   map[int]chan []byte
	pendingMu sync.Mutex

	diagnosticsHandler func(uri string, diagnostics []Diagnostic)
	diagMu             sync.Mutex
}

func NewPLSClient(workspaceDir string) (*PLSClient, error) {
	cmd := exec.Command("gopls")
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
			logs.Info("gopls stderr:", strings.TrimRight(string(buf[:n]), "\n"))
		}
	}()

	client := &PLSClient{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		reader:  bufio.NewReaderSize(stdout, 10*1024*1024),
		pending: make(map[int]chan []byte),
	}

	go client.listenResponses()

	logs.Info("gopls 启动成功, 工作目录:", workspaceDir, "PID:", cmd.Process.Pid)
	return client, nil
}

func (c *PLSClient) SetDiagnosticsHandler(handler func(uri string, diagnostics []Diagnostic)) {
	c.diagMu.Lock()
	c.diagnosticsHandler = handler
	c.diagMu.Unlock()
}

func (c *PLSClient) Initialize(rootURI string) error {
	logs.Info("gopls Initialize rootURI:", rootURI)
	params := map[string]interface{}{
		"processId": nil,
		"rootUri":   rootURI,
		"capabilities": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"completion": map[string]interface{}{
					"completionItem": map[string]interface{}{
						"snippetSupport":    true,
						"deprecatedSupport": true,
						"resolveSupport": map[string]interface{}{
							"properties": []string{"documentation", "detail", "additionalTextEdits"},
						},
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
								"quickfix", "refactor", "refactor.extract",
								"refactor.inline", "refactor.rewrite",
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
					"dynamicRegistration": false,
					"willSave":            false,
					"willSaveWaitUntil":   false,
					"didSave":             true,
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
	logs.Info("gopls Initialize 完成")
	return nil
}

func (c *PLSClient) Completion(fileURI string, line, column int, triggerKind int, triggerChar string) ([]CompletionItem, error) {
	logs.Info("gopls Completion 请求: uri=", fileURI, "line=", line, "column=", column, "triggerKind=", triggerKind, "triggerChar=", triggerChar)
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

	logs.Info("gopls Completion 返回", len(completionList.Items), "个建议项, isIncomplete=", completionList.IsIncomplete)

	// If result is incomplete, resolve items to get full data (additionalTextEdits etc.)
	if completionList.IsIncomplete && len(completionList.Items) > 0 {
		for i := range completionList.Items {
			if len(completionList.Items[i].AdditionalTextEdits) == 0 {
				resolved, err := c.ResolveCompletionItem(completionList.Items[i])
				if err == nil && resolved != nil {
					if len(resolved.AdditionalTextEdits) > 0 {
						completionList.Items[i].AdditionalTextEdits = resolved.AdditionalTextEdits
					}
					if resolved.Documentation != nil {
						completionList.Items[i].Documentation = resolved.Documentation
					}
					if resolved.Detail != nil {
						completionList.Items[i].Detail = resolved.Detail
					}
				}
			}
		}
	}

	return completionList.Items, nil
}

func (c *PLSClient) SignatureHelp(fileURI string, line, column int) (*SignatureHelpResult, error) {
	logs.Info("gopls SignatureHelp 请求: uri=", fileURI, "line=", line, "column=", column)
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

	logs.Info("gopls SignatureHelp 返回", len(result.Signatures), "个签名, activeSignature=", result.ActiveSignature, "activeParameter=", result.ActiveParameter)
	return &result, nil
}

func (c *PLSClient) ResolveCompletionItem(item CompletionItem) (*CompletionItem, error) {
	logs.Info("gopls ResolveCompletionItem: label=", item.Label)
	resp, err := c.sendRequest("completionItem/resolve", item)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &item, nil
	}
	var resolved CompletionItem
	if err := json.Unmarshal(resp, &resolved); err != nil {
		return &item, nil
	}
	return &resolved, nil
}

func (c *PLSClient) CodeAction(fileURI string, startLine, startChar, endLine, endChar int, kinds []string, diagnostics []Diagnostic) ([]CodeAction, error) {
	logs.Info("gopls CodeAction 请求: uri=", fileURI, "kinds=", kinds, "diagnostics=", len(diagnostics))

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

	logs.Info("gopls CodeAction 返回", len(actions), "个操作")
	return actions, nil
}

func (c *PLSClient) DidOpen(fileURI, languageID, content string, version int) error {
	logs.Info("gopls DidOpen: uri=", fileURI, "lang=", languageID, "version=", version, "contentLen=", len(content))
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
	logs.Info("gopls DidChange: uri=", fileURI, "version=", version, "contentLen=", len(content))
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
	logs.Info("gopls DidSave: uri=", fileURI, "textLen=", len(text))
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri": fileURI,
		},
		"text": text,
	}

	return c.sendNotification("textDocument/didSave", params)
}

func (c *PLSClient) DidClose(fileURI string) error {
	logs.Info("gopls DidClose: uri=", fileURI)
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
	}

	return c.sendNotification("textDocument/didClose", params)
}

func (c *PLSClient) sendRequest(method string, params interface{}) ([]byte, error) {
	c.mu.Lock()
	c.requestID++
	id := c.requestID
	c.mu.Unlock()

	respChan := make(chan []byte, 1)
	c.pendingMu.Lock()
	c.pending[id] = respChan
	c.pendingMu.Unlock()

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}

	data, err := json.Marshal(request)
	if err != nil {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	c.mu.Lock()
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		c.mu.Unlock()
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, err
	}
	if _, err := c.stdin.Write(data); err != nil {
		c.mu.Unlock()
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, err
	}
	c.mu.Unlock()

	select {
	case resp, ok := <-respChan:
		if !ok {
			return nil, fmt.Errorf("gopls 连接已关闭")
		}
		return resp, nil
	case <-time.After(10 * time.Second):
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		return nil, fmt.Errorf("gopls 请求超时: %s", method)
	}
}

func (c *PLSClient) sendNotification(method string, params interface{}) error {
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return err
	}

	_, err = c.stdin.Write(data)
	return err
}

func (c *PLSClient) listenResponses() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
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
			break
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
			continue
		}

		if resp.ID > 0 {
			c.pendingMu.Lock()
			ch, ok := c.pending[resp.ID]
			c.pendingMu.Unlock()

			if ok {
				if resp.Error != nil {
					errData, _ := json.Marshal(resp.Error)
					select {
					case ch <- errData:
					default:
					}
				} else {
					select {
					case ch <- []byte(resp.Result):
					default:
					}
				}
			}
		}

		if resp.Method == "textDocument/publishDiagnostics" && len(resp.Params) > 0 {
			c.handleDiagnostics(resp.Params)
		}
	}
}

func (c *PLSClient) handleDiagnostics(params json.RawMessage) {
	var notification struct {
		URI         string       `json:"uri"`
		Diagnostics []Diagnostic `json:"diagnostics"`
	}

	if err := json.Unmarshal(params, &notification); err != nil {
		return
	}

	c.diagMu.Lock()
	handler := c.diagnosticsHandler
	c.diagMu.Unlock()

	if handler != nil {
		handler(notification.URI, notification.Diagnostics)
	}
}

func (c *PLSClient) Close() {
	c.stdin.Close()
	c.cmd.Wait()
}

func goplsMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- PLS Type Definitions ---

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

func (c *PLSClient) Definition(fileURI string, line, column int) ([]Location, error) {
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
