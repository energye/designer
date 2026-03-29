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

package bean

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/winres"
	"path/filepath"
)

var (
	// GPath 全局 Path 完整项目路径, 打开项目时设置. C:/YouProjectXxx/xxx.egp > C:/YouProjectXxx
	GPath string
	// GProject 全局项目配置, 在创建或加载项目时设置
	GProject *TProject
)

// TProject 项目信息 xxx.egp 配置文件
type TProject struct {
	Name               string       `json:"name"`                 // 项目名 "Your Application"
	EGPName            string       `json:"egp_name"`             // 项目配置文件名 "xxx.egp"
	Package            string       `json:"package"`              // 项目(应用)包名 package "app"
	Main               string       `json:"main"`                 // 主程序入口文件或相对文件目录名 "main.go"
	GUIRenderFramework string       `json:"gui_render_framework"` // GUI 渲染框架  lcl, webview, cef
	UIForms            []TUIForm    `json:"ui_forms"`             // 窗体信息
	ActiveUIForm       int          `json:"active_ui_form"`       // 当前激活设计的窗体Id
	BuildOption        TBuildOption `json:"build_option"`         // 构建配置
	AppOption          TAppOption   `json:"app_option"`           // 应用配置
	Data               any          `json:"-"`                    // 其它数据
}

// TUIForm 窗体信息
type TUIForm struct {
	Id         int    `json:"id"`           // 设计窗体Id
	Name       string `json:"name"`         // 窗体名
	UIFile     string `json:"ui_file"`      // UI文件名
	GOFile     string `json:"go_file"`      // UI Go 文件名
	GOUserFile string `json:"go_user_file"` // UI Go 用户文件名
	UpdateTime string `json:"date_time"`    // 更新时间
}

// TBuildOption 构建配置
type TBuildOption struct {
	// 基础配置
	PlatformWindows  bool   `json:"windows"`
	PlatformMacOS    bool   `json:"macos"`
	PlatformLinux    bool   `json:"linux"`
	ArchAmd64        bool   `json:"arch_amd64"`
	Arch386          bool   `json:"arch_386"`
	ArchArm64        bool   `json:"arch_arm64"`
	ArchArm          bool   `json:"arch_arm"`
	ArchLoong64      bool   `json:"arch_loong64"`
	UIWin32          bool   `json:"ui_win32"`
	UICocoa          bool   `json:"ui_cocoa"`
	UIGtk2           bool   `json:"ui_gtk2"`
	UIGtk3           bool   `json:"ui_gtk3"`
	Output           string `json:"output"`
	BuildFileName    string `json:"build_file_name"`
	BuildModeDebug   bool   `json:"build_mode_debug"`
	BuildModeRelease bool   `json:"build_mode_release"`
	GoArgs           string `json:"go_args"`
	CodeObfuscation  bool   `json:"code_obfuscation"`
	DisableDebug     bool   `json:"disable_debug"`
	// 打包配置
	PackageName              string   `json:"package_name"`
	WinMsi                   bool     `json:"win_msi"`
	WinExe                   bool     `json:"win_exe"`
	WinDefaultInstall        string   `json:"win_default_install"`
	WinAssociateFileList     []string `json:"win_associate_file_list"`
	WinAssociateProtocolList []string `json:"win_associate_protocol_list"`
	WinSign                  TSign    `json:"win_sign"`
	NSIS                     TNSIS    `json:"nsis"`
	MacDMG                   bool     `json:"mac_dmg"`
	MacPKG                   bool     `json:"mac_pkg"`
	MacSign                  TSign    `json:"mac_sign"`
	MacCommonLib             bool     `json:"mac_common_lib"`
	LinuxDEB                 bool     `json:"linux_deb"`
	Depends                  string   `json:"depends"`
}

type TSign struct {
	Enable bool     `json:"enable"`
	Cert   []string `json:"cert"`
}

type TNSIS struct {
	WelcomeBanner string `json:"welcome_banner"`
	HeaderBanner  string `json:"header_banner"`
	License       string `json:"license"`
}

// TAppOption 应用配置
type TAppOption struct {
	Title     string      `json:"title"`
	Id        string      `json:"id"`
	Desc      string      `json:"desc"`
	Version   string      `json:"version"`
	Copyright string      `json:"copyright"`
	Icon      TAppIcon    `json:"icon"` // 应用图标, 作用在窗口和执行文件上
	Lang      string      `json:"lang"` // 语言
	Windows   TAppWindows `json:"windows"`
	MacOS     TAppMacOS   `json:"macos"`
	Linux     TAppLinux   `json:"linux"`
}

// TAppIcon 应用图标 png 格式
// 标准大小为 1024x1024
// 不同平台需要做对应的绽放处理
type TAppIcon struct {
	Data []byte `json:"data"`
	W    int32  `json:"w"`
	H    int32  `json:"h"`
}

