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

func URIToFilePath(uri string) string {
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

func FilePathToURI(filePath string) string {
	uri := filepath.ToSlash(filePath)
	if tool.IsWindows {
		return "file:///" + uri
	}
	return "file://" + uri
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
