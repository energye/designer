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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/dast"
	"path/filepath"
)

// GetLCL 获取映射的类型值
func GetLCL(name string) any {
	lclPath := config.Config.FrameworkDirForLCL()
	srcLCLTypes := filepath.Join(lclPath, "lcl", "types", "lcl.go")
	val := dast.GetConstValue(srcLCLTypes, name)
	return val
}
