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

package lib

import (
	"archive/zip"
	"bytes"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api/libname"
	"path/filepath"
)

// ExtractLibrary 从内置资源中提取库文件到指定输出路径
//
//		outputPath: 库文件的输出目录路径
//		libPath: 提取后的库文件完整路径
//	 说明:
func ExtractLibrary(outputPath string) (libPath string) {
	libPath = filepath.Join(outputPath, libname.GetDLLName())
	if tool.IsExist(libPath) {
		return
	}
	libByte, e := lib.ReadFile(path)
	err.CheckErr(e)
	zipReader, e := zip.NewReader(bytes.NewReader(libByte), int64(len(libByte)))
	err.CheckErr(e)
	for _, file := range zipReader.File {
		_, e := tool.ExtractFile(file, outputPath)
		err.CheckErr(e)
		break // 只有一个文件
	}
	return
}
