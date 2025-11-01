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
	"github.com/energye/lcl/types/colors"
)

// 设计 - 组件属性

type InspectorComponentProperty struct {
	box              lcl.IPanel           // 组件属性盒子
	filter           lcl.ITreeFilterEdit  // 组件属性过滤框
	currentComponent *TDesigningComponent // 当前正在设计的组件
}

func (m *InspectorComponentProperty) init(leftBoxWidth int32) {
	componentPropertyTitle := lcl.NewLabel(m.box)
	componentPropertyTitle.SetParent(m.box)
	componentPropertyTitle.SetCaption("属性")
	componentPropertyTitle.Font().SetStyle(types.NewSet(types.FsBold))

	componentPropertyTitle.SetTop(5)
	componentPropertyTitle.SetLeft(5)

	m.filter = lcl.NewTreeFilterEdit(m.box)
	m.filter.SetParent(m.box)
	m.filter.SetTop(2)
	m.filter.SetLeft(30)
	m.filter.SetWidth(leftBoxWidth - m.filter.Left())
	m.filter.SetAlign(types.AlCustom)
	m.filter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
}

// 属性列表虚拟树配置方法
func vstConfig(tree lcl.ILazVirtualStringTree) {
	tree.SetBorderStyleToBorderStyle(types.BsNone)
	tree.SetAlign(types.AlClient)
	tree.SetLineStyle(types.LsSolid)
	tree.SetDefaultNodeHeight(28)
	tree.SetIndent(8)

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
	header.Font().SetColor(colors.ClGray)
	columns := header.Columns()
	columns.Clear()
	propNameCol := columns.AddToVirtualTreeColumn()
	propNameCol.SetText("名")
	propNameCol.SetWidth(125)
	propNameCol.SetAlignment(types.TaLeftJustify)
	//propNameCol.SetOptions(propNameCol.Options().Include(types.CoDisableAnimatedResize))

	propValueCol := columns.AddToVirtualTreeColumn()
	propValueCol.SetText("值")
	//propValueCol.SetWidth(leftBoxWidth - 125)
	propValueCol.SetAlignment(types.TaLeftJustify)
	propValueCol.SetOptions(propValueCol.Options().Include(types.CoAutoSpring))
}
