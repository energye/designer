package designer

import (
	"github.com/energye/lcl/types"
	"unsafe"
)

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
	// 调整组件树节点显示顺序
	switch changeLevel {
	case CLevelFront:
		parentNode := m.node.Parent()
		if parentNode == nil || !parentNode.IsValid() {
			return
		}
		lastChild := parentNode.GetLastChild()
		if lastChild != nil && lastChild.IsValid() {
			m.node.MoveTo(lastChild, types.NaInsertBehind)
		}
	case CLevelBack:
		parentNode := m.node.Parent()
		if parentNode == nil || !parentNode.IsValid() {
			return
		}
		firstChild := parentNode.GetFirstChild()
		if firstChild != nil && firstChild.IsValid() {
			m.node.MoveTo(firstChild, types.NaInsert)
		}
	case CLevelForwardOne:
		nextNode := m.node.GetNextSibling()
		if nextNode != nil && nextNode.IsValid() {
			m.node.MoveTo(nextNode, types.NaInsertBehind)
		}
	case CLevelBackOne:
		prevNode := m.node.GetPrevSibling()
		if prevNode != nil && prevNode.IsValid() {
			m.node.MoveTo(prevNode, types.NaInsert)
		}
	}
	// 调整组件对象顺序
	switch changeLevel {
	case CLevelFront:

	case CLevelBack:

	case CLevelForwardOne:
		nextComp := m.NextSibling()
		if nextComp != nil {
			m.MoveTo(nextComp, types.NaInsertBehind)
		}
	case CLevelBackOne:
		prevComp := m.PrevSibling()
		if prevComp != nil {
			m.MoveTo(prevComp, types.NaInsert)
		}
	}
	go triggerUIGeneration(m)
}
