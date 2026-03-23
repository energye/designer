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
	_ "github.com/energye/designer/resources"
	"github.com/energye/designer/resources/frameworks"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"os"
	"runtime"
	"strings"
)

// go build -ldflags="-H windowsgui -s -w" -trimpath -o build/designer.exe
// go build -ldflags="-H windowsgui -s -w -buildid=" -trimpath -o build/designer.exe
// go build -ldflags="-H windowsgui" -trimpath -o build/designer.exe
// go build --tags liball -ldflags="-H windowsgui" -trimpath -o build/designer.exe
func main() {
	//go tool pprof http://localhost:8080/debug/pprof/profile?seconds=15
	//go http.ListenAndServe(":8080", nil)
	logs.Level = logs.LevelDebug
	//logs.Level = logs.LevelInfo
	//logs.Level = logs.LevelError
	libname.LibName = frameworks.ExtractLibrary()
	lcl.Init(nil, nil)
	logs.Debug(strings.Join(os.Args, " "))
	// 运行设计器
	designer.Run()
}

func init() {
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
		// os.Setenv("--universal-binary", "universal")
	}
}
