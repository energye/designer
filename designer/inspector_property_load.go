package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
)

// 设计 - 组件的属性和事件数据加载

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
		m.loadPropertyList(component.propertyList)

		// 加载事件列表
		m.loadEventList(component.eventList)
	}
}

// 加载属性列表
func (m *InspectorComponentProperty) loadPropertyList(propertyList []*lcl.ComponentProperties) {
	//data := &vtedit.TEditLinkNodeData{Type: vtedit.PdtText, Name: "TextEdit", StringValue: "Value"}
	//vtedit.AddPropertyNodeData(m.propertyTree, 0, data)
	compProp := config.ComponentProperty
	for i, property := range propertyList {
		if compProp.IsExclude(property.Name) {
			logs.Debug("排除属性:", property.ToJSON())
			continue
		}
		logs.Debug("加载属性:", property.ToJSON())
		// 自定义属性, 使用会覆蓋掉
		// 返回数组
		if customProps := compProp.GetCustomPropertyList(property.Name); customProps != nil {
			if len(customProps) == 1 {
				// 数组只有一个元素，规则为直接作用在当前属性上
				customProperty := &customProps[0]
				propertyList[i] = customProperty // 更新到组件属性
				property = propertyList[i]
			} else {

			}
		}
		// 根据组件属性构建节点数据
		newData := vtedit.NewEditLinkNodeData(property)
		vtedit.AddPropertyNodeData(m.propertyTree, 0, newData)
	}
}

// 加载事件列表
func (m *InspectorComponentProperty) loadEventList(eventList []*lcl.ComponentProperties) {

}

func (m *InspectorComponentProperty) addNodeData() {

}
