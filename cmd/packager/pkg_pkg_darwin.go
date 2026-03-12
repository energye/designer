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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api/libname"
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
	event.ConsoleWriteInfo("Package - ProjectPath", bean.GPath)
	event.ConsoleWriteInfo("Package - ProjectName", bean.GProject.Name)
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	outputFilename := filepath.Join(output, option.BuildFileName)
	event.ConsoleWriteInfo("Package - GUI-Render:", proj.GUIRenderFramework, "Output:", outputFilename)
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
	event.ConsoleWriteInfo("Package - ProjectName", bean.GProject.Name, "Success")
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
	event.ConsoleWriteInfo("Package - Copy App Info.plist")
	resourcesPath := filepath.Join(bean.GPath, "resources", "metadata", "Info.plist")
	if !tool.IsExist(resourcesPath) {
		event.ConsoleWriteError("Package - Info.plist not exist", resourcesPath)
		return false
	}
	infoPlistData, err := os.ReadFile(resourcesPath)
	if err != nil {
		event.ConsoleWriteError("Package - Copy App Info.plist ReadFile:", err.Error())
		return false
	}
	appRoot := appRootDir()
	contents := filepath.Join(appRoot, appContents)
	copyInfoPlistPath := filepath.Join(contents, "Info.plist")
	err = os.WriteFile(copyInfoPlistPath, infoPlistData, 0666)
	if err != nil {
		event.ConsoleWriteError("Package - Copy App Info.plist WriteFile:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("Package - Copy App Info.plist success")
	return true
}

func copyAppPkgInfo() bool {
	pkgInfo := []byte{0x41, 0x50, 0x50, 0x4C, 0x3F, 0x3F, 0x3F, 0x3F, 0x0D, 0x0A}
	event.ConsoleWriteInfo("Package - Copy App PkgInfo")
	appRoot := appRootDir()
	contents := filepath.Join(appRoot, appContents)
	copyPkgInfoPath := filepath.Join(contents, "PkgInfo")
	err := os.WriteFile(copyPkgInfoPath, pkgInfo, 0666)
	if err != nil {
		event.ConsoleWriteError("Package - Copy App PkgInfo WriteFile:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("Package - Copy App PkgInfo success")
	return true
}

func copyFiles() bool {
	event.ConsoleWriteInfo("Package - Copy App Files")
	//libDll := lib.Libs().Get(libname.GetDLLName())
	cfg := config.Config
	proj := bean.GProject
	appRoot := appRootDir()
	frameworksRuntime := filepath.Join(cfg.FrameworkDir, "runtime")
	// libenergy-xxx-xxx-xxx.dylib
	srcRuntimeFilePath := filepath.Join(frameworksRuntime, libname.GetDLLName())
	if proj.BuildOption.MacCommonLib {
		srcRuntimeFilePath = filepath.Join(frameworksRuntime, libname.DarwinUniversalBinaryName)
	}
	if !tool.IsExist(srcRuntimeFilePath) {
		event.ConsoleWriteError("Package - energy runtime lib not exist.", srcRuntimeFilePath)
		return false
	}
	// app binary
	output := proj.BuildOption.Output
	if !filepath.IsAbs(proj.BuildOption.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	srcAppBinary := filepath.Join(output, proj.BuildOption.BuildFileName)
	if !tool.IsExist(srcAppBinary) {
		event.ConsoleWriteError("Package - app binary not exist.", srcAppBinary)
		return false
	}
	// icon.icns
	srcAppIconIcns := filepath.Join(bean.ResourceEmbedPath(), "icon.icns")
	if !tool.IsExist(srcAppIconIcns) {
		event.ConsoleWriteError("Package - app icon.icns not exist.", srcAppIconIcns)
		return false
	}

	contents := filepath.Join(appRoot, appContents)
	contentsFrameworks := filepath.Join(contents, appContentsFrameworks)
	contentsMacOS := filepath.Join(contents, appContentsMacOS)
	contentsResources := filepath.Join(contents, appContentsResources)

	outputRuntimeFilePath := filepath.Join(contentsFrameworks, libname.GetDLLName())
	if proj.BuildOption.MacCommonLib {
		outputRuntimeFilePath = filepath.Join(contentsFrameworks, libname.DarwinUniversalBinaryName)
	}
	event.ConsoleWriteInfo("Package - Copy Frameworks")
	event.ConsoleWriteInfo("Package - Copy Runtime lib", outputRuntimeFilePath)
	if err := tool.CopyFile(srcRuntimeFilePath, outputRuntimeFilePath); err != nil {
		event.ConsoleWriteError(err.Error())
		return false
	}

	outputAppBinary := filepath.Join(contentsMacOS, proj.AppOption.MacOS.PList.CFBundleExecutable)
	event.ConsoleWriteInfo("Package - Copy MacOS")
	event.ConsoleWriteInfo("Package - Copy App binary", outputAppBinary)
	if err := tool.CopyFile(srcAppBinary, outputAppBinary); err != nil {
		event.ConsoleWriteError(err.Error())
		return false
	}

	outputAppIcns := filepath.Join(contentsResources, proj.AppOption.MacOS.PList.CFBundleIconFile)
	event.ConsoleWriteInfo("Package - Copy Resources")
	event.ConsoleWriteInfo("Package - Copy icns", outputAppIcns)
	if err := tool.CopyFile(srcAppIconIcns, outputAppIcns); err != nil {
		event.ConsoleWriteError(err.Error())
		return false
	}
	event.ConsoleWriteInfo("Package - Copy App Files. Success")
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
