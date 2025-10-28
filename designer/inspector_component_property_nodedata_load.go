package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
)

// 设计 - 组件的设计属性和设计事件数据加载

// 加载组件属性列表
func (m *TDesigningComponent) loadPropertyList() {
	configCompProp := config.ComponentProperty
	for i, nodeData := range m.PropertyList {
		if configCompProp.IsExclude(nodeData.EditNodeData.Name) {
			logs.Debug("排除属性:", nodeData.EditNodeData.Metadata.ToJSON())
			continue
		}
		logs.Debug("加载属性:", nodeData.EditNodeData.Metadata.ToJSON())
		if !nodeData.IsFinal {
			// 自定义属性, 使用会覆蓋掉
			// 返回数组
			if customProps := configCompProp.GetCustomPropertyList(nodeData.EditNodeData.Name); customProps != nil {
				if len(customProps) == 1 {
					// 数组只有一个元素，规则为直接作用在当前属性上
					customProperty := vtedit.NewEditLinkNodeData(&customProps[0])
					newEditNodeData := &vtedit.TEditNodeData{IsFinal: true, EditNodeData: customProperty,
						OriginNodeData: customProperty.Clone(), AffiliatedComponent: m}
					m.PropertyList[i] = newEditNodeData // 更新到组件属性
					nodeData = m.PropertyList[i]
					newEditNodeData.Build()
				} else {
					// 自定义属性添加？？
				}
			}
			nodeData.IsFinal = true
		}
		// 属性节点数据添加到树
		vtedit.AddPropertyNodeData(m.propertyTree, 0, nodeData)
	}
}

// 加载组件事件列表
func (m *TDesigningComponent) loadEventList() {

}
