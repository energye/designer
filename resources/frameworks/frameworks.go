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
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/tool/exec"
	"os"
	"path/filepath"
)

// 启用的模块依赖配置
const (
	EnableLCL = true
	EnableCEF = true
	EnableWV  = true
)

var (
	// Path 框架目录
	Path = filepath.Join(exec.Dir, "frameworks")
	// RuntimePath 运行时库目录
	RuntimePath = filepath.Join(Path, "runtime")
)

var (
	src = filepath.Join(Path, "src")
	// LCLLocalPath LCL 框架源码路径
	LCLLocalPath = filepath.Join(src, "lcl")
	// CEFLocalPath CEF 框架源码路径
	CEFLocalPath = filepath.Join(src, "cef")
	// WVLocalPath WebView 框架源码路径, darwin/linux/windows
	WVLocalPath = filepath.Join(src, "wv")
)

// ExtractLibrary 解压设计器运行时库 libenergy
// 这个函数作为解压过程的入口点
func ExtractLibrary() {
	_ = os.MkdirAll(RuntimePath, os.ModePerm)
	// 释放LCL框架文件
	lib.ExtractLibrary(RuntimePath)
}

// ExtractLCL 根据enable参数决定是否执行 LCL 库提取操作
func ExtractLCL(enable bool) {
	if enable {
		_ = os.MkdirAll(src, os.ModePerm)
		extractLCL(LCLLocalPath)
	}
}

// ExtractCEF 根据enable参数决定是否执行 CEF 库提取操作
func ExtractCEF(enable bool) {
	if enable {
		_ = os.MkdirAll(src, os.ModePerm)
		extractLCL(CEFLocalPath)
	}
}

// ExtractWV 根据enable参数决定是否执行 WebView 库提取操作
func ExtractWV(enable bool) {
	if enable {
		_ = os.MkdirAll(src, os.ModePerm)
		extractLCL(WVLocalPath)
	}
}
