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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 布局下 设计器

var (
	leftBoxWidth int32 = 290
)

type BottomBox struct {
	box      lcl.IPanel
	leftBox  lcl.IPanel    // 左侧-面板组件对象查看器
	splitter lcl.ISplitter // 分割线
	rightBox lcl.IPanel    // 右侧-窗体设计器
	console  *TConsole     // 底部输出
}

func (m *TAppWindow) createBottomBox() *BottomBox {
	box := &BottomBox{}
	box.box = lcl.NewPanel(m)
	box.box.SetBevelOuter(types.BvNone)
	box.box.SetDoubleBuffered(true)
	box.box.SetTop(toolbarHeight)
	box.box.SetWidth(m.Width())
	box.box.SetHeight(m.Height() - box.box.Top())
	box.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	SetComponentDefaultColor(box.box)
	box.box.SetParent(m)
	//box.box.SetColor(bottomColor)
	m.box = box

	// 工具栏-分隔线
	box.splitter = lcl.NewSplitter(box.box)
	box.splitter.SetAlign(types.AlLeft)
	box.splitter.SetWidth(splitterWidth)
	box.splitter.SetMinSize(50)
	box.splitter.SetResizeStyle(types.RsNone)
	box.splitter.SetParent(box.box)

	// 左侧-面板组件对象查看器
	box.leftBox = lcl.NewPanel(box.box)
	box.leftBox.SetBevelOuter(types.BvNone)
	box.leftBox.SetDoubleBuffered(true)
	box.leftBox.SetWidth(leftBoxWidth)
	box.leftBox.SetHeight(box.box.Height())
	box.leftBox.Constraints().SetMinWidth(50)
	box.leftBox.SetAlign(types.AlLeft)
	SetComponentDefaultColor(box.leftBox)
	box.leftBox.SetParent(box.box)

	// 右侧-窗体设计器
	box.rightBox = lcl.NewPanel(box.box)
	box.rightBox.SetBevelOuter(types.BvNone)
	box.rightBox.SetDoubleBuffered(true)
	box.rightBox.SetAlign(types.AlClient)
	SetComponentDefaultColor(box.rightBox)
	box.rightBox.SetParent(box.box)

	// 创建对象查看器
	inspector = box.createInspectorLayout()

	// 创建窗体设计器
	designer = box.createFromDesignerLayout()
	box.createConsole()

	AddOnShow(func() {
		//// 2. 添加到组件树
		//// 2.1. 显示之后创建一个默认的设计面板
		//defaultForm := designer.addDesignerFormTab()
		//// 2.1. 加载属性到设计器
		//// 此步骤会初始化并填充设计组件实例
		//defaultForm.FormRoot.LoadPropertyToInspector()
		//// 2.2. 添加到组件树
		//defaultForm.AddFormNode()
		//
		//go triggerUIGeneration(defaultForm.FormRoot)
	})

	return box
}
