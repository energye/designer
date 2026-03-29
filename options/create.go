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

package options

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"os"
	"path/filepath"
	"strings"
)

// 项目创建, 在指定目录创建新项目
// 检查目录是否为空

// 运行创建项目窗口
func runCreate() {
	// 显示创建项目窗口
	lcl.RunOnMainThreadAsync(func(id uint32) {
		form := NewCreateProjectForm()
		// 为什么 debug 时显示不出来呢？
		form.ShowModal()
		//form.Show()
		logs.Debug("runCreate NewCreateProjectForm show end")
	})
}

// 检查创建条件
func checkCreate(dir string) bool {
	logs.Debug("运行创建项目 目录:", dir)
	if !tool.IsExist(dir) {
		event.ConsoleWriteError("目录不存在:", dir)
		return false
	}
	de, err := os.ReadDir(dir)
	if err != nil {
		event.ConsoleWriteError("读取目录失败:", err.Error())
		return false
	}
	var (
		isNotEmpty bool   // 当前目录是否为空
		existEgp   string // 目录的egp文件
		isCreate   bool   // 是否创建
	)
	for _, entry := range de {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.LastIndex(name, consts.EGPExt) != -1 {
			existEgp = entry.Name()
			break
		}
		isNotEmpty = true // 非空目录
	}
	// 已存在项目 egp 文件, 提示覆盖
	if existEgp != "" {
		msg := fmt.Sprintf("当前目录已存在项目配置 %s\n是否覆盖？", existEgp)
		event.ConsoleWriteWarn("当前目录已存在项目配置", existEgp, "是否覆盖？")
		isCreate = api.MessageDlg(msg, types.MtCustom, types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
		if !isCreate {
			event.ConsoleWriteInfo("取消创建项目")
			return false
		}
		// 覆盖并创建项目, 删除已存在的 xx.egp 文件
		existEGPPath := filepath.Join(dir, existEgp)
		event.ConsoleWriteWarn("创建并覆盖项目配置文件:", existEGPPath)
		err = os.Remove(existEGPPath)
		if err != nil {
			event.ConsoleWriteError("删除项目配置文件错误件:", err.Error())
			return false
		}
	} else if isNotEmpty {
		// 目录非空并且没有项目配置文件 egp, 提示是否在当前目录创建项目
		event.ConsoleWriteWarn("当前目录非空是否创建？")
		isCreate = api.MessageDlg("当前目录非空是否创建？", types.MtCustom, types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
		if !isCreate {
			event.ConsoleWriteInfo("取消创建项目")
			return false
		}
	}
	return true
}

// 运行创建项目
func doRunCreate(name, dir, guiRenderFramework string) bool {
	// 开始创建项目
	event.ConsoleWriteInfo("开始创建项目", name)
	newProject := new(bean.TProject)
	newProject.Name = name
	newProject.EGPName = name + consts.EGPExt
	newProject.Main = "main.go"
	newProject.GUIRenderFramework = guiRenderFramework
	newProject.Package = consts.AppPackageName
	newProject.InitAppOption()   // 初始化应用配置数据
	newProject.InitBuildOption() // 初始化构建配置
	// 创建并写入项目配置文件
	if err := WriteEGPConfig(dir, newProject); err != nil {
		event.ConsoleWriteError("创建项目, 写入项目配置失败:", err.Error())
		SetGlobalProject("", nil)
		return false
	} else {
		// 设置项目目录
		SetGlobalProject(dir, newProject)
		// 创建项目目录结构和文件
		createProjectDir()
		// 创建 windows 应用程序清单配置
		saveOrUpdateWindowsManifest()
		// 创建 macOS 应用程序清单配置
		saveOrUpdateMacOSPList()
		// 创建 本地语言文件
		createAppLocalizations()
		// 更新应用图标
		updateWindowICON()
		// 创建项目成功
		event.ConsoleWriteInfo("创建项目成功", newProject.Name, dir)
		return true
	}
}

// 写入项目配置文件
func WriteEGPConfig(path string, project *bean.TProject) error {
	if project == nil {
		return errors.New("项目配置为空")
	}
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	egpFilePath := filepath.Join(path, project.EGPName)
	err = os.WriteFile(egpFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
