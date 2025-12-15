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

package options

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/tool/command"
	"path/filepath"
)

var cmdGoRoot string

func init() {
	event.On(event.Project, func(trigger event.TTrigger) {
		logs.Debug("项目管理事件 Payload:", trigger.Payload)
		payload, ok := trigger.Payload.(event.TPayload)
		if ok {
			switch payload.Type {
			case event.ProjectCreate:
				// 创建项目
				runCreate()
			case event.ProjectLoad:
				// 加载项目或UI
				dir := payload.Data.(string)
				runLoad(dir)
			case event.ProjectUpdateForm:
				formTab := payload.Data.(*designer.FormTab)
				runUpdate(formTab)
			case event.ProjectConfig:
				// 项目(应用)配置
				runAppConfig()
			case event.EnvConfig:
				// 项目(环境)配置
				runEnvConfig()
			case event.BuildConfig:
				// 项目(构建)配置
				runBuildConfig()
			}
		}
	}, func() {
		logs.Println("停止项目配置更新生成器")
	})
	// 获取Go Root目录
	go func() {
		result := false
		cmd := command.NewCMD()
		cmd.IsPrint = false
		cmd.HideWindow = true
		cmd.Console = func(data string, level command.Level) {
			if !result {
				logs.Debug(level, data)
				_, err := filepath.Abs(data)
				if err == nil {
					cmdGoRoot = data
				}
			}
			result = true
		}
		cmd.Command("go", "env", "GOROOT")
	}()
}
