package vtedit

// 查看器的数据类型

// PropertyDataType 属性数据组件类型
// 有哪些？TODO 0:按钮 1:复选框 2:下拉框 3:进度条 4:微调框 5:日期选择器
type PropertyDataType int32

const (
	PdtText PropertyDataType = iota
	PdtRadiobutton
	PdtCheckBox
	PdtCheckBoxList
	PdtComboBox
	PdtTree
	PdtClassDialog
)

type TNodeData struct {
	Name          string           // 属性名
	Index         int32            // 值索引 值是数组类型时，选中的索引
	Checked       bool             // 选中列表 值是数组类型时，是否选中
	StringValue   string           // 属性值 string
	DoubleValue   float64          // 属性值 double
	BoolValue     bool             // 属性值 boolean
	CheckBoxValue []TNodeData      // 属性值 checkbox
	ComboBoxValue []TNodeData      // 属性值 combobox
	Type          PropertyDataType // 属性值类型 普通文本, 单选框, 多选框, 下拉框, 菜单(子菜单)
}
