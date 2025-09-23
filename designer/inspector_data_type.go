package designer

import "github.com/energye/lcl/types"

// 查看器的数据类型

// PropertyDataType 属性数据组件类型
// 有哪些？TODO 0:按钮 1:复选框 2:下拉框 3:进度条 4:微调框 5:日期选择器
type PropertyDataType int32

const (
	PdtText PropertyDataType = iota
	PdtRadiobutton
	PdtCheckBox
	PdtComboBox
	PdtTree
	PdtClassDialog
)

var treePropertyNodeDatas = make(map[types.PVirtualNode]*TTreePropertyNodeData)

type TTreePropertyNodeData struct {
	Name      string           // 属性名
	Value     string           // 属性值
	ValueList []string         // 属性值列表
	Type      PropertyDataType // 属性值类型 普通文本, 单选框, 多选框, 下拉框, 菜单(子菜单)
}

type ValueList struct {
	Label string
	Value string
}

func GetPropertyNodeData(nodeKey types.PVirtualNode) *TTreePropertyNodeData {
	if data, ok := treePropertyNodeDatas[nodeKey]; ok {
		return data
	}
	return nil
}
