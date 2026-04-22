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
	"github.com/energye/lcl/types/colors"
)

// 设计 - 组件属性

type InspectorComponentProperty struct {
	propBox           lcl.IPanel           // 组件属性盒子
	propFilterBox     lcl.IPanel           // 组件属性过滤盒子
	propComponentProp lcl.IPanel           // 组件属性
	filter            lcl.ITreeFilterEdit  // 组件属性过滤框
	currentComponent  *TDesigningComponent // 当前正在设计的组件
}

func (m *InspectorComponentProperty) init(leftBoxWidth int32) {
	m.propFilterBox = lcl.NewPanel(m.propBox)
	m.propFilterBox.SetBevelOuter(types.BvNone)
	m.propFilterBox.SetDoubleBuffered(true)
	m.propFilterBox.SetAlign(types.AlTop)
	if tool.IsLinux() {
		// Linux 编辑框高度差异
		m.propFilterBox.SetHeight(45)
	} else {
		m.propFilterBox.SetHeight(35)
	}
	m.propFilterBox.SetParent(m.propBox)

	m.filter = lcl.NewTreeFilterEdit(m.propFilterBox)
	m.filter.SetTop(3)
	m.filter.SetLeft(3)
	m.filter.SetWidth(m.propFilterBox.Width() - (m.filter.Left() + 3))
	m.filter.SetAlign(types.AlCustom)
	m.filter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	m.filter.SetTextHint("搜索属性")
	m.filter.SetParent(m.propFilterBox)

	m.propComponentProp = lcl.NewPanel(m.propBox)
	m.propComponentProp.SetBevelOuter(types.BvNone)
	m.propComponentProp.SetDoubleBuffered(true)
	m.propComponentProp.SetAlign(types.AlClient)
	m.propComponentProp.SetParent(m.propBox)
}

// 属性列表虚拟树配置方法
func vstConfig(tree lcl.ILazVirtualStringTree) {
	tree.SetBorderStyleToBorderStyle(types.BsNone)
	tree.SetAlign(types.AlClient)
	tree.SetLineStyle(types.LsSolid)
	SetComponentDefaultColor(tree)
	tree.SetDefaultNodeHeight(28)
	tree.SetIndent(8)
	//tree.Font().SetSize(8)
	tree.ScrollBarOptions().SetScrollBars(types.SsVertical)

	// options
	propTreeOptions := tree.TreeOptions()
	propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Exclude(types.ToShowTreeLines))
	propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Include(types.ToShowVertGridLines, types.ToShowHorzGridLines))
	propTreeOptions.SetSelectionOptions(propTreeOptions.SelectionOptions().Include(types.ToFullRowSelect, types.ToLevelSelectConstraint))
	propTreeOptions.SetMiscOptions(propTreeOptions.MiscOptions().Include(types.ToEditable, types.ToEditOnClick, types.ToEditOnDblClick))

	// 颜色
	propColors := tree.Colors()
	propColors.SetFocusedSelectionColor(colors.RGBToColor(43, 169, 241))
	propColors.SetUnfocusedSelectionColor(colors.RGBToColor(43, 169, 241))

	// header
	header := tree.Header()
	header.SetOptions(header.Options().Include(types.HoVisible, types.HoAutoSpring)) //types.HoAutoResize
	header.Font().SetStyle(header.Font().Style().Include(types.FsBold))
	columns := header.Columns()
	columns.Clear()
	propNameCol := columns.AddToVirtualTreeColumn()
	propNameCol.SetText("Name")
	propNameCol.SetAlignment(types.TaLeftJustify)
	propNameCol.SetWidth(125)
	propNameCol.SetMinWidth(50)
	//propNameCol.SetOptions(propNameCol.Options().Include(types.CoDisableAnimatedResize))

	propValueCol := columns.AddToVirtualTreeColumn()
	propValueCol.SetText("Value")
	propValueCol.SetAlignment(types.TaLeftJustify)
	propValueCol.SetOptions(propValueCol.Options().Include(types.CoAutoSpring))
	if tool.IsLinux() || tool.IsDarwin() {
		width := int32(135) // tree.Width() - 65
		propValueCol.SetWidth(width)
		propValueCol.SetMinWidth(50)
	}
}
