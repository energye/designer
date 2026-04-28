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

package designer

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
)

// 项目源码

var gProjectSrcTree = &TProjectSrcTree{}

type TProjectSrcTree struct {
}

func (m *TProjectSrcTree) scanProjectSrc() {
	//bean.GPath
}

func initProjectSrcEvent() {
	logs.Println("启动项目 SRC Tree 监听")
	event.On(event.ListenProjectSrcFileChange, func(trigger event.TTrigger) {
		payload, ok := trigger.Payload.(event.TPayload)
		if ok {
			switch payload.Type {
			case event.ProjectSrcScan:
				gProjectSrcTree.scanProjectSrc()
			}
			if tool.IsMainThread() {
			} else {
				lcl.RunOnMainThreadAsync(func(id uint32) {
				})
			}
		}
	}, func() {
		logs.Println("停止项目 SRC Tree 监听")
	})
}
