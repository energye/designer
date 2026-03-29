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

//go:build windows

package packager

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/winres"
	"github.com/energye/designer/resources/app"
	"github.com/energye/designer/resources/frameworks/lib"
	"path/filepath"
	"runtime"
	"strings"
)

// makensis.exe

const makensis = "makensis.exe"

// TWinAssociateFiles  关联文件
type TWinAssociateFiles struct {
	Ext         string // 要关联的文件后缀（不带点）
	FileClass   string // 注册表唯一类名（自定义，不能重复） 软件名+后缀+File
	Description string // 文件类型描述（鼠标悬浮时显示的文字）
	Icon        string // 文件显示的图标路径 icon.ico "$INSTDIR\\图标名.ico"
	SrcIcon     string // 文件显示的图标路径 icon.ico File "图标名.ico"
	CommandText string // 右键菜单显示的文字 "Open with EnergyTool"
}

func (m *TWinAssociateFiles) IsEmpty() bool {
	return m.Ext == "" || m.FileClass == "" || m.Description == "" || m.Icon == "" || m.CommandText == ""
}

// TWinAssociateProtocols  关联协议
type TWinAssociateProtocols struct {
	Scheme      string // 协议名 如 myapp, 调用格式：myapp://xxx
	Description string // 描述 系统显示的协议名称
}

func (m *TWinAssociateProtocols) IsEmpty() bool {
	return m.Scheme == "" || m.Description == ""
}

