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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

// Go代码生成

func init() {
	event.On(event.GenCode, func(trigger event.TTrigger) {
		// 触发代码生成事件
		if uiFilePath, ok := trigger.Payload.(string); ok {
			err := runGenerateCode(uiFilePath)
			if err != nil {
				logs.Error("代码生成错误:", err.Error())
			}
		}
	}, func() {
		logs.Info("停止代码生成器")
	})
}
