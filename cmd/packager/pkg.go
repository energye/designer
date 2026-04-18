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
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/lcl/tool/command"
	"os"
	"text/template"
)

// Run 执行打包命令的入口函数
func Run() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - project GProject is nil")
		return false
	}

	// 构建环境 arch
	defaultGOARCH := os.Getenv("GOARCH")
	defer os.Setenv("GOARCH", defaultGOARCH)

	// 平台打包
	platformPackage()

	// 构建全平台
	// 如果禁用CGO并启用了其它平台构建
	build.All()

	return true
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

func RunCMD(dir, name string, args ...string) error {
	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.HideWindow = true
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.Command(name, args...)
	return nil
}

func WriteFile(filePath string, data []byte) error {
	_ = os.Remove(filePath)
	err := os.WriteFile(filePath, data, 0644)
	return err
}
