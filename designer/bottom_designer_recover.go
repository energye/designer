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
	"github.com/energye/designer/uigen/bean"
	"os"
)

// 恢复 FormTab
// 从 UI 布局文件恢复
// 触发功能:
//  1. 打开 xxx.ui 布局文件
//	1.1 恢复成功后, 提示更新到项目配置？？TODO 不同目录是否可以？还是复制到此项目目录？
//  2. 打开项目配置文件 xxx.egp, 根据 ui_forms 字段恢复所有窗体
//	2.1 恢复所有窗体对象到设计器

// RecoverDesignerFormTab 恢复设计窗体
func RecoverDesignerFormTab(uiFilePaths ...string) {
	for _, uiFilePath := range uiFilePaths {
		data, err := os.ReadFile(uiFilePath)
		if err != nil {
			logs.Error("恢复设计窗体, 读取UI布局文件错误:", err.Error())
			event.ConsoleWriteError("恢复设计窗体, 读取UI布局文件错误:", err.Error())
			continue
		}
		uiComponent := bean.TUIComponent{}
		err = json.Unmarshal(data, &uiComponent)
		if err != nil {
			logs.Error("恢复设计窗体, 解析窗体布局错误:", err.Error())
			event.ConsoleWriteError("恢复设计窗体, 解析窗体布局错误:", err.Error())
			continue
		}

	}
}
