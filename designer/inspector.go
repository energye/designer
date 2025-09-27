package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"sort"
	"strings"
)

var (
	inspector *Inspector
)

// 组件树和对象查看器
type Inspector struct {
	boxSplitter       lcl.ISplitter               // 分割线
	componentTree     *InspectorComponentTree     // 组件树
	componentProperty *InspectorComponentProperty // 组件属性
}

// 返回查看器实例
func GetInspector() *Inspector {
	return inspector
}

// 加载组件
// 属性, 事件
func (m *Inspector) LoadComponent(component *DesigningComponent) {
	properties := lcl.DesigningComponent().GetComponentProperties(component.object)
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
	}
	// 排序
	sort.Slice(propertyList, func(i, j int) bool {
		return strings.ToLower(propertyList[i].Name) < strings.ToLower(propertyList[j].Name)
	})
	sort.Slice(eventList, func(i, j int) bool {
		return strings.ToLower(eventList[i].Name) < strings.ToLower(eventList[j].Name)
	})
	// 测试输出
	{
		for _, prop := range propertyList {
			fmt.Printf("%+v\n", prop)
		}
		for _, event := range eventList {
			fmt.Printf("%+v\n", event)
		}
	}
	m.componentProperty.Load(propertyList, eventList, component)
}
