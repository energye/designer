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

package codegen

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/project"
	"os"
	"path/filepath"
	"time"
)

var (
	ticker   *time.Ticker
	fileInfo *tool.HashMap[string, *TListenFileInfo]
)

type TListenFileInfo struct {
	Id      int       // 窗体 id
	Path    string    // 文件路径
	ModTime time.Time // 修改时间
	Exist   bool      // 是否存在
}

// 监听文件改变任务
func runListenFileChange() {
	if ticker != nil {
		logs.Debug("停止运行监听文件改变任务")
		ticker.Stop()
		ticker = nil
	}
	logs.Debug("启动监听文件改变任务")
	// 创建(重置)
	fileInfo = tool.NewHashMap[string, *TListenFileInfo]()
	ticker = time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			detectFileChange()
		}
	}()
}

// 检测文件
func detectFileChange() {
	projPath := project.Path()
	proj := project.Project()
	if proj == nil || projPath == "" {
		return
	}
	// 重置文件存在状态, 默认 false
	fileInfo.Iterate(func(key string, value *TListenFileInfo) bool {
		value.Exist = false
		return false
	})

	// 重新定位文件信息
	appCodePath := filepath.Join(projPath, proj.Package)
	for _, form := range project.Project().UIForms {
		userFile := filepath.Join(appCodePath, form.GOUserFile)
		if fi := fileInfo.Get(form.Name); fi != nil {
			fi.Exist = true
		} else {
			fileInfo.Add(form.Name, &TListenFileInfo{Id: form.Id, Path: userFile, ModTime: time.Now(), Exist: true})
		}
	}

	// 待删除不存在的文件列表
	var removeNotExistList []string
	// 验证文件状态
	fileInfo.Iterate(func(formName string, file *TListenFileInfo) bool {
		if !file.Exist {
			// 不存在的添加到待删除列表
			removeNotExistList = append(removeNotExistList, formName)
			return false
		}
		fi, err := os.Stat(file.Path)
		if err != nil {
			return false
		}
		// 使用修改时间判断
		if !fi.ModTime().Equal(file.ModTime) {
			file.ModTime = fi.ModTime()
			// 当文件改变后处理对应逻辑
			// 功能:
			//   获得窗体引用接收者的所有方法
			//   其它待添加...
			logs.Debug("listenFileChange 文件被修改 窗体名:", formName, "窗体ID", file.Id, "文件:", file.Path)
		}
		return false
	})

	for _, name := range removeNotExistList {
		logs.Debug("listenFileChange 删除不存在的文件:", name)
		fileInfo.Remove(name)
	}
}
