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
// 有哪些？TODO 0:按钮 1:复选框 2:下拉框 3:进度条 4:微调框 5:日期选择器
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
	ClassInstance uintptr                  // 属性值 class 实例
	CheckBoxValue []*TEditLinkNodeData     // 属性值 checkbox
	ComboBoxValue []*TEditLinkNodeData     // 属性值 combobox
	Type          PropertyDataType         // 属性值类型 普通文本, 单选框, 多选框, 下拉框, 菜单(子菜单)
}

func (m *TEditLinkNodeData) IsModify(originNodeData *TEditLinkNodeData) bool {
	switch m.Type {
	case PdtCheckBox:
		return m.Checked != originNodeData.Checked
	case PdtText:
		return m.StringValue != originNodeData.StringValue
	case PdtInt, PdtInt64:
		return m.IntValue != originNodeData.IntValue
	case PdtFloat:
		return m.FloatValue != originNodeData.FloatValue
	case PdtCheckBoxList, PdtClass:
		return m.StringValue != originNodeData.StringValue
	case PdtComboBox:
		return m.StringValue != originNodeData.StringValue
	case PdtColorSelect:
		return m.IntValue != originNodeData.IntValue
	}
	return false
}

func (m *TEditLinkNodeData) EditValue() string {
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
		Name:          m.Name,
		Index:         m.Index,
		Checked:       m.Checked,
		StringValue:   m.StringValue,
		FloatValue:    m.FloatValue,
		BoolValue:     m.BoolValue,
		IntValue:      m.IntValue,
		Type:          m.Type,
		ClassInstance: m.ClassInstance,
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
	UpdateComponentProperty(nodeData *TEditNodeData)
}

// 编辑的节点数据
type TEditNodeData struct {
	IsFinal             bool                // 标记是否最终对象, 用于完整的数据
	EditNodeData        *TEditLinkNodeData  // 编辑数据
	OriginNodeData      *TEditLinkNodeData  // 原始数据
	AffiliatedNode      types.PVirtualNode  // 所属属性树节点
	AffiliatedComponent IDesigningComponent // 所属组件对象
	Child               []*TEditNodeData    // 子节点
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

func ResetPropertyNodeData() {
	propertyTreeDataList = make(map[types.PVirtualNode]*TEditNodeData)
}

// 添加数据到指定节点
func AddPropertyNodeData(tree lcl.ILazVirtualStringTree, parent types.PVirtualNode, data *TEditNodeData) types.PVirtualNode {
	node := tree.AddChild(parent, 0)
	// 节点设置到节点数据
	data.AffiliatedNode = node
	// 设置到数据列表, 增加绑定关系
	propertyTreeDataList[node] = data
	if data.EditNodeData.Type == PdtCheckBoxList {
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
	} else if data.EditNodeData.Type == PdtClass {
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

// 构建节点数据
func (m *TEditNodeData) Build() {
	// 构建类字段属性, 做为子节点
	if m.EditNodeData.ClassInstance != 0 {
		object := lcl.AsObject(m.EditNodeData.ClassInstance)
		var properties []lcl.ComponentProperties
		properties = lcl.DesigningComponent().GetComponentProperties(object)
		logs.Debug("TkClass LoadComponent", object.ToString(), "Count:", len(properties))
		for _, prop := range properties {
			if prop.Kind == "tkMethod" {
				// tkMethod 事件函数
				continue
			}
			newProp := prop
			newEditLinkNodeData := NewEditLinkNodeData(&newProp)
			newEditNodeData := &TEditNodeData{EditNodeData: newEditLinkNodeData, OriginNodeData: newEditLinkNodeData.Clone(), AffiliatedComponent: m.AffiliatedComponent}
			m.Child = append(m.Child, newEditNodeData)
			newEditNodeData.Build()
		}
	} else {
		// 其它？？
	}
}

// 从设计属性更新到组件属性
func (m *TEditNodeData) FormInspectorPropertyToComponentProperty() {
	if m.EditNodeData != nil {
		logs.Debug("TEditLinkNodeData FormInspectorPropertyToComponentProperty property-name:", m.EditNodeData.Name)
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			m.AffiliatedComponent.UpdateComponentProperty(m)
		})
	}
}

// 从组件属性更新到设计属性
func (m *TEditNodeData) FormComponentPropertyToInspectorProperty() {
	if m.EditNodeData != nil {
		logs.Debug("TEditLinkNodeData FormComponentPropertyToInspectorProperty property-name:", m.EditNodeData.Name)
	}
}

// 是否被修改
func (m *TEditNodeData) IsModify() bool {
	return m.EditNodeData.IsModify(m.OriginNodeData)
}

// 返回编辑字符串值
func (m *TEditNodeData) EditValue() string {
	return m.EditNodeData.EditValue()
}

// 返回编辑字符串值
func (m *TEditNodeData) SetEditValue(value any) {
	m.EditNodeData.SetEditValue(value)
}
