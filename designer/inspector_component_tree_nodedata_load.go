package designer

import "unsafe"

// 设计 - 组件设计树数据加载

// 删除当前节点
func (m *DesigningComponent) Remove() {
	//owner:=m.owner
	//m.owner=nil
}

func (m *DesigningComponent) instance() uintptr {
	return uintptr(unsafe.Pointer(m))
}

// 向当前组件节点添加子组件节点
func (m *DesigningComponent) AddChild(child *DesigningComponent) {
	inspector.componentTree.AddComponentNode(m, child)
}
