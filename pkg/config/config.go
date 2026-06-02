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
	"github.com/energye/designer/resources/metadata"
	toolExec "github.com/energye/lcl/tool/exec"
	"os"
	"path/filepath"
	"runtime"
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
type chromium map[string]string

// 设计器窗体配置
type designerConfig struct {
	Title         string              `json:"title"`         // 设计器标题
	Version       string              `json:"version"`       // 设计器版本
	Dependencies  dependencies        `json:"dependencies"`  // 核心依赖列表: "模块路径": "版本号" => github.com/energye/energy/v3@latest or github.com/energye/energy/v3@v3.0.0
	Chromium      chromium            `json:"chromium"`      // CEF chromium 可用支持版本列表
	WindowLayout  StorageWindowLayout `json:"window"`        // 设计器窗体信息
	ComponentTabs ComponentTabs       `json:"componentTabs"` // 设计器加载组件
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
	WindowLayout   StorageWindowLayout `json:"window"`          // 窗口配置
	FrameworkDir   string              `json:"framework"`       // 框架目录
	Chromium       TChromium           `json:"chromium"`        // CEF 框架
	Registry       string              `json:"registry"`        // 远程服务配置地址
	Proxy          string              `json:"proxy"`           // 代理地址
	LastProject    string              `json:"last_project"`    // 最后打开项目
	HistoryProject []string            `json:"history_project"` // 历史项目列表
	Env            map[string]*TEnv    `json:"env"`             // 环境配置
	EnvLang        string              `json:"env_lang"`        // 环境配置-语言
}

// TEnv 环境配置
type TEnv struct {
	GoRoot            []string `json:"go_root"`              // Go SDK
	GoRootSelectIndex int32    `json:"go_root_select_index"` // Go SDK select index
}

type TChromium struct {
	Dir     string   `json:"dir"`     // CEF 根目录, 默认 ~/.energy/chromium
	Version string   `json:"version"` // 当前使用的 CEF 版本
	List    []string `json:"list"`    // 已安装的 CEF 版本列表, 如 ["110", "120"]
}

// DefaultDir CEF 默认目录
func (m *TChromium) DefaultDir() string {
	return filepath.Join(energyDir, "chromium")
}

func (m *TChromium) SetDir(dir string) {
	if m.Dir != dir && tool.IsExist(dir) {
		m.Dir = dir
		UpdateConfig()
	}
}

func (m *TChromium) SetVersion(version string) {
	if version != "" && m.Version != version {
		m.Version = version
		exist := false
		for _, v := range m.List {
			if v == version {
				exist = true
				break
			}
		}
		if !exist {
			m.List = append(m.List, version)
			sort.Strings(m.List)
		}
		UpdateConfig()
	}
}

// CEFLibraryName 返回当前平台的 CEF 核心库文件名
func CEFLibraryName() string {
	switch runtime.GOOS {
	case "windows":
		return "libcef.dll"
	case "linux":
		return "libcef.so"
	case "darwin":
		return "libcef.dylib"
	default:
		return "libcef.dll"
	}
}

// IsCEFInstalled 检查 CEF 是否已完整安装
// 判断 Dir/Version/ 目录下是否存在平台的 libcef 库文件
func (m *TChromium) IsCEFInstalled() bool {
	if m.Dir == "" || m.Version == "" {
		return false
	}
	libPath := filepath.Join(m.Dir, m.Version, CEFLibraryName())
	return tool.IsExist(libPath)
}

// CEFVersionDir 返回指定版本的 CEF 安装目录
func (m *TChromium) CEFVersionDir(version string) string {
	return filepath.Join(m.Dir, version)
}

func (m *TConfig) FrameworkRuntimePath() string {
	runtimeDir := filepath.Join(m.FrameworkDir, "runtime")
	return runtimeDir
}

func (m chromium) Get(name string) string {
	return m[name]
}

func (m dependencies) Get(modPath string) string {
	return m[modPath]
}

func (m goEnv) Get(name string) string {
	value := m[name]
	if value == "" {
		value = os.Getenv(name)
	}
	return value
}

func (m goEnv) Set(name, value string) {
	m[name] = value
	_ = os.Setenv(name, value)
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
	event.ConsoleWriteInfo("Project Environment Configuration, name:", envName, "Go SDK:", goRoot)
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
	Config = &TConfig{WindowLayout: DesignerConfig.WindowLayout}
	defer func() {
		// 初始化 i18n
		_, _ = metadata.GI18n.Get(Config.EnvLang)

	}()

	if !tool.IsExist(configPath) {
		// 不存在创建 config.json
		Config.WindowLayout = DesignerConfig.WindowLayout
		Config.FrameworkDir = filepath.Join(energyDir)
		Config.WindowLayout.InitDefaultMenuView()
		Config.WindowLayout.InitDefaultContentLayout()
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

	Config.WindowLayout.InitDefaultMenuView()
	Config.WindowLayout.InitDefaultContentLayout()
	DesignerConfig.WindowLayout = Config.WindowLayout

	// 框架目录为空或无效重新设置
	if Config.FrameworkDir == "" || !tool.IsExist(Config.FrameworkDir) {
		Config.FrameworkDir = filepath.Join(energyDir)
		go UpdateConfig()
	}
}
