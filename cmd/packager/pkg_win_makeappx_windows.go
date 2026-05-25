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
	"github.com/energye/designer/resources/frameworks/lib"
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

func (m *Package) packageAppx() bool {
	if !checkToolCMD(makeappx) {
		event.ConsoleWriteError("Package - check", makeappx, " Not Installed")
		return false
	}
	var (
		libEnergy   = lib.GetDLLName()
		libWebview2 string
	)
	GOARCH := os.Getenv("GOARCH")
	processorArchitecture := "x64"
	switch GOARCH {
	case "amd64":
		libWebview2 = "WebView2Loader-amd64.dll"
	case "386":
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
	if m.AppendPlatform {
		msixPackageName += "_" + lib.GOOS()
	}
	if m.AppendArch {
		msixPackageName += "_" + lib.GOARCH()
	}
	if m.CustomSuffix != "" {
		msixPackageName += "_" + m.CustomSuffix
	}
	msixPackageName += ".msix"

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
	assetsPath := bean.ResourceAssetsPath()
	embedPath := bean.ResourceEmbedPath()
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
	//     ├── Square44x44Logo.png
	//     ├── Square150x150Logo.png
	//     ├── Wide310x150Logo.png
	//     ├── SplashScreen.png
	//     ├── PropertiesLogo.png
	//     ├── AssociateFileIcon.png
	//     └── AssociateProtocolLogo.png
	// | 文件               		| 尺寸      |
	// | ---------------------- | -------   |
	// | Square44x44Logo   		| 44x44     |
	// | Square150x150Logo 		| 150x150   |
	// | Wide310x150Logo   		| 310x150   |
	// | SplashScreen      		| 620x300   |
	// | PropertiesLogo    		| 50x50     |
	// | AssociateFileIcon      | 256 x 256 |
	// | AssociateProtocolLogo  | 88 x 88   |

	publisher := ""
	if buildOption.WinSign.Enable {
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
	// 创建 app 打包目录
	energyMsixAppRootDir := filepath.Join(output, "energy_msix_"+GOARCH)
	_ = os.RemoveAll(energyMsixAppRootDir)
	err := os.MkdirAll(energyMsixAppRootDir, 0755)
	if err != nil {
		event.ConsoleWriteError("Package - Failed to create directory for energy_msix_app:", err.Error())
		return false
	}
	msixPackagePath := filepath.Join(output, msixPackageName)
	_ = os.Remove(msixPackagePath)

	assetsDir := filepath.Join(energyMsixAppRootDir, "Assets")
	err = os.MkdirAll(assetsDir, 0755)
	if err != nil {
		event.ConsoleWriteError("Package - Failed to create directory for Assets:", err.Error())
		return false
	}
	copyAppBinary := filepath.Join(energyMsixAppRootDir, buildFileName)
	dstLibName := packageLibName()
	copyLibEnergyBinary := filepath.Join(energyMsixAppRootDir, dstLibName)

	// 应用二进制程序使用的图片
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
	// Assets dir
	PropertiesLogo := &assetsPng{Name: "PropertiesLogo.png", W: 50, H: 50}
	Square44x44Logo := &assetsPng{Name: "Square44x44Logo.png", W: 44, H: 44}
	Square150x150Logo := &assetsPng{Name: "Square150x150Logo.png", W: 150, H: 150}
	Wide310x150Logo := &assetsPng{Name: "Wide310x150Logo.png", W: 310, H: 150}
	SplashScreen := &assetsPng{Name: "SplashScreen.png", W: 620, H: 300}
	AssociateFileIcon := &assetsPng{Name: "AssociateFileIcon.png", W: 256, H: 256}
	AssociateProtocolLogo := &assetsPng{Name: "AssociateProtocolLogo.png", W: 88, H: 88}

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
	data["AssociateFiles"] = parserAssociateFile(buildOption.WinAssociateFileList)
	data["AssociateFileInfoTip"] = "Use " + appOption.Title + " Open" // TODO 需换成自动语言
	data["AssociateProtocols"] = parserAssociateProtocol(buildOption.WinAssociateProtocolList)

	// 处理 appx/Assets 图片资源
	handleAssets := func(logo, customLogoName string, useSrcIconPng bool, assPng *assetsPng) bool {
		if srcLogoPath := filepath.Join(assetsPath, customLogoName); tool.IsExist(srcLogoPath) && customLogoName != "" {
			// 使用自定义图
			err = tool.CopyFile(srcLogoPath, filepath.Join(assetsDir, customLogoName))
			if err != nil {
				event.ConsoleWriteError("Package - Custom Resize And Save", srcLogoPath, err.Error())
				return false
			}
			data[logo] = customLogoName
		} else if useSrcIconPng {
			// 使用应用二进制图
			newIconPng := resize.Resize(assPng.W, assPng.H, srcIconPngSrcImg, resize.Lanczos3)
			err = saveAccessPNG(newIconPng, filepath.Join(assetsDir, assPng.Name))
			if err != nil {
				event.ConsoleWriteError("Package - Use src icon Resize And Save", assPng.Name, err.Error())
				return false
			}
			data[logo] = assPng.Name
		} else {
			// 使用创建空透明图
			newIconPng, err := createTransparentPNG(int(assPng.W), int(assPng.H))
			if err != nil {
				event.ConsoleWriteError("Package - Create Transparent PNG", assPng.Name, err.Error())
				return false
			}
			err = saveAccessPNG(newIconPng, filepath.Join(assetsDir, assPng.Name))
			if err != nil {
				event.ConsoleWriteError("Package - Save", assPng.Name, err.Error())
				return false
			}
			data[logo] = assPng.Name
		}
		return true
	}
	if !handleAssets("PropertiesLogo", bean.GProject.BuildOption.WinAppx.PropertiesLogo, true, PropertiesLogo) {
		return false
	}
	if !handleAssets("Square44x44Logo", bean.GProject.BuildOption.WinAppx.Square44x44Logo, true, Square44x44Logo) {
		return false
	}
	if !handleAssets("Square150x150Logo", bean.GProject.BuildOption.WinAppx.Square150x150Logo, true, Square150x150Logo) {
		return false
	}
	if !handleAssets("Wide310x150Logo", bean.GProject.BuildOption.WinAppx.Wide310x150Logo, false, Wide310x150Logo) {
		return false
	}
	if !handleAssets("SplashScreen", bean.GProject.BuildOption.WinAppx.SplashScreen, false, SplashScreen) {
		return false
	}
	if !handleAssets("AssociateFileIcon", bean.GProject.BuildOption.WinAppx.AssociateFileIcon, true, AssociateFileIcon) {
		return false
	}
	if !handleAssets("AssociateProtocolLogo", bean.GProject.BuildOption.WinAppx.AssociateProtocolLogo, true, AssociateProtocolLogo) {
		return false
	}

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
	}

	// 签名文件 signtool
	// 应用二进制文件 和 libenergy.dll
	signWindowsBinary(copyAppBinary)
	signWindowsBinary(copyLibEnergyBinary)

	// 执行 makeappx 构建安装包命令
	packageArgs := []string{"pack", "/d", energyMsixAppRootDir, "/p", msixPackageName}
	event.ConsoleWriteInfo("Package:", strings.Join(packageArgs, " "))

	cmd := command.NewCMD()
	cmd.HideWindow = true
	cmd.Dir = output
	cmd.Console = func(data string, level command.Level) {
		if level == command.LError {
			event.ConsoleWriteError(data)
		} else {
			event.ConsoleWriteInfo(data)
		}
	}
	err = cmd.CommandContext(m.Context, makeappx, packageArgs...)
	if err != nil {
		event.ConsoleWriteError("Package - RunCMD", makeappx, err.Error())
		return false
	}
	signWindowsBinary(msixPackagePath)
	return true
}

type assetsPng struct {
	Name string
	W, H uint
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
