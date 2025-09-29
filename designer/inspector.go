package designer

import (
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
	if component == nil {
		logs.Error("加载组件属性/事件失败, 设计组件为空")
		return
	}
	if component.propertyList == nil {
		object := component.object
		var properties []lcl.ComponentProperties
		properties = lcl.DesigningComponent().GetComponentProperties(object)
		logs.Debug("LoadComponent Count:", len(properties))
		// 拆分 属性和事件
		var (
			eventList    []*lcl.ComponentProperties
			propertyList []*lcl.ComponentProperties
		)
		for _, prop := range properties {
			newProp := prop
			if prop.Kind == "tkMethod" {
				eventList = append(eventList, &newProp)
			} else {
				propertyList = append(propertyList, &newProp)
			}
			//logs.Debug("  ", toJSON(prop))
		}
		// 排序
		sort.Slice(eventList, func(i, j int) bool {
			return strings.ToLower(eventList[i].Name) < strings.ToLower(eventList[j].Name)
		})
		sort.Slice(propertyList, func(i, j int) bool {
			return strings.ToLower(propertyList[i].Name) < strings.ToLower(propertyList[j].Name)
		})
		component.eventList = eventList
		component.propertyList = propertyList
	}
	m.componentProperty.Load(component)
}
