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
	"embed"
	"fmt"
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/cmd/dflag"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/icns"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/pkgs/win"
	"github.com/energye/lcl/tool/command"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	AppContents           = "Contents"
	AppContentsFrameworks = "Frameworks"
	AppContentsMacOS      = "MacOS"
	AppContentsResources  = "Resources"
)

//go:embed dmg
var images embed.FS

type TMacAssociateFiles struct {
	CFBundleTypeExtensions string // 扩展名 不带点
	CFBundleTypeName       string // 名称
	CFBundleTypeRole       string // 角色 Editor/Viewer
	LSHandlerRank          string // 优先级 Owner/Default
	CFBundleTypeIconFile   string // 图标 png/.icns
	CFBundleTypeMimeType   string // MIME(允许空)
}

func (m *TMacAssociateFiles) IsEmpty() bool {
	return m.CFBundleTypeExtensions == "" || m.CFBundleTypeName == "" ||
		m.CFBundleTypeRole == ""
}

type TMacAssociateProtocols struct {
	Scheme      string // 协议名 如 myapp, 调用格式：myapp://xxx
	Description string // 描述 系统显示的协议名称
}

func (m *TMacAssociateProtocols) IsEmpty() bool {
	return m.Scheme == "" || m.Description == ""
}

func (m *Package) platformPackage() {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build - project GProject is nil")
		return
	}
	event.ConsoleWriteInfo("CMD-package-run", "GOOS:", lib.GOOS(), "GOARCH:", lib.GOARCH())
	if !build.Run() {
		return
	}
	event.ConsoleWriteInfo("CMD-package-run")
	if m.packager() {
		event.ConsoleWriteInfo("Package Successfully")
	}
}

func (m *Package) packager() bool {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Package - GProject is nil")
		return false
	}
	event.ConsoleWriteInfo("Package - project check config options")
	option := proj.BuildOption
	if !m.createAppBundle() {
		return false
	}
	if option.MacSign.Enable {
		if !cert() {
			return false
		}
	} else {
		event.ConsoleWriteWarn("Package - Not Enabled cert")
	}

	if option.MacPKG {
		if !m.pkg() {
			return false
		}
	} else {
		event.ConsoleWriteWarn("Package - Not Enabled PKG")
	}
	if option.MacDMG {
		if !m.dmg() {
			return false
		}
	} else {
		event.ConsoleWriteWarn("Package - Not Enabled DMG")
	}
	return true
}

func (m *Package) createAppBundle() bool {
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
	event.ConsoleWriteInfo("App Bundle", bean.GProject.Name, "END")
	return true
}

func (m *Package) pkg() bool {
	event.ConsoleWriteInfo("Package - PKG")
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	app := appRootDir()
	appName := option.PackageName + ".app"
	pkgName := option.PackageName
	if m.AppendPlatform {
		pkgName += "_" + lib.GOOS()
	}
	if proj.BuildOption.MacCommonLib {
		pkgName += "_universal"
	} else {
		if m.AppendArch {
			pkgName += "_" + lib.GOARCH()
		}
	}
	if m.CustomSuffix != "" {
		pkgName += "_" + m.CustomSuffix
	}
	pkgName += ".pkg"

	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.HideWindow = true
	cmd.Dir = output
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.Command("rm", "-rf", pkgName)
	// pkgbuild --root demo.app --identifier com.demo.demo --version 1.0.0 --install-location /Applications/demo.app demo.pkg
	args := []string{"--root", app,
		"--identifier", proj.AppOption.MacOS.PList.CFBundleIdentifier,
		"--version", proj.AppOption.MacOS.PList.CFBundleVersion,
		"--install-location", fmt.Sprintf("/Applications/%s", appName), pkgName}
	cmd.Command("pkgbuild", args...)
	event.ConsoleWriteInfo("Package - PKG END")
	return true
}

func (m *Package) dmg() bool {
	event.ConsoleWriteInfo("Package - DMG")

	whichCmd := exec.Command("which", "create-dmg")
	err := whichCmd.Run()
	if err != nil {
		msg := "create-dmg command was not found.\nUse 'brew install create-dmg' Install"
		event.ConsoleWriteError("Package - DMG", err.Error(), "\n", msg)
		if lcl.Application != nil {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				text := api.PasStr(msg)
				title := api.PasStr("ENERGY Designer")
				lcl.Application.MessageBox(text, title, win.MB_OK+win.MB_ICONINFORMATION)
			})
		}
		return false
	}

	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	appName := option.PackageName + ".app"
	dmgName := option.PackageName
	if m.AppendPlatform {
		dmgName += "_" + lib.GOOS()
	}
	if proj.BuildOption.MacCommonLib {
		dmgName += "_universal"
	} else {
		if m.AppendArch {
			dmgName += "_" + lib.GOARCH()
		}
	}
	if m.CustomSuffix != "" {
		dmgName += "_" + m.CustomSuffix
	}
	dmgName += ".dmg"

	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.HideWindow = true
	cmd.Dir = output
	cmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	cmd.Command("rm", "-rf", dmgName)
	args := []string{
		"--volname", option.PackageName,
		"--window-pos", "200", "120",
		"--window-size", "500", "350",
		"--icon-size", "80",
		"--icon", appName, "100", "130",
		"--app-drop-link", "350", "130",
		//"--background", "bg.png",
		dmgName,
		appName,
	}
	cmd.Command("create-dmg", args...)
	event.ConsoleWriteInfo("Package - DMG END")
	return true
}

