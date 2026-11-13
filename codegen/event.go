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
	"sync"

	"github.com/energye/designer/designer"
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

// Go代码生成

// generateMu 串行化代码生成：事件队列、构建前 flush 可能并发触发，
// 同时写 xxx.ui.go / xxx.go 会产生半成品文件。
var generateMu sync.Mutex

// Generate 同步执行一次代码生成。生成会整文件覆盖 .go 源码，
// 编辑器有未保存内容时先挂起，保存完成后自动重放。
func Generate(uiGenData designer.TUIGenerationData, type_ event.Type) {
	if !editor.AllowGenerate(func() {
		Generate(uiGenData, type_)
	}) {
		event.ConsoleWriteWarn("代码编辑器有未保存内容：已先自动保存，保存完成后继续生成代码")
		return
	}

	generateMu.Lock()
	defer generateMu.Unlock()
	switch type_ {
	case event.CodeGenUI:
		// 根据UI布局文件生成代码
		if err := runGenerateCode(uiGenData); err != nil {
			logs.Error("代码生成错误-UI布局:", err.Error())
		}
	case event.CodeGenEvent:
		// 更新用户代码-绑定事件
		if err := runUpdateEvent(uiGenData); err != nil {
			logs.Error("代码生成错误-绑定事件:", err.Error())
		}
	case event.CodeGenSelf:
		// 更新用户代码-自引用（当窗体名被修改时）
		if err := runUpdateSelf(uiGenData); err != nil {
			logs.Error("代码生成错误-自引用:", err.Error())
		}
	}
}

func Init() {
	// 代码生成事件监听
	event.On(event.GenCode, func(trigger event.TTrigger) {
		if payload, ok := trigger.Payload.(event.TPayload); ok {
			uiGenData, ok := payload.Data.(designer.TUIGenerationData)
			if ok {
				Generate(uiGenData, payload.Type)
			}
		}
	}, func() {
		logs.Println("停止代码生成器")
	})

	// 文件改变监听
	// 运行时机, 每次都会重置监斩
	//   打开项目后
	//   新建项目后
	event.On(event.ListenFileChange, func(trigger event.TTrigger) {
		runListenFileChange()
	}, func() {
		logs.Println("停止文件改变监听")
		stopListenFileChange()
	})
}
