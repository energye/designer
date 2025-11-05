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
	"bytes"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/keys"
	"strconv"
)

// CheckBox

type TCheckBoxEditLink struct {
	*TBaseEditLink
	checkbox  lcl.ICheckBox
	alignment types.TAlignment
	stopping  bool
}

func NewCheckBoxEditLink(bindData *TEditNodeData) *TCheckBoxEditLink {
	link := new(TCheckBoxEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TCheckBoxEditLink) Create() {
	logs.Debug("TCheckBoxEditLink CreateEdit")
	m.checkbox = lcl.NewCheckBox(nil)
	m.checkbox.SetVisible(false)
	m.checkbox.SetCaption("(false)")
	m.checkbox.SetDoubleBuffered(true)
	m.checkbox.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		logs.Debug("TCheckBoxEditLink OnKeyDown key:", *key)
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})
	m.checkbox.SetOnChange(func(sender lcl.IObject) {
		m.checkbox.SetCaption("(" + strconv.FormatBool(m.checkbox.Checked()) + ")")
		m.BindData.EditNodeData.Checked = m.checkbox.Checked()
		m.BindData.FormInspectorPropertyToComponentProperty()
		logs.Debug("TCheckBoxEditLink OnChange checked:", m.BindData.EditNodeData.Checked)
		node := m.Node.ToGo()
		parentNode := node.Parent
		if pData := GetPropertyNodeData(parentNode); pData != nil {
			// 只当前节点的父节点是 PdtCheckBoxList 类型时才修改父节点的显示单元格
			if pData.EditNodeData.Type != consts.PdtCheckBoxList {
				return
			}
			dataList := pData.EditNodeData.CheckBoxValue
			buf := bytes.Buffer{}
			buf.WriteString("[")
			i := 0
			for _, item := range dataList {
				if item.Checked {
					if i > 0 {
						buf.WriteString(",")
					}
					buf.WriteString(item.Name)
					i++
				}
			}
			buf.WriteString("]")
			pData.EditNodeData.StringValue = buf.String()
			//pData.EditNodeData.IsModify = pData.EditNodeData.StringValue != pData.OriginNodeData.StringValue
			logs.Debug("TCheckBoxEditLink OnChange ParentNode-text:", pData.EditNodeData.StringValue)
			m.VTree.InvalidateNode(parentNode)
		}
	})
}

func (m *TCheckBoxEditLink) BeginEdit() bool {
	if !m.stopping {
		m.checkbox.SetVisible(true)
	}
	return true
}

func (m *TCheckBoxEditLink) CancelEdit() bool {
	logs.Debug("TCheckBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.checkbox.SetVisible(false)
		if m.VTree != nil {
			m.VTree.CancelEditNode()
		}
	}
	return true
}

func (m *TCheckBoxEditLink) EndEdit() bool {
	value := m.checkbox.Checked()
	logs.Debug("TCheckBoxEditLink EndEdit", "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.EditNodeData.Checked = value
		m.checkbox.SetVisible(false)
		if m.VTree != nil {
			m.VTree.EndEditNode()
		}
	}
	return true
}

func (m *TCheckBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TCheckBoxEditLink PrepareEdit")
	if m.checkbox == nil || m.checkbox.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	value := m.VTree.Text(m.Node, m.Column)
	if v, err := strconv.ParseBool(value); err == nil {
		m.BindData.EditNodeData.Checked = v
	}
	m.checkbox.SetParent(m.VTree)
	m.checkbox.HandleNeeded()
	m.checkbox.SetChecked(m.BindData.EditNodeData.Checked)
	return true
}

func (m *TCheckBoxEditLink) GetBounds() (R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R = columnRect
	return
}

func (m *TCheckBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TCheckBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.checkbox, msg)
}

func (m *TCheckBoxEditLink) SetBounds(R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left + 5
	R.Top = columnRect.Top + 3
	//R.SetHeight(columnRect.Height())
	//R.SetWidth(columnRect.Width())
	m.checkbox.SetBoundsRect(R)
	logs.Debug("TCheckBoxEditLink SetBounds", R)
}

func (m *TCheckBoxEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TCheckBoxEditLink Destroy")
	m.checkbox.Free()
}
