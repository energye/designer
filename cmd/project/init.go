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

package project

import (
	"context"
	"runtime"
	"strings"

	"github.com/energye/designer/cmd/cef"
	"github.com/energye/designer/cmd/clui"
	"github.com/energye/designer/cmd/dflag"
	"github.com/energye/designer/options/bean"
)

func RunInit(args dflag.Args, defaultPath string) {
	console := clui.New()
	name := argValue(args, "name")
	if name == "" {
		var err error
		name, err = console.Input("Project name", "")
		if err != nil {
			console.Error(err.Error())
			return
		}
	}
	path := argValue(args, "path")
	if path == "" {
		path = defaultPath
	}
	uiName := strings.ToUpper(argValue(args, "ui"))
	if uiName == "" {
		uiName = strings.ToUpper(argValue(args, "framework"))
	}
	if uiName == "" {
		items := []string{bean.GUIRenderFramework_LCL, bean.GUIRenderFramework_WV, bean.GUIRenderFramework_CEF}
		idx, err := console.Select("Select UI", items, 0)
		if err != nil {
			console.Error(err.Error())
			return
		}
		uiName = items[idx]
	}
	framework := bean.GUIRenderFramework(uiName)

	frameworkVersion := ""
	if framework == bean.GUIRenderFramework_CEF {
		version := argValue(args, "cef-version")
		if version == "" {
			versions := cef.Versions()
			idx, err := console.Select("Select CEF version", versions, 0)
			if err != nil {
				console.Error(err.Error())
				return
			}
			version = versions[idx]
		}
		osName := argValue(args, "os")
		if osName == "" {
			osName = runtime.GOOS
		}
		arch := argValue(args, "arch")
		if arch == "" {
			arch = runtime.GOARCH
		}
		cefDir := argValue(args, "cef-dir")
		console.Info("Preparing CEF:", version, osName, arch)
		progressSink := clui.NewProgressSink(console)
		result, err := cef.EnsureInstalled(context.Background(), cef.InstallOptions{
			Dir:     cefDir,
			Version: version,
			OS:      osName,
			Arch:    arch,
			OnProgress: func(progress cef.Progress) {
				progressSink.Update(progress.Message, progress.Current, progress.Total)
			},
		})
		progressSink.Finish()
		if err != nil {
			console.Error(err.Error())
			return
		}
		frameworkVersion = result.OSArchVersion
	}

	result, err := Create(CreateOptions{
		Name:               name,
		Dir:                path,
		GUIRenderFramework: framework,
		FrameworkVersion:   frameworkVersion,
		OnConflict: func(conflict Conflict) ConflictDecision {
			ok, _ := console.Confirm(conflict.Message, true)
			if ok {
				return ConflictOverwrite
			}
			return ConflictCancel
		},
	})
	if err != nil {
		console.Error(err.Error())
		return
	}
	console.Success("Project created successfully:", result.Project.Name, "->", result.Dir)
	console.Info("Project config:", result.EGPPath)
	if err = UpdateGoModDependencies(context.Background(), GoModUpdateOptions{
		Dir: result.Dir,
		OnOutput: func(message string) {
			console.Info(message)
		},
	}); err != nil {
		console.Error("go mod update failed:", err.Error())
	}
}

func argValue(args dflag.Args, names ...string) string {
	for _, name := range names {
		if args.Contains(name) {
			return strings.TrimSpace(args.Get(name))
		}
	}
	return ""
}