func packageNSIS() bool {
	if !checkToolCMD(makensis) {
		event.ConsoleWriteInfo("Package - check nsis Not Installed")
		return false
	}
	var (
		libEnergy   string
		libWebview2 string
	)
	switch runtime.GOARCH {
	case "amd64":
		libEnergy = "libenergy-amd64.dll"
		libWebview2 = "WebView2Loader-amd64.dll"
	case "386":
		libEnergy = "libenergy-386.dll"
		libWebview2 = "WebView2Loader-386.dll"
	case "arm64":
		event.ConsoleWriteInfo("Package - Currently, windows arm64 arch is not support")
		return false
	}
	proj := bean.GProject
	buildOption := proj.BuildOption
	appOption := proj.AppOption

	buildFileName := buildOption.BuildFileName
	if filepath.Ext(buildFileName) != ".exe" {
		buildFileName += ".exe"
	}
	packageName := buildOption.PackageName
	if filepath.Ext(packageName) != ".exe" {
		packageName += ".exe"
	}
	output := buildOption.Output
	if !filepath.IsAbs(buildOption.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	if proj.GUIRenderFramework == "WV" {
		webview2Setup := lib.Libs().Get(lib.PathWV2Setup)
		if webview2Setup == nil {
			event.ConsoleWriteInfo("Package - Failed to obtain MicrosoftEdgeWebview2Setup.exe")
			return false
		}
		err := webview2Setup.Release(output)
		if err != nil {
			event.ConsoleWriteError("Package - Release MicrosoftEdgeWebview2Setup.exe", err.Error())
			return false
		}
		defer func() {
			//_ = os.Remove(filepath.Join(output, "MicrosoftEdgeWebview2Setup.exe"))
		}()
	}

	var (
		appCompanyName = ""
		appProductName = ""
	)
	appID := appOption.Id // CompanyName.productName.AppName
	if ids := strings.Split(appID, "."); len(ids) >= 2 {
		appCompanyName = ids[0]
		appProductName = ids[1]
	}
	nsisExecLevel := ""
	runLevel := winres.ExecutionLevel(appOption.Windows.Manifest.RunLevel)
	switch runLevel {
	case winres.AsInvoker:
		nsisExecLevel = "user"
	case winres.HighestAvailable:
		nsisExecLevel = "highest"
	case winres.RequireAdministrator:
		nsisExecLevel = "admin"
	}

	installNsisScriptTemp := app.Packager("windows/install-nsis.nsi")
	installToolsScriptTemp := app.Packager("windows/install-tools.nsh")
	embedPath := bean.ResourceEmbedPath()
	iconIcoFilePath := filepath.Join(embedPath, "icon.ico")
	frameworkRuntime := config.Config.FrameworkRuntimePath()
	libEnergyPath := filepath.Join(frameworkRuntime, libEnergy)
	libWebView2LoaderPath := filepath.Join(frameworkRuntime, libWebview2)

	data := map[string]any{}
	data["BuildName"] = buildFileName                                // 应用运行二进制名
	data["BuildFileNamePath"] = filepath.Join(output, buildFileName) // 二进制文件目录
	data["InstallFileName"] = packageName                            // 安装包名
	data["CompanyName"] = appCompanyName                             // 企业名
	data["ProductName"] = appProductName                             // 产品名
	data["ShortCutName"] = appOption.Title                           // 快捷方试名
	data["IsShortcut"] = buildOption.WinDesktopShortcut              // 是否快捷方试名
	data["IsStartMenu"] = buildOption.WinAddStartMenu                // 是否开始菜单
	data["FileVersion"] = appOption.Version                          //
	data["ProductVersion"] = appOption.Version                       //
	data["FileDescription"] = appOption.Desc                         //
	data["Copyright"] = appOption.Copyright                          //
	if buildOption.WinDefaultInstall != "" {
		data["DefaultInstall"] = buildOption.WinDefaultInstall // 自定义安装目录
	}
	data["RuntimeLibEnergy"] = libEnergyPath // runtime lib energy dll
	// 判断是否使用的 webview2
	if proj.GUIRenderFramework == "WV" {
		data["RuntimeWebView2Loader"] = libWebView2LoaderPath           // runtime lib webview2  dll
		data["RuntimeWebView2Setup"] = "MicrosoftEdgeWebview2Setup.exe" // webview2 setup exe
	}
	data["NSISIcon"] = iconIcoFilePath   // 安装包程序图标
	data["NSISUnIcon"] = iconIcoFilePath // 安装包卸载程序图标
	data["NSISLanguage"] = "SimpChinese" // 中文: SimpChinese, 英文: English, 语言在 NSIS_HOME/Contrib/Language files
	if licensePath := filepath.Join(bean.ResourcePath(), buildOption.License); buildOption.License != "" && tool.IsExist(licensePath) {
		data["NSISLicense"] = licensePath // (license.txt) 文件路径
	}
	data["NSISRequestExecutionLevel"] = nsisExecLevel // run_level NSISRequestExecutionLevel
	data["AssociateFiles"] = paserAssociateFile(buildOption.WinAssociateFileList)
	data["AssociateProtocols"] = paserAssociateProtocol(buildOption.WinAssociateProtocolList)

	installToolsScript, err := RenderTemplate(data, string(installToolsScriptTemp))
	if err != nil {
		event.ConsoleWriteError("Package - check nsis RenderTemplate:", err.Error())
		return false
	}
	nsisInstallScriptPath := filepath.Join(output, "install-nsis.nsi")
	nsisToolScriptPath := filepath.Join(output, "install-tools.nsh")

	defer func() {
		//_ = os.Remove(nsisInstallScriptPath)
		//_ = os.Remove(nsisToolScriptPath)
	}()
	utf8Bom := []byte{0xEF, 0xBB, 0xBF}
	installNsisScriptTemp = append(utf8Bom, installNsisScriptTemp...)
	err = WriteFile(nsisInstallScriptPath, installNsisScriptTemp)
	if err != nil {
		event.ConsoleWriteError("Package - WriteFile", err.Error())
		return false
	}
	installToolsScript = append(utf8Bom, installToolsScript...)
	err = WriteFile(nsisToolScriptPath, installToolsScript)
	if err != nil {
		event.ConsoleWriteError("Package - WriteFile", err.Error())
		return false
	}
	// 执行 makensis 构建安装包命令
	err = RunCMD(output, "makensis", "install-nsis.nsi")
	if err != nil {
		event.ConsoleWriteError("Package - RunCMD makensis", err.Error())
		return false
	}
	return true
}

func paserAssociateFile(associateFileList []string) (associateFiles []TWinAssociateFiles) {
	embedPath := bean.ResourceEmbedPath()
	for _, line := range associateFileList {
		associates := strings.Split(line, "|")
		if len(associates) >= 5 {
			icon := strings.TrimSpace(associates[3])
			srcIcon := icon
			if !filepath.IsAbs(srcIcon) {
				srcIcon = filepath.Join(embedPath, srcIcon)
			} else {
				_, icon = filepath.Split(icon)
			}
			associate := TWinAssociateFiles{
				Ext:         strings.TrimSpace(associates[0]),
				FileClass:   strings.TrimSpace(associates[1]),
				Description: strings.TrimSpace(associates[2]),
				Icon:        icon,
				SrcIcon:     srcIcon,
				CommandText: strings.TrimSpace(associates[4]),
			}
			if !associate.IsEmpty() {
				associateFiles = append(associateFiles, associate)
			}
		}
	}
	return
}

func paserAssociateProtocol(associateProtocolList []string) (associateFiles []TWinAssociateProtocols) {
	for _, line := range associateProtocolList {
		associates := strings.Split(line, "|")
		if len(associates) >= 2 {
			associate := TWinAssociateProtocols{
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

func copyTmpFile() {
	proj := bean.GProject
	buildOption := proj.BuildOption
	output := buildOption.Output
	if !filepath.IsAbs(buildOption.Output) {
		output = filepath.Join(bean.GPath, output)
	}
}
