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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	toolExec "github.com/energye/lcl/tool/exec"
	"github.com/energye/lcl/types"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

// 设计器配置

const (
	LCLModPath    = "github.com/energye/lcl"
	CEFModPath    = "github.com/energye/cef"
	WVModPath     = "github.com/energye/wv"
	ENERGYModPath = "github.com/energye/energy/v3"
)

var (
	// DesignerConfig 设计器窗体全局配置(布局窗体信息等)
	DesignerConfig *designerConfig
	// Config energy 配置文件
	Config *TConfig
	// 用户目录
	homeDir = toolExec.HomeDir
	// energy 配置目录
	energyDir = filepath.Join(homeDir, ".energy")
	// energy 配置文件路径
	configPath = filepath.Join(energyDir, "config.json")
	// 全局当前 Go 环境配置
	GGoEnv goEnv
)

type dependencies map[string]string
type goEnv map[string]string

// 设计器窗体配置
type designerConfig struct {
	Title         string        `json:"title"`         // 设计器标题
	Version       string        `json:"version"`       // 设计器版本
	Dependencies  dependencies  `json:"dependencies"`  // 核心依赖列表: "模块路径": "版本号" => github.com/energye/energy/v3@latest or github.com/energye/energy/v3@v3.0.0
	Window        Window        `json:"window"`        // 设计器窗体信息
	ComponentTabs ComponentTabs `json:"componentTabs"` // 设计器加载组件
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
	SynEdit    Tab `json:"synedit"`
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
	Window         Window           `json:"window"`          // 窗口配置
	FrameworkDir   string           `json:"framework"`       // 框架目录
	Mod            TMod             `json:"mod"`             // 模块配置
	Registry       string           `json:"registry"`        // 远程服务配置地址
	Proxy          string           `json:"proxy"`           // 代理地址
	LastProject    string           `json:"last_project"`    // 最后打开项目
	HistoryProject []string         `json:"history_project"` // 历史项目列表
	Env            map[string]*TEnv `json:"env"`             // 环境配置
}

// TEnv 环境配置
type TEnv struct {
	GoRoot            []string `json:"go_root"`              // Go SDK
	GoRootSelectIndex int32    `json:"go_root_select_index"` // Go SDK select index
}

type TMod struct {
}

// FrameworkDirForSrcVersion 返回源码版本的框架目录路径。
// 它通过将 FrameworkDir 与 "src" 和当前版本连接来构建路径。
//
//   - string: 构建的源码版本框架目录路径，如果未设置 FrameworkDir 则返回空字符串
func (m *TConfig) FrameworkDirForSrcVersion() (string, string) {
	if m.FrameworkDir != "" {
		return filepath.Join(m.FrameworkDir, "src"), DesignerConfig.Version
	}
	return "", ""
}

func (m *TConfig) FrameworkDirForLCL() string {
	src, version := m.FrameworkDirForSrcVersion()
	return filepath.Join(src, "lcl@"+version)
}

func (m *TConfig) FrameworkDirForLCLRelativePath() string {
	_, version := m.FrameworkDirForSrcVersion()
	return "../lcl@" + version
}

func (m *TConfig) FrameworkDirForCEF() string {
	src, version := m.FrameworkDirForSrcVersion()
	return filepath.Join(src, "cef@"+version)
}

func (m *TConfig) FrameworkDirForCEFRelativePath() string {
	_, version := m.FrameworkDirForSrcVersion()
	return "../cef@" + version
}

func (m *TConfig) FrameworkDirForWV() string {
	src, version := m.FrameworkDirForSrcVersion()
	return filepath.Join(src, "wv@"+version)
}

func (m *TConfig) FrameworkDirForWVRelativePath() string {
	_, version := m.FrameworkDirForSrcVersion()
	return "../wv@" + version
}

func (m *TConfig) FrameworkDirForENERGY() string {
	src, version := m.FrameworkDirForSrcVersion()
	return filepath.Join(src, "energy@"+version)
}

func (m *TConfig) FrameworkDirForENERGYRelativePath() string {
	_, version := m.FrameworkDirForSrcVersion()
	return "../energy@" + version
}

func (m dependencies) Get(modPath string) string {
	return m[modPath]
}

func (m goEnv) Get(name string) string {
	return m[name]
}

