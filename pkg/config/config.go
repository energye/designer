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

package config

import (
	"encoding/json"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/tool/homedir"
	"github.com/energye/lcl/types"
	"os"
	"path/filepath"
)

// 设计器配置

var (
	// FormConfig 设计器窗体全局配置(布局窗体信息等)
	FormConfig *formConfig
	// Config energy 配置文件
	Config *TConfig
	// 用户目录
	homeDir string
)

// 设计器窗体配置
type formConfig struct {
	Title         string        `json:"title"`
	Version       string        `json:"version"`
	Window        Window        `json:"window"`
	ComponentTabs ComponentTabs `json:"componentTabs"`
}

// 设计器窗口配置
type Window struct {
	X           int32              `json:"x"`
	Y           int32              `json:"y"`
	Width       int32              `json:"width"`
	Height      int32              `json:"height"`
	WindowState types.TWindowState `json:"window_state"`
}

// 设计器组件标签页
type ComponentTabs struct {
	Standard   Tab `json:"standard"`
	Additional Tab `json:"additional"`
	Common     Tab `json:"common"`
	Dialogs    Tab `json:"dialogs"`
	Misc       Tab `json:"misc"`
	System     Tab `json:"system"`
	LazControl Tab `json:"lazcontrol"`
	WebView    Tab `json:"webview"`
}

// 设计器组件标签页项
type Tab struct {
	Cn        string   `json:"cn"`
	En        string   `json:"en"`
	Component []string `json:"component"`
}

// TConfig energy 配置文件
type TConfig struct {
	Window       Window `json:"window"`    // 窗口配置
	FrameworkDir string `json:"framework"` // 框架目录
}

// UpdateWindow 更新窗体配置
// 在窗体改变大小时调用, 窗体关闭时
func UpdateWindow(x, y, w, h int32, windowState types.TWindowState) {
	if windowState == types.WsNormal {
		Config.Window.X = x
		Config.Window.Y = y
		Config.Window.Width = w
		Config.Window.Height = h
	}
	Config.Window.WindowState = windowState
	energyDir := filepath.Join(homeDir, ".energy")
	configPath := filepath.Join(energyDir, "config.json")
	data, e := json.MarshalIndent(Config, "", "\t")
	err.CheckErr(e)
	e = os.WriteFile(configPath, data, os.ModePerm)
	err.CheckErr(e)
}

// UpdateFrameworkDir 更新设计器的框架存放目录配置
func UpdateFrameworkDir(frameworkDir string) bool {
	if !tool.IsExist(frameworkDir) {
		return false
	}
	Config.FrameworkDir = frameworkDir
	energyDir := filepath.Join(homeDir, ".energy")
	configPath := filepath.Join(energyDir, "config.json")
	data, e := json.MarshalIndent(Config, "", "\t")
	err.CheckErr(e)
	e = os.WriteFile(configPath, data, os.ModePerm)
	err.CheckErr(e)
	return true
}

func init() {
	FormConfig = &formConfig{}
	e := json.Unmarshal(resources.Config(), FormConfig)
	err.CheckErr(e)
	// 从 energy 配置文件读取
	homeDir, e = homedir.Dir()
	err.CheckErr(e)
	energyDir := filepath.Join(homeDir, ".energy")
	if !tool.IsDir(energyDir) {
		// 非目录删除文件
		e = os.Remove(energyDir)
		err.CheckErr(e)
	}
	// 创建 energy 目录
	_ = os.Mkdir(energyDir, os.ModePerm)

	// config.json
	Config = &TConfig{}
	configPath := filepath.Join(energyDir, "config.json")
	if !tool.IsExist(configPath) {
		// 不存在创建 config.json
		Config.Window = FormConfig.Window
		data, e := json.MarshalIndent(Config, "", "\t")
		err.CheckErr(e)
		e = os.WriteFile(configPath, data, os.ModePerm)
		err.CheckErr(e)
		return
	}
	// 存在读取 config.json
	data, e := os.ReadFile(configPath)
	err.CheckErr(e)
	e = json.Unmarshal(data, Config)
	err.CheckErr(e)
	FormConfig.Window = Config.Window
}
