// Copyright 漏 yanghy. All Rights Reserved.
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
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
)

func runCreate() {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		form := NewCreateProjectForm()
		form.ShowModal()
		logs.Debug("runCreate NewCreateProjectForm show end")
		if !form.submitted || form.createOptions == nil {
			return
		}
		create := *form.createOptions
		// 第一个模态窗口已完全退出, 这里顺序打开 CEF 版本窗口;
		// 不在按钮事件回调里 ShowModal, 避免嵌套模态循环。
		if create.GUIRenderFramework == bean.GUIRenderFramework_CEF {
			chromiumForm := NewChromiumDirForm()
			// 新项目尚未创建, 不能把当前打开的项目传给安装流程, 否则会改写旧项目的 framework_version
			chromiumForm.Project = nil
			chromiumForm.ShowModal()
			// 用户未点确认或安装未完成, 中止创建
			if !chromiumForm.Confirmed || chromiumForm.Version == "" ||
				!config.Config.Chromium.IsCEFInstalled(chromiumForm.Version) {
				return
			}
			create.FrameworkVersion = chromiumForm.Version
		}
		runCreateProject(create)
	})
}
