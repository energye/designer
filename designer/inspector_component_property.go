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

// 设计 - 组件属性

// 属性列表虚拟树配置方法
func newVirtualStringTree(owner lcl.IWinControl) lcl.ILazVirtualStringTree {
	tree := lcl.NewLazVirtualStringTree(owner)
	tree.SetBorderStyleToBorderStyle(types.BsNone)
	tree.SetAlign(types.AlClient)
	tree.SetLineStyle(types.LsSolid)
	SetComponentDefaultColor(tree)
	tree.SetDefaultNodeHeight(28)
	tree.SetIndent(8)
	//tree.Font().SetSize(8)
	//tree.ScrollBarOptions().SetScrollBars(types.SsBoth)
	tree.ScrollBarOptions().SetScrollBars(types.SsVertical)

	// options
	propTreeOptions := tree.TreeOptions()
	propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Exclude(types.ToShowTreeLines))
	propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Include(types.ToShowVertGridLines, types.ToShowHorzGridLines))
	propTreeOptions.SetSelectionOptions(propTreeOptions.SelectionOptions().Include(types.ToFullRowSelect, types.ToLevelSelectConstraint))
	propTreeOptions.SetMiscOptions(propTreeOptions.MiscOptions().Include(types.ToEditable, types.ToEditOnClick, types.ToEditOnDblClick))

	// 颜色
	propColors := tree.Colors()
	propColors.SetFocusedSelectionColor(0xF9D99E)
	propColors.SetUnfocusedSelectionColor(0xF9D99E)

	// header
	header := tree.Header()
	header.SetOptions(header.Options().Include(types.HoVisible, types.HoAutoSpring)) //types.HoAutoResize
	header.Font().SetStyle(header.Font().Style().Include(types.FsBold))
	columns := header.Columns()
	columns.Clear()
	propNameCol := columns.AddToVirtualTreeColumn()
	propNameCol.SetText("Name")
	propNameCol.SetAlignment(types.TaLeftJustify)

	propValueCol := columns.AddToVirtualTreeColumn()
	propValueCol.SetText("Value")
	propValueCol.SetAlignment(types.TaLeftJustify)
	propValueCol.SetOptions(propValueCol.Options().Include(types.CoAutoSpring))

	return tree
}
