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

package frameworks

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"os"
	"path/filepath"
)

//go:embed cef/cef.zip
var cef embed.FS

// 释放 CEF 框架源码库
func extractCEF(outputPath string) {
	// 存在 go.mod 文件则不进行解压
	gomod := filepath.Join(outputPath, "go.mod")
	if tool.IsExist(gomod) {
		return
	}
	zipData, e := cef.ReadFile("cef/cef.zip")
	err.CheckErr(e)
	zipReader, e := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	err.CheckErr(e)
	for _, file := range zipReader.File {
		_, e := tool.ExtractFile(file, outputPath, "")
		err.CheckErr(e)
	}
	replace := fmt.Sprintf(`replace (
	github.com/energye/lcl => %v
)`, config.Config.FrameworkDirForLCLRelativePath())
	data, e := renderModLocalTemplate("github.com/energye/cef", replace)
	err.CheckErr(e)
	e = os.WriteFile(gomod, data, 0644)
	err.CheckErr(e)
	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.HideWindow = true
	cmd.Dir = outputPath
	cmd.Console = func(data string, level command.Level) {
		logs.Debug("go mod tidy:", data)
	}
	cmd.Command("go", "mod", "tidy")
}

// CEF 从嵌入的 cef.zip 文件中读取指定文件的内容
//
//	targetFileName - 要读取的目标文件名
func CEF(targetFileName string) (data []byte, err error) {
	zipData, err := cef.ReadFile("cef/cef.zip")
	if err != nil {
		return nil, err
	}
	return readFileForZIPData(zipData, targetFileName)
}
