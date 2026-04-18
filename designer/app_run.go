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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/locales"
	"strings"
	"time"
)

var (
	associateProtocol string
	associateFile     string
	universalLink     string
)

func Run() {
	logs.Println("TID:", api.MainThreadId(), "ENERGY Designer RUN")
	//locales.SwitchLCLLang("de")
	locales.SwitchLCLLang("zh_CN")
	lcl.Application.Initialize()
	lcl.Application.SetTitle(config.DesignerConfig.Title)
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&MainWindow)

	beforeRun()

	lcl.Application.Run()
	logs.Println("ENERGY Designer RUN END.")
}

// 是否停止加载关联项目
var isAssociateStopLoading = false

// 主动停止关联项目加载
func stopAutoAssociateProjectLoad() {
	logs.Info("stopAutoAssociateProjectLoad")
	isAssociateStopLoading = true
}

// 自动尝试关联项目加载
// 1. 尝试获取关联文件
// 2. 尝试从最后一次打开项目
func autoAssociateProjectLoad() {
	logs.Info("autoAssociateProjectLoad")
	isAssociateStopLoading = false
	var filePath string
	// 轮训获取 openFile
	for i := 0; i < 10; i++ {
		if isAssociateStopLoading {
			break
		}
		filePath = OpenFile()
		if filePath != "" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !isAssociateStopLoading {
		logs.Info("autoAssociateProjectLoad filePath:", filePath)
		loadProject(filePath)
	}
	isAssociateStopLoading = false
	logs.Info("autoAssociateProjectLoad end")
}

func loadProject(filePath string) {
	isEgp := strings.HasSuffix(filePath, consts.EGPExt)
	if isEgp {
		// 自动打开 energy 项目
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: filePath}})
	} else if config.Config.LastProject != "" && tool.IsExist(config.Config.LastProject) {
		// 自动打开 最后一次打开的项目
		event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: config.Config.LastProject}})
	}
}

func OpenFile() string {
	return associateFile
}

func OpenURL() string {
	return associateProtocol
}

func UniversalLink() string {
	return universalLink
}
