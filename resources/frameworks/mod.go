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

import "github.com/energye/designer/pkg/tool"

// 模块本地依赖模板
//
//	 使用:
//		  在 cef, wv 库源码依赖配置
const modLocalTemplate = `module {{.Module}}

go 1.20

{{.Replace}}
`

// renderModLocalTemplate 渲染本地模块模板
func renderModLocalTemplate(module string, replace string) ([]byte, error) {
	return tool.RenderTemplate(modLocalTemplate, map[string]any{
		"Module":  module,
		"Replace": replace,
	})
}
