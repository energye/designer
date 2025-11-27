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

//go:build windows
// +build windows

package frameworks

import (
	"archive/zip"
	"bytes"
	"embed"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"path/filepath"
)

//go:embed wv/wv.zip
var wv embed.FS

// 释放 WV 框架源码库
func extractWV(outputPath string) {
	// 存在 go.mod 文件则不进行解压
	if tool.IsExist(filepath.Join(outputPath, "go.mod")) {
		return
	}
	// 提到父目录, 因为解压时包含 wv 目录
	outputPath = filepath.Join(outputPath, "../")
	data, e := lcl.ReadFile("wv/wv.zip")
	err.CheckErr(e)
	zipReader, e := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	err.CheckErr(e)
	for _, file := range zipReader.File {
		_, e := tool.ExtractFile(file, outputPath)
		err.CheckErr(e)
	}
}
