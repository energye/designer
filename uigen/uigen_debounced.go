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
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// UI 布局文件生成

var (
	debounceTimers = make(map[string]*time.Timer)
	debounceMutex  sync.Mutex
	debounceDelay  = 500 * time.Millisecond
)

// UI布局文件生成
func DebouncedGenerate(formTab *designer.FormTab) {
	debounceMutex.Lock()
	defer debounceMutex.Unlock()
	formName := formTab.Name
	// 取消之前的定时器
	if timer, exists := debounceTimers[formName]; exists {
		timer.Stop()
	}

	// 创建新的定时器
	timer := time.AfterFunc(debounceDelay, func() {
		debounceMutex.Lock()
		delete(debounceTimers, formName)
		debounceMutex.Unlock()

		formId := strings.ToLower(formName)
		uiFilePath := filepath.Join(project.Path, "forms", formId+".ui")

		// 执行UI生成
		err := GenerateUIFile(formTab.FormRoot, uiFilePath)
		if err != nil {
			logs.Error("UI布局文件生成错误:", err.Error())
		} else {
			// 触发代码生成
			triggerCodeGeneration(uiFilePath)
		}
	})

	debounceTimers[formName] = timer
}
