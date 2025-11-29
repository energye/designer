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

package codegen

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/project"
	"os"
	"path/filepath"
)

// 更新用户代码文件
// 此处功能为:
//   1. 事件绑定/删除/修改
//	 2. 自引用

// 执行更新自引用
// 被执行时说明当前窗体名称被修改. 窗体 name 修改过程 => file name (xxx.ui.go xxx.go xxx.ui) => go self code
// ast 用户文件
//
//	功能: 修改窗体绑定的自引用结构名. 例如 窗体结构 TForm1 > TForm2,  原: func (m *Form1) 新名: func (m *TForm2)
//	仅修改 xxx.go 用户文件, 其它文件需手动修改
func doUpdateSelf(uiGenData designer.TUIGenerationData) {
	formTab := uiGenData.Component.FormTab
	goUserFilePath := filepath.Join(project.Path(), project.Project().Package, formTab.GOUserFile())
	newFormName := "T" + formTab.FormRoot.Name()
	oldFormName := "T" + uiGenData.Component.FormTab.OldFormName
	logs.Debug("NewFormName:", newFormName, "OldFormName:", oldFormName, "FilePath:", goUserFilePath)
	// 查找 OldFormName 的类型函数名
	// 修改 OldFormName 类型函数为 newFormName 的函数
	// 写入文件
	newCode, isUpdate, err := dast.UpdateMethodRecv(goUserFilePath, oldFormName, newFormName)
	if err != nil {
		logs.Error("执行更新自引用-UpdateMethodRecv:", err.Error())
		return
	}
	if !isUpdate {
		return
	}
	st, err := os.Stat(goUserFilePath)
	if err != nil {
		logs.Error("执行更新自引用-Stat:", err.Error())
		return
	}
	if err = os.WriteFile(goUserFilePath, newCode, st.Mode()); err != nil {
		logs.Error("执行更新自引用-WriteFile:", err.Error())
		return
	}
}

// 执行更新绑定事件
// 被执行时说明组件事件被修改. 组件属性-事件 => file code (xxx.ui.go xxx.go xxx.ui) => go event code
// ast 用户文件
//
//	xxx.ui.go => 绑定事件
//	xxx.go => 事件定义
//	xxx.ui => 事件记录
func doUpdateEvent(uiGenData designer.TUIGenerationData) {
	formTab := uiGenData.Component.FormTab
	goUserFilePath := filepath.Join(project.Path(), project.Project().Package, formTab.GOUserFile())
	newFormName := "T" + formTab.FormRoot.Name()
	logs.Debug("FormName:", newFormName, goUserFilePath)

}