// APP_NAME="MyApp.app"
// CERT="Developer ID Application: 你的名字 (团队ID)"
// echo "签名 dylib..."
// codesign -f -s "$CERT" "$APP_NAME/Contents/Frameworks/"your.dylib
// echo "签名主程序..."
// codesign -f -s "$CERT" --options runtime "$APP_NAME"
// echo "开发自签名"
// codesign --force --deep --sign - "$APP_NAME"
// echo "验证签名..."
// codesign -vvv --deep --strict "$APP_NAME"
func cert() bool {
	event.ConsoleWriteInfo("Package - Cert")
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	certCommandList := option.MacSign.Cert
	if len(certCommandList) > 0 {
		appName := option.PackageName + ".app"
		runtimeFile := lib.GetDLLName()
		if proj.BuildOption.MacCommonLib {
			runtimeFile = lib.DarwinUniversalBinaryName
		}
		cmd := command.NewCMD()
		cmd.IsPrint = false
		cmd.HideWindow = true
		cmd.Dir = output
		cmd.Console = func(data string, level command.Level) {
			event.ConsoleWriteInfo(data)
		}
		for _, codesignCMD := range certCommandList {
			// 必须是 codesign 命令
			if strings.Index(codesignCMD, "codesign") == 0 {
				// 替换变量 $APP_NAME $ENERGY.DYLIB
				codesignCMD = strings.ReplaceAll(codesignCMD, "$APP_NAME", appName)
				codesignCMD = strings.ReplaceAll(codesignCMD, "$ENERGY.DYLIB", runtimeFile)
				args := dflag.ParseCommandLine(codesignCMD)
				if len(args) > 1 {
					event.ConsoleWriteInfo(codesignCMD)
					cmd.Command("codesign", args[1:]...)
				}
			}
		}
	}
	event.ConsoleWriteInfo("Package - Cert END")
	return true
}

func appRootDir() string {
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	appName := option.PackageName + ".app"
	appRoot := filepath.Join(output, appName)
	return appRoot
}

func createApp() bool {
	appRoot := appRootDir()
	event.ConsoleWriteInfo("App Bundle - Create App Dir", appRoot)
	err := os.RemoveAll(appRoot)
	if err != nil {
		event.ConsoleWriteError("Package - Remove App:", err.Error())
		return false
	}
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
	event.ConsoleWriteInfo("App Bundle - Create END")
	return true
}

func copyAppInfoPList() bool {
	event.ConsoleWriteInfo("App Bundle - Copy App Info.plist")

	appRoot := appRootDir()
	appContentsPath := filepath.Join(appRoot, AppContents)
	appResourcesPath := filepath.Join(appContentsPath, AppContentsResources)

	pListInfoTemplate := app.Packager("darwin/Info.plist")
	if pListInfoTemplate == nil {
		event.ConsoleWriteError("macOS 应用配置-保存配置 info.plist 模板获取失败, 模板内容为 nil")
		return false
	}
	associateFiles := parserAssociateFile(bean.GProject.BuildOption.MacAssociateFileList)
	associateProtocols := parserAssociateProtocol(bean.GProject.BuildOption.MacAssociateProtocolList)
	// 处理 associateFiles.CFBundleTypeIconFile 图标处理
	// 如果配置 png 转 icns 复制到 xxx.app/Contents/Resources/ 目录
	for _, associateFile := range associateFiles {
		iconFile := associateFile.CFBundleTypeIconFile
		associateFile.CFBundleTypeIconFile = "" // 设置为空, 在下面赋值
		ext := filepath.Ext(iconFile)
		if ext != ".png" && ext != ".icns" || !tool.IsExist(iconFile) {
			continue
		}
		// 为了保证唯一，重新命名
		newIcnsName := "AssociateFile_" + associateFile.CFBundleTypeExtensions + ".icns"
		outIcnsFilePath := filepath.Join(appResourcesPath, newIcnsName)
		if ext == ".png" {
			// png 转 icns, 并保存到 xxx.app/Contents/Resources/xxx.icns
			if err := icns.PngToIcns(iconFile, outIcnsFilePath); err == nil {
				associateFile.CFBundleTypeIconFile = newIcnsName
			}
		} else {
			// 复制 icns, 保存到 xxx.app/Contents/Resources/xxx.icns
			if err := tool.CopyFile(iconFile, outIcnsFilePath); err == nil {
				associateFile.CFBundleTypeIconFile = newIcnsName
			}
		}
	}
	data := make(map[string]any)
	data["PList"] = bean.GProject.AppOption.MacOS.PList
	data["AssociateFiles"] = associateFiles
	data["AssociateProtocols"] = associateProtocols
	pListInfoData, err := tool.RenderTemplate(string(pListInfoTemplate), data)
	if err != nil {
		event.ConsoleWriteError("macOS 应用配置-保存配置 info.plist 内容渲染失败:", err.Error())
		return false
	}
	// copy to app contents
	copyInfoPlistPath := filepath.Join(appContentsPath, "Info.plist")
	err = os.WriteFile(copyInfoPlistPath, pListInfoData, 0644)
	if err != nil {
		event.ConsoleWriteError("App Bundle - Copy App Info.plist WriteFile:", err.Error())
		return false
	}
	event.ConsoleWriteInfo("App Bundle - Copy App Info.plist END")
	return true
}

