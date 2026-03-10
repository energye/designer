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
	"fmt"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"os"
	"path/filepath"
	"strings"
)

// runLoad 运行加载项目/UI文件目录
// filePath: 要加载的文件路径
func runLoad(filePath string) {
	logs.Debug("运行加载项目/UI 文件目录:", filePath)
	Load(filePath)
}

// Load 加载指定路径的文件，根据文件扩展名判断文件类型并调用相应的加载函数
// filePath: 要加载的文件路径
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
			event.ConsoleWriteError("文件非项目配置文件(.egp)或UI布局文件(.ui)")
			SetGlobalProject("", nil)
			return
		}
	}
}

// LoadProject 加载项目配置文件并初始化项目
// path: 项目路径
// egpFilePath: 项目配置文件路径
func LoadProject(path, egpFilePath string) {
	event.ConsoleWriteInfo("开始加载项目配置文件:", egpFilePath)
	data, err := os.ReadFile(egpFilePath)
	if err != nil {
		event.ConsoleWriteError("读取项目配置文件失败:", err.Error())
		SetGlobalProject("", nil)
		return
	}
	loadProject := &bean.TProject{}
	err = json.Unmarshal(data, loadProject)
	if err != nil {
		event.ConsoleWriteError("解析项目配置文件失败:", err.Error())
		SetGlobalProject("", nil)
		return
	}
	event.ConsoleWriteInfo("加载项目成功", loadProject.Name)
	SetGlobalProject(path, loadProject)
	// 触发文件修改监听事件
	event.Emit(event.TTrigger{Name: event.ListenFileChange})
	// 重置设计器
	designer.ResetDesigner()
	// 恢复设计器窗体
	designer.RecoverDesignerFormTab(bean.GPath, loadProject, nil)
	// 加载完后设置窗口标题
	designer.UpdateDesignerTitle(fmt.Sprintf("%v (%v)", loadProject.Name, bean.GPath))
	// 更新打开项目历史记录
	designer.UpdateHistoryProject(egpFilePath)
}

// LoadUI 加载UI布局文件
// uiFilePath: UI布局文件的完整路径
func LoadUI(uiFilePath string) {
	event.ConsoleWriteInfo("开始加载UI布局文件:", uiFilePath)
	if bean.GPath == "" || bean.GProject == nil {
		event.ConsoleWriteError("不允许加载的UI布局, 当前项目未创建")
		return
	}
	path, uiFileName := filepath.Split(uiFilePath)
	// 匹配 ui 文件是否属于当前项目
	if !strings.HasPrefix(path, bean.GPath) {
		event.ConsoleWriteError("不允许加载的UI布局, 不属于当前项目:", uiFilePath)
		return
	}
	var loadUIForm *bean.TUIForm
	for _, uiForm := range bean.GProject.UIForms {
		if uiForm.UIFile == uiFileName {
			loadUIForm = &uiForm
			break
		}
	}
	if loadUIForm == nil {
		event.ConsoleWriteError("UI布局, 在当前项目未匹配到, 无法加载:", uiFilePath)
		return
	}
	event.ConsoleWriteInfo("开始加载UI布局文件:", uiFilePath)
	// 恢复设计器窗体
	designer.RecoverDesignerFormTab(bean.GPath, bean.GProject, loadUIForm)
}
