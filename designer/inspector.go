package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
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
func (m *Inspector) LoadComponent(component lcl.IObject) {
	properties := lcl.DesigningComponent().GetComponentProperties(component)
	for _, prop := range properties {
		fmt.Println(prop)
	}
}
