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
	"github.com/energye/energy/v3/ipc"
)

func OpenFileInEditor(filePath string) {
	ipc.Emit("open-file", filePath)
}

func CloseFileInEditor(filePath string) {
	ipc.Emit("close-file", filePath)
}

func SaveCurrentFile() {
	ipc.Emit("save-current-file", "")
}

func GetOpenedFiles() []string {
	var files []string

	type FileData struct {
		File    string `json:"file"`
		Content string `json:"content"`
	}

	for filePath := range GetAllOpenedFiles() {
		files = append(files, filePath)
	}

	return files
}

func GetAllOpenedFiles() map[string]*FileState {
	if gCurrentEditor == nil {
		return make(map[string]*FileState)
	}
	return gCurrentEditor.fileManager.GetAllFiles()
}

var gCurrentEditor *TWebviewEditor

func SetCurrentEditor(editor *TWebviewEditor) {
	gCurrentEditor = editor
}

func GetCurrentEditor() *TWebviewEditor {
	return gCurrentEditor
}
