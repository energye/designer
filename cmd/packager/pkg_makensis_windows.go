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
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
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

	installNsisTemp := app.Packager("windows/install-nsis.nsi")
	installToolsTemp := app.Packager("windows/install-tools.nsh")

	buildFileName := buildOption.BuildFileName
	if filepath.Ext(buildFileName) != ".exe" {
		buildFileName += ".exe"
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

	data := map[string]any{}
	data["BuildName"] = buildFileName         // 应用运行二进制名
	data["InstallFileName"] = appOption.Title // 安装包名
	data["CompanyName"] = appCompanyName      // 企业名
	data["ProductName"] = appProductName      // 产品名
	data["ShortCutName"] = appOption.Title    // 快捷方试名

	fmt.Println(string(installNsisTemp))
	fmt.Println(string(installToolsTemp))

	return true
}
