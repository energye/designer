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

// 生成回调事件

type CodeGenerationCallback func(uiFilePath string)

var (
	// Go代码生成回调
	onCodeGeneration CodeGenerationCallback
)

// 设置Go代码生成回调
func SetCodeGenerationCallback(callback CodeGenerationCallback) {
	onCodeGeneration = callback
}

// 触发Go代码生成事件
func triggerCodeGeneration(uiFilePath string) {
	if onCodeGeneration != nil {
		onCodeGeneration(uiFilePath)
	}
}
