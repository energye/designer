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

package project

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/project/bean"
)

// 项目文件 xxx.egp 配置文件, 存在项目根目录
// 功能说明
// 1. 创建项目
// 2. 加载项目
// 3. 更新项目信息
// 4. 更新项目的窗体信息
// 5. 更新其它配置和选项
// 6. UI 布局文件加载, 恢复到设计器
// 6.1 UI 布局同步更新: 新增/修改/删除

var (
	// 全局 Path 完整项目路径, 打开项目时设置. C:/YouProjectXxx/xxx.egp > C:/YouProjectXxx
	gPath string
	// 全局项目配置, 在创建或加载项目时设置
	gProject *bean.TProject
)

// SetGlobalProject 设置全局项目路径和项目对象
// path: 项目路径
// project: 项目对象指针
func SetGlobalProject(path string, project *bean.TProject) {
	gPath = path
	gProject = project
	if path == "" || project == nil {
		designer.SetEnableFuncComponent(false)
	} else {
		designer.SetEnableFuncComponent(true)
	}
}

// 返回当前项目路径
func Path() string {
	return gPath
}

// 返回当前项目对象
func Project() *bean.TProject {
	return gProject
}
