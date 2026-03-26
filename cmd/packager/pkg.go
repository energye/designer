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

package packager

import (
	"bytes"
	"github.com/energye/designer/event"
	"text/template"
)

// Run 执行打包命令的入口函数
func Run() bool {
	event.ConsoleWriteInfo("CMD-package-run")
	return packager()
}

// AppBundle 创建 macOS 应用程序包
func AppBundle() bool {
	return createAppBundle()
}

func RenderTemplate(data any, templateText string) ([]byte, error) {
	tmpl, err := template.New("package-render-template").Parse(templateText)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err = tmpl.Execute(&out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
