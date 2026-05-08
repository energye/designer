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

type formFileWatcher struct {
	mu       sync.RWMutex
	fileInfo map[string]time.Time
	stopCh   chan struct{}
}

var (
	gFormWatcher *formFileWatcher
	watcherMu    sync.RWMutex
)

func StartFormFileWatcher() {
	watcherMu.RLock()
	if gFormWatcher != nil {
		watcherMu.RUnlock()
		return
	}
	watcherMu.RUnlock()

	watcherMu.Lock()
	defer watcherMu.Unlock()
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

func StopFormFileWatcher() {
	watcherMu.Lock()
	defer watcherMu.Unlock()
	if gFormWatcher == nil {
		return
	}
	close(gFormWatcher.stopCh)
	gFormWatcher = nil
}

// UpdateSavedModTime 编辑器保存文件后调用，同步 ModTime 防止误判
func UpdateSavedModTime(filePath string) {
	watcherMu.RLock()
	w := gFormWatcher
	watcherMu.RUnlock()
	if w == nil {
		return
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return
	}
	w.mu.Lock()
	w.fileInfo[filePath] = fi.ModTime()
	w.mu.Unlock()
}

func (w *formFileWatcher) run() {
	ticker := time.NewTicker(time.Second)
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

	var filesToStat []string
	for _, form := range bean.GProject.UIForms {
		for _, fileName := range []string{form.GOFile, form.GOUserFile} {
			if fileName == "" {
				continue
			}
			filesToStat = append(filesToStat, filepath.Join(codePath, fileName))
		}
	}

	currentFiles := make(map[string]time.Time)
	for _, fp := range filesToStat {
		fi, err := os.Stat(fp)
		if err != nil {
			continue
		}
		currentFiles[fp] = fi.ModTime()
	}

	w.mu.Lock()
	var changedFiles []string
	for fp, modTime := range currentFiles {
		oldModTime, existed := w.fileInfo[fp]
		if !existed || !modTime.Equal(oldModTime) {
			changedFiles = append(changedFiles, fp)
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
