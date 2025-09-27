package vtedit

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
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
	if tool.Equal(m.metadata.Name, "ShowHint") {
		println()
	}
	kind := PropertyKind(m.metadata.Kind)
	switch kind {
	case TkEnumeration:
		m.Type = PdtComboBox
		m.Name = m.metadata.Name
		m.StringValue = m.metadata.Value
		options := strings.Split(m.metadata.Options, ",")
		sort.Strings(options)
		for _, option := range options {
			item := &TEditLinkNodeData{StringValue: option}
			m.ComboBoxValue = append(m.ComboBoxValue, item)
		}
	case TkSet:
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
	case TkBool:
		m.Type = PdtCheckBox
		m.Name = m.metadata.Name
		m.Checked = m.metadata.Value == "1"
	case TkAString:
		m.Type = PdtText
		m.Name = m.metadata.Name
		m.StringValue = m.metadata.Value
	case TkChar:
		m.Type = PdtText
		m.Name = m.metadata.Name
		m.StringValue = ""
	case TkInteger:
		m.Type = PdtInt
		m.Name = m.metadata.Name
		v, _ := strconv.Atoi(m.metadata.Value)
		m.IntValue = v
		// TModalResult TCursor TGraphicsColor
	case TkInt64:
		m.Type = PdtInt64
		m.Name = m.metadata.Name
		v, _ := strconv.Atoi(m.metadata.Value)
		m.IntValue = v
	default:
		m.Type = PdtText
		m.Name = m.metadata.Name
		logs.Warn("未识别的元数据类型:", m.metadata.ToJSON())
		return
	}
}
