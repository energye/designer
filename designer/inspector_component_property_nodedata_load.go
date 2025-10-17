package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
)

// 设计 - 组件的设计属性和设计事件数据加载

// 清空树
func (m *InspectorComponentProperty) Clear() {
	vtedit.ResetPropertyNodeData()
	m.propertyTree.Clear()
	m.eventTree.Clear()
}

// 加载属性和事件列表
func (m *InspectorComponentProperty) Load(component *DesigningComponent) {
	if component != m.currentComponent {
		m.currentComponent = component
		// 清空树数据
		m.Clear()

		// 加载属性列表
		m.loadPropertyList(component)

		// 加载事件列表
		m.loadEventList(component)
	}
}

// 加载属性列表
func (m *InspectorComponentProperty) loadPropertyList(component *DesigningComponent) {
	configCompProp := config.ComponentProperty
	for i, nodeData := range component.propertyList {
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
						OriginNodeData: customProperty.Clone(), AffiliatedComponent: component}
					component.propertyList[i] = newEditNodeData // 更新到组件属性
					nodeData = component.propertyList[i]
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

// 加载事件列表
func (m *InspectorComponentProperty) loadEventList(component *DesigningComponent) {

}
