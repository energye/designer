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

package project

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

func init() {
	event.On(event.Project, func(trigger event.TTrigger) {
		logs.Debug("项目事件 Payload:", trigger.Payload)
		payload, ok := trigger.Payload.(event.TPayload)
		if ok {
			switch payload.Type {
			case event.ProjectCreate:

			}
		}
	}, func() {
		logs.Info("停止项目配置更新生成器")
	})
}
