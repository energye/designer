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
)

// go install golang.org/x/tools/gopls@latest

type LSPClient struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	scanner   *bufio.Scanner
	mu        sync.Mutex
	requestID int
}

func NewLSPClient(workspaceDir string) (*LSPClient, error) {
	cmd := exec.Command("gopls")
	cmd.Dir = workspaceDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	client := &LSPClient{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		scanner: bufio.NewScanner(stdout),
	}

	// 启动响应监听
	go client.listenResponses()

	return client, nil
}

func (c *LSPClient) Initialize(rootURI string) error {
	c.mu.Lock()
	c.requestID++
	id := c.requestID
	c.mu.Unlock()

	params := map[string]interface{}{
		"processId": nil,
		"rootUri":   rootURI,
		"capabilities": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"completion": map[string]interface{}{
					"completionItem": map[string]interface{}{
						"snippetSupport": true,
					},
				},
			},
		},
	}

	_, err := c.sendRequest("initialize", id, params)
	return err
}

func (c *LSPClient) Completion(fileURI string, line, column int) ([]CompletionItem, error) {
	c.mu.Lock()
	c.requestID++
	id := c.requestID
	c.mu.Unlock()

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": fileURI,
		},
		"position": map[string]int{
			"line":      line,
			"character": column,
		},
	}

	resp, err := c.sendRequest("textDocument/completion", id, params)
	if err != nil {
		return nil, err
	}

	var result CompletionResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (c *LSPClient) DidOpen(fileURI, languageID, content string) error {
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":        fileURI,
			"languageId": languageID,
			"version":    1,
			"text":       content,
		},
	}

	return c.sendNotification("textDocument/didOpen", params)
}

func (c *LSPClient) sendRequest(method string, id int, params interface{}) ([]byte, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return nil, err
	}

	if _, err := c.stdin.Write(data); err != nil {
		return nil, err
	}

	// TODO: 等待响应(需要实现响应匹配逻辑)
	return nil, nil
}

func (c *LSPClient) sendNotification(method string, params interface{}) error {
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

func (c *LSPClient) listenResponses() {
	for c.scanner.Scan() {
		line := c.scanner.Text()
		if strings.HasPrefix(line, "Content-Length: ") {
			length, _ := strconv.Atoi(strings.TrimPrefix(line, "Content-Length: "))
			c.scanner.Scan() // 跳过空行

			buf := make([]byte, length)
			c.stdout.Read(buf)

			// 处理响应
			var resp map[string]interface{}
			json.Unmarshal(buf, &resp)

			// TODO: 根据 resp["id"] 匹配请求
		}
	}
}

func (c *LSPClient) Close() {
	c.stdin.Close()
	c.cmd.Wait()
}

type CompletionItem struct {
	Label      string `json:"label"`
	Kind       int    `json:"kind"`
	InsertText string `json:"insertText,omitempty"`
}

type CompletionResponse struct {
	Items []CompletionItem `json:"items"`
}
