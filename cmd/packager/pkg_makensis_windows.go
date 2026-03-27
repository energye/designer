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
	"github.com/energye/designer/pkg/winres"
	"github.com/energye/designer/resources/app"
	"path/filepath"
	"strings"
)

// makensis.exe

const makensis = "makensis.exe"

func packageNSIS() bool {
	if !checkToolCMD(makensis) {
		event.ConsoleWriteInfo("Package - check nsis Not Installed")
		///return false
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

	// 创建临时目录
	//tmpPath := filepath.Join(output, "tmp")
	//err := os.MkdirAll(tmpPath, 0755)
	//if err != nil {
	//	return false
	//}

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
	data["DefaultInstall"] = buildOption.WinDefaultInstall           // 安装目录
	data["NSISIcon"] = iconIcoFilePath                               //
	data["NSISUnIcon"] = iconIcoFilePath                             //
	data["NSISLanguage"] = "SimpChinese"                             // 中文: SimpChinese, 英文: English, 语言在 NSIS_HOME/Contrib/Language files
	data["NSISLicense"] = ""                                         // (license.txt) 文件路径
	data["NSISRequestExecutionLevel"] = nsisExecLevel                // run_level NSISRequestExecutionLevel

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
		return false
	}
	installToolsScript = append(utf8Bom, installToolsScript...)
	err = WriteFile(nsisToolScriptPath, installToolsScript)
	if err != nil {
		return false
	}
	return true
}
