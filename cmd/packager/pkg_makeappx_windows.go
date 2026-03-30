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
	"bytes"
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/resize"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/app"
	"github.com/energye/lcl/tool/command"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const makeappx = "MakeAppx.exe"

func packageAppx() bool {
	var (
		libEnergy   string
		libWebview2 string
	)
	GOARCH := os.Getenv("GOARCH")
	processorArchitecture := "x64"
	switch GOARCH {
	case "amd64":
		libEnergy = "libenergy-amd64.dll"
		libWebview2 = "WebView2Loader-amd64.dll"
	case "386":
		libEnergy = "libenergy-386.dll"
		libWebview2 = "WebView2Loader-386.dll"
		processorArchitecture = "x86"
	case "arm64":
		event.ConsoleWriteWarn("Package - Currently, windows arm64 arch is not support")
		return false
	}

	proj := bean.GProject
	buildOption := proj.BuildOption
	appOption := proj.AppOption

	buildFileName := buildOption.BuildFileName
	if filepath.Ext(buildFileName) != ".exe" {
		buildFileName += ".exe"
	}
	msixPackageName := buildOption.PackageName
	if exeIdx := strings.LastIndex(msixPackageName, ".msix"); exeIdx != -1 {
		msixPackageName = msixPackageName[:exeIdx] + "_" + GOARCH + ".msix"
	} else {
		msixPackageName = msixPackageName + "_" + GOARCH + ".msix"
	}
	output := buildOption.Output
	if !filepath.IsAbs(buildOption.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	var (
		appCompanyName = ""
	)
	appID := appOption.Id // CompanyName.productName.AppName
	if ids := strings.Split(appID, "."); len(ids) >= 2 {
		appCompanyName = ids[0]
	}
	frameworkRuntime := config.Config.FrameworkRuntimePath()
	libEnergyPath := filepath.Join(frameworkRuntime, libEnergy)
	libWebView2LoaderPath := filepath.Join(frameworkRuntime, libWebview2)
	binaryFileNamePath := filepath.Join(output, buildFileName)

	// 创建 app package 目录, 并复制文件
	// xxx_msix/
	// │
	// ├── AppxManifest.xml
	// ├── xxx.exe
	// ├── 其它 DLL / 资源
	// │
	// └── Assets/
	//     ├── StoreLogo.png
	//     ├── Square44x44Logo.png
	//     ├── Square150x150Logo.png
	//     ├── Wide310x150Logo.png
	//     ├── SplashScreen.png
	//     ├── FileIcon.png
	//     └── ProtocolLogo.png
	// | 文件               | 尺寸      |
	// | ----------------- | -------   |
	// | Square44x44Logo   | 44x44     |
	// | Square71x71Logo   | 71x71     |
	// | Square150x150Logo | 150x150   |
	// | Wide310x150Logo   | 310x150   |
	// | SplashScreen      | 620x300   |
	// | StoreLogo         | 50x50     |
	// | FileIcon          | 256 x 256 |
	// | ProtocolLogo      | 88 x 88   |

	publisher := ""
	if buildOption.WinSign.Enable {
		// signtool /dump /v codesign.pfx
		cmdInfo := getSignCMDInfo()
		if cmdInfo == nil {
			event.ConsoleWriteError("Package - 已开启签名, 但获取签名命令配置失败")
			return false
		}
		if cmdInfo.Type == "auto" {
			publisher = publisherFromInstalledCert()
		} else {
			publisher = publisherFromPFX(cmdInfo.File, cmdInfo.Password)
		}
	}

	// 模板填充数据
	data := map[string]any{}
	data["BinaryName"] = buildFileName                    // 应用运行二进制名
	data["CompanyName"] = appCompanyName                  // 企业名
	data["ProductIdentity"] = appOption.Id                // app 唯一id, 内部标识
	data["Publisher"] = publisher                         // Publisher，"signtool /dump /v cert.pfx" : CN=MyCompany, O=MyOrg, C=CN
	data["ProcessorArchitecture"] = processorArchitecture // ProcessorArchitecture，x86 x64 arm arm64
	data["DisplayName"] = appOption.Title                 // 显示名
	data["ProductVersion"] = appOption.Version            // app 版本
	data["Description"] = appOption.Desc                  //
	data["AssociateFiles"] = paserAssociateFile(buildOption.WinAssociateFileList)
	data["AssociateFileInfoTip"] = "使用 " + appOption.Title + " 打开" // TODO 需换成自动语言
	data["AssociateProtocols"] = paserAssociateProtocol(buildOption.WinAssociateProtocolList)

	// 创建 app 打包目录
	energyMsixAppRootDir := filepath.Join(output, "energy_msix_"+GOARCH)
	err := os.MkdirAll(energyMsixAppRootDir, 0755)
	if err != nil {
		event.ConsoleWriteError("Package - Failed to create directory for energy_msix_app:", err.Error())
		return false
	}
	assetsDir := filepath.Join(energyMsixAppRootDir, "Assets")
	err = os.MkdirAll(assetsDir, 0755)
	if err != nil {
		event.ConsoleWriteError("Package - Failed to create directory for Assets:", err.Error())
		return false
	}
	copyAppBinary := filepath.Join(energyMsixAppRootDir, buildFileName)
	copyLibEnergyBinary := filepath.Join(energyMsixAppRootDir, libEnergy)

	{
		// copy file
		// app 二进制
		err = tool.CopyFile(binaryFileNamePath, copyAppBinary)
		if err != nil {
			event.ConsoleWriteError("Package - Copy app.binary.exe file:", err.Error())
			return false
		}
		// libenergy-xxx.dll
		err = tool.CopyFile(libEnergyPath, copyLibEnergyBinary)
		if err != nil {
			event.ConsoleWriteError("Package - Copy lib.energy.dll file:", err.Error())
			return false
		}
		// webview2loader.dll
		if proj.GUIRenderFramework == "WV" {
			// runtime lib webview2  dll
			err = tool.CopyFile(libWebView2LoaderPath, filepath.Join(energyMsixAppRootDir, libWebview2))
			if err != nil {
				event.ConsoleWriteError("Package - Copy Webview2Loader.dll file:", err.Error())
				return false
			}
		}
		// AppxManifest.xml
		installAppxMaifestTemp := app.Packager("windows/install-appx-manifest.xml")
		installAppxMaifest, err := RenderTemplate(data, string(installAppxMaifestTemp))
		if err != nil {
			event.ConsoleWriteError("Package - Render AppxManifest.xml template:", err.Error())
			return false
		}
		installAppxMaifestPath := filepath.Join(energyMsixAppRootDir, "AppxManifest.xml")
		err = WriteFile(installAppxMaifestPath, installAppxMaifest)
		if err != nil {
			event.ConsoleWriteError("Package - Create AppxManifest.xml:", err.Error())
			return false
		}
		// Assets dir
		pngFiles := []assetsPng{
			{Name: "StoreLogo.png", W: 50, H: 50},
			{Name: "Square44x44Logo.png", W: 44, H: 44},
			{Name: "Square150x150Logo.png", W: 150, H: 150},
			{Name: "Wide310x150Logo.png", W: 310, H: 150},
			{Name: "SplashScreen.png", W: 6200, H: 300},
			{Name: "FileIcon.png", W: 256, H: 256},
			{Name: "ProtocolLogo.png", W: 88, H: 88},
		}
		embedPath := bean.ResourceEmbedPath()
		resourcesPath := bean.ResourcePath()
		srcIconPng := filepath.Join(embedPath, "icon.png")
		srcIconPngFile, err := os.Open(srcIconPng)
		if err != nil {
			event.ConsoleWriteError("Package - Open icon.png:", err.Error())
			return false
		}
		defer srcIconPngFile.Close()
		srcIconPngSrcImg, _, err := image.Decode(srcIconPngFile)
		if err != nil {
			event.ConsoleWriteError("Package - Decode icon.png:", err.Error())
			return false
		}
		for _, pngFile := range pngFiles {
			// 项目目录 resources/assets 查找这个文件, 如果有直接复制到 Assets 打包目录
			srcPng := filepath.Join(resourcesPath, "assets", pngFile.Name)
			if tool.IsExist(srcPng) {
				err = tool.CopyFile(srcPng, filepath.Join(assetsDir, pngFile.Name))
				if err != nil {
					event.ConsoleWriteError("Package - Copy Assets:", err.Error())
					return false
				}
			} else {
				// 项目目录 resources/assets 没有这个文件
				if pngFile.Name == "FileIcon.png" {
					// 使用 icon.png 256x256
					err = tool.CopyFile(srcIconPng, filepath.Join(assetsDir, pngFile.Name))
					if err != nil {
						event.ConsoleWriteError("Package - Copy FileIcon.png:", err.Error())
						return false
					}
				} else if pngFile.Name == "ProtocolLogo.png" {
					// 使用 icon.png 88x88
					newIconPng := resize.Resize(88, 88, srcIconPngSrcImg, resize.Lanczos3)
					err = saveAccessPNG(newIconPng, filepath.Join(assetsDir, "ProtocolLogo.png"))
					if err != nil {
						event.ConsoleWriteError("Package - Resize And Save ProtocolLogo.png:", err.Error())
						return false
					}
				} else {
					// 创建空文件代替
					newIconPng, err := createTransparentPNG(pngFile.W, pngFile.H)
					if err != nil {
						event.ConsoleWriteError("Package - Create Empty", pngFile.Name, err.Error())
						return false
					}
					err = saveAccessPNG(newIconPng, filepath.Join(assetsDir, pngFile.Name))
					if err != nil {
						event.ConsoleWriteError("Package - Save", pngFile.Name, err.Error())
						return false
					}
				}
			}
		}
	}

	// 签名文件 signtool
	// 应用二进制文件 和 libenergy.dll
	signWindowsBinary(copyAppBinary)
	signWindowsBinary(copyLibEnergyBinary)

	// 执行 makeappx 构建安装包命令
	packageArgs := []string{"pack", "/d", energyMsixAppRootDir, "/p", msixPackageName}
	event.ConsoleWriteInfo("Package:", strings.Join(packageArgs, " "))
	err = RunCMD(output, makeappx, packageArgs...)
	if err != nil {
		event.ConsoleWriteError("Package - RunCMD", makeappx, err.Error())
		return false
	}
	msixPackagePath := filepath.Join(output, msixPackageName)
	signWindowsBinary(msixPackagePath)
	return true
}

type assetsPng struct {
	Name string
	W, H int
}

func createTransparentPNG(width, height int) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.Transparent)
		}
	}
	return img, nil
}

