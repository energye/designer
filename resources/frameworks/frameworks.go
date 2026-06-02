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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"os"
)

// ExtractLibrary 解压设计器运行时库 libenergy
// 这个函数作为解压过程的入口点
func ExtractLibrary() string {
	frameworksDir := config.Config.FrameworkDir
	runtimeDir := config.Config.FrameworkRuntimePath()
	if config.Config.FrameworkDir != "" && tool.IsExist(frameworksDir) {
		// 确保运行时库被释放
		_ = os.MkdirAll(runtimeDir, os.ModePerm)
		// 释放LCL框架文件
		return lib.ExtractLibrary(runtimeDir)
	} else if !tool.IsDarwin {
		// windows, linux 在未设置框架目录时，将运行时库释放到临时目录
		tempDir := os.TempDir()
		return lib.ExtractLibrary(tempDir)
	}
	return ""
}
