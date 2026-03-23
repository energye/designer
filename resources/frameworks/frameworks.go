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
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/tool/exec"
	"io"
	"os"
	"path/filepath"
)

// 启用的模块依赖配置
const (
	EnableLCL    = true
	EnableCEF    = true
	EnableWV     = true
	EnableENERGY = true
)

var (
	// Path 框架目录
	Path = filepath.Join(exec.AppDir(), "frameworks")
	// RuntimePath 运行时库目录
	RuntimePath = filepath.Join(Path, "runtime")
)

// ExtractLibrary 解压设计器运行时库 libenergy
// 这个函数作为解压过程的入口点
func ExtractLibrary() string {
	frameworksDir := config.Config.FrameworkDir
	runtimeDir := filepath.Join(frameworksDir, "runtime")
	if config.Config.FrameworkDir != "" && tool.IsExist(runtimeDir) {
		// 确保运行时库被释放
		_ = os.MkdirAll(runtimeDir, os.ModePerm)
		// 释放LCL框架文件
		return lib.ExtractLibrary(runtimeDir)
	} else if !tool.IsDarwin {
		// windows, linux 在未设置框架目录时，将运行时库释放到临时目录
		tempDir := os.TempDir()
		return lib.ExtractLibrary(tempDir)
	}
	//_ = os.MkdirAll(RuntimePath, os.ModePerm)
	// 释放LCL框架文件
	//return lib.ExtractLibrary(RuntimePath)
	return ""
}

// ExtractFrameworks 提取所有启用的框架
// 该函数会根据启用的标志位，依次提取LCL、CEF和WV框架
func ExtractFrameworks(fn func(message string), ok func()) {
	go func() {
		ExtractLCL(EnableLCL, fn)
		ExtractCEF(EnableCEF, fn)
		ExtractWV(EnableWV, fn)
		ExtractENERGY(EnableENERGY, fn)
		if ok != nil {
			ok()
		}
	}()
}

// ExtractLCL 根据enable参数决定是否执行 LCL 库提取操作
func ExtractLCL(enable bool, fn func(message string)) {
	if enable {
		fn("提取 LCL 组件库")
		// LocalPath LCL 框架源码路径
		LocalPath := config.Config.FrameworkDirForLCL()
		_ = os.MkdirAll(LocalPath, os.ModePerm)
		extractLCL(LocalPath)
	}
}

// ExtractCEF 根据enable参数决定是否执行 CEF 库提取操作
func ExtractCEF(enable bool, fn func(message string)) {
	if enable {
		fn("提取 CEF 组件库")
		LocalPath := config.Config.FrameworkDirForCEF()
		_ = os.MkdirAll(LocalPath, os.ModePerm)
		extractCEF(LocalPath)
	}
}

// ExtractWV 根据enable参数决定是否执行 WebView 库提取操作
func ExtractWV(enable bool, fn func(message string)) {
	if enable {
		fn("提取 Webview 组件库")
		LocalPath := config.Config.FrameworkDirForWV()
		_ = os.MkdirAll(LocalPath, os.ModePerm)
		extractWV(LocalPath)
	}
}

// ExtractENERGY 根据enable参数决定是否执行 energy 库提取操作
func ExtractENERGY(enable bool, fn func(message string)) {
	if enable {
		fn("提取 ENERGY GUI 框架库")
		LocalPath := config.Config.FrameworkDirForENERGY()
		_ = os.MkdirAll(LocalPath, os.ModePerm)
		extractENERGY(LocalPath)
	}
}

// readFileForZIPData 从ZIP数据中读取指定文件的内容
//
//	zipData: ZIP格式的字节数据
//	targetFileName: 要读取的目标文件名
func readFileForZIPData(zipData []byte, targetFileName string) (data []byte, err error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}
	for _, file := range zipReader.File {
		if !file.FileInfo().IsDir() && file.Name == targetFileName {
			srcFile, err := file.Open()
			if err != nil {
				return nil, err
			}
			data, err = io.ReadAll(srcFile)
			if err != nil {
				return nil, fmt.Errorf("读取文件[%s]内容失败: %w", targetFileName, err)
			}
			_ = srcFile.Close()
			break
		}
	}
	return
}
