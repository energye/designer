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
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/energye/designer/cmd/cef"
	"github.com/energye/designer/cmd/dflag"
	"github.com/energye/designer/options/bean"
)

func RunInit(args dflag.Args, defaultPath string) {
	reader := bufio.NewReader(os.Stdin)
	name := argValue(args, "name")
	if name == "" {
		name = promptText(reader, "Project name")
	}
	path := argValue(args, "path")
	if path == "" {
		path = defaultPath
	}
	ui := strings.ToUpper(argValue(args, "ui"))
	if ui == "" {
		ui = strings.ToUpper(argValue(args, "framework"))
	}
	if ui == "" {
		ui = promptSelect(reader, "Select UI", []string{bean.GUIRenderFramework_LCL, bean.GUIRenderFramework_WV, bean.GUIRenderFramework_CEF})
	}
	framework := bean.GUIRenderFramework(ui)

	frameworkVersion := ""
	if framework == bean.GUIRenderFramework_CEF {
		version := argValue(args, "cef-version")
		if version == "" {
			version = promptSelect(reader, "Select CEF version", cef.Versions())
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
		fmt.Println("Preparing CEF:", version, osName, arch)
		result, err := cef.EnsureInstalled(context.Background(), cef.InstallOptions{
			Dir:     cefDir,
			Version: version,
			OS:      osName,
			Arch:    arch,
			OnProgress: func(progress cef.Progress) {
				if progress.Total > 0 {
					fmt.Printf("%s: %s (%d%%)\n", progress.Kind, progress.Message, progress.Current*100/progress.Total)
				} else if progress.Message != "" {
					fmt.Println(progress.Message)
				}
			},
		})
		if err != nil {
			fmt.Println("[ERROR]", err.Error())
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
			if askYesNo(reader, conflict.Message+" [y/N] ") {
				return ConflictOverwrite
			}
			return ConflictCancel
		},
	})
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return
	}
	fmt.Println("Project created successfully:", result.Project.Name, "->", result.Dir)
	fmt.Println("Project config:", result.EGPPath)
}

func argValue(args dflag.Args, names ...string) string {
	for _, name := range names {
		if args.Contains(name) {
			return strings.TrimSpace(args.Get(name))
		}
	}
	return ""
}

func promptText(reader *bufio.Reader, label string) string {
	for {
		fmt.Print(label + ": ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text != "" {
			return text
		}
	}
}

func promptSelect(reader *bufio.Reader, label string, items []string) string {
	for {
		fmt.Println(label + ":")
		for i, item := range items {
			fmt.Printf("  %d. %s\n", i+1, item)
		}
		fmt.Print("Select number: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		idx, err := strconv.Atoi(text)
		if err == nil && idx > 0 && idx <= len(items) {
			return items[idx-1]
		}
	}
}

func askYesNo(reader *bufio.Reader, label string) bool {
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes"
}
