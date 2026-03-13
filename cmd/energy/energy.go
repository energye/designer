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
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/cmd/dflag"
	"github.com/energye/designer/cmd/packager"
	"github.com/energye/designer/cmd/project"
	"github.com/energye/designer/cmd/run"
	"os"
)

/*
energy run -path /you/app/path
energy build -path /you/app/path
energy package -path /you/app/path
*/
func main() {
	projectPath := func(args dflag.Args) string {
		if args.Contains("path") {
			path := args.Get("path")
			return path
		}
		wd, _ := os.Getwd()
		return wd
	}
	cmd := dflag.New()
	cmd.Add(&dflag.Command{
		Name: "run",
		Long: "energy run, 运行应用",
		Run: func(args dflag.Args) {
			path := projectPath(args)
			project.LoadProject(path)
			if build.Run() {
				// 运行
				run.Run(nil)
			}
		},
	})
	cmd.Add(&dflag.Command{
		Name: "build",
		Long: "energy build, 构建应用二进制程序",
		Run: func(args dflag.Args) {
			path := projectPath(args)
			project.LoadProject(path)
			build.Run()
		},
	})
	cmd.Add(&dflag.Command{
		Name: "package",
		Long: "energy package, 制作应用安装包",
		Run: func(args dflag.Args) {
			path := projectPath(args)
			project.LoadProject(path)
			if !build.Run() {
				return
			}
			if !packager.Run() {
				return
			}
		},
	})
	cmd.Parse()
}
