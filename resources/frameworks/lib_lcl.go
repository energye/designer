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
)

//go:embed lcl/lcl.zip
var lcl embed.FS

// 释放 LCL 框架库
func extractLCL(outputPath string) {
	if tool.IsExist(outputPath) {
		return
	}
	data, e := lcl.ReadFile("lcl/lcl.zip")
	err.CheckErr(e)
	zipReader, e := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	err.CheckErr(e)
	for _, file := range zipReader.File {
		_, e := tool.ExtractFile(file, Path)
		err.CheckErr(e)
	}
}
