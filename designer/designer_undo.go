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
// See the License for the specific language governing permissions and
// limitations under the License.

package designer

import (
	"sort"
	"time"

	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	uiBean "github.com/energye/designer/uigen/bean"
)

// 设计 - 撤销/恢复(操作快照)
//
// 快照 = 整棵设计组件树的序列化(.ui 同源数据结构, uigen.ExportUIComponentTree)
// + 各容器内组件的相对 Z 序。撤销/恢复 = 清空重建: Free 现有组件树后按快照
// 经 GetDesignerComponent + recoverDesignerComponentProperty 重建, 与
// .ui 布局恢复共用同一管线, 属性(含组件名)逐组件还原。
//
// 捕获点(变更前调用 saveUndoPoint, 变更后状态即"当前"):
//   - 拖动/缩放: OnMouseDown/beginResize 时暂存快照, 真正发生位移时提交
//     (pendingUndo, 普通点击不产生撤销点)
//   - 删除/剪切/粘贴/属性修改/对齐/同尺寸/分布/Z 序/方向键移动
//
// 恢复期间 undoRestoring 置位, 抑制属性更新等路径的撤销点捕获。

const (
	// designerUndoMaxDepth 撤销栈深度上限
	designerUndoMaxDepth = 50
	// designerUndoCoalesceInterval 方向键移动的撤销点合并间隔(连按一次撤销)
	designerUndoCoalesceInterval = 400 * time.Millisecond
)

// designerSnapshot 设计状态快照
type designerSnapshot struct {
	root   uiBean.TUIComponent // FormRoot 序列化(属性+子树)
	zOrder []designerZOrderGroup
}

// designerZOrderGroup 容器内设计组件的相对 Z 序(底层→顶层)
type designerZOrderGroup struct {
	parent string   // 容器组件名, "" = 设计画布(FormRoot 直属)
	comps  []string // 组件名列表, 按当前 LCL 子控件索引升序(底层→顶层)
}

// captureSnapshot 捕获当前设计状态快照
func (m *FormTab) captureSnapshot() *designerSnapshot {
	if m.FormRoot == nil {
		return nil
	}
	return &designerSnapshot{
		root:   ExportUIComponentTree(m.FormRoot),
		zOrder: m.captureZOrder(),
	}
}

// captureZOrder 捕获每个容器内设计组件的相对 Z 序
func (m *FormTab) captureZOrder() []designerZOrderGroup {
	type entry struct {
		name  string
		index int32
	}
	groups := map[*TDesigningComponent][]entry{}
	var walk func(comp *TDesigningComponent)
	walk = func(comp *TDesigningComponent) {
		for _, child := range comp.Child {
			if ctrl := child.Control(); ctrl != nil && ctrl.IsValid() {
				if p := ctrl.Parent(); p != nil && p.IsValid() {
					groups[comp] = append(groups[comp], entry{name: child.Name(), index: p.GetControlIndex(ctrl)})
				}
			}
			walk(child)
		}
	}
	if m.FormRoot != nil {
		walk(m.FormRoot)
	}
	var result []designerZOrderGroup
	for parent, list := range groups {
		sort.Slice(list, func(i, j int) bool { return list[i].index < list[j].index })
		names := make([]string, 0, len(list))
		for _, item := range list {
			names = append(names, item.name)
		}
		parentName := ""
		if parent != m.FormRoot {
			parentName = parent.Name()
		}
		result = append(result, designerZOrderGroup{parent: parentName, comps: names})
	}
	return result
}

// restoreZOrder 按快照恢复各容器内设计组件的相对 Z 序:
// 每组按底层→顶层依次 BringToFront(对齐线/框选等瞬态子控件无固定 Z 序要求)。
func (m *FormTab) restoreZOrder(groups []designerZOrderGroup) {
	for _, g := range groups {
		var parent *TDesigningComponent
		if g.parent == "" {
			parent = m.FormRoot
		} else {
			parent = m.findComponentByName(m.FormRoot, g.parent)
		}
		if parent == nil {
			continue
		}
		for _, name := range g.comps {
			comp := m.findComponentByName(parent, name)
			if comp == nil {
				continue
			}
			if ctrl := comp.Control(); ctrl != nil && ctrl.IsValid() {
				ctrl.BringToFront()
			}
		}
	}
}

// findComponentByName 在组件树中按名称查找组件
func (m *FormTab) findComponentByName(root *TDesigningComponent, name string) *TDesigningComponent {
	if root == nil {
		return nil
	}
	var walk func(comp *TDesigningComponent) *TDesigningComponent
	walk = func(comp *TDesigningComponent) *TDesigningComponent {
		if comp.Name() == name {
			return comp
		}
		for _, child := range comp.Child {
			if found := walk(child); found != nil {
				return found
			}
		}
		return nil
	}
	return walk(root)
}

// clearDesignComponents 清空设计组件树(Free 全部子组件), 保留窗体根节点。
// 同时从组件实例列表移除, 避免实例地址复用后的陈旧映射。
func (m *FormTab) clearDesignComponents() {
	if m.FormRoot == nil {
		return
	}
	children := make([]*TDesigningComponent, len(m.FormRoot.Child))
	copy(children, m.FormRoot.Child)
	for _, comp := range children {
		m.removeCompInstances(comp)
		comp.Remove()
	}
	m.FormRoot.Child = nil
}

