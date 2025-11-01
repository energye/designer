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
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strconv"
)

// 查看器的数据类型

// PropertyDataType 属性数据组件类型
type PropertyDataType int32

const (
	PdtText PropertyDataType = iota
	PdtInt
	PdtInt64
	PdtFloat
	PdtRadiobutton
	PdtCheckBox
	PdtCheckBoxDraw
	PdtCheckBoxList
	PdtComboBox
	PdtClassDialog
	PdtColorSelect
	PdtClass
)

// 节点数据
type TEditLinkNodeData struct {
	Metadata      *lcl.ComponentProperties // 组件属性元数据
	Name          string                   // 属性名
	Index         int32                    // 值索引 值是数组类型时，选中的索引
	Checked       bool                     // 选中列表 值是数组类型时，是否选中
	StringValue   string                   // 属性值 string
	FloatValue    float64                  // 属性值 float64
	BoolValue     bool                     // 属性值 bool
	IntValue      int                      // 属性值 int
	Class         TPropClass               // 属性值 class 实例
	CheckBoxValue []*TEditLinkNodeData     // 属性值 checkbox
	ComboBoxValue []*TEditLinkNodeData     // 属性值 combobox
	Type          PropertyDataType         // 属性值类型 普通文本, 单选框, 多选框, 下拉框, 菜单(子菜单)
}

type TPropClass struct {
	Instance uintptr // 属性值 class 实例
	Count    int32   // 属性值 class 属性数量
}

// 编辑数据返回字符串值
func (m *TEditLinkNodeData) EditStringValue() string {
	switch m.Type {
	case PdtText:
		return m.StringValue
	case PdtInt, PdtInt64:
		return strconv.Itoa(m.IntValue)
	case PdtFloat:
		val := strconv.FormatFloat(m.FloatValue, 'f', 2, 64)
		return val
	case PdtCheckBox:
		return strconv.FormatBool(m.Checked)
	case PdtCheckBoxList:
		return m.StringValue
	case PdtComboBox:
		return m.StringValue
	case PdtColorSelect:
		return fmt.Sprintf("0x%X", m.IntValue)
	case PdtClass:
		return m.StringValue
	default:
		return ""
	}
}

// 编辑数据返回原始类型值
func (m *TEditLinkNodeData) EditValue() any {
	switch m.Type {
	case PdtText:
		return m.StringValue
	case PdtInt, PdtInt64:
		return m.IntValue
	case PdtFloat:
		return m.FloatValue
	case PdtCheckBox:
		return m.Checked
	case PdtCheckBoxList:
		return m.StringValue
	case PdtComboBox:
		return m.StringValue
	case PdtColorSelect:
		return m.IntValue
	case PdtClass:
		return m.StringValue
	default:
		return ""
	}
}

func (m *TEditLinkNodeData) SetEditValue(value any) {
	switch m.Type {
	case PdtText:
		m.StringValue = value.(string)
	case PdtInt, PdtInt64:
		m.IntValue = int(value.(int32))
	case PdtFloat:
		m.FloatValue = value.(float64)
	case PdtCheckBox:
		m.Checked = value.(bool)
	case PdtCheckBoxList:
		m.StringValue = value.(string)
	case PdtComboBox:
		m.StringValue = value.(string)
	case PdtColorSelect:
		m.IntValue = int(value.(uint32))
	}
}

func (m *TEditLinkNodeData) Clone() *TEditLinkNodeData {
	if m == nil {
		return nil
	}
	clone := &TEditLinkNodeData{
		Name:        m.Name,
		Index:       m.Index,
		Checked:     m.Checked,
		StringValue: m.StringValue,
		FloatValue:  m.FloatValue,
		BoolValue:   m.BoolValue,
		IntValue:    m.IntValue,
		Type:        m.Type,
		Class:       TPropClass{Instance: m.Class.Instance, Count: m.Class.Count},
	}
	if m.Metadata != nil {
		cloneMetadata := *m.Metadata
		clone.Metadata = &cloneMetadata
	}

	// 深拷贝CheckBoxValue切片（递归克隆每个元素）
	if m.CheckBoxValue != nil {
		clone.CheckBoxValue = make([]*TEditLinkNodeData, len(m.CheckBoxValue))
		for i, item := range m.CheckBoxValue {
			clone.CheckBoxValue[i] = item.Clone()
		}
	}

	// 深拷贝ComboBoxValue切片（递归克隆每个元素）
	if m.ComboBoxValue != nil {
		clone.ComboBoxValue = make([]*TEditLinkNodeData, len(m.ComboBoxValue))
		for i, item := range m.ComboBoxValue {
			clone.ComboBoxValue[i] = item.Clone()
		}
	}
	return clone
}

// 设计组件接口
type IDesigningComponent interface {
	UpdateComponentPropertyToObject(nodeData *TEditNodeData)
}

// 编辑的节点数据
type TEditNodeData struct {
	Parent              *TEditNodeData      // 父属性节点
	Child               []*TEditNodeData    // 子属性节点
	IsFinal             bool                // 标记是否最终对象, 用于完整的数据
	EditNodeData        *TEditLinkNodeData  // 编辑数据
	OriginNodeData      *TEditLinkNodeData  // 原始数据
	AffiliatedNode      types.PVirtualNode  // 所属属性树节点
	AffiliatedComponent IDesigningComponent // 所属组件对象
}

var (
	// 组件属性数据列表, key: 节点指针 value: 节点数据
	propertyTreeDataList = make(map[types.PVirtualNode]*TEditNodeData)
)

