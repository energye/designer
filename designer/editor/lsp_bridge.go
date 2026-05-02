// Copyright © yanghy.. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package editor

import (
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"os"
	"path/filepath"
	"strings"
)

// JSCompletionItem 发送给前端的补全项
type JSCompletionItem struct {
	Label               string       `json:"label"`
	Kind                int          `json:"kind"`
	Detail              string       `json:"detail,omitempty"`
	Documentation       string       `json:"documentation,omitempty"`
	SortText            string       `json:"sortText,omitempty"`
	FilterText          string       `json:"filterText,omitempty"`
	InsertText          string       `json:"insertText,omitempty"`
	InsertTextFormat    int          `json:"insertTextFormat,omitempty"`
	AdditionalTextEdits []JSTextEdit `json:"additionalTextEdits,omitempty"`
	Preselect           bool         `json:"preselect,omitempty"`
	Deprecated          bool         `json:"deprecated,omitempty"`
}

type JSTextEdit struct {
	Range   JSRange `json:"range"`
	NewText string  `json:"newText"`
}

type JSRange struct {
	Start JSPosition `json:"start"`
	End   JSPosition `json:"end"`
}

type JSPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type JSParameter struct {
	Label         string `json:"label"`
	Documentation string `json:"documentation,omitempty"`
}

type JSSignature struct {
	Label         string        `json:"label"`
	Documentation string        `json:"documentation,omitempty"`
	Parameters    []JSParameter `json:"parameters,omitempty"`
}

type JSSignatureHelpResult struct {
	Signatures      []JSSignature `json:"signatures"`
	ActiveSignature int           `json:"activeSignature"`
	ActiveParameter int           `json:"activeParameter"`
}

type JSWorkspaceEdit struct {
	Changes map[string][]JSTextEdit `json:"changes"`
}

type JSCodeAction struct {
	Title       string           `json:"title"`
	Kind        string           `json:"kind,omitempty"`
	IsPreferred bool             `json:"isPreferred,omitempty"`
	Edit        *JSWorkspaceEdit `json:"edit,omitempty"`
}

// FileData 文件数据结构，用于IPC文件读写
type FileData struct {
	File     string `json:"file"`
	Content  string `json:"content"`
	Language string `json:"language"`
	ModTime  int64  `json:"modTime"`
	ReadOnly bool   `json:"readOnly"`
}

func uriToFilePath(uri string) string {
	if tool.IsWindows {
		if !strings.HasPrefix(uri, "file:///") {
			return ""
		}
		path := strings.TrimPrefix(uri, "file:///")
		return strings.ReplaceAll(path, "/", "\\")
	}
	if !strings.HasPrefix(uri, "file://") {
		return ""
	}
	return strings.TrimPrefix(uri, "file://")
}

// isFileReadOnly checks if a file should be opened as read-only.
// Files in the Go module cache (GOMODCACHE) or vendor directory are read-only.
// Files that cannot be written to are also read-only.
func isFileReadOnly(filePath string) bool {
	// Check actual file write permission
	if !isWritable(filePath) {
		return true
	}
	// Check if file is in Go module cache
	modCache := os.Getenv("GOMODCACHE")
	if modCache != "" {
		absPath, _ := filepath.Abs(filePath)
		if strings.HasPrefix(filepath.ToSlash(absPath), filepath.ToSlash(modCache)) {
			return true
		}
	}
	// Check if file is in project vendor directory
	codePath := bean.CodePath()
	if codePath != "" {
		vendorPath := filepath.Join(codePath, "vendor")
		absPath, _ := filepath.Abs(filePath)
		if strings.HasPrefix(filepath.ToSlash(absPath), filepath.ToSlash(vendorPath)) {
			return true
		}
	}
	return false
}

// isWritable checks if a file is writable by attempting to open it for writing
func isWritable(filePath string) bool {
	f, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func filePathToURI(filePath string) string {
	uri := filepath.ToSlash(filePath)
	if tool.IsWindows {
		return "file:///" + uri
	}
	return "file://" + uri
}

// isTextFile checks if a file is a text file that can be edited.
// Returns false for binary files, images, archives, etc.
func isTextFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	// Images
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".ico", ".svg", ".webp", ".tiff", ".tif", ".icns":
		return false
	// Audio/Video
	case ".mp3", ".mp4", ".wav", ".avi", ".mkv", ".flac", ".ogg", ".wma", ".wmv", ".mov":
		return false
	// Archives
	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar", ".tgz", ".zst":
		return false
	// Binaries/Executables
	case ".exe", ".dll", ".so", ".dylib", ".bin", ".dat", ".o", ".a", ".lib", ".pdb":
		return false
	// Fonts
	case ".ttf", ".otf", ".woff", ".woff2", ".eot":
		return false
	// Documents
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return false
	// Database
	case ".db", ".sqlite", ".mdb":
		return false
	// Certificate/Key (often binary)
	case ".p12", ".pfx", ".der", ".crt", ".pem":
		return true // PEM is text-based
	}
	return true
}

func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".jsx":
		return "javascript"
	case ".tsx":
		return "typescript"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".scss":
		return "scss"
	case ".less":
		return "less"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".md":
		return "markdown"
	case ".py":
		return "python"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".hpp", ".cc":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".php":
		return "php"
	case ".rb":
		return "ruby"
	case ".sql":
		return "sql"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	case ".sh":
		return "shell"
	case ".bat":
		return "batch"
	default:
		return "plaintext"
	}
}
