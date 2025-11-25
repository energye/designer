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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/winres"
)

// TProject 项目信息 xxx.egp 配置文件
type TProject struct {
	Name         string       `json:"name"`           // 项目名 "Your Application"
	EGPName      string       `json:"egp_name"`       // 项目配置文件名 "xxx.egp"
	Package      string       `json:"package"`        // 项目(应用)包名 package "app"
	Main         string       `json:"main"`           // 主程序入口文件或相对文件目录名 "main.go"
	UIForms      []TUIForm    `json:"ui_forms"`       // 窗体信息
	ActiveUIForm int          `json:"active_ui_form"` // 当前激活设计的窗体Id
	BuildOption  TBuildOption `json:"build_option"`   // 构建配置
	EnvOption    TEnvOption   `json:"env_option"`     // 环境配置
	AppOption    TAppOption   `json:"app_option"`     // 应用配置
	Data         any          `json:"-"`              // 其它数据
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
	GoArgument string `json:"go_argument"` // 构建参数
	Output     string `json:"output"`      // 构建输出目录
}

// TEnvOption 环境配置
type TEnvOption struct {
	GoRoot string `json:"go_root"` // Go 安装目录
}

// TAppOption 应用配置
type TAppOption struct {
	Title   string      `json:"app_title"`
	Id      string      `json:"id"`
	Desc    string      `json:"desc"`
	Version string      `json:"version"`
	AppIcon []byte      `json:"app_icon"`
	Lang    string      `json:"lang"` // 语言
	Windows TAppWindows `json:"windows"`
	MacOS   TAppMacOS   `json:"macos"`
	Linux   TAppLinux   `json:"linux"`
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
}

type TAppLinux struct {
}

var (
	CompatibilityOSList = tool.NewArrayMap[winres.SupportedOS, string]()
	DPIList             = tool.NewArrayMap[winres.DPIAwareness, string]()
	RunLevelList        = tool.NewArrayMap[winres.ExecutionLevel, string]()
)

func (m *TProject) InitAppOption() {
	m.AppOption.Title = "MyEnergyApp"
	m.AppOption.Id = "CompanyName.ProductName.AppName"
	m.AppOption.Desc = "Your application description."
	m.AppOption.Version = "1.0.0"
	m.AppOption.Lang = "zh_CN"

	// windows 默认值
	m.AppOption.Windows.Manifest.CompatibilityOS = int32(winres.WinVistaAndAbove)
	m.AppOption.Windows.Manifest.DPI = int32(winres.DPIAware)
	m.AppOption.Windows.Manifest.RunLevel = int32(winres.AsInvoker)
	m.AppOption.Windows.Manifest.HighResolutionScrollingAware = true
	m.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware = true
	m.AppOption.Windows.Manifest.LongPathAware = true
	m.AppOption.Windows.Manifest.GDIScaling = true
	m.AppOption.Windows.Manifest.UseCommonControlsV6 = true

	// macos 默认值

	// linux 默认值
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
}
