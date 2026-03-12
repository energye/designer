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
	"io/fs"
	"os"
	"path/filepath"
)

const (
	AppContents           = "Contents"
	AppContentsFrameworks = "Frameworks"
	AppContentsMacOS      = "MacOS"
	AppContentsResources  = "Resources"
)

func packager() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - GProject is nil")
		return false
	}
	event.ConsoleWriteInfo("Package - project check config options")
	option := proj.BuildOption
	if !option.MacDMG && !option.MacPKG {
		event.ConsoleWriteWarn("Package - project not package format DMG PKG")
		return false
	}
	if !createAppBundle() {
		return false
	}
	if option.MacPKG {
		if !pkg() {
			return false
		}
	}
	if option.MacDMG {
		if !dmg() {
			return false
		}
	}
	return true
}

func createAppBundle() bool {
	event.ConsoleWriteInfo("App Bundle", bean.GPath)
	event.ConsoleWriteInfo("App Bundle", bean.GProject.Name)
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	outputFilename := filepath.Join(output, option.BuildFileName)
	event.ConsoleWriteInfo("App Bundle - GUI-Render:", proj.GUIRenderFramework, "Output:", outputFilename)
	if !createApp() {
		return false
	}
	if !copyAppInfoPList() {
		return false
	}
	if !copyAppPkgInfo() {
		return false
	}
	if !copyFiles() {
		return false
	}
	event.ConsoleWriteInfo("App Bundle", bean.GProject.Name, "Success")
	return true
}

func pkg() bool {
	event.ConsoleWriteInfo("Package - PKG")
	return true
}

func dmg() bool {
	event.ConsoleWriteInfo("Package - DMG")
	return true
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
	event.ConsoleWriteInfo("App Bundle - Create App Dir", appRoot)
	//err := os.RemoveAll(appRoot)
	//if err != nil {
	//	event.ConsoleWriteError("Package - remove app root:", err.Error())
	//	return false
	//}
	// Contents
	contents := filepath.Join(appRoot, AppContents)
	if !tool.IsExist(contents) {
		if err := os.MkdirAll(contents, 0755); err != nil {
			event.ConsoleWriteError("App Bundle - unable to create directory:", err.Error())
			return false
		}
	}
	// Contents/Frameworks
	contentsFrameworks := filepath.Join(contents, AppContentsFrameworks)
	if !tool.IsExist(contentsFrameworks) {
		if err := os.MkdirAll(contentsFrameworks, 0755); err != nil {
			event.ConsoleWriteError("App Bundle - unable to create directory:", err.Error())
			return false
		}
	}
	// Contents/MacOS
	contentsMacOS := filepath.Join(contents, AppContentsMacOS)
	if !tool.IsExist(contentsMacOS) {
		if err := os.MkdirAll(contentsMacOS, 0755); err != nil {
			event.ConsoleWriteError("App Bundle - unable to create directory:", err.Error())
			return false
		}
	}
	// Contents/Resources
	contentsResources := filepath.Join(contents, AppContentsResources)
	if !tool.IsExist(contentsResources) {
		if err := os.MkdirAll(contentsResources, 0755); err != nil {
			event.ConsoleWriteError("App Bundle - unable to create directory:", err.Error())
			return false
		}
	}
	event.ConsoleWriteInfo("App Bundle - Create Success")
	return true
}