func (m goEnv) Set(name, value string) {
	m[name] = value
	_ = os.Setenv("name", value)
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
}

// UpdateFrameworkDir 更新设计器的框架存放目录配置
func UpdateFrameworkDir(frameworkDir string) bool {
	if !tool.IsExist(frameworkDir) {
		return false
	}
	Config.FrameworkDir = frameworkDir
	return true
}

// UpdateLastProject 更新设计器的框架存放目录配置
func UpdateLastProject(projectEGPPath string) bool {
	if !tool.IsExist(projectEGPPath) {
		return false
	}
	Config.LastProject = projectEGPPath
	return true
}

// UpdateHistoryProject 更新设计器打开的历史项目
func UpdateHistoryProject(projectEGPPath string) bool {
	if !tool.IsExist(projectEGPPath) {
		return false
	}
	projectEGPPath = filepath.ToSlash(projectEGPPath)
	isAdd := true
	for _, hp := range Config.HistoryProject {
		hp = filepath.ToSlash(hp)
		if tool.Equal(projectEGPPath, hp) {
			isAdd = false
			break
		}
	}
	if isAdd {
		Config.HistoryProject = append(Config.HistoryProject, projectEGPPath)
		sort.Strings(Config.HistoryProject)
	}
	return true
}

// UpdateConfig 更新设计器配置到配置文件, 在修改了 Config 后调用
func UpdateConfig() bool {
	data, e := json.MarshalIndent(Config, "", "\t")
	err.CheckErr(e)
	e = os.WriteFile(configPath, data, os.ModePerm)
	err.CheckErr(e)
	return true
}

// UpdateEnvGoRoot 更新环境配置
func UpdateEnvGoRoot(envName string, goRoot string) {
	event.ConsoleWriteInfo("项目环境配置, name:", envName, "Go SDK:", goRoot)
	if Config.Env == nil {
		Config.Env = make(map[string]*TEnv)
	}
	if env := Config.Env[envName]; env != nil {
		selectIndex := -1
		for i, item := range env.GoRoot {
			if item == goRoot {
				selectIndex = i
				break
			}
		}
		if selectIndex == -1 {
			env.GoRoot = append(env.GoRoot, goRoot)
			env.GoRootSelectIndex = int32(len(env.GoRoot)) - 1
		} else {
			env.GoRootSelectIndex = int32(selectIndex)
		}
	} else {
		env = &TEnv{
			GoRoot:            []string{goRoot},
			GoRootSelectIndex: 0,
		}
		Config.Env[envName] = env
	}
}

// Path 返回 energy designer 环境目录
// 该目录在当前用户目录/.energy
func Path() string {
	return energyDir
}

func init() {
	DesignerConfig = &designerConfig{}
	e := json.Unmarshal(resources.Config(), DesignerConfig)
	err.CheckErr(e)
	// 从 energy 配置文件读取
	if !tool.IsDir(energyDir) {
		// 非目录删除文件
		_ = os.Remove(energyDir)
	}
	// 创建 energy 目录
	_ = os.Mkdir(energyDir, os.ModePerm)

	// config.json
	Config = &TConfig{Window: DesignerConfig.Window}
	if !tool.IsExist(configPath) {
		// 不存在创建 config.json
		Config.Window = DesignerConfig.Window
		Config.FrameworkDir = filepath.Join(energyDir)
		data, e := json.MarshalIndent(Config, "", "\t")
		err.CheckErr(e)
		e = os.WriteFile(configPath, data, 0644)
		err.CheckErr(e)
		return
	}
	// 存在读取 config.json
	data, e := os.ReadFile(configPath)
	err.CheckErr(e)
	e = json.Unmarshal(data, Config)
	err.CheckErr(e)
	DesignerConfig.Window = Config.Window

	// 框架目录为空或无效重新设置
	if Config.FrameworkDir == "" || !tool.IsExist(Config.FrameworkDir) {
		Config.FrameworkDir = filepath.Join(energyDir)
		go UpdateConfig()
	}

	// 加载 Go 环境信息
	goEnvCmd := exec.Command("go", "env", "-json")
	goEnvData, err := goEnvCmd.Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(goEnvData, &GGoEnv)
	if err != nil {
		return
	}
}
