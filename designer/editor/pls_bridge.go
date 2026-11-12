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
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
)

// URIToFilePath 把 gopls 回传的 file URI 还原成本地路径。
// 必须走 url.Parse 解码 %XX（中文、空格、# 等），否则前端按 filePath
// 查 Monaco model 会静默丢弃诊断与跳转结果。Windows 盘符统一大写。
func URIToFilePath(uri string) string {
	return uriToFilePath(uri, tool.IsWindows)
}

func uriToFilePath(uri string, windows bool) string {
	if uri == "" {
		return ""
	}
	if !strings.HasPrefix(strings.ToLower(uri), "file://") {
		// 少数实现会直接回传本地路径。
		return nativeFilePath(uri, windows)
	}
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	p := u.Path
	if p == "" {
		p = u.Opaque
	}
	if p == "" {
		return ""
	}
	return nativeFilePath(p, windows)
}

func nativeFilePath(p string, windows bool) string {
	p = strings.ReplaceAll(p, "\\", "/")
	if !windows {
		return p
	}
	p = strings.TrimPrefix(p, "/")
	p = strings.ReplaceAll(p, "/", "\\")
	if len(p) >= 2 && p[1] == ':' && isASCIILetter(p[0]) {
		p = strings.ToUpper(p[:1]) + p[1:]
	}
	return p
}

// FilePathToURI 按 RFC 8089 生成 file URI：路径分隔符、盘符大小写与
// 特殊字符全部由 url.URL 统一处理，保证与 gopls 回传的 URI 可互转。
func FilePathToURI(filePath string) string {
	return filePathToURI(filePath, tool.IsWindows)
}

func filePathToURI(filePath string, windows bool) string {
	p := filePath
	if windows {
		p = strings.ReplaceAll(p, "\\", "/")
	} else {
		p = filepath.ToSlash(p)
	}
	// LSP/gopls 在 Windows 上使用小写盘符，主动归一避免同一文件被当成两份文档。
	if windows && len(p) >= 2 && p[1] == ':' && isASCIILetter(p[0]) {
		p = strings.ToLower(p[:1]) + p[1:]
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return (&url.URL{Scheme: "file", Path: p}).String()
}

// CanonicalPathKey 生成路径比较用的 key，与前端 normalizePathKey 保持一致。
// 用于诊断/CodeAction 结果和文档版本表的匹配，忽略盘符大小写与斜杠差异。
func CanonicalPathKey(filePath string) string {
	return canonicalPathKey(filePath, tool.IsWindows)
}

func canonicalPathKey(p string, windows bool) string {
	p = strings.ReplaceAll(p, "\\", "/")
	if strings.HasPrefix(strings.ToLower(p), "file://") {
		if u, err := url.Parse(p); err == nil {
			p = u.Path
			if windows && len(p) >= 3 && p[0] == '/' && p[2] == ':' {
				p = p[1:]
			}
		}
	}
	if len(p) >= 2 && p[1] == ':' && isASCIILetter(p[0]) {
		p = strings.ToLower(p[:1]) + p[1:]
	}
	for strings.Contains(p, "//") {
		p = strings.ReplaceAll(p, "//", "/")
	}
	return p
}

func isASCIILetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// IsFileReadOnly checks if a file should be opened as read-only.
// Files in the Go module cache (GOMODCACHE) or vendor directory are read-only.
// Files that cannot be written to are also read-only.
func IsFileReadOnly(filePath string) bool {
	if !IsWritable(filePath) {
		return true
	}
	modCache := os.Getenv("GOMODCACHE")
	if modCache != "" {
		absPath, _ := filepath.Abs(filePath)
		if strings.HasPrefix(filepath.ToSlash(absPath), filepath.ToSlash(modCache)) {
			return true
		}
	}
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

// IsWritable checks if a file is writable by attempting to open it for writing
func IsWritable(filePath string) bool {
	f, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

// IsTextFile checks if a file is a text file that can be edited.
// Returns false for binary files, images, archives, etc.
func IsTextFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".ico", ".svg", ".webp", ".tiff", ".tif", ".icns":
		return false
	case ".mp3", ".mp4", ".wav", ".avi", ".mkv", ".flac", ".ogg", ".wma", ".wmv", ".mov":
		return false
	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar", ".tgz", ".zst":
		return false
	case ".exe", ".dll", ".so", ".dylib", ".bin", ".dat", ".o", ".a", ".lib", ".pdb":
		return false
	case ".ttf", ".otf", ".woff", ".woff2", ".eot":
		return false
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return false
	case ".db", ".sqlite", ".mdb":
		return false
	case ".p12", ".pfx", ".der", ".crt", ".pem":
		return true
	}
	return true
}

func DetectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go", ".mod":
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
	case ".json", ".egp":
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
