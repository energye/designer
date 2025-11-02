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

package preview

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

// 预览事件实例
var preview = &TPreview{trigger: make(chan event.TEventTrigger, 1), cancel: make(chan bool, 1)}

// TPreview 预览
type TPreview struct {
	trigger chan event.TEventTrigger // 触发生成事件
	cancel  chan bool                // 取消生成事件
}

// Start 启动项目配置文件更新
func (m *TPreview) Start() {
	for {
		select {
		case trigger := <-m.trigger:
			if trigger.GenType == event.GtProject {
				runPreview()
			}
		case <-m.cancel:
			logs.Info("停止预览事件处理器失败")
			return
		}
	}
}

func init() {
	event.Preview = event.NewEvent(preview.trigger, preview.cancel)
	go preview.Start()
}
