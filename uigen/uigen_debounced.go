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
	debounceDelay  = 500 * time.Millisecond
)

// UI布局文件生成
func runDebouncedGenerate(formTab *designer.FormTab) {
	debounceMutex.Lock()
	defer debounceMutex.Unlock()
	formName := formTab.Id
	// 取消之前的定时器
	if timer, exists := debounceTimers[formName]; exists {
		timer.Stop()
	}

	// 创建新的定时器
	timer := time.AfterFunc(debounceDelay, func() {
		tempFormTab := formTab
		debounceMutex.Lock()
		delete(debounceTimers, formName)
		debounceMutex.Unlock()

		// 尝试更新文件名
		tryRenameFileName(tempFormTab)

		uiFilePath := filepath.Join(project.Path(), project.Project().Package, tempFormTab.UIFile())
		// 执行UI生成
		err := generateUIFile(tempFormTab.FormRoot, uiFilePath)
		if err != nil {
			logs.Error("UI布局文件生成错误:", err.Error())
		} else {
			// 触发代码生成事件
			triggerCodeGeneration(tempFormTab)
			// 触发更新项目管理的窗体信息事件
			triggerProjectUpdate(tempFormTab)
		}
	})

	debounceTimers[formName] = timer
}

// 尝试更新文件名
// 如果窗体名称被改变, 修改文件名
//
//	修改文件:
//	xxx.ui
//	xxx.ui.go
//	xxx.go
func tryRenameFileName(tempFormTab *designer.FormTab) {
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
	if uiForm != nil && uiForm.UIFile != uiFileName {
		// 修改 xxx.ui 布局文件名
		oldUIFilePath := filepath.Join(project.Path(), project.Project().Package, uiForm.UIFile)
		newUIFilePath := filepath.Join(project.Path(), project.Project().Package, uiFileName)
		if err := os.Rename(oldUIFilePath, newUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return
		}
		uiForm.UIFile = uiFileName

		// 修改 xxx.ui.go 布局文件名
		oldGoUIFilePath := filepath.Join(project.Path(), project.Project().Package, uiForm.GOFile)
		newGoUIFilePath := filepath.Join(project.Path(), project.Project().Package, tempFormTab.GOFile())
		if err := os.Rename(oldGoUIFilePath, newGoUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return
		}
		uiForm.GOFile = tempFormTab.GOFile()

		// 修改 xxx.go 用户代码文件名
		oldGoUIUserFilePath := filepath.Join(project.Path(), project.Project().Package, uiForm.GOUserFile)
		newGoUIUserFilePath := filepath.Join(project.Path(), project.Project().Package, tempFormTab.GOUserFile())
		if err := os.Rename(oldGoUIUserFilePath, newGoUIUserFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return
		}
		uiForm.GOUserFile = tempFormTab.GOUserFile()
	}
}
