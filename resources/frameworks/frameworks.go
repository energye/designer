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
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/exec"
	"path/filepath"
)

//go:embed lcl/lcl.zip
var lcl embed.FS

// 是否已释放依赖框架
var gIsExtractFramework bool

// extractLCL 从嵌入的zip文件中提取LCL框架文件
// 读取lcl/lcl.zip文件，创建zip读取器，并将所有文件解压到可执行文件目录下的frameworks文件夹中
func extractLCL() {
	data, e := lcl.ReadFile("lcl/lcl.zip")
	err.CheckErr(e)
	zipReader, e := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	err.CheckErr(e)
	outPath := filepath.Join(exec.Dir, "frameworks")
	for _, file := range zipReader.File {
		e := tool.ExtractFile(file, outPath)
		err.CheckErr(e)
	}
}
