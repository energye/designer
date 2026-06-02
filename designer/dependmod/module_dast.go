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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"os"
	"path/filepath"
)

// 初始化模块类型信息
// 功能依赖 frameworks/src 源码:
//   lcl cef wv(windows, darwin, linux)
// 获得类型信息, 回调函数类型元数据
// 作用:
//   事件绑定修改源码文件时使用
// 调用时机:
//   应用启动期间
//

type modCacheDir struct {
	lclDir string
	cefDir string
	wvDir  string
	engDir string
}

// 从模块缓存 初始化模块类型信息
func initModuleTypeInfoFormModCache(dir *modCacheDir) {
	logs.Println("初始化模块类型信息")
	event.ConsoleWriteInfo("   LCL:", dir.lclDir)
	event.ConsoleWriteInfo("   CEF:", dir.cefDir)
	event.ConsoleWriteInfo("    WV:", dir.wvDir)
	event.ConsoleWriteInfo("ENERGY:", dir.engDir)

	lclSRCEventDef := filepath.Join(dir.lclDir, "lcl", "callback_event_def.go")
	cefSRCEventDef := filepath.Join(dir.cefDir, "cef", "callback_event_def.go")
	wvWindowsSRCEventDef := filepath.Join(dir.wvDir, "windows", "callback_event_def.go")
	wvDarwinSRCEventDef := filepath.Join(dir.wvDir, "darwin", "callback_event_def.go")
	wvLinuxSRCEventDef := filepath.Join(dir.wvDir, "linux", "callback_event_def.go")

	{
		// LCL 模块的事件回调函数类型
		lclSRCEventDefData, err := os.ReadFile(lclSRCEventDef)
		if err != nil {
			logs.Error("initModuleTypeInfoFormModCache", err.Error())
			return
		}
		GLCLFuncTypeAliases = dast.GetAllFuncTypeAliasesByCode(lclSRCEventDef, lclSRCEventDefData)
		if GLCLFuncTypeAliases != nil {
			GLCLFuncTypeAliases.Mod = consts.ModLCL
			GLCLFuncTypeAliases.Imports.Add(GLCLFuncTypeAliases.Mod, consts.DmLCL)
		}
	}
	{
		// CEF 模块的事件回调函数类型
		cefSRCEventDefData, err := os.ReadFile(cefSRCEventDef)
		if err != nil {
			logs.Error("initModuleTypeInfoFormModCache", err.Error())
			return
		}
		if GCEFFuncTypeAliases = dast.GetAllFuncTypeAliasesByCode(cefSRCEventDef, cefSRCEventDefData); GCEFFuncTypeAliases != nil {
			GCEFFuncTypeAliases.Mod = consts.ModCEF
			GCEFFuncTypeAliases.Imports.Add(GCEFFuncTypeAliases.Mod, consts.DmCEF)
		}
	}
	{
		// WV 模块的事件回调函数类型
		// Windows
		wvWindowsSRCEventDefData, err := os.ReadFile(wvWindowsSRCEventDef)
		if err != nil {
			logs.Error("initModuleTypeInfoFormModCache", err.Error())
			return
		}
		GWVWindowsFuncTypeAliases = dast.GetAllFuncTypeAliasesByCode(wvWindowsSRCEventDef, wvWindowsSRCEventDefData)
		if GWVWindowsFuncTypeAliases != nil {
			GWVWindowsFuncTypeAliases.Mod = consts.ModWVWindows
			GWVWindowsFuncTypeAliases.Imports.Add(GWVWindowsFuncTypeAliases.Mod, consts.DmWVWindows)
		}
		//  macOS
		wvDarwinSRCEventDefData, err := os.ReadFile(wvDarwinSRCEventDef)
		if err != nil {
			logs.Error("initModuleTypeInfoFormModCache", err.Error())
			return
		}
		GWVDarwinFuncTypeAliases = dast.GetAllFuncTypeAliasesByCode(wvDarwinSRCEventDef, wvDarwinSRCEventDefData)
		if GWVDarwinFuncTypeAliases != nil {
			GWVDarwinFuncTypeAliases.Mod = consts.ModWVDarwin
			GWVDarwinFuncTypeAliases.Imports.Add(GWVDarwinFuncTypeAliases.Mod, consts.DmWVMacOS)
		}
		// Linux
		wvLinuxSRCEventDefData, err := os.ReadFile(wvLinuxSRCEventDef)
		if err != nil {
			logs.Error("initModuleTypeInfoFormModCache", err.Error())
			return
		}
		GWVLinuxFuncTypeAliases = dast.GetAllFuncTypeAliasesByCode(wvLinuxSRCEventDef, wvLinuxSRCEventDefData)
		if GWVLinuxFuncTypeAliases != nil {
			GWVLinuxFuncTypeAliases.Mod = consts.ModWVLinux
			GWVLinuxFuncTypeAliases.Imports.Add(GWVLinuxFuncTypeAliases.Mod, consts.DmWVLinux)
		}
	}
	logs.Println("初始化模块类型信息 结束")
}
