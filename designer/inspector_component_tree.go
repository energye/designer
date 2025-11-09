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
	"github.com/energye/lcl/tool"
	"github.com/energye/lcl/types"
)

// 设计 - 组件树

var (
	gTreeId int // 维护组件树全局数据id
)

// 获取下一个树数据ID
func nextTreeDataId() (id int) {
	id = gTreeId
	gTreeId++
	return
}

// 查看器组件树
type InspectorComponentTree struct {
	treeBox           lcl.IPanel          // 组件树盒子
	treeFilterBox     lcl.IPanel          // 组件树过滤盒子
	treeComponentTree lcl.IPanel          // 组件树
	treeFilter        lcl.ITreeFilterEdit // 组件树过滤框
}

func (m *InspectorComponentTree) init(leftBoxWidth int32) {
	m.treeFilterBox = lcl.NewPanel(m.treeBox)
	m.treeFilterBox.SetBevelOuter(types.BvNone)
	m.treeFilterBox.SetDoubleBuffered(true)
	m.treeFilterBox.SetAlign(types.AlTop)
	if tool.IsLinux() {
		// Linux 编辑框高度差异
		m.treeFilterBox.SetHeight(45)
	} else {
		m.treeFilterBox.SetHeight(35)
	}
	m.treeFilterBox.SetParent(m.treeBox)

	componentTreeTitle := lcl.NewLabel(m.treeFilterBox)
	componentTreeTitle.SetCaption("组件")
	componentTreeTitle.Font().SetStyle(types.NewSet(types.FsBold))
	componentTreeTitle.SetTop(8)
	componentTreeTitle.SetLeft(5)
	componentTreeTitle.SetParent(m.treeFilterBox)

	m.treeFilter = lcl.NewTreeFilterEdit(m.treeFilterBox)
	m.treeFilter.SetTop(3)
	if tool.IsLinux() {
		m.treeFilter.SetLeft(40)
	} else {
		m.treeFilter.SetLeft(30)
	}
	m.treeFilter.SetWidth(m.treeFilterBox.Width() - m.treeFilter.Left())
	m.treeFilter.SetAlign(types.AlCustom)
	m.treeFilter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	m.treeFilter.SetParent(m.treeFilterBox)

	m.treeComponentTree = lcl.NewPanel(m.treeBox)
	m.treeComponentTree.SetBevelOuter(types.BvNone)
	m.treeComponentTree.SetDoubleBuffered(true)
	m.treeComponentTree.SetAlign(types.AlClient)
	m.treeComponentTree.SetParent(m.treeBox)

}
