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
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/api/libname"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

var (
	Path string
)

func init() {
	tempDir := os.TempDir()
	outPath := filepath.Join(tempDir, libname.GetDLLName())
	libByte, e := lib.ReadFile(path)
	err.CheckErr(e)
	zipReader, e := zip.NewReader(bytes.NewReader(libByte), int64(len(libByte)))
	err.CheckErr(e)

	for _, file := range zipReader.File {
		e := extractFile(file, outPath)
		err.CheckErr(e)
		break // 只有一个文件
	}
	Path = outPath
	logs.Info("Lib Path:", outPath)
}

func extractFile(zipFile *zip.File, targetFile string) error {
	srcFile, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		if runtime.GOOS == "windows" {
			if pathErr, ok := err.(*os.PathError); ok {
				if errno, ok := pathErr.Err.(syscall.Errno); ok && errno == 32 {
					logs.Error("File is busy, skipping extraction:", targetFile)
					return nil
				}
			}
		}
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}
