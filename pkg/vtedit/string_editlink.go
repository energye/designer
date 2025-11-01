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

package vtedit

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
)

// 文本编辑框

type TStringEditLink struct {
	*TBaseEditLink
	edit      lcl.IEdit
	bounds    types.TRect
	alignment types.TAlignment
	stopping  bool
}

func NewStringEditLink(bindData *TEditNodeData) *TStringEditLink {
	link := new(TStringEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TStringEditLink) Create() {
	logs.Debug("TStringEditLink Create")
	m.edit = lcl.NewEdit(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetDoubleBuffered(true)
	m.edit.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		logs.Debug("TStringEditLink OnKeyDown key:", *key)
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})
}

func (m *TStringEditLink) SetReadOnly(v bool) {
	m.edit.SetReadOnly(v)
}

func (m *TStringEditLink) BeginEdit() bool {
	logs.Debug("TStringEditLink BeginEdit")
	if !m.stopping {
		m.edit.Show()
		m.edit.SelectAll()
		m.edit.SetFocus()
	}
	return true
}

func (m *TStringEditLink) CancelEdit() bool {
	logs.Debug("TStringEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.Hide()
		if m.VTree != nil {
			m.VTree.CancelEditNode()
		}
	}
	return true
}

func (m *TStringEditLink) EndEdit() bool {
	value := m.edit.Text()
	logs.Debug("TStringEditLink EndEdit", "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.EditNodeData.StringValue = value
		m.edit.Hide()
		if m.VTree != nil {
			m.VTree.EndEditNode()
		}
	}
	return true
}

func (m *TStringEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TStringEditLink PrepareEdit")
	if m.edit == nil || !m.edit.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	// 节点的初始大小、字体和文本。
	m.edit.Font().SetColor(colors.ClWindowText)
	m.edit.SetParent(m.VTree)
	m.edit.HandleNeeded()
	m.edit.SetText(m.BindData.EditNodeData.StringValue)
	if column <= -1 {
		m.edit.SetBiDiMode(m.VTree.BiDiMode())
		m.alignment = m.VTree.Alignment()
	} else {
		columns := m.VTree.Header().Columns()
		m.edit.SetBiDiMode(columns.ItemsWithColumnIndexToVirtualTreeColumn(column).BiDiMode())
		m.alignment = columns.ItemsWithColumnIndexToVirtualTreeColumn(column).Alignment()
	}

	if m.edit.BiDiMode() != types.BdLeftToRight {
		switch m.alignment {
		case types.TaLeftJustify:
			m.alignment = types.TaRightJustify
		case types.TaRightJustify:
			m.alignment = types.TaLeftJustify
		}
	}
	return true
}

func (m *TStringEditLink) GetBounds() types.TRect {
	return m.edit.BoundsRect()
}

func (m *TStringEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TStringEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TStringEditLink) SetBounds(R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)
	logs.Debug("TStringEditLink SetBounds", R)
}

func (m *TStringEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TStringEditLink Destroy")
	m.edit.Free()
}