// TAppWindows 应用配置-Windows
type TAppWindows struct {
	Manifest struct {
		CompatibilityOS                   int32 `json:"compatibility_os"`
		DPI                               int32 `json:"dpi"`
		RunLevel                          int32 `json:"run_level"`
		UIAccess                          bool  `json:"ui_access"`
		AutoElevate                       bool  `json:"auto_elevate"`
		DisableTheming                    bool  `json:"disable_theming"`
		DisableWindowFiltering            bool  `json:"disable_window_filtering"`
		HighResolutionScrollingAware      bool  `json:"high_resolution_scrolling_aware"`
		UltraHighResolutionScrollingAware bool  `json:"ultra_high_resolution_scrolling_aware"`
		LongPathAware                     bool  `json:"long_path_aware"`
		PrinterDriverIsolation            bool  `json:"printer_driver_isolation"`
		GDIScaling                        bool  `json:"gdi_scaling"`
		SegmentHeap                       bool  `json:"segment_heap"`
		UseCommonControlsV6               bool  `json:"use_common_controls_v6"`
	} `json:"manifest"`
}

type TAppMacOS struct {
	PList struct {
		CFBundleExecutable          string   `json:"cf_bundle_executable"`
		CFBundleName                string   `json:"cf_bundle_name"`
		CFBundleDisplayName         string   `json:"cf_bundle_display_name"`
		CFBundleLocalizations       []string `json:"cf_bundle_localizations"`
		CFBundleIdentifier          string   `json:"cf_bundle_identifier"`
		CFBundleVersion             string   `json:"cf_bundle_version"`
		CFBundleShortVersionString  string   `json:"cf_bundle_short_version_string"`
		CFBundleGetInfoString       string   `json:"cf_bundle_get_info_string"`
		CFBundleIconFile            string   `json:"cf_bundle_icon_file"`
		NSHumanReadableCopyright    string   `json:"cf_bundle_human_readable_copyright"`
		LSUIElementIndex            int32    `json:"ls_ui_element"`
		LSMinimumSystemVersionIndex int32    `json:"ls_minimum_system_version"`
		LSUIElement                 bool     `json:"-"`
		LSMinimumSystemVersion      string   `json:"-"`
	}
}

type TAppLinux struct {
}

var (
	CompatibilityOSList        = tool.NewArrayMap[winres.SupportedOS, string]()
	DPIList                    = tool.NewArrayMap[winres.DPIAwareness, string]()
	RunLevelList               = tool.NewArrayMap[winres.ExecutionLevel, string]()
	LSUIElementList            = tool.NewArrayMap[MacOSUIElementList, string]()
	LSMinimumSystemVersionList = tool.NewArrayMap[LSMinimumSystemVersion, string]()
	GUIRenderFrameworks        = tool.NewArrayMap[string, string]()
)

// InitAppOption 初始化应用程序选项，设置默认值
// 该函数用于为 TProject 结构体的 AppOption 字段设置初始默认配置，
// 包括通用的应用信息以及针对不同操作系统的特定默认设置。
func (m *TProject) InitAppOption() {
	m.AppOption.Title = "MyEnergyApp"
	m.AppOption.Id = "CompanyName.productName.AppName"
	m.AppOption.Desc = "Your Application Description."
	m.AppOption.Copyright = "Copyright (C) YYYY-YYYY Your Company Name. All rights reserved."
	m.AppOption.Version = "1.0.0.0"
	m.AppOption.Lang = "zh_CN"
	//m.AppOption.Icon = resources.Images("icons/window-icon_256x256.png") // 默认内置图标

	{
		// windows 默认值
		m.AppOption.Windows.Manifest.CompatibilityOS = int32(winres.WinVistaAndAbove)
		m.AppOption.Windows.Manifest.DPI = int32(winres.DPIAware)
		m.AppOption.Windows.Manifest.RunLevel = int32(winres.AsInvoker)
		m.AppOption.Windows.Manifest.HighResolutionScrollingAware = true
		m.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware = true
		m.AppOption.Windows.Manifest.LongPathAware = true
		m.AppOption.Windows.Manifest.GDIScaling = true
		m.AppOption.Windows.Manifest.UseCommonControlsV6 = true
	}
	{
		// macos 默认值
		m.AppOption.MacOS.PList.CFBundleExecutable = m.Name
		m.AppOption.MacOS.PList.CFBundleName = m.Name
		m.AppOption.MacOS.PList.CFBundleDisplayName = m.AppOption.Title
		m.AppOption.MacOS.PList.CFBundleLocalizations = []string{m.AppOption.Lang}
		m.AppOption.MacOS.PList.CFBundleIdentifier = m.AppOption.Id
		m.AppOption.MacOS.PList.CFBundleVersion = m.AppOption.Version
		m.AppOption.MacOS.PList.CFBundleShortVersionString = m.AppOption.Version
		m.AppOption.MacOS.PList.CFBundleGetInfoString = m.AppOption.Desc
		m.AppOption.MacOS.PList.CFBundleIconFile = m.Name + ".icns"
		m.AppOption.MacOS.PList.NSHumanReadableCopyright = m.AppOption.Copyright
		m.AppOption.MacOS.PList.LSUIElementIndex = int32(MacOSUIElementListNo)
		m.AppOption.MacOS.PList.LSUIElement = false
		m.AppOption.MacOS.PList.LSMinimumSystemVersionIndex = int32(LSMinimumSystemVersion_10_15)
		m.AppOption.MacOS.PList.LSMinimumSystemVersion = "10.15"
	}
	{
		// linux 默认值
	}
}

