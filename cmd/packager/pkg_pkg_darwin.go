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

//go:build darwin

package packager

import (
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"os"
	"path/filepath"
)

const (
	appContents           = "Contents"
	appContentsFrameworks = "Frameworks"
	appContentsMacOS      = "MacOS"
	appContentsResources  = "Resources"
)

func packager() {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - GProject is nil")
		return
	}
	event.ConsoleWriteInfo("Package - project check config options")
	option := proj.BuildOption
	if !option.MacDMG && !option.MacPKG {
		event.ConsoleWriteWarn("Package - project not package format DMG PKG")
		return
	}
	if option.MacPKG {
		event.ConsoleWriteInfo("Package - PKG")
		pkg()
	}
	if option.MacDMG {
		event.ConsoleWriteInfo("Package - DMG")
		dmg()
	}

}

func pkg() {
	event.ConsoleWriteInfo("Package - ProjectPath:", bean.GPath)
	event.ConsoleWriteInfo("Package - ProjectName:", bean.GProject.Name)
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	outputFilename := filepath.Join(output, option.BuildFileName)
	event.ConsoleWriteInfo("Package - GUI-Render:", proj.GUIRenderFramework, outputFilename)
	switch proj.GUIRenderFramework {
	case bean.GUIRenderFramework_LCL, bean.GUIRenderFramework_WV:
	case bean.GUIRenderFramework_CEF:
	}

}

func dmg() {

}

func createApp() error {
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	appRoot := ""
	event.ConsoleWriteInfo("Package - Create App Dir: " + appRoot)
	buildOutDir := ""
	appDir := filepath.Join(buildOutDir, appRoot)
	os.Remove(appDir)
	// Contents
	contents := filepath.Join(appDir, appContents)
	if err := os.MkdirAll(contents, 0755); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}
	// Contents/Frameworks
	contentsFrameworks := filepath.Join(contents, appContentsFrameworks)
	if err := os.MkdirAll(contentsFrameworks, 0755); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}
	// Contents/MacOS
	contentsMacOS := filepath.Join(contents, appContentsMacOS)
	if err := os.MkdirAll(contentsMacOS, 0755); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}
	// Contents/Resources
	contentsResources := filepath.Join(contents, appContentsResources)
	if err := os.MkdirAll(contentsResources, 0755); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}
	return nil
}
