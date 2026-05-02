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

func OpenFileInEditor(filePath string, readOnly ...bool) {
	if gCurrentEditor != nil {
		gCurrentEditor.OpenFile(filePath, readOnly...)
	}
}

func CloseFileInEditor(filePath string) {
	if gCurrentEditor != nil {
		gCurrentEditor.CloseFile(filePath)
	}
}

func SaveCurrentFile() {
	if gCurrentEditor != nil {
		gCurrentEditor.SaveCurrentFile()
	}
}

func GetAllOpenedFiles() map[string]*FileState {
	if gCurrentEditor == nil {
		return make(map[string]*FileState)
	}
	return gCurrentEditor.FileManager().GetAllFiles()
}

// NotifyFileChanged 通知gopls文件已被外部修改（如codegen/uigen），
// 使用 DidClose+DidOpen 强制 gopls 重新读取文件内容。
func NotifyFileChanged(filePath string) {
	notifyFileChanged(filePath)
}

// NotifyFilesChanged 批量通知gopls多个文件已被外部修改
func NotifyFilesChanged(filePaths []string) {
	for _, filePath := range filePaths {
		NotifyFileChanged(filePath)
	}
}

var gCurrentEditor IEditor

func SetCurrentEditor(ed IEditor) {
	gCurrentEditor = ed
}

func GetCurrentEditor() IEditor {
	return gCurrentEditor
}