func copyAppInfoPList() bool {
	event.ConsoleWriteInfo("App Bundle - Copy App Info.plist")
	resourcesPath := filepath.Join(bean.GPath, "resources", "metadata", "Info.plist")
	if !tool.IsExist(resourcesPath) {
		event.ConsoleWriteError("App Bundle - Info.plist not exist", resourcesPath)
		return false
	}
	infoPlistData, err := os.ReadFile(resourcesPath)
	if err != nil {
		event.ConsoleWriteError("App Bundle - Copy App Info.plist ReadFile:", err.Error())
		return false
	}
	appRoot := appRootDir()
	contents := filepath.Join(appRoot, AppContents)
	copyInfoPlistPath := filepath.Join(contents, "Info.plist")
	err = os.WriteFile(copyInfoPlistPath, infoPlistData, 0666)
	if err != nil {
		event.ConsoleWriteError("App Bundle - Copy App Info.plist WriteFile:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("App Bundle - Copy App Info.plist Success")
	return true
}

func copyAppPkgInfo() bool {
	pkgInfo := []byte{0x41, 0x50, 0x50, 0x4C, 0x3F, 0x3F, 0x3F, 0x3F, 0x0D, 0x0A}
	event.ConsoleWriteInfo("App Bundle - Copy App PkgInfo")
	appRoot := appRootDir()
	contents := filepath.Join(appRoot, AppContents)
	copyPkgInfoPath := filepath.Join(contents, "PkgInfo")
	if !tool.IsExist(copyPkgInfoPath) {
		err := os.WriteFile(copyPkgInfoPath, pkgInfo, 0666)
		if err != nil {
			event.ConsoleWriteError("Package - Copy App PkgInfo WriteFile:", err.Error())
			return false
		}
	}
	event.ConsoleWriteInfo("App Bundle - Copy App PkgInfo success")
	return true
}

func copyFiles() bool {
	event.ConsoleWriteInfo("App Bundle - Copy App Files")
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
		event.ConsoleWriteError("App Bundle - energy runtime lib not exist.", srcRuntimeFilePath)
		return false
	}
	// app binary
	output := proj.BuildOption.Output
	if !filepath.IsAbs(proj.BuildOption.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	srcAppBinary := filepath.Join(output, proj.BuildOption.BuildFileName)
	if !tool.IsExist(srcAppBinary) {
		event.ConsoleWriteError("App Bundle - app binary not exist.", srcAppBinary)
		return false
	}
	// icon.icns
	srcAppIconIcns := filepath.Join(bean.ResourceEmbedPath(), "icon.icns")
	if !tool.IsExist(srcAppIconIcns) {
		event.ConsoleWriteError("App Bundle - app icon.icns not exist.", srcAppIconIcns)
		return false
	}

	contents := filepath.Join(appRoot, AppContents)
	contentsFrameworks := filepath.Join(contents, AppContentsFrameworks)
	contentsMacOS := filepath.Join(contents, AppContentsMacOS)
	contentsResources := filepath.Join(contents, AppContentsResources)

	outputRuntimeFilePath := filepath.Join(contentsFrameworks, libname.GetDLLName())
	if proj.BuildOption.MacCommonLib {
		outputRuntimeFilePath = filepath.Join(contentsFrameworks, libname.DarwinUniversalBinaryName)
	}
	event.ConsoleWriteInfo("App Bundle - Copy Frameworks")
	event.ConsoleWriteInfo("App Bundle - Copy Runtime lib", outputRuntimeFilePath)
	if err := tool.CopyFile(srcRuntimeFilePath, outputRuntimeFilePath); err != nil {
		event.ConsoleWriteError(err.Error())
		return false
	}

	outputAppBinary := filepath.Join(contentsMacOS, proj.AppOption.MacOS.PList.CFBundleExecutable)
	event.ConsoleWriteInfo("App Bundle - Copy MacOS")
	event.ConsoleWriteInfo("App Bundle - Copy App binary", outputAppBinary)
	if err := tool.CopyFile(srcAppBinary, outputAppBinary); err != nil {
		event.ConsoleWriteError(err.Error())
		return false
	}

	outputAppIcns := filepath.Join(contentsResources, proj.AppOption.MacOS.PList.CFBundleIconFile)
	event.ConsoleWriteInfo("App Bundle - Copy Resources")
	event.ConsoleWriteInfo("App Bundle - Copy icns", outputAppIcns)
	if err := tool.CopyFile(srcAppIconIcns, outputAppIcns); err != nil {
		event.ConsoleWriteError(err.Error())
		return false
	}

	// Contents/Resources/xxx.lproj
	event.ConsoleWriteInfo("App Bundle - Copy Localizations")
	srcResourceMetadataPath := bean.ResourceMetadataPath()
	for _, local := range bean.GProject.AppOption.MacOS.PList.CFBundleLocalizations {
		srcLocal := filepath.Join(srcResourceMetadataPath, local+".lproj")
		dstLocal := filepath.Join(contentsResources, local+".lproj")
		err := filepath.WalkDir(srcLocal, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			_, fileName := filepath.Split(path)
			dstPath := filepath.Join(dstLocal, fileName)
			return tool.CopyFile(path, dstPath)
		})
		if err != nil {
			event.ConsoleWriteError("App Bundle - Copy Localizations", err.Error())
		}
	}

	event.ConsoleWriteInfo("App Bundle - Copy App Files. Success")
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
