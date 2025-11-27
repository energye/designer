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

package dependmod

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"path/filepath"
)

// 初始化模块类型信息
// 功能依赖 frameworks/src 源码:
//   lcl cef wv(windows, darwin, linux)
// 获得类型信息, 回调函数类型元数据
// 作用:
//   事件绑定修改源码文件时使用
// 调用时机:
//   1. 设计器打开后 2. 项目创建后
//     检查 frameworks/src 源码是否已安装, 然后加载
//

// InitDependencyModule 初始化模块类型信息
// 每次调用都会重置
func InitDependencyModule() {
	logs.Println("初始化模块类型信息")
	go initModuleTypeInfo()
}

// 初始化模块类型信息
func initModuleTypeInfo() {
	frameDir := config.Config.FrameworkDir
	// LCL 模块的事件回调函数类型
	lclSRCEventDef := filepath.Join(frameDir, "src", "lcl", "lcl", "callback_event_def.go")
	gLCLFuncTypeAliases = dast.GetAllFuncTypeAliases(lclSRCEventDef)

	// CEF 模块的事件回调函数类型
	cefSRCEventDef := filepath.Join(frameDir, "src", "cef", "cef", "callback_event_def.go")
	gCEFFuncTypeAliases = dast.GetAllFuncTypeAliases(cefSRCEventDef)

	// WV 模块的事件回调函数类型
	// Windows
	wvWindowsSRCEventDef := filepath.Join(frameDir, "src", "wv", "windows", "callback_event_def.go")
	gWVWindowsFuncTypeAliases = dast.GetAllFuncTypeAliases(wvWindowsSRCEventDef)
	//  macOS
	wvDarwinSRCEventDef := filepath.Join(frameDir, "src", "wv", "darwin", "callback_event_def.go")
	gWVDarwinFuncTypeAliases = dast.GetAllFuncTypeAliases(wvDarwinSRCEventDef)
	// Linux
	wvLinuxSRCEventDef := filepath.Join(frameDir, "src", "wv", "linux", "callback_event_def.go")
	gWVLinuxFuncTypeAliases = dast.GetAllFuncTypeAliases(wvLinuxSRCEventDef)
}