func (m *TProject) InitBuildOption() {
	// 基础配置
	m.BuildOption.PlatformWindows = true
	m.BuildOption.PlatformMacOS = true
	m.BuildOption.PlatformLinux = true
	m.BuildOption.ArchAmd64 = true
	m.BuildOption.Arch386 = true
	m.BuildOption.ArchArm64 = true
	m.BuildOption.ArchArm = true
	m.BuildOption.ArchLoong64 = false
	m.BuildOption.UIWin32 = true
	m.BuildOption.UICocoa = true
	m.BuildOption.UIGtk2 = true
	m.BuildOption.UIGtk3 = false
	m.BuildOption.Output = "./build"
	m.BuildOption.BuildFileName = m.Name
	m.BuildOption.BuildModeDebug = true
	m.BuildOption.BuildModeRelease = false
	m.BuildOption.GoArgs = ""
	m.BuildOption.CodeObfuscation = false
	m.BuildOption.DisableDebug = false
	// 打包配置
	m.BuildOption.PackageName = m.Name
	m.BuildOption.WinMsi = false
	m.BuildOption.WinExe = true
	m.BuildOption.WinDefaultInstall = ""
	m.BuildOption.MacDMG = false
	m.BuildOption.MacPKG = true
	m.BuildOption.MacSign.Cert = []string{
		// 默认
		`codesign -f -s "-" "$APP_NAME/Contents/Frameworks/$ENERGY.DYLIB"`,
		`codesign -f -s "-" --options runtime "$APP_NAME"`,
	}
	m.BuildOption.MacCommonLib = false
	m.BuildOption.LinuxDEB = true
	m.BuildOption.Depends = ""
}

// 模板调用 返回当前项目的所有窗体名称
func (m *TProject) GoFormNames() string {
	buf := tool.Buffer{}
	for _, form := range m.UIForms {
		buf.WriteString("&", form.Name, ", ")
	}
	return buf.String()
}

func init() {
	CompatibilityOSList.Add(winres.WinVistaAndAbove, "Windows Vista")
	CompatibilityOSList.Add(winres.Win7AndAbove, "Windows 7")
	CompatibilityOSList.Add(winres.Win8AndAbove, "Windows 8")
	CompatibilityOSList.Add(winres.Win81AndAbove, "Windows 8.1")
	CompatibilityOSList.Add(winres.Win10AndAbove, "Windows 10")
	CompatibilityOSList.Add(winres.Win11AndAbove, "Windows 11")

	DPIList.Add(winres.DPIAware, "System (开启)")
	DPIList.Add(winres.DPIUnaware, "UnAware (关闭)")
	DPIList.Add(winres.DPIPerMonitor, "PerMonitor (true/PM)")
	DPIList.Add(winres.DPIPerMonitorV2, "PerMonitorV2 (true/PM-V2)")

	RunLevelList.Add(winres.AsInvoker, "AsInvoker (当前用户)")
	RunLevelList.Add(winres.HighestAvailable, "HighestAvailable (最高可用权限)")
	RunLevelList.Add(winres.RequireAdministrator, "RequireAdministrator (要求管理员)")

	LSUIElementList.Add(MacOSUIElementListNo, "false (常规前台应用)")
	LSUIElementList.Add(MacOSUIElementListYes, "true (后台应用, 无 Dock 图标)")

	LSMinimumSystemVersionList.Add(LSMinimumSystemVersion_10_15, "10.15 (Intel)")
	LSMinimumSystemVersionList.Add(LSMinimumSystemVersion_11_0, "11.0 (Apple Silicon)")

	GUIRenderFrameworks.Add(GUIRenderFramework_LCL, "LCL (Native - Lazarus Component Library)")
	GUIRenderFrameworks.Add(GUIRenderFramework_WV, "WV (Web - WebView2, WebKit2)")
	GUIRenderFrameworks.Add(GUIRenderFramework_CEF, "CEF (Web - Chromium Embedded Framework)")
}

// 返回当前项目布局文件存放目录
func LayoutsPath() string {
	return filepath.Join(GPath, consts.LayoutsDir)
}

// 返回当前项目代码存放目录
func CodePath() string {
	return filepath.Join(GPath, GProject.Package)
}

// 返回当前项目资源路径
func ResourcePath() string {
	return filepath.Join(GPath, "resources")
}

// 返回当前应用元数据资源路径
func ResourceMetadataPath() string {
	return filepath.Join(ResourcePath(), "metadata")
}

// 返回当前项目内置资源目录
func ResourceEmbedPath() string {
	return filepath.Join(ResourcePath(), "embed")
}
