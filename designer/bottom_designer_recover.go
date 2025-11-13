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

package designer

import (
	"encoding/json"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	projBean "github.com/energye/designer/project/bean"
	uiBean "github.com/energye/designer/uigen/bean"
	"github.com/energye/lcl/lcl"
	"os"
	"path/filepath"
)

// 恢复 FormTab
// 从 UI 布局文件恢复
// 触发功能:
//  1. 打开 xxx.ui 布局文件
//	1.1 恢复成功后, 提示更新到项目配置？？TODO 不同目录是否可以？还是复制到此项目目录？
//  2. 打开项目配置文件 xxx.egp, 根据 ui_forms 字段恢复所有窗体
//	2.1 恢复所有窗体对象到设计器

// RecoverDesignerFormTab 恢复设计窗体
// 只恢复当前项目下的窗体
// path: 当前项目路径
// project: 项目对象
// loadUIForm: 要加载的 UI 窗体对象, 如果 nil 表示加载所有窗体, 否则只加载当前这个窗体
func RecoverDesignerFormTab(path string, project *projBean.TProject, loadUIForm *projBean.TUIForm) {
	if loadUIForm != nil {
		// 只加载这个窗体
	} else {
		// 加载所有
		uiForms := project.UIForms
		for _, uiForm := range uiForms {
			uiFilePath := filepath.Join(path, project.Package, uiForm.UIFile)
			data, err := os.ReadFile(uiFilePath)
			if err != nil {
				logs.Error("恢复设计窗体, 读取UI布局文件错误:", err.Error())
				event.ConsoleWriteError("恢复设计窗体, 读取UI布局文件错误:", err.Error())
				continue
			}
			uiComponent := &uiBean.TUIComponent{}
			err = json.Unmarshal(data, uiComponent)
			if err != nil {
				logs.Error("恢复设计窗体, 解析窗体布局错误:", err.Error())
				event.ConsoleWriteError("恢复设计窗体, 解析窗体布局错误:", err.Error())
				continue
			}
			lcl.RunOnMainThreadAsync(func(id uint32) {
				formTab := designer.addDesignerFormTab()
				formTab.Id = uiForm.Id
				designer.ActiveFormTab(formTab)
				// 1. 加载属性到设计器
				// 此步骤会初始化并填充设计组件实例
				formTab.FormRoot.LoadPropertyToInspector()
				// 2. 添加到组件树
				formTab.AddFormNode()
				// 3. 放置设计组件
				//formTab.placeComponent()
			})
		}

	}
}
