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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
)

// 事件选择方法下拉框

const defaultText = "(none)"

type TEventComboBoxEditLink struct {
	*TBaseEditLink
	combobox lcl.IComboBox
	bounds   types.TRect
	stopping bool
}

func NewEventComboBoxEditLink(bindData *TEditNodeData) *TEventComboBoxEditLink {
	link := new(TEventComboBoxEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.CreateEdit()
	return link
}

func (m *TEventComboBoxEditLink) CreateEdit() {
	logs.Debug("TEventComboBoxEditLink CreateEdit")
	m.combobox = lcl.NewComboBox(nil)
	m.combobox.SetVisible(false)
	m.combobox.SetBorderStyle(types.BsSingle)
	m.combobox.SetAutoSize(false)
	m.combobox.SetDoubleBuffered(true)
	m.combobox.SetOnChange(func(sender lcl.IObject) {
		text := m.combobox.Text()
		if text == defaultText {
			m.combobox.SetText("")
			m.combobox.SetItemIndex(-1)
			return
		}
	})
	m.combobox.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		logs.Debug("TEventComboBoxEditLink OnKeyDown key:", *key)
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})
	items := m.combobox.Items()

	comboBoxValue := []*TEditLinkNodeData{
		{StringValue: defaultText},
	}
	methods := m.BindData.AffiliatedComponent.GetRecvMethods()
	if methods != nil {
		// TODO 过滤匹配出符合的事件类型方法
		for _, method := range methods {
			comboBoxValue = append(comboBoxValue, &TEditLinkNodeData{StringValue: method.Name})
		}
	}
	m.BindData.EditNodeData.ComboBoxValue = comboBoxValue
	for _, item := range m.BindData.EditNodeData.ComboBoxValue {
		items.Add(item.StringValue)
	}
}

func (m *TEventComboBoxEditLink) SetValue(index int32, value string) {
	if index == -1 && value == "" {
		// 删除
		m.BindData.EditNodeData.Index = -1
		m.BindData.EditNodeData.EventState = consts.EsDelete
		return
	}
	if index == -1 && value != "" {
		// -1 新增
		m.BindData.EditNodeData.ComboBoxValue = append(m.BindData.EditNodeData.ComboBoxValue, &TEditLinkNodeData{StringValue: value})
		m.BindData.EditNodeData.Index = int32(len(m.BindData.EditNodeData.ComboBoxValue)) - 1
		m.BindData.EditNodeData.EventState = consts.EsAdd
	} else if index == 0 || value == "" {
		// 0 (none) 删除
		m.BindData.EditNodeData.Index = -1
		m.BindData.EditNodeData.EventState = consts.EsDelete
	} else {
		// 1-n 更新
		m.BindData.EditNodeData.Index = index
		m.BindData.EditNodeData.EventState = consts.EsUpdate
	}
}

func (m *TEventComboBoxEditLink) BeginEdit() bool {
	logs.Debug("TEventComboBoxEditLink BeginEdit")
	if !m.stopping {
		m.combobox.Show()
		m.combobox.SetFocus()
	}
	return true
}

func (m *TEventComboBoxEditLink) CancelEdit() bool {
	logs.Debug("TEventComboBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.combobox.Hide()
		if m.VTree != nil {
			m.VTree.CancelEditNode()
		}
	}
	return true
}

func (m *TEventComboBoxEditLink) EndEdit() bool {
	value := m.combobox.Text()
	logs.Debug("TEventComboBoxEditLink EndEdit", "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.SetValue(m.combobox.ItemIndex(), m.combobox.Text())
		m.combobox.Hide()
		if m.VTree != nil {
			m.VTree.EndEditNode()
		}
	}
	return true
}

func (m *TEventComboBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TEventComboBoxEditLink PrepareEdit")
	if m.combobox == nil || !m.combobox.IsValid() {
		m.CreateEdit()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	text := ""
	m.VTree.GetTextInfo(node, column, m.combobox.Font(), &m.bounds, &text)
	logs.Debug("  PrepareEdit GetTextInfo:", m.bounds, text)
	m.combobox.Font().SetColor(colors.ClWindowText)
	m.combobox.SetParent(m.VTree)
	m.combobox.HandleNeeded()
	if text != "" && text != defaultText {
		items := m.combobox.Items()
		index := int32(-1)
		for i := int32(0); i < items.Count(); i++ {
			if items.Strings(i) == text {
				index = i
				break
			}
		}
		if index != -1 {
			m.combobox.SetItemIndex(index)
			m.BindData.EditNodeData.Index = index
		}
	}
	return true
}

func (m *TEventComboBoxEditLink) GetBounds() types.TRect {
	return m.combobox.BoundsRect()
}

func (m *TEventComboBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TEventComboBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.combobox, msg)
}

func (m *TEventComboBoxEditLink) SetBounds(R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetWidth(columnRect.Width())
	m.combobox.SetBoundsRect(R)
	logs.Debug("TEventComboBoxEditLink SetBounds", R)
}

func (m *TEventComboBoxEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TEventComboBoxEditLink Destroy")
	m.combobox.Free()
}
