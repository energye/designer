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

package build

import (
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"os"
	"runtime"
)

// Run 执行构建命令的入口函数
func Run() bool {
	event.ConsoleWriteInfo("CMD-build-run")
	return build()
}

func xBuildPackVar() (result []string) {
	if bean.GProject == nil {
		return nil
	}
	GOARCH := os.Getenv("GOARCH")
	if GOARCH == "" {
		GOARCH = runtime.GOARCH
	}
	identity := fmt.Sprintf("%s_%s_%s", bean.GProject.Name, bean.GProject.AppOption.Id, GOARCH)
	//packMap := make(map[string]string)
	//packMap["name"] = bean.GProject.Name
	//packMap["id"] = bean.GProject.AppOption.Id
	//packMap["version"] = bean.GProject.AppOption.Version
	//packMap["arch"] = GOARCH
	//data, _ := json.Marshal(packMap)
	result = append(result, "-X github.com/energye/energy/v3/application/pack.JSON="+identity)
	return
}
