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

package main

import (
	"github.com/energye/designer/designer"
	_ "github.com/energye/designer/internal"
	"github.com/energye/designer/pkg/logs"
	_ "github.com/energye/designer/pkg/syso"
	"github.com/energye/designer/resources/lib"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	logs.Level = logs.LevelDebug
	logs.Level = logs.LevelError
	{
		// 这是一段测试时用的代码
		libname.LibName = func() string {
			wd, _ := os.Getwd()
			return filepath.Join(wd, "../", "gen", "gout", libname.GetDLLName())
		}()
		if !tool.IsExist(libname.LibName) {
			libname.LibName = lib.Path
		}
	}
	logs.Debug(strings.Join(os.Args, " "))
	lcl.Init(nil, nil)
	// 运行设计器
	designer.Run()
}
