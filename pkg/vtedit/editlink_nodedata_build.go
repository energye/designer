package vtedit

import (
	"github.com/energye/designer/pkg/logs"
	"sort"
	"strconv"
	"strings"
)

// 查看器的数据类型 构建

type PropertyKind string

const (
	TkClass       PropertyKind = "tkClass"
	TkEnumeration PropertyKind = "tkEnumeration"
	TkSet         PropertyKind = "tkSet"
	TkBool        PropertyKind = "tkBool"
	TkAString     PropertyKind = "tkAString"
	TkChar        PropertyKind = "tkChar"
	TkInteger     PropertyKind = "tkInteger"
	TkInt64       PropertyKind = "tkInt64"
)

type PropertyType string

func (m *TEditLinkNodeData) Build() {
	kind := PropertyKind(m.metadata.Kind)
	switch kind {
	case TkEnumeration: // 枚举 单选, 使用下拉框
		m.Type = PdtComboBox
		m.Name = m.metadata.Name
		m.StringValue = m.metadata.Value
		options := strings.Split(m.metadata.Options, ",")
		sort.Strings(options)
		for _, option := range options {
			item := &TEditLinkNodeData{StringValue: option}
			m.ComboBoxValue = append(m.ComboBoxValue, item)
		}
	case TkSet: // 集合 多选, 使用子菜单列表
		m.Type = PdtCheckBoxList
		m.Name = m.metadata.Name
		values := strings.Split(m.metadata.Value, ",")
		options := strings.Split(m.metadata.Options, ",")
		sort.Strings(options)
		for _, option := range options {
			checkBox := &TEditLinkNodeData{Type: PdtCheckBox, Name: option, Checked: false}
			for _, value := range values {
				if option == value {
					checkBox.Checked = true
					break
				}
			}
			m.CheckBoxValue = append(m.CheckBoxValue, checkBox)
		}
	case TkBool: // 布尔类型
		m.Type = PdtCheckBox
		m.Name = m.metadata.Name
		m.Checked = m.metadata.Value == "1"
	case TkAString: // 字符串
		m.Type = PdtText
		m.Name = m.metadata.Name
		m.StringValue = m.metadata.Value
	case TkChar: // 密码
		m.Type = PdtText
		m.Name = m.metadata.Name
		m.StringValue = ""
	case TkInteger: // 数字
		m.Type = PdtInt
		m.Name = m.metadata.Name
		v, _ := strconv.Atoi(m.metadata.Value)
		m.IntValue = v
		// TModalResult TCursor TGraphicsColor
		switch m.metadata.Type {
		case "TGraphicsColor": // 颜色
			m.Type = PdtColorSelect
		case "TCursor": // 指针样式

		case "TModalResult": // 模态返回值

		}
	case TkInt64: // 数字 64
		m.Type = PdtInt64
		m.Name = m.metadata.Name
		v, _ := strconv.Atoi(m.metadata.Value)
		m.IntValue = v
	case TkClass: // 类
		m.Type = PdtClass
		m.Name = m.metadata.Name
		m.StringValue = m.metadata.Value
		// 获取类实例 属性

	default: // 未识别类型
		m.Type = PdtText
		m.Name = m.metadata.Name
		logs.Warn("未识别的元数据类型:", m.metadata.ToJSON())
		return
	}
}
