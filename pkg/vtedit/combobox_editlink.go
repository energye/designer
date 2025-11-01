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

// 下拉框

type TComboBoxEditLink struct {
	*TBaseEditLink
	combobox lcl.IComboBox
	bounds   types.TRect
	text     string
	stopping bool
}

func NewComboBoxEditLink(bindData *TEditNodeData) *TComboBoxEditLink {
	link := new(TComboBoxEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.CreateEdit()
	return link
}

func (m *TComboBoxEditLink) CreateEdit() {
	logs.Debug("TComboBoxEditLink CreateEdit")
	m.combobox = lcl.NewComboBox(nil)
	m.combobox.SetVisible(false)
	m.combobox.SetBorderStyle(types.BsSingle)
	m.combobox.SetAutoSize(false)
	m.combobox.SetDoubleBuffered(true)
	m.combobox.SetReadOnly(true)
	m.combobox.SetOnChange(func(sender lcl.IObject) {
		m.BindData.EditNodeData.Index = m.combobox.ItemIndex()
		m.BindData.EditNodeData.StringValue = m.combobox.Text()
		logs.Debug("TComboBoxEditLink OnChange index:", m.BindData.EditNodeData.Index, "text:", m.BindData.EditNodeData.StringValue)
		m.BindData.FormInspectorPropertyToComponentProperty()
	})
	m.combobox.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		logs.Debug("TComboBoxEditLink OnKeyDown key:", *key)
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})
	items := m.combobox.Items()
	for _, item := range m.BindData.EditNodeData.ComboBoxValue {
		items.Add(item.StringValue)
	}
}

func (m *TComboBoxEditLink) BeginEdit() bool {
	logs.Debug("TComboBoxEditLink BeginEdit")
	if !m.stopping {
		m.combobox.Show()
		m.combobox.SetFocus()
	}
	return true
}

func (m *TComboBoxEditLink) CancelEdit() bool {
	logs.Debug("TComboBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.combobox.Hide()
		if m.VTree != nil {
			m.VTree.CancelEditNode()
		}
	}
	return true
}

func (m *TComboBoxEditLink) EndEdit() bool {
	value := m.combobox.Text()
	logs.Debug("TComboBoxEditLink EndEdit", "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.EditNodeData.Index = m.combobox.ItemIndex()
		m.BindData.EditNodeData.StringValue = m.combobox.Text()
		m.combobox.Hide()
		if m.VTree != nil {
			m.VTree.EndEditNode()
		}
	}
	return true
}

func (m *TComboBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TComboBoxEditLink PrepareEdit")
	if m.combobox == nil || !m.combobox.IsValid() {
		m.CreateEdit()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	m.VTree.GetTextInfo(node, column, m.combobox.Font(), &m.bounds, &m.text)
	logs.Debug("  PrepareEdit GetTextInfo:", m.bounds, m.text)
	m.combobox.Font().SetColor(colors.ClWindowText)
	m.combobox.SetParent(m.VTree)
	m.combobox.HandleNeeded()
	m.combobox.SetText(m.text)
	return true
}

func (m *TComboBoxEditLink) GetBounds() types.TRect {
	return m.combobox.BoundsRect()
}

func (m *TComboBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TComboBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.combobox, msg)
}

func (m *TComboBoxEditLink) SetBounds(R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetWidth(columnRect.Width())
	m.combobox.SetBoundsRect(R)
	logs.Debug("TComboBoxEditLink SetBounds", R)
}

func (m *TComboBoxEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TComboBoxEditLink Destroy")
	m.combobox.Free()
}
