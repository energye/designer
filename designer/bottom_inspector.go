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
	"widget/wg"
)

// 对象查看器

var (
	componentTreeHeight int32 = 222
)

func (m *BottomBox) createInspectorLayout() *Inspector {
	ins := new(Inspector)
	// 面板 对象查看器分隔
	{
		ins.boxSplitter = lcl.NewSplitter(m.leftBox)
		ins.boxSplitter.SetParent(m.leftBox)
		ins.boxSplitter.SetAlign(types.AlTop)
		ins.boxSplitter.SetResizeStyle(types.RsNone)

		tree := new(InspectorComponentTree)
		tree.treeBox = lcl.NewPanel(m.leftBox)
		tree.treeBox.SetBevelColor(wg.LightenColor(bgDarkColor, 0.3))
		tree.treeBox.SetBevelOuter(types.BvLowered)
		tree.treeBox.SetDoubleBuffered(true)
		tree.treeBox.SetWidth(m.leftBox.Width())
		tree.treeBox.SetHeight(componentTreeHeight)
		tree.treeBox.Constraints().SetMinWidth(50)
		tree.treeBox.Constraints().SetMinHeight(50)
		tree.treeBox.SetAlign(types.AlTop)
		SetComponentDefaultColor(tree.treeBox)
		tree.treeBox.SetParent(m.leftBox)
		ins.componentTree = tree

		property := new(InspectorComponentProperty)
		property.propBox = lcl.NewPanel(m.leftBox)
		property.propBox.SetBevelColor(wg.LightenColor(bgDarkColor, 0.3))
		property.propBox.SetBevelOuter(types.BvLowered)
		property.propBox.SetDoubleBuffered(true)
		property.propBox.SetWidth(m.leftBox.Width())
		property.propBox.Constraints().SetMinWidth(50)
		property.propBox.Constraints().SetMinHeight(50)
		property.propBox.SetAlign(types.AlClient)
		property.propBox.SetParent(m.leftBox)
		ins.componentProperty = property
		//ins.componentPropertyBox.SetColor(colors.Cl3DDkShadow)
	}
	// 组件树
	{
		ins.componentTree.init(m.leftBox.Width())
	}
	// 组件属性
	{
		ins.componentProperty.init(m.leftBox.Width())
	}
	return ins
}
