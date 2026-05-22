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
	"github.com/energye/designer/codegen"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	_ "github.com/energye/designer/internal"
	"github.com/energye/designer/options"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/preview"
	_ "github.com/energye/designer/resources"
	"github.com/energye/designer/resources/frameworks"
	"github.com/energye/designer/uigen"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"os"
	"runtime"
	"strings"
)

// Linux 环境: sudo apt install gdb make binutils build-essential libx11-dev libgtk2.0-dev libgdk-pixbuf2.0-dev libcairo2-dev libpango1.0-dev

// go build -ldflags="-H windowsgui -s -w" -trimpath -o build/designer.exe
// go build -ldflags="-H windowsgui -s -w -buildid=" -trimpath -o build/designer.exe
// go build -ldflags="-H windowsgui" -trimpath -o build/designer.exe
// go build --tags liball -ldflags="-H windowsgui" -trimpath -o build/designer.exe
// =================================================================================================
// go mod
// go env -w GOPROXY=https://goproxy.cn,direct
// go env -w GOSUMDB=sum.golang.google.cn
// go env -w GONOSUMDB=github.com/energye/*
//
// go get github.com/energye/lcl@v1.0.4
// go get github.com/energye/wv@v1.0.5
// go mod tidy
// =================================================================================================
func main() {
	api.SetDebug(true)
	//go tool pprof http://localhost:8080/debug/pprof/profile?seconds=15
	//go http.ListenAndServe(":8080", nil)
	logs.Level = logs.LevelDebug
	//logs.Level = logs.LevelInfo
	//logs.Level = logs.LevelError
	libname.LibName = frameworks.ExtractLibrary()
	lcl.Init()
	logs.Debug(strings.Join(os.Args, " "))
	// 运行设计器
	designer.Run()
}

func init() {
	if tool.IsLinux {
		libname.UseWS = "gtk3"
	}
	if tool.IsDarwin {
		// libname.EnableUniversalBinary = true
	}
	// 初始化Go环境变量, macOS linux
	config.InitGoEnv()
	// 初始化事件系统
	event.Init()
	// 初始化预览功能
	preview.Init()
	// 初始化项目配置和选项管理
	options.Init()
	// 初始化 UI 代码生成触发器
	uigen.Init()
	// 初始化代码生成引擎
	codegen.Init()
	// 初始化设计器核心功能
	designer.Init()

	if runtime.GOOS == "darwin" {
		isAmd64 := runtime.GOARCH == "amd64"
		isArm64 := runtime.GOARCH == "arm64"
		if isAmd64 {
			os.Setenv("MACOSX_DEPLOYMENT_TARGET", "10.15")
			os.Setenv("CGO_CFLAGS", "-mmacosx-version-min=10.15")
			os.Setenv("CGO_LDFLAGS", "-mmacosx-version-min=10.15")
		}
		if isArm64 {
			os.Setenv("MACOSX_DEPLOYMENT_TARGET", "11.0")
			os.Setenv("CGO_CFLAGS", "-mmacosx-version-min=11.0")
			os.Setenv("CGO_LDFLAGS", "-mmacosx-version-min=11.0")
		}
	}
}
