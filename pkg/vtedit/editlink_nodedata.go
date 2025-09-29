package vtedit

import (
	"bytes"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
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
	metadata      *lcl.ComponentProperties // 组件属性元数据
	Name          string                   // 属性名
	Index         int32                    // 值索引 值是数组类型时，选中的索引
	Checked       bool                     // 选中列表 值是数组类型时，是否选中
	StringValue   string                   // 属性值 string
	FloatValue    float64                  // 属性值 float64
	BoolValue     bool                     // 属性值 bool
	IntValue      int                      // 属性值 int
	CheckBoxValue []*TEditLinkNodeData     // 属性值 checkbox
	ComboBoxValue []*TEditLinkNodeData     // 属性值 combobox
	Type          PropertyDataType         // 属性值类型 普通文本, 单选框, 多选框, 下拉框, 菜单(子菜单)
}

var (
	// 组件属性数据列表, key: 节点指针 value: 节点数据
	propertyTreeDataList = make(map[types.PVirtualNode]*TEditLinkNodeData)
)

func NewEditLinkNodeData(prop *lcl.ComponentProperties) *TEditLinkNodeData {
	m := &TEditLinkNodeData{metadata: prop}
	m.Build()
	return m
}

func ResetPropertyNodeData() {
	propertyTreeDataList = make(map[types.PVirtualNode]*TEditLinkNodeData)
}

// 添加数据到指定节点
func AddPropertyNodeData(tree lcl.ILazVirtualStringTree, parent types.PVirtualNode, data *TEditLinkNodeData) types.PVirtualNode {
	node := tree.AddChild(parent, 0)
	// 设置到数据列表, 增加绑定关系
	propertyTreeDataList[node] = data
	if data.Type == PdtCheckBoxList {
		dataList := data.CheckBoxValue
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
			AddPropertyNodeData(tree, node, item)
		}
		buf.WriteString("]")
		data.StringValue = buf.String()
	}
	return node
}

// 获取节点属性数据
func GetPropertyNodeData(node types.PVirtualNode) *TEditLinkNodeData {
	if data, ok := propertyTreeDataList[node]; ok {
		return data
	}
	return nil
}

// 通知更新组件属性
func (m *TEditLinkNodeData) UpdateComponentProperties() {
	if m.metadata != nil {
		logs.Debug("TEditLinkNodeData UpdateComponentProperties")
	}
}
