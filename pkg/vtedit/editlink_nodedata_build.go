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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"math/bits"
	"os"
	"sort"
	"strconv"
)

// 查看器的数据类型 构建

// 构建节点数据
func (m *TEditLinkNodeData) Build() {
	kind := consts.PropertyKind(m.Metadata.Kind)
	switch kind {
	case consts.TkEnumeration: // 枚举 单选, 使用下拉框
		m.Type = consts.PdtComboBox
		m.Name = m.Metadata.Name
		m.StringValue = tool.FirstToUpper(m.Metadata.Value)
		options := tool.Split(m.Metadata.Options, ",")
		sort.Strings(options)
		for _, option := range options {
			option = tool.FirstToUpper(option)
			item := &TEditLinkNodeData{StringValue: option}
			m.ComboBoxValue = append(m.ComboBoxValue, item)
		}
	case consts.TkSet: // 集合 多选, 使用子菜单列表
		m.Type = consts.PdtCheckBoxList
		m.Name = m.Metadata.Name
		values := tool.Split(m.Metadata.Value, ",")
		options := tool.Split(m.Metadata.Options, ",")
		sort.Strings(options)
		for _, option := range options {
			option = tool.FirstToUpper(option)
			checkBox := &TEditLinkNodeData{Type: consts.PdtCheckBox, Name: option, Checked: false}
			for _, value := range values {
				if tool.Equal(option, value) {
					checkBox.Checked = true
					break
				}
			}
			m.CheckBoxValue = append(m.CheckBoxValue, checkBox)
		}
	case consts.TkBool: // 布尔类型
		m.Type = consts.PdtCheckBox
		m.Name = m.Metadata.Name
		m.Checked = m.Metadata.Value == "1"
	case consts.TkAString: // 字符串
		m.Type = consts.PdtText
		m.Name = m.Metadata.Name
		m.StringValue = m.Metadata.Value
	case consts.TkChar: // 密码
		m.Type = consts.PdtUint16
		m.Name = m.Metadata.Name
		m.StringValue = "" // 默认 0 == 空
	case consts.TkInteger: // 数字
		m.Type = consts.PdtInt
		m.Name = m.Metadata.Name
		v, _ := strconv.Atoi(m.Metadata.Value)
		m.IntValue = v
		// TModalResult TCursor TGraphicsColor
		switch m.Metadata.Type {
		case "TGraphicsColor": // 颜色
			m.Type = consts.PdtColorSelect
		case "TCursor": // 指针样式-在配置文件转换
		case "TModalResult": // 模态返回值-在配置文件转换
		}
	case consts.TkInt64: // 数字 64
		m.Type = consts.PdtInt64
		m.Name = m.Metadata.Name
		v, _ := strconv.Atoi(m.Metadata.Value)
		m.IntValue = v
	case consts.TkClass: // 类
		m.Type = consts.PdtClass
		m.Name = m.Metadata.Name
		// 获取类实例 属性
		classInstance, err := strconv.ParseUint(m.Metadata.Value, 10, bits.UintSize)
		if err != nil {
			logs.Error("获取类实例失败:", err.Error())
			os.Exit(1)
		}
		m.Class = TPropClass{Instance: uintptr(classInstance)}
		m.StringValue = "(" + m.Metadata.Type + ")"
	default: // 未识别类型
		m.Type = consts.PdtText // todo 使用文本
		m.Name = m.Metadata.Name
		logs.Warn("未识别的元数据类型:", m.Metadata.ToJSON())
		return
	}
}

// 构建节点数据
func (m *TEditNodeData) Build() {
	// 构建类字段属性, 做为子节点
	if m.EditNodeData.Type == consts.PdtClass {
		if m.EditNodeData.Class.Instance != 0 {
			object := lcl.AsObject(m.EditNodeData.Class.Instance)
			methods := tool.GetObjectMethodNames(object)
			if methods == nil {
				logs.Error("获取当前组件对象属性错误, 获取对象方法列表为空, 组件名:", m.Name())
			}
			var properties []lcl.ComponentProperties
			properties = lcl.DesigningComponent().GetComponentProperties(object)
			m.EditNodeData.Class.Count = int32(len(properties))
			logs.Debug("TkClass LoadComponent", object.ToString(), "Count:", len(properties))
			for _, prop := range properties {
				newProp := prop
				tool.FixPropInfo(methods, &newProp)
				if newProp.Kind == "tkMethod" {
					continue // tkMethod 事件函数
				}
				newEditLinkNodeData := NewEditLinkNodeData(&newProp)
				newEditNodeData := &TEditNodeData{EditNodeData: newEditLinkNodeData, OriginNodeData: newEditLinkNodeData.Clone(),
					AffiliatedComponent: m.AffiliatedComponent, Parent: m}
				m.Child = append(m.Child, newEditNodeData)
				newEditNodeData.Build()
			}
		} else {
			logs.Debug("TEditNodeData Build Class 实例是'0', 属性名:", m.EditNodeData.Name)
		}
	} else {
		// 其它？？
	}
}
