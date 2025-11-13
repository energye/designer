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

package project

import (
	"encoding/json"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project/bean"
	"os"
	"path/filepath"
	"strings"
)

func runLoad(filePath string) {
	logs.Debug("运行加载项目/UI 文件目录:", filePath)
	Load(filePath)
}

// 加载项目
func Load(filePath string) {
	if tool.IsExist(filePath) {
		path, file := filepath.Split(filePath)
		// 加载文件
		// 项目配置文件
		isEgp := strings.ToLower(filepath.Ext(file)) == consts.EGPExt
		// UI 布局文件
		isUI := strings.ToLower(filepath.Ext(file)) == consts.UIExt
		if isEgp {
			LoadProject(path, filePath)
		} else if isUI {
			LoadUI(filePath)
		} else {
			logs.Warn("文件非项目配置文件(.egp)或UI布局文件(.ui)")
			event.ConsoleWriteError("文件非项目配置文件(.egp)或UI布局文件(.ui)")
			SetGlobalProject("", nil)
			return
		}
	}
}

func LoadProject(path, egpFilePath string) {
	logs.Info("开始加载项目配置文件:", egpFilePath)
	event.ConsoleWriteInfo("开始加载项目配置文件:", egpFilePath)
	data, err := os.ReadFile(egpFilePath)
	if err != nil {
		logs.Error("读取项目配置文件失败:", err)
		event.ConsoleWriteError("读取项目配置文件失败:", err.Error())
		SetGlobalProject("", nil)
		return
	}
	loadProject := &bean.TProject{}
	err = json.Unmarshal(data, loadProject)
	if err != nil {
		logs.Error("解析项目配置文件失败:", err)
		event.ConsoleWriteError("解析项目配置文件失败:", err.Error())
		SetGlobalProject("", nil)
		return
	}
	event.ConsoleWriteInfo("加载项目成功", loadProject.Name)
	SetGlobalProject(path, loadProject)
	// 恢复设计器窗体
	designer.RecoverDesignerFormTab(gPath, loadProject, nil)
}

func LoadUI(uiFilePath string) {
	logs.Info("开始加载UI布局文件:", uiFilePath)
	event.ConsoleWriteInfo("开始加载UI布局文件:", uiFilePath)
	if gPath == "" || gProject == nil {
		logs.Error("不允许加载的UI布局, 当前项目未创建")
		event.ConsoleWriteError("不允许加载的UI布局, 当前项目未创建")
		return
	}
	path, uiFileName := filepath.Split(uiFilePath)
	// 匹配 ui 文件是否属于当前项目
	if !strings.HasPrefix(gPath, path) {
		logs.Error("不允许加载的UI布局, 不属于当前项目:", uiFilePath)
		event.ConsoleWriteError("不允许加载的UI布局, 不属于当前项目:", uiFilePath)
		return
	}
	var loadUIForm *bean.TUIForm
	for _, uiForm := range gProject.UIForms {
		if uiForm.UIFile == uiFileName {
			loadUIForm = &uiForm
			break
		}
	}
	if loadUIForm == nil {
		logs.Error("UI布局, 在当前项目未匹配到, 无法加载:", uiFilePath)
		event.ConsoleWriteError("UI布局, 在当前项目未匹配到, 无法加载:", uiFilePath)
		return
	}
	event.ConsoleWriteInfo("开始加载UI布局文件:", uiFilePath)
	// 恢复设计器窗体
	designer.RecoverDesignerFormTab(gPath, gProject, loadUIForm)
}
