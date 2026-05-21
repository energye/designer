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
	"errors"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/resources/frameworks"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/tool/command"
	"os"
	"strings"
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
	var lastError string
	cmd.Console = func(data string, level command.Level) {
		if level == command.LError {
			lastError = data
			event.ConsoleWriteError(data)
		} else {
			event.ConsoleWriteInfo(data)
		}
	}
	cmd.Command(name, args...)
	if cmd.Cmd != nil && cmd.Cmd.ProcessState != nil && !cmd.Cmd.ProcessState.Success() {
		if lastError != "" {
			return errors.New(lastError)
		}
		return errors.New(name + " exited with non-zero code")
	}
	return nil
}

func WriteFile(filePath string, data []byte) error {
	_ = os.Remove(filePath)
	err := os.WriteFile(filePath, data, 0644)
	return err
}

// packageLibName 返回打包时使用的运行时库文件名
// Debug 模式保留架构前缀 (如 libenergy-amd64-gtk3.so)
// Release 模式去掉架构前缀 (如 libenergy-gtk3.so)，与 SONAME 一致
func packageLibName() string {
	libName := lib.GetDLLName()
	proj := bean.GProject
	if proj != nil && proj.BuildOption.BuildModeRelease {
		libName = sonameFromLib(libName)
	}
	return libName
}

// sonameFromLib 从库文件名中去掉架构前缀，得到 SONAME
// Linux:   "libenergy-amd64-gtk3.so" -> "libenergy-gtk3.so"     (3段，去掉中间的arch)
// Windows: "libenergy-amd64.dll"     -> "libenergy.dll"          (2段，去掉arch)
// macOS:   "libenergy-amd64.dylib"   -> "libenergy.dylib"        (2段，去掉arch)
func sonameFromLib(libName string) string {
	parts := strings.SplitN(libName, "-", 3)
	switch len(parts) {
	case 3:
		// Linux: libenergy-{arch}-{ws}.so -> libenergy-{ws}.so
		return parts[0] + "-" + parts[2]
	case 2:
		// Windows/macOS: libenergy-{arch}.dll -> libenergy.dll
		return parts[0] + "." + strings.SplitN(parts[1], ".", 2)[1]
	default:
		return libName
	}
}
