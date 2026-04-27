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

type ContentLayoutInspector struct {
	searchEdit lcl.ITreeFilterEdit // 组件搜索框
	topBox     lcl.IPanel
	title      lcl.ILabel
	box        lcl.IPanel
}

func initContentLayoutInspector(owner *ContentLayout) *ContentLayoutInspector {
	m := &ContentLayoutInspector{}
	m.searchEdit = lcl.NewTreeFilterEdit(owner.inspectorPanel)
	m.searchEdit.SetTextHint("搜索属性")
	m.searchEdit.SetAlign(types.AlTop)
	m.searchEdit.SetAutoSelect(false)
	borderSpacing := m.searchEdit.BorderSpacing()
	borderSpacing.SetLeft(3)
	borderSpacing.SetRight(3)
	borderSpacing.SetTop(3)
	borderSpacing.SetBottom(3)
	m.searchEdit.SetParent(owner.inspectorPanel)

	m.topBox = lcl.NewPanel(owner.inspectorPanel)
	m.topBox.SetBorderStyleToBorderStyle(types.BsNone)
	m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	m.topBox.SetParent(owner.inspectorPanel)

	title := lcl.NewLabel(m.topBox)
	title.SetCaption("对象检查器")
	title.SetLeft(5)
	title.SetTop(5)
	font := title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	title.SetParent(m.topBox)

	m.box = lcl.NewPanel(owner.inspectorPanel)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlClient)
	m.box.SetWidth(owner.inspectorPanel.Width())
	m.box.SetParent(owner.inspectorPanel)

	return m
}