func saveAccessPNG(srcPng image.Image, outputPath string) error {
	newPngFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer newPngFile.Close()
	err = png.Encode(newPngFile, srcPng)
	return err
}

// 根据指定证书获取 Subject
func publisherFromPFX(pfxPath, pfxPassword string) string {
	cmdStr := fmt.Sprintf(`
$pfxPath = "%s";
$password = "%s";
$securePass = ConvertTo-SecureString $password -AsPlainText -Force;
$cert = (New-Object System.Security.Cryptography.X509Certificates.X509Certificate2);
$cert.Import($pfxPath, $securePass, [System.Security.Cryptography.X509Certificates.X509KeyStorageFlags]::PersistKeySet);
$cert.Subject;
	`, pfxPath, pfxPassword)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", cmdStr)
	cmd.SysProcAttr = command.HideWindow(true)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return ""
	}
	publisher := strings.TrimSpace(stdout.String())
	return publisher
}

// 根据已安装证书获取 Subject
func publisherFromInstalledCert() string {
	cmdStr := `
	Get-ChildItem -Path Cert:\CurrentUser\My | 
	Where-Object { $_.HasPrivateKey -and $_.EnhancedKeyUsageList -match 'CodeSigning' } | 
	Select-Object -ExpandProperty Subject -First 1
	`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", cmdStr)
	cmd.SysProcAttr = command.HideWindow(true)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return ""
	}
	publisher := strings.TrimSpace(stdout.String())
	if publisher == "" {
		return ""
	}
	return publisher
}
