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
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/locales"
)

func Run() {
	logs.Info("ENERGY Designer RUN")
	//locales.SwitchLCLLang("de")
	locales.SwitchLCLLang("zh_CN")
	lcl.Application.Initialize()
	lcl.Application.SetTitle(config.Config.Title)
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&mainWindow)
	lcl.Application.Run()
	logs.Info("ENERGY Designer RUN END.")
}
