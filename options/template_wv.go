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

// main.go 文件代码模板
const runWVCodeTemplate = `// ==============================================================================
// Application startup portal
// ==============================================================================

package main

import (
	"embed"
	"github.com/energye/energy/v3/application"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/energy/v3/wv"
	"{{.Name}}/app"
	_ "{{.Name}}/resources"
)

//go:embed web
var web embed.FS

func main() {
	// Global Initialization
	wvApp := wv.Init()
	wvApp.SetOptions(application.Options{
		DefaultURL: "fs://energy/index.html",
	})
	// Local resource loading
	wvApp.SetLocalLoad(application.LocalLoad{
		Scheme:     "fs",
		Domain:     "energy",
		ResRootDir: "web",
		FS:         web,
	})
	// IPC
	ipc.On("counter:change", func(context ipc.IContext) {
		data := context.Data().([]any)
		context.Result(data[0])
	})
	wv.Run(app.Forms...)
}
`