// removeCompInstances 递归从组件实例列表移除组件(含非可视包裹实例)
func (m *FormTab) removeCompInstances(comp *TDesigningComponent) {
	if comp == nil {
		return
	}
	m.RemoveComponentFormList(comp.Instance())
	if comp.objectNonWrap != nil {
		m.RemoveComponentFormList(comp.objectNonWrap.Instance())
	}
	for _, child := range comp.Child {
		m.removeCompInstances(child)
	}
}

// saveUndoPoint 捕获撤销点(在任何变更应用之前调用)
func (m *FormTab) saveUndoPoint() {
	if m.undoRestoring || m.State == FtsClose || m.FormRoot == nil {
		return
	}
	m.commitUndoPoint(m.captureSnapshot())
}

// commitUndoPoint 提交一份已捕获的快照到撤销栈(拖动/缩放延迟提交用)
func (m *FormTab) commitUndoPoint(snap *designerSnapshot) {
	if snap == nil || m.undoRestoring {
		return
	}
	m.undoStack = append(m.undoStack, snap)
	if len(m.undoStack) > designerUndoMaxDepth {
		m.undoStack = m.undoStack[len(m.undoStack)-designerUndoMaxDepth:]
	}
	m.redoStack = nil
	m.lastUndoTime = time.Now()
	updateUndoRedoUI()
}

// undoPushDue 方向键移动的撤销点合并判定: 换键或距上次提交超过间隔才需要新快照
func (m *FormTab) undoPushDue(key int32) bool {
	now := time.Now()
	due := key != m.lastUndoKey || now.Sub(m.lastUndoTime) > designerUndoCoalesceInterval
	m.lastUndoKey = key
	return due
}

// Undo 撤销: 恢复到上一个撤销点, 当前状态进入恢复栈
func (m *FormTab) Undo() {
	if m.undoRestoring || len(m.undoStack) == 0 {
		return
	}
	cur := m.captureSnapshot()
	snap := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]
	m.restoreSnapshot(snap)
	if cur != nil {
		m.redoStack = append(m.redoStack, cur)
		if len(m.redoStack) > designerUndoMaxDepth {
			m.redoStack = m.redoStack[len(m.redoStack)-designerUndoMaxDepth:]
		}
	}
	updateUndoRedoUI()
	logs.Debug("撤销完成, 剩余撤销点:", len(m.undoStack))
}

// Redo 恢复: 重做下一个恢复点, 当前状态回到撤销栈
func (m *FormTab) Redo() {
	if m.undoRestoring || len(m.redoStack) == 0 {
		return
	}
	cur := m.captureSnapshot()
	snap := m.redoStack[len(m.redoStack)-1]
	m.redoStack = m.redoStack[:len(m.redoStack)-1]
	m.restoreSnapshot(snap)
	if cur != nil {
		m.undoStack = append(m.undoStack, cur)
		if len(m.undoStack) > designerUndoMaxDepth {
			m.undoStack = m.undoStack[len(m.undoStack)-designerUndoMaxDepth:]
		}
	}
	updateUndoRedoUI()
	logs.Debug("恢复完成, 剩余恢复点:", len(m.redoStack))
}

// CanUndo 是否可撤销
func (m *FormTab) CanUndo() bool {
	return !m.undoRestoring && len(m.undoStack) > 0
}

// CanRedo 是否可恢复
func (m *FormTab) CanRedo() bool {
	return !m.undoRestoring && len(m.redoStack) > 0
}

// restoreSnapshot 按快照重建设计组件树:
// 清空现有组件 → 创建子树(属性暂存) → 恢复属性(含组件名/几何) → 恢复 Z 序。
func (m *FormTab) restoreSnapshot(snap *designerSnapshot) {
	if snap == nil || m.FormRoot == nil {
		return
	}
	prevRestoring := m.undoRestoring
	m.undoRestoring = true
	defer func() { m.undoRestoring = prevRestoring }()
	m.ClearSelection()
	m.clearDesignComponents()
	ProjectTreeBeginUpdate()
	defer ProjectTreeEndUpdate()
	m.FormRoot.RecoverProperty = snap.root.Properties
	for _, child := range snap.root.Child {
		m.createRecoverSubtree(child, m.FormRoot)
	}
	m.RecoverComponentPropertyValue()
	// 延迟设计模式: 属性恢复完成后统一开启(与 .ui 布局恢复一致)
	EnableDeferredDesignMode(m)
	m.restoreZOrder(snap.zOrder)
	// 重建后回到窗体根节点(属性面板/组件树同步), 并全量重新生成 UI
	m.SwitchComponentEditing(m.FormRoot)
	triggerUIGeneration(m.FormRoot, nil, event.CodeGenUI)
	logs.Debug("恢复快照完成, 组件数:", len(snap.root.Child))
}

// updateUndoRedoUI 刷新撤销/恢复的按钮可用状态(工具栏)
func updateUndoRedoUI() {
	tab := activeFormTab()
	canUndo := tab != nil && tab.IsDesigner && tab.CanUndo()
	canRedo := tab != nil && tab.IsDesigner && tab.CanRedo()
	if MainWindow.toolBtnLayout != nil && MainWindow.toolBtnLayout.toolbarBtn != nil {
		MainWindow.toolBtnLayout.toolbarBtn.undoBtn.SetEnabled(canUndo)
		MainWindow.toolBtnLayout.toolbarBtn.redoBtn.SetEnabled(canRedo)
	}
}
