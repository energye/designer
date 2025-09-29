package vtedit

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"math/bits"
	"os"
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

// 构建节点数据
func (m *TEditLinkNodeData) Build() {
	kind := PropertyKind(m.metadata.Kind)
	switch kind {
	case TkEnumeration: // 枚举 单选, 使用下拉框
		m.Type = PdtComboBox
		m.Name = m.metadata.Name
		m.StringValue = tool.FirstToUpper(m.metadata.Value)
		options := strings.Split(m.metadata.Options, ",")
		sort.Strings(options)
		for _, option := range options {
			option = tool.FirstToUpper(option)
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
			option = tool.FirstToUpper(option)
			checkBox := &TEditLinkNodeData{Type: PdtCheckBox, Name: option, Checked: false}
			for _, value := range values {
				if tool.Equal(option, value) {
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
		case "TCursor": // 指针样式-在配置文件转换
		case "TModalResult": // 模态返回值-在配置文件转换
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
		classInstance, err := strconv.ParseUint(m.metadata.Value, 10, bits.UintSize)
		if err != nil {
			logs.Error("获取类实例失败:", err.Error())
			os.Exit(1)
		}
		// 转换 object 获取对象属性
		object := lcl.AsObject(classInstance)
		var properties []lcl.ComponentProperties
		properties = lcl.DesigningComponent().GetComponentProperties(object)
		logs.Debug("TkClass LoadComponent Count:", len(properties))

	default: // 未识别类型
		m.Type = PdtText
		m.Name = m.metadata.Name
		logs.Warn("未识别的元数据类型:", m.metadata.ToJSON())
		return
	}
}
