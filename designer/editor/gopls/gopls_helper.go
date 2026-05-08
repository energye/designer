// Copyright © yanghy. All Rights Reserved.
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

package gopls

import (
	"path/filepath"
	"strings"
)

var defaultCustomIgnoreList = []string{"node_modules", "vendor", "bin", "testdata", "dist"}

// ShouldIgnoreFile 判断文件路径是否被忽略
//
//   - filePath: 文件的完整路径 (例如: "/project/src/main.go")
//   - customIgnoreDirs: 自定义需要忽略的目录名列表 (例如: []string{"node_modules", "bin"})
func ShouldIgnoreFile(filePath string, customIgnoreDirs []string) bool {
	if filepath.Ext(filePath) != ".go" {
		return true
	}
	normalizedPath := filepath.ToSlash(filePath)
	// "xxx/.xxx"
	if strings.Contains(normalizedPath, "/.") {
		return true
	}
	customIgnoreDirs = append(customIgnoreDirs, defaultCustomIgnoreList...)
	for _, dir := range customIgnoreDirs {
		pattern := "/" + dir + "/"
		if strings.Contains(normalizedPath, pattern) {
			return true
		}
	}
	return false
}
