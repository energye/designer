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
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
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
	if !createApp() {
		return
	}
	if !copyAppInfoPList() {
		return
	}
	if !copyAppPkgInfo() {
		return
	}
	if !copyFiles() {
		return
	}
	switch proj.GUIRenderFramework {
	case bean.GUIRenderFramework_LCL, bean.GUIRenderFramework_WV:
	case bean.GUIRenderFramework_CEF:
	}

}

func dmg() {

}

func appRootDir() string {
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	packageName := option.PackageName + ".app"
	appRoot := filepath.Join(output, packageName)
	return appRoot
}

func createApp() bool {
	appRoot := appRootDir()
	event.ConsoleWriteInfo("Package - Create App Dir", appRoot)
	err := os.RemoveAll(appRoot)
	if err != nil {
		event.ConsoleWriteError("Package - remove app root:", err.Error())
		return false
	}
	// Contents
	contents := filepath.Join(appRoot, appContents)
	if err := os.MkdirAll(contents, 0755); err != nil {
		event.ConsoleWriteError("Package - unable to create directory:", err.Error())
		return false
	}
	// Contents/Frameworks
	contentsFrameworks := filepath.Join(contents, appContentsFrameworks)
	if err := os.MkdirAll(contentsFrameworks, 0755); err != nil {
		event.ConsoleWriteError("Package - unable to create directory:", err.Error())
		return false
	}
	// Contents/MacOS
	contentsMacOS := filepath.Join(contents, appContentsMacOS)
	if err := os.MkdirAll(contentsMacOS, 0755); err != nil {
		event.ConsoleWriteError("Package - unable to create directory:", err.Error())
		return false
	}
	// Contents/Resources
	contentsResources := filepath.Join(contents, appContentsResources)
	if err := os.MkdirAll(contentsResources, 0755); err != nil {
		event.ConsoleWriteError("Package - unable to create directory:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("Package - Create App Dir success")
	return true
}

func copyAppInfoPList() bool {
	event.ConsoleWriteInfo("Package - Copy app Info.plist")
	resourcesPath := filepath.Join(bean.GPath, "resources", "metadata", "Info.plist")
	if !tool.IsExist(resourcesPath) {
		event.ConsoleWriteError("Package - Info.plist not exist", resourcesPath)
		return false
	}
	infoPlistData, err := os.ReadFile(resourcesPath)
	if err != nil {
		event.ConsoleWriteError("Package - Copy app Info.plist ReadFile:", err.Error())
		return false
	}
	appRoot := appRootDir()
	copyInfoPlistPath := filepath.Join(appRoot, "Info.plist")
	err = os.WriteFile(copyInfoPlistPath, infoPlistData, 0666)
	if err != nil {
		event.ConsoleWriteError("Package - Copy app Info.plist WriteFile:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("Package - Copy app Info.plist success")
	return true
}

func copyAppPkgInfo() bool {
	pkgInfo := []byte{0x41, 0x50, 0x50, 0x4C, 0x3F, 0x3F, 0x3F, 0x3F, 0x0D, 0x0A}
	event.ConsoleWriteInfo("Package - Copy app PkgInfo")
	appRoot := appRootDir()
	copyPkgInfoPath := filepath.Join(appRoot, "PkgInfo")
	err := os.WriteFile(copyPkgInfoPath, pkgInfo, 0666)
	if err != nil {
		event.ConsoleWriteError("Package - Copy app PkgInfo WriteFile:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("Package - Copy app PkgInfo success")
	return true
}

func copyFiles() bool {
	event.ConsoleWriteInfo("Package - Copy app file")
	//libDll := lib.Libs().Get(libname.GetDLLName())

	//appRoot := appRootDir()
	event.ConsoleWriteInfo("Package - Copy frameworks")

	return true
}

// pkgbuild --root demo.app --identifier com.demo.demo --version 1.0.0 --install-location /Applications/demo.app demo.pkg
//func pkgbuild(c *command.Config, proj *project.Project, appRoot string) error {
//	proj.AppType = project.AtApp
//	proj.ProjectPath = projectPath
//	buildOutDir := assets.BuildOutPath(proj)
//	cmdWorkDir := filepath.Join(buildOutDir, "darwin")
//	term.Logger.Info("Generate app pkgbuild", term.Logger.Args("cmd work dir", cmdWorkDir))
//	// remove xxx.pkg
//	os.Remove(filepath.Join(cmdWorkDir, fmt.Sprintf("%s.pkg", getAppName(c, proj))))
//	cmd := cmd.NewCMD()
//	//cmd.IsPrint = false
//	cmd.Dir = cmdWorkDir
//	cmd.MessageCallback = func(bytes []byte, err error) {
//		msg := string(bytes)
//		if msg != "" {
//			println(msg)
//		}
//	}
//	app := fmt.Sprintf("%s.app", getAppName(c, proj))
//	pkg := fmt.Sprintf("%s.pkg", getAppName(c, proj))
//	var args = []string{"--root", app,
//		"--identifier", proj.PList.BundleIdentifier,
//		"--version", proj.PList.BundleVersion,
//		"--install-location", fmt.Sprintf("/Applications/%s", app), pkg}
//	cmd.Command("pkgbuild", args...)
//	cmd.Close()
//	// remove xxx.app
//	os.RemoveAll(filepath.Join(cmdWorkDir, app))
//	return nil
//}
