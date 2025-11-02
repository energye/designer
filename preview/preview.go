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
	"github.com/energye/designer/project"
	"github.com/energye/lcl/tool/command"
)

// 执行应用程序的预览功能
// 根据项目配置预览当前项目
func runPreview() {
	cmd := command.NewCMD()
	cmd.Dir = project.Path
	cmd.MessageCallback = func(bytes []byte, err error) {

	}
	cmd.Command("go", "run", "main.go")
}
