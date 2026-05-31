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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
)

// 创建设计器布局

var (
	toolBarHeight                = int32(30)
	switchAutoContentLayoutAlign = true // 一个开关, 用于自动布局
)

func init() {
	if tool.IsLinux {
		return
	}
	switchAutoContentLayoutAlign = false
}

// 初始化设计器布局 v2
func (m *TAppWindow) initDesignerLayoutV2() {
	// 初始化主窗口用到的组件
	m.initWindowComponent()
	// 主菜单
	m.initMainMenu()
	// 顶侧布局 - 工具栏
	m.toolBtnLayout = initToolBtnLayout(m)
	// 底部布局 - 左: 组件库, 左中: 项目查看, 中: 中间画布(自适应), 右: 属性, 下: 日志控制
	m.contentLayout = initBottomBox(m)
	// 设计器
	designer = m.contentLayout.initFromDesignerLayout()
}

// 初始化主窗口用到的组件
func (m *TAppWindow) initWindowComponent() {
	m.openDialog = lcl.NewOpenDialog(m)
	m.openDialog.SetName("OpenDialog")
	m.saveDialog = lcl.NewSaveDialog(m)
	m.saveDialog.SetName("SaveDialog")
	m.selectDirectoryDialog = lcl.NewSelectDirectoryDialog(m)
	m.selectDirectoryDialog.SetName("SelectDirectoryDialog")
}
