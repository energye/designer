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
	"os"
	"sync"
	"time"
)

type FileState struct {
	Path         string
	ModTime      time.Time
	IsDirty      bool
	LastSaveTime time.Time
}

type FileManager struct {
	files map[string]*FileState
	lock  sync.RWMutex
}

func NewFileManager() *FileManager {
	return &FileManager{
		files: make(map[string]*FileState),
	}
}

func (fm *FileManager) RegisterFile(filePath string) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	if _, exists := fm.files[filePath]; exists {
		return
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		return
	}

	fm.files[filePath] = &FileState{
		Path:         filePath,
		ModTime:      fi.ModTime(),
		IsDirty:      false,
		LastSaveTime: time.Now(),
	}
}

func (fm *FileManager) UnregisterFile(filePath string) {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	delete(fm.files, filePath)
}

func (fm *FileManager) SetDirty(filePath string, isDirty bool) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	if state, exists := fm.files[filePath]; exists {
		state.IsDirty = isDirty
		if !isDirty {
			state.LastSaveTime = time.Now()
		}
	}
}

func (fm *FileManager) UpdateModTime(filePath string, modTime time.Time) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	if state, exists := fm.files[filePath]; exists {
		state.ModTime = modTime
	}
}

func (fm *FileManager) GetFileState(filePath string) (*FileState, bool) {
	fm.lock.RLock()
	defer fm.lock.RUnlock()

	state, exists := fm.files[filePath]
	return state, exists
}

func (fm *FileManager) CheckExternalChanges() []string {
	fm.lock.RLock()
	defer fm.lock.RUnlock()

	var changed []string
	for filePath, state := range fm.files {
		fi, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		if !fi.ModTime().Equal(state.ModTime) {
			changed = append(changed, filePath)
		}
	}

	return changed
}

func (fm *FileManager) GetAllFiles() map[string]*FileState {
	fm.lock.RLock()
	defer fm.lock.RUnlock()

	result := make(map[string]*FileState)
	for k, v := range fm.files {
		result[k] = v
	}
	return result
}