// 创建一个编辑节点数据
func NewEditLinkNodeData(prop *lcl.ComponentProperties) *TEditLinkNodeData {
	m := &TEditLinkNodeData{Metadata: prop}
	m.Build()
	return m
}

//func ResetPropertyNodeData() {
//	//propertyTreeDataList = make(map[types.PVirtualNode]*TEditNodeData)
//}

// 添加数据到指定节点
func AddPropertyNodeData(tree lcl.ILazVirtualStringTree, parent types.PVirtualNode, data *TEditNodeData) types.PVirtualNode {
	node := tree.AddChild(parent, 0)
	// 节点设置到节点数据
	data.AffiliatedNode = node
	// 设置到数据列表, 增加绑定关系
	propertyTreeDataList[node] = data
	if data.Type() == PdtCheckBoxList {
		// 复选框列表
		dataList := data.EditNodeData.CheckBoxValue
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
			newItemData := &TEditNodeData{EditNodeData: item, OriginNodeData: item.Clone(), AffiliatedComponent: data.AffiliatedComponent}
			AddPropertyNodeData(tree, node, newItemData)
		}
		buf.WriteString("]")
		data.EditNodeData.StringValue = buf.String()
	} else if data.Type() == PdtClass {
		for _, nodeData := range data.Child {
			AddPropertyNodeData(tree, node, nodeData)
		}
	}
	return node
}

// 获取节点属性数据
func GetPropertyNodeData(node types.PVirtualNode) *TEditNodeData {
	if data, ok := propertyTreeDataList[node]; ok {
		return data
	}
	return nil
}

// 判断节点对象是否存在
func IsExistNodeData(node types.PVirtualNode) bool {
	if node == 0 {
		return false
	}
	_, ok := propertyTreeDataList[node]
	return ok
}

// 删除节点属性数据
func DelPropertyNodeData(node types.PVirtualNode) {
	delete(propertyTreeDataList, node)
}

func (m *TEditNodeData) Type() PropertyDataType {
	return m.EditNodeData.Type
}
func (m *TEditNodeData) Class() TPropClass {
	return m.EditNodeData.Class
}
func (m *TEditNodeData) Name() string {
	return m.EditNodeData.Name
}

// 从设计属性更新到组件属性
func (m *TEditNodeData) FormInspectorPropertyToComponentProperty() {
	if m.EditNodeData != nil {
		logs.Debug("TEditLinkNodeData FormInspectorPropertyToComponentProperty property-name:", m.EditNodeData.Name)
		//go lcl.RunOnMainThreadAsync(func(id uint32) {
		m.AffiliatedComponent.UpdateComponentPropertyToObject(m)
		//})
	}
}

// 从组件属性更新到设计属性
func (m *TEditNodeData) FormComponentPropertyToInspectorProperty() {
	if m.EditNodeData != nil {
		logs.Debug("TEditLinkNodeData FormComponentPropertyToInspectorProperty property-name:", m.EditNodeData.Name)
	}

	//m.AffiliatedNode
}

// 是否被修改
func (m *TEditNodeData) IsModify() bool {
	switch m.Type() {
	case PdtCheckBox:
		return m.EditNodeData.Checked != m.OriginNodeData.Checked
	case PdtText:
		return m.EditNodeData.StringValue != m.OriginNodeData.StringValue
	case PdtInt, PdtInt64:
		return m.EditNodeData.IntValue != m.OriginNodeData.IntValue
	case PdtFloat:
		return m.EditNodeData.FloatValue != m.OriginNodeData.FloatValue
	case PdtCheckBoxList:
		// CheckBox 判断集合是否修改
		for i, item := range m.EditNodeData.CheckBoxValue {
			if item.Checked != m.OriginNodeData.CheckBoxValue[i].Checked {
				return true
			}
		}
	case PdtComboBox:
		return m.EditNodeData.StringValue != m.OriginNodeData.StringValue
	case PdtColorSelect:
		return m.EditNodeData.IntValue != m.OriginNodeData.IntValue
	case PdtClass:
		// 类实例, 需要判断类下的属性是否被修改
		for _, child := range m.Child {
			if child.IsModify() {
				return true
			}
		}
	}
	return false
}

// 获取修改class的子节点
func (m *TEditNodeData) GetModifyClassChildNodeData() *TEditNodeData {
	if m.Type() == PdtClass {
		for _, child := range m.Child {
			if child.IsModify() {
				if child.Type() == PdtClass {
					return child.GetModifyClassChildNodeData()
				} else {
					return child
				}
			}
		}
	}
	return nil
}

// 编辑数据返回字符串值
func (m *TEditNodeData) EditStringValue() string {
	return m.EditNodeData.EditStringValue()
}

// 编辑数据返回原始类型值
func (m *TEditNodeData) EditValue() any {
	return m.EditNodeData.EditValue()
}

// 返回编辑字符串值
func (m *TEditNodeData) SetEditValue(value any) {
	m.EditNodeData.SetEditValue(value)
}

// 获得类的路径 Txxx.Txxx.Txxx ...
func (m *TEditNodeData) Paths() []string {
	// todo 1: 可能存在的问题, 某父对象不是class一定是错误的
	var paths []string
	pData := m.Parent
	for pData != nil {
		if pData.Type() == PdtClass { // todo 1
			paths = append(paths, pData.Name())
		} else {
			// 不正确, 直接退出
			panic("递归遍历属性节点对象路径错误, 对象非class类型, 节点必须为class类型: " + pData.Name())
		}
		pData = pData.Parent
	}
	return paths
}
