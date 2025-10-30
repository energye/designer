package designer

import "unsafe"

// 设计 - 组件设计树数据加载

// 删除当前节点
func (m *TDesigningComponent) Remove() {
	//owner:=m.owner
	//m.owner=nil
}

func (m *TDesigningComponent) instance() uintptr {
	return uintptr(unsafe.Pointer(m))
}

// 向当前组件节点添加子组件节点
func (m *TDesigningComponent) AddChild(child *TDesigningComponent) {
	m.formTab.AddComponentNode(m, child)
}

// 设置当前设计组件为选中状态
func (m *TDesigningComponent) SetSelected() {
	m.node.SetSelected(true)
}

// 重新排序当前组件树节点
func (m *TDesigningComponent) Order(changeLevel ChangeLevel) {
	switch changeLevel {
	case CLevelFront:
	case CLevelBack:
	case CLevelForwardOne:
	case CLevelBackOne:
	}
}
