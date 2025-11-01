// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"github.com/energye/lcl/types"
	"unsafe"
)

// 设计 - 组件设计树数据加载

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
		parentComp := m.Parent()
		if parentComp == nil {
			return
		}
		lastChild := parentComp.LastChild()
		if lastChild != nil {
			m.MoveTo(lastChild, types.NaInsertBehind)
		}
	case CLevelBack:
		parentComp := m.Parent()
		if parentComp == nil {
			return
		}
		firstChild := parentComp.FirstChild()
		if firstChild != nil {
			m.MoveTo(firstChild, types.NaInsert)
		}
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
	// TODO 排序所有控件显示 Z 序, 显示有些问题, laz也有
	//control := m.WinControl()
	//parent := control.Parent()
	//if parent != nil && parent.IsValid() {
	//	for i := 0; i < int(parent.ControlCount()); i++ {
	//		child := parent.Controls(int32(i))
	//		println("ControlIndex:", parent.GetControlIndex(child), child.Name())
	//	}
	//}
	go triggerUIGeneration(m)
}
