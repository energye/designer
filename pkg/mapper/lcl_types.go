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

package mapper

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"path/filepath"
)

// lcl 源码模块常量
var _LCLTypeFiles = []string{
	"lcl.go",
	"cursors.go",
	"consts.go",
}

// GetLCL 获取映射的类型值
func GetLCL(name string) any {
	if v := cacheValue.Get(name); v != nil {
		return v
	}
	modCachePath := config.GGoEnv.Get("GOMODCACHE")
	dependencies := config.DesignerConfig.Dependencies
	lclModVersion := dependencies.Get(config.LCLModPath)
	lclMod := fmt.Sprintf("%s@%s", config.LCLModPath, lclModVersion)
	lclPath := filepath.Join(modCachePath, lclMod)
	for _, file := range _LCLTypeFiles {
		filePath := filepath.Join(lclPath, "types", file)
		val := dast.GetConstValue(filePath, name)
		if val != nil {
			cacheValue.Add(name, val)
			return val
		}
	}
	logs.Error("mapper.GetLCL 获取常量值失败, 常量名:", name)
	return nil
}
