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

package uigen

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/project"
	projBean "github.com/energye/designer/project/bean"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UI 布局文件生成

var (
	debounceTimers = make(map[int]*time.Timer)
	debounceMutex  sync.Mutex
	debounceDelay  = 500 * time.Millisecond // 最大更新间隔
)

// UI布局文件生成
func runDebouncedGenerate(uiGenData designer.TUIGenerationData, type_ event.Type) {
	debounceMutex.Lock()
	defer debounceMutex.Unlock()
	tempUIGenData := uiGenData
	formTab := tempUIGenData.Component.FormTab
	formId := formTab.Id
	// 取消之前的定时器
	if timer, exists := debounceTimers[formId]; exists {
		timer.Stop()
	}
	// 创建新的定时器
	timer := time.AfterFunc(debounceDelay, func() {
		debounceMutex.Lock()
		delete(debounceTimers, formId)
		debounceMutex.Unlock()
		formName := tempUIGenData.Component.FormTab.FormRoot.Name()
		if formName == "" {
			// 如果正在加载设计时, 同时点击关闭设计窗口, 获取Name为空, 跳过生成
			return
		}

		// 尝试更新文件名
		isUpdateSelf := tryRenameFileName(formTab)

		uiFilePath := filepath.Join(project.LayoutsPath(), formTab.UIFile())
		// 执行UI生成
		err := generateUIFile(formTab.FormRoot, uiFilePath)
		if err != nil {
			logs.Error("UI布局文件生成错误:", err.Error())
		} else {
			if isUpdateSelf {
				// 触发代码生成事件 - 自引用 > 用户代码
				triggerCodeGeneration(tempUIGenData, event.CodeGenSelf)
			} else if type_ == event.CodeGenEvent {
				// 触发代码生成事件 - 绑定事件 > 用户代码
				triggerCodeGeneration(tempUIGenData, event.CodeGenEvent)
			}
			// 触发代码生成事件 - UI 布局, 全量更新
			triggerCodeGeneration(tempUIGenData, event.CodeGenUI)
			// 触发更新项目管理的窗体信息事件
			triggerProjectUpdate(formTab)
		}
	})

	debounceTimers[formId] = timer
}

// 尝试更新文件名
// 如果窗体名称被改变, 修改文件名
//
//	修改文件:
//	xxx.ui
//	xxx.ui.go
//	xxx.go
func tryRenameFileName(tempFormTab *designer.FormTab) bool {
	// ui 布局文件名
	uiFileName := tempFormTab.UIFile()

	// 验证UI布局文件名
	var uiForm *projBean.TUIForm
	for _, form := range project.Project().UIForms {
		if form.Id == tempFormTab.Id {
			uiForm = &form
			break
		}
	}
	// 判断 UI 文件名是否相同, 如果不相同则更改当前窗体的文件名
	if uiForm != nil && uiForm.UIFile != uiFileName {
		// 修改 xxx.ui 布局文件名
		oldUIFilePath := filepath.Join(project.LayoutsPath(), uiForm.UIFile)
		newUIFilePath := filepath.Join(project.LayoutsPath(), uiFileName)
		if err := os.Rename(oldUIFilePath, newUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return false
		}
		uiForm.UIFile = uiFileName

		// 修改 xxx.ui.go 布局代码文件名
		oldGoUIFilePath := filepath.Join(project.CodePath(), uiForm.GOFile)
		newGoUIFilePath := filepath.Join(project.CodePath(), tempFormTab.GOFile())
		if err := os.Rename(oldGoUIFilePath, newGoUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return false
		}
		uiForm.GOFile = tempFormTab.GOFile()

		// 修改 xxx.go 用户代码文件名
		oldGoUIUserFilePath := filepath.Join(project.CodePath(), uiForm.GOUserFile)
		newGoUIUserFilePath := filepath.Join(project.CodePath(), tempFormTab.GOUserFile())
		if err := os.Rename(oldGoUIUserFilePath, newGoUIUserFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return false
		}
		uiForm.GOUserFile = tempFormTab.GOUserFile()
		return true
	}
	return false
}
