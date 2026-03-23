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

package dependmod

import (
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
)

const (
	lclPath    = "github.com/energye/lcl"
	cefPath    = "github.com/energye/cef"
	wvPath     = "github.com/energye/wv"
	energyPath = "github.com/energye/energy/v3"
)

// 根据 designer/resources/config.json 配置依赖模块下载模块
func downloadMod() {
	logs.Println("根据 designer/resources/config.json 配置依赖模块下载模块")
	formCfg := config.FormConfig
	dependencies := formCfg.Dependencies
	lclVer := dependencies[lclPath]
	cefVer := dependencies[cefPath]
	wvVer := dependencies[wvPath]
	energyVer := dependencies[energyPath]
	dependenciesInfo := fmt.Sprintf("Dependencies LCL: %s, CEF: %s, WV: %s, ENERGY: %s", lclVer, cefVer, wvVer, energyVer)
	event.Emit(event.TTrigger{Name: event.Console, Payload: event.TPayload{Type: event.ConsoleInfo, Data: dependenciesInfo}})
}
