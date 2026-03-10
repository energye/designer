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
	"github.com/energye/designer/cmd/build"
	"github.com/energye/designer/cmd/packager"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

// 构建打包
func configBuildPackage() {
	logs.Debug("构建配置-打包")
	event.ConsoleWriteClear()
	build.Run()
	packager.Run()
}
