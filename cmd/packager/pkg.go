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
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/resources/frameworks"
	"github.com/energye/lcl/tool/command"
	"os"
	"text/template"
)

type Package struct {
	AppendPlatform bool
	AppendArch     bool
	CustomSuffix   string
}

func Default() *Package {
	return &Package{}
}

// Run 执行打包命令的入口函数
func Run(pack *Package) bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - project GProject is nil")
		return false
	}
	frameworks.ExtractLibrary()
	if pack == nil {
		pack = Default()
	}
	// 平台打包
	pack.platformPackage()

	return true
}

// AppBundle 创建 macOS 应用程序包
func AppBundle(pack *Package) bool {
	if pack == nil {
		pack = Default()
	}
	return pack.createAppBundle()
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
		if level == command.LError {
			event.ConsoleWriteError(data)
		} else {
			event.ConsoleWriteInfo(data)
		}
	}
	cmd.Command(name, args...)
	return nil
}

func WriteFile(filePath string, data []byte) error {
	_ = os.Remove(filePath)
	err := os.WriteFile(filePath, data, 0644)
	return err
}
