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
	"github.com/energye/lcl/types"
)

// 布局下 设计器

var (
	leftBoxWidth int32 = 290
)

// BottomBox 底部布局
//
//	左: 组件库
//	左中: 项目查看
//	中: 中间画布(自适应)
//	右: 属性
//	下: 日志控制

type BottomLayout struct {
	box      lcl.IPanel
	leftBox  lcl.IPanel    // 左侧-面板组件对象查看器
	splitter lcl.ISplitter // 分割线
	rightBox lcl.IPanel    // 右侧-窗体设计器
	console  *TConsole     // 底部输出
}

func (m *TAppWindow) initBottomBox() *BottomLayout {
	bottomLayout := &BottomLayout{}
	bottomLayout.box = lcl.NewPanel(m)
	bottomLayout.box.SetBevelOuter(types.BvNone)
	bottomLayout.box.SetDoubleBuffered(true)
	bottomLayout.box.SetTop(toolbarHeight)
	bottomLayout.box.SetWidth(m.Width())
	bottomLayout.box.SetHeight(m.Height() - bottomLayout.box.Top())
	bottomLayout.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	SetComponentDefaultColor(bottomLayout.box)
	bottomLayout.box.SetParent(m)
	//box.box.SetColor(bottomColor)
	m.bottomLayout = bottomLayout

	// 工具栏-分隔线
	bottomLayout.splitter = lcl.NewSplitter(bottomLayout.box)
	bottomLayout.splitter.SetAlign(types.AlLeft)
	bottomLayout.splitter.SetWidth(splitterWidth)
	bottomLayout.splitter.SetMinSize(50)
	if tool.IsWindows {
		bottomLayout.splitter.SetResizeStyle(types.RsNone)
	}
	bottomLayout.splitter.SetParent(bottomLayout.box)

	// 左侧-面板组件对象查看器
	bottomLayout.leftBox = lcl.NewPanel(bottomLayout.box)
	bottomLayout.leftBox.SetBevelOuter(types.BvNone)
	bottomLayout.leftBox.SetDoubleBuffered(true)
	bottomLayout.leftBox.SetWidth(leftBoxWidth)
	bottomLayout.leftBox.SetHeight(bottomLayout.box.Height())
	bottomLayout.leftBox.Constraints().SetMinWidth(255)
	bottomLayout.leftBox.SetAlign(types.AlLeft)
	SetComponentDefaultColor(bottomLayout.leftBox)
	bottomLayout.leftBox.SetParent(bottomLayout.box)

	// 右侧-窗体设计器
	bottomLayout.rightBox = lcl.NewPanel(bottomLayout.box)
	bottomLayout.rightBox.SetBevelOuter(types.BvNone)
	bottomLayout.rightBox.SetDoubleBuffered(true)
	bottomLayout.rightBox.SetAlign(types.AlClient)
	SetComponentDefaultColor(bottomLayout.rightBox)
	bottomLayout.rightBox.SetParent(bottomLayout.box)

	// 创建对象查看器
	inspector = bottomLayout.createInspectorLayout()

	// 创建窗体设计器
	designer = bottomLayout.createFromDesignerLayout()

	// 创建消息控制台
	bottomLayout.createConsole()

	return bottomLayout
}
