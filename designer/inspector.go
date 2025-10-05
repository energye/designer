package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
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

// 返回查看器实例
func GetInspector() *Inspector {
	return inspector
}

// 加载组件
// 属性, 事件
// 参数: component 当前正在设计的组件
func (m *Inspector) LoadComponentProps(component *DesigningComponent) {
	if component == nil {
		logs.Error("加载组件属性/事件失败, 设计组件为空")
		return
	}
	// 属性列表为空时获取属性列表
	component.GetProps()
	// 加载属性列表和事件列表
	go lcl.RunOnMainThreadAsync(func(id uint32) {
		m.componentProperty.Load(component)
		logs.Debug("加载组件属性完成", component.object.ToString())
		selectNode := component.compPropTreeState.selectNode
		if vtedit.IsExistNodeData(selectNode) {
			// 恢复上次选择的编辑节点
			m.componentProperty.propertyTree.SetSelected(selectNode, true)
			m.componentProperty.propertyTree.ScrollIntoViewWithPVirtualNodeBoolX2(selectNode, true, true)
		}
	})
}

// 加载组件
// 组件树
// 参数: root 当前正在设计的根节点
// 参数: component 当前正在设计的组件
func (m *Inspector) LoadComponentTree(root, component *DesigningComponent) {
	m.componentTree.Clear()
	var iterateTreeNode func(node *DesigningComponent)
	iterateTreeNode = func(node *DesigningComponent) {
		if node == nil {
			return
		}
		//for _, item := range node.child {
		//
		//}
	}
	iterateTreeNode(root)
}
