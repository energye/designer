package designer

import (
	"github.com/energye/lcl/lcl"
)

var (
	inspector *Inspector
)

// 组件树和对象查看器
type Inspector struct {
	boxSplitter       lcl.ISplitter               // 分割线
	componentTree     *InspectorComponentTree     // 组件树实例
	componentProperty *InspectorComponentProperty // 组件属性实现
}