func copyAppPkgInfo() bool {
	pkgInfo := []byte{0x41, 0x50, 0x50, 0x4C, 0x3F, 0x3F, 0x3F, 0x3F, 0x0D, 0x0A}
	event.ConsoleWriteInfo("App Bundle - Copy App PkgInfo")
	appRoot := appRootDir()
	contents := filepath.Join(appRoot, AppContents)
	copyPkgInfoPath := filepath.Join(contents, "PkgInfo")
	if !tool.IsExist(copyPkgInfoPath) {
		err := os.WriteFile(copyPkgInfoPath, pkgInfo, 0644)
		if err != nil {
			event.ConsoleWriteError("Package - Copy App PkgInfo WriteFile:", err.Error())
			return false
		}
	}
	event.ConsoleWriteInfo("App Bundle - Copy App PkgInfo END")
	return true
}

func copyFiles() bool {
	event.ConsoleWriteInfo("App Bundle - Copy App Files")
	proj := bean.GProject
	appRoot := appRootDir()
	frameworksRuntime := config.Config.FrameworkRuntimePath()
	// libenergy-[arch].dylib
	srcRuntimeFilePath := filepath.Join(frameworksRuntime, lib.GetDLLName())
	if proj.BuildOption.MacCommonLib {
		srcRuntimeFilePath = filepath.Join(frameworksRuntime, lib.DarwinUniversalBinaryName)
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

	outputRuntimeFilePath := ""
	outputCurrRuntimeFilePath := filepath.Join(contentsFrameworks, lib.GetDLLName())
	outputUnivRuntimeFilePath := filepath.Join(contentsFrameworks, lib.DarwinUniversalBinaryName)
	if proj.BuildOption.MacCommonLib {
		// 通用二进制
		outputRuntimeFilePath = outputUnivRuntimeFilePath
		_ = os.Remove(outputCurrRuntimeFilePath)
	} else {
		outputRuntimeFilePath = outputCurrRuntimeFilePath
		_ = os.Remove(outputUnivRuntimeFilePath)
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

	event.ConsoleWriteInfo("App Bundle - Copy App Files. END")
	return true
}

func parserAssociateFile(associateFileList []string) (associateFiles []*TMacAssociateFiles) {
	assetsPath := bean.ResourceAssetsPath()
	for _, line := range associateFileList {
		associates := strings.Split(line, "|")
		if len(associates) >= 6 {
			icon := strings.TrimSpace(associates[4])
			if !filepath.IsAbs(icon) {
				icon = filepath.Join(assetsPath, icon)
			}
			associate := &TMacAssociateFiles{
				CFBundleTypeExtensions: strings.TrimSpace(associates[0]),
				CFBundleTypeName:       strings.TrimSpace(associates[1]),
				CFBundleTypeRole:       strings.TrimSpace(associates[2]),
				LSHandlerRank:          strings.TrimSpace(associates[3]),
				CFBundleTypeIconFile:   icon,
				CFBundleTypeMimeType:   strings.TrimSpace(associates[5]),
			}
			if !associate.IsEmpty() {
				associateFiles = append(associateFiles, associate)
			}
		}
	}
	return
}

func parserAssociateProtocol(associateProtocolList []string) (associateFiles []TMacAssociateProtocols) {
	for _, line := range associateProtocolList {
		associates := strings.Split(line, "|")
		if len(associates) >= 2 {
			associate := TMacAssociateProtocols{
				Scheme:      strings.TrimSpace(associates[0]),
				Description: strings.TrimSpace(associates[1]),
			}
			if !associate.IsEmpty() {
				associateFiles = append(associateFiles, associate)
			}
		}
	}
	return
}
