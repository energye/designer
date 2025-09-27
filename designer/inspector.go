package designer

import (
	"encoding/json"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"sort"
	"strings"
)

var (
	inspector *Inspector
)

// 组件树和对象查看器
type Inspector struct {
	boxSplitter        lcl.ISplitter                         // 分割线
	componentTree      *InspectorComponentTree               // 组件树
	componentProperty  *InspectorComponentProperty           // 组件属性
	objectPropertyList map[uintptr][]lcl.ComponentProperties // 组件的属性列表, 删除时同步删除
}

// 返回查看器实例
func GetInspector() *Inspector {
	return inspector
}

// 加载组件
// 属性, 事件
func (m *Inspector) LoadComponent(component *DesigningComponent) {
	toJSON := func(cp lcl.ComponentProperties) string {
		str, _ := json.Marshal(cp)
		return string(str)
	}
	_ = toJSON
	object := component.object
	var properties []lcl.ComponentProperties
	if propList, ok := m.objectPropertyList[object.Instance()]; ok {
		properties = propList
	} else {
		properties = lcl.DesigningComponent().GetComponentProperties(object)
		logs.Debug("LoadComponent Count:", len(properties))
		m.objectPropertyList[object.Instance()] = properties
	}
	// 拆分 属性和事件
	var (
		propertyList []lcl.ComponentProperties
		eventList    []lcl.ComponentProperties
	)
	for _, prop := range properties {
		if prop.Kind == "tkMethod" {
			eventList = append(eventList, prop)
		} else {
			propertyList = append(propertyList, prop)
		}
		//logs.Debug("  ", toJSON(prop))
	}
	// 排序
	sort.Slice(propertyList, func(i, j int) bool {
		return strings.ToLower(propertyList[i].Name) < strings.ToLower(propertyList[j].Name)
	})
	sort.Slice(eventList, func(i, j int) bool {
		return strings.ToLower(eventList[i].Name) < strings.ToLower(eventList[j].Name)
	})
	// 测试输出
	//{
	//	for _, prop := range propertyList {
	//		logs.Debug(toJSON(prop))
	//	}
	//	for _, event := range eventList {
	//		logs.Debug(toJSON(event))
	//	}
	//}
	m.componentProperty.Load(propertyList, eventList, component)
}
