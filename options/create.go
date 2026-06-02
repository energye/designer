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
	logs.Debug("Run to create project - directory:", dir)
	if !tool.IsExist(dir) {
		event.ConsoleWriteError("Directory does not exist:", dir)
		return false
	}
	de, err := os.ReadDir(dir)
	if err != nil {
		event.ConsoleWriteError("Failed to read directory:", err.Error())
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
		msg := fmt.Sprintf("Project configuration already exists %s\nOverwrite？", existEgp)
		event.ConsoleWriteWarn("Project config exists in current directory.", existEgp, "Overwrite?")
		isCreate = api.MessageDlg(msg, types.MtCustom, types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
		if !isCreate {
			event.ConsoleWriteInfo("Cancel project creation")
			return false
		}
		// 覆盖并创建项目, 删除已存在的 xx.egp 文件
		existEGPPath := filepath.Join(dir, existEgp)
		event.ConsoleWriteWarn("Create and overwrite project configuration:", existEGPPath)
		err = os.Remove(existEGPPath)
		if err != nil {
			event.ConsoleWriteError("Failed to delete project configuration file:", err.Error())
			return false
		}
	} else if isNotEmpty {
		// 目录非空并且没有项目配置文件 egp, 提示是否在当前目录创建项目
		event.ConsoleWriteWarn("Directory not empty. Create anyway?")
		isCreate = api.MessageDlg("Directory is not empty. Create?", types.MtCustom, types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
		if !isCreate {
			event.ConsoleWriteInfo("Cancel Project Creation")
			return false
		}
	}
	return true
}

type CreateProject struct {
	Name               string
	Dir                string
	GuiRenderFramework bean.GUIRenderFramework
	FrameworkVersion   string
}

// 运行创建项目
func doRunCreate(create *CreateProject) bool {
	// 开始创建项目
	event.ConsoleWriteInfo("Start creating project", create.Name)
	newProject := new(bean.TProject)
	newProject.Name = create.Name
	newProject.EGPName = create.Name + consts.EGPExt
	newProject.Main = "main.go"
	newProject.GUIRenderFramework = create.GuiRenderFramework
	newProject.FrameworkVersion = create.FrameworkVersion
	newProject.Package = consts.AppPackageName
	newProject.InitAppOption()   // 初始化应用配置数据
	newProject.InitBuildOption() // 初始化构建配置
	// 创建并写入项目配置文件
	if err := WriteEGPConfig(create.Dir, newProject); err != nil {
		event.ConsoleWriteError("Failed to create project and write configuration:", err.Error())
		SetGlobalProject("", nil)
		return false
	} else {
		// 设置项目目录
		SetGlobalProject(create.Dir, newProject)
		// 创建项目目录结构和文件
		createProjectDir()
		// 创建 windows 应用程序清单配置
		saveOrUpdateWindowsManifest()
		// 创建 本地语言文件
		createAppLocalizations()
		// 更新应用图标
		updateWindowICON()
		// 创建项目成功
		event.ConsoleWriteInfo("Project created successfully", newProject.Name, create.Dir)
		return true
	}
}

// 写入项目配置文件
func WriteEGPConfig(path string, project *bean.TProject) error {
	if project == nil {
		return errors.New("项目配置为空")
	}
	if project.CheckLinuxWSGTK3() {
		// Linux: 如果渲染框架为 WV 或 CEF 则强制启用 GTK3
		project.BuildOption.UIGtk3 = true
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
