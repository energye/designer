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

package editor

import (
	"github.com/energye/designer/options/bean"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// formFileWatcher 监控项目所有窗体关联文件的变更
// FileManager 只跟踪 Monaco 中打开的文件，而 ui.go 文件通常不在 Monaco 中打开，
// 所以 formFileWatcher 独立扫描所有窗体文件，确保 gopls 感知外部修改。
//
// 职责划分:
//   - checkFileChanges: 检测 FileManager 注册的文件 → 通知 Monaco + gopls
//   - formFileWatcher: 检测项目窗体文件 → 仅通知 gopls（Monaco 未打开这些文件）
type formFileWatcher struct {
	mu       sync.RWMutex
	fileInfo map[string]time.Time // filePath -> ModTime
	stopCh   chan struct{}
}

var gFormWatcher *formFileWatcher

func startFormFileWatcher() {
	if gFormWatcher != nil {
		return
	}
	gFormWatcher = &formFileWatcher{
		fileInfo: make(map[string]time.Time),
		stopCh:   make(chan struct{}),
	}
	gFormWatcher.scanAndRecord()
	go gFormWatcher.run()
}

func stopFormFileWatcher() {
	if gFormWatcher == nil {
		return
	}
	close(gFormWatcher.stopCh)
	gFormWatcher = nil
}

// UpdateSavedModTime 编辑器保存文件后调用，同步 ModTime 防止误判
func updateSavedModTime(filePath string) {
	if gFormWatcher == nil {
		return
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return
	}
	gFormWatcher.mu.Lock()
	gFormWatcher.fileInfo[filePath] = fi.ModTime()
	gFormWatcher.mu.Unlock()
}

func (w *formFileWatcher) run() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.checkChanges()
		}
	}
}

func (w *formFileWatcher) scanAndRecord() {
	codePath := bean.CodePath()
	if codePath == "" || bean.GProject == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, form := range bean.GProject.UIForms {
		for _, fileName := range []string{form.GOFile, form.GOUserFile} {
			if fileName == "" {
				continue
			}
			fp := filepath.Join(codePath, fileName)
			fi, err := os.Stat(fp)
			if err != nil {
				continue
			}
			w.fileInfo[fp] = fi.ModTime()
		}
	}
}

func (w *formFileWatcher) checkChanges() {
	codePath := bean.CodePath()
	if codePath == "" || bean.GProject == nil {
		return
	}

	var changedFiles []string
	currentFiles := make(map[string]time.Time)

	w.mu.Lock()
	for _, form := range bean.GProject.UIForms {
		for _, fileName := range []string{form.GOFile, form.GOUserFile} {
			if fileName == "" {
				continue
			}
			fp := filepath.Join(codePath, fileName)
			fi, err := os.Stat(fp)
			if err != nil {
				continue
			}
			currentFiles[fp] = fi.ModTime()

			oldModTime, existed := w.fileInfo[fp]
			if !existed || !fi.ModTime().Equal(oldModTime) {
				changedFiles = append(changedFiles, fp)
			}
		}
	}
	for fp, modTime := range currentFiles {
		w.fileInfo[fp] = modTime
	}
	for fp := range w.fileInfo {
		if _, exists := currentFiles[fp]; !exists {
			delete(w.fileInfo, fp)
		}
	}
	w.mu.Unlock()

	for _, filePath := range changedFiles {
		notifyFileChanged(filePath)
	}
}
