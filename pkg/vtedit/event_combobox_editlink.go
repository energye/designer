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
	"github.com/energye/designer/designer/dependmod"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
	"strings"
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
	logs.Debug("TEventComboBoxEditLink-CreateEdit")
	m.combobox = lcl.NewComboBox(nil)
	m.combobox.SetVisible(false)
	m.combobox.SetBorderStyle(types.BsSingle)
	m.combobox.SetAutoSize(false)
	m.combobox.SetDoubleBuffered(true)
	m.combobox.SetShowHint(true)
	m.combobox.SetHint("双击创建绑定事件")
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
	m.combobox.SetOnDblClick(func(sender lcl.IObject) {
		compName := m.BindData.AffiliatedComponent.GetName()
		propName := m.BindData.Name()
		bindEventName := compName + propName
		logs.Debug("TEventComboBoxEditLink-CreateEdit 自动创建绑定事件名, 组件名:", compName, "属性名:", propName, "=> 事件名:", bindEventName)
		// 验证当前列表是否已有该事件
		items := m.combobox.Items()
		count := m.combobox.Items().Count()
		itemIndex := int32(-1)
		for i := int32(0); i < count; i++ {
			if items.Strings(i) == bindEventName {
				itemIndex = i
				break
			}
		}
		logs.Debug("TEventComboBoxEditLink-CreateEdit itemIndex:", itemIndex)
		if itemIndex == -1 {
			// 添加
			m.combobox.SetItemIndex(-1)
			m.combobox.SetText(bindEventName)
			items.Add(bindEventName)
		} else {
			// 索引
			m.combobox.SetText(bindEventName)
			m.combobox.SetItemIndex(itemIndex)
		}
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.VTree.EndEditNode()
		})
	})

	// 动态添加下拉项, 需要过滤匹配
	items := m.combobox.Items()
	comboBoxValue := []*TEditLinkNodeData{{StringValue: defaultText}}
	methods := m.BindData.AffiliatedComponent.GetRecvMethods()
	if methods != nil {
		// 过滤匹配出符合的事件类型方法
		funcTypeAliases := dependmod.GetFuncTypeAliases(m.BindData.AffiliatedComponent.GetMod())
		if funcTypeAliases != nil {
			lowerType := strings.ToLower(m.BindData.EditNodeData.Metadata.Type)
			if funcType := funcTypeAliases.Funcs.Get(lowerType); funcType != nil {
				// 得到参数和返回值列表与 GetRecvMethods 匹配过滤
				funcInfo := &dast.TFuncInfo{
					Params:  dast.ParseFields(funcType.Params),
					Results: dast.ParseFields(funcType.Results),
				}
				var tempMethods []*dast.TFuncInfo
				for _, method := range methods {
					if dast.IsFuncInfoMatch(method, funcInfo) {
						tempMethods = append(tempMethods, method)
					}
				}
				methods = tempMethods
			}
		}
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
	value = strings.TrimSpace(value)
	m.BindData.EditNodeData.SetEditValue(value)
	isEmpty := value == "" || value == defaultText
	if isEmpty {
		// -1 删除
		m.BindData.EditNodeData.Index = -1
		m.BindData.EditNodeData.EventState = consts.EsDelete
		return
	}
	isExist := false
	for _, item := range m.BindData.EditNodeData.ComboBoxValue {
		if item.StringValue == value {
			isExist = true
			break
		}
	}
	if !isEmpty && !isExist {
		// > 0 新增
		m.BindData.EditNodeData.ComboBoxValue = append(m.BindData.EditNodeData.ComboBoxValue, &TEditLinkNodeData{StringValue: value})
		m.BindData.EditNodeData.Index = int32(len(m.BindData.EditNodeData.ComboBoxValue)) - 1
		m.BindData.EditNodeData.EventState = consts.EsAdd
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
	itemIndex := m.combobox.ItemIndex()
	value := m.combobox.Text()
	logs.Debug("TEventComboBoxEditLink-EndEdit", "itemIndex:", itemIndex, "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.SetValue(itemIndex, value)
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
		count := items.Count()
		for i := int32(0); i < count; i++ {
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
