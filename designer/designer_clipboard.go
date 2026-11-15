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
	"strings"

	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	uiBean "github.com/energye/designer/uigen/bean"
	"github.com/energye/lcl/types"
)

// 设计 - 组件剪切/复制/粘贴(设计器内部剪贴板)
//
// 剪贴板内容为组件子树的序列化数据(与 .ui 布局文件同源结构,
// uigen.ExportUIComponentTree): 属性含 Name/几何/已修改项。
// 粘贴时经 GetDesignerComponent 重建 + RecoverComponentPropertyValue
// 应用属性, 与 .ui 布局恢复共用同一管线; 组件名粘贴时自动重新生成
// (GetComponentCaptionName), 不会与现有组件冲突。

var (
	// designClipboardItems 设计器内部剪贴板: 顶层组件序列化子树列表
	designClipboardItems []uiBean.TUIComponent
	// designClipboardPos 各顶层组件复制时在设计画布客户区的绝对位置
	// (与 items 一一对应; 序列化 Left/Top 是相对原父容器的,
	// 跨容器粘贴必须按绝对位置换算, 见 PasteComponents)
	designClipboardPos []types.TPoint
	// designClipboardPasteSeq 级联粘贴偏移计数(鼠标不在画布内时使用)
	designClipboardPasteSeq int32
)

// designClipboardHasContent 设计器剪贴板是否有内容
func designClipboardHasContent() bool {
	return len(designClipboardItems) > 0
}

// clipboardTopComponents 复制/剪切的顶层组件集合:
// 排除窗体根节点与被其它待复制组件包含的子组件, 包含非可视组件。
func (m *FormTab) clipboardTopComponents() []*TDesigningComponent {
	var selected []*TDesigningComponent
	for _, comp := range m.SelectedComponents() {
		if comp == m.FormRoot || comp.ComponentType == consts.CtForm {
			continue
		}
		if m.HasSelectedAncestor(comp) {
			continue
		}
		selected = append(selected, comp)
	}
	return selected
}

// componentCanvasPos 组件在设计画布客户区的绝对位置(跨容器层级换算)
func (m *FormTab) componentCanvasPos(comp *TDesigningComponent) types.TPoint {
	var pos types.TPoint
	ctrl := comp.Control()
	root := m.FormRoot.Control()
	if ctrl == nil || !ctrl.IsValid() || root == nil || !root.IsValid() {
		return pos
	}
	co, ro := ctrl.ClientOrigin(), root.ClientOrigin()
	return types.TPoint{X: co.X - ro.X, Y: co.Y - ro.Y}
}

// CopySelectedComponents 复制选中组件到设计器剪贴板
func (m *FormTab) CopySelectedComponents() bool {
	targets := m.clipboardTopComponents()
	if len(targets) == 0 {
		return false
	}
	items := make([]uiBean.TUIComponent, 0, len(targets))
	pos := make([]types.TPoint, 0, len(targets))
	for _, comp := range targets {
		item := ExportUIComponentTree(comp)
		// 剥离 Name 属性: 粘贴时组件名由 GetComponentCaptionName 规则
		// 重新生成, 不得沿用被复制组件的名字(避免与原组件重名)
		stripClipboardNames(&item)
		items = append(items, item)
		pos = append(pos, m.componentCanvasPos(comp))
	}
	designClipboardItems = items
	designClipboardPos = pos
	designClipboardPasteSeq = 0
	logs.Debug("复制组件:", len(items))
	return true
}

// stripClipboardNames 递归移除剪贴板组件树中的 Name 属性。
// Name 缺失时属性恢复阶段跳过改名, 组件保持创建时自动生成的名字;
// TDesigningComponent.Name() 回退读取 m.name(自动名), 组件树节点显示一致。
func stripClipboardNames(item *uiBean.TUIComponent) {
	item.Properties = stripNameProps(item.Properties)
	for i := range item.Child {
		stripClipboardNames(&item.Child[i])
	}
}

// stripNameProps 移除属性列表中的 Name 属性(不区分大小写)
func stripNameProps(props []uiBean.TProperty) []uiBean.TProperty {
	result := make([]uiBean.TProperty, 0, len(props))
	for _, prop := range props {
		if strings.EqualFold(prop.Name, "name") {
			continue
		}
		result = append(result, prop)
	}
	return result
}

// CutSelectedComponents 剪切选中组件: 复制到剪贴板后删除
func (m *FormTab) CutSelectedComponents() bool {
	targets := m.clipboardTopComponents()
	if len(targets) == 0 {
		return false
	}
	if !m.CopySelectedComponents() {
		return false
	}
	m.saveUndoPoint()
	parents := m.removeDesignComponents(targets)
	if len(parents) == 1 {
		for parent := range parents {
			m.SwitchComponentEditing(parent)
		}
	}
	return true
}

// PasteComponents 粘贴剪贴板组件到设计画布:
// 鼠标在画布内时以鼠标位置为粘贴点, 否则级联偏移。
func (m *FormTab) PasteComponents() bool {
	if !designClipboardHasContent() {
		return false
	}
	m.saveUndoPoint()
	tx, ty := m.pasteTargetPos()
	// 整组保持相对布局平移到粘贴点: 原包围盒左上角(画布绝对坐标) → 目标位置。
	// 深拷贝剪贴板数据: 位置改写/Caption 同步都是就地修改,
	// 浅拷贝会污染剪贴板本体(下次粘贴从被改写的数据出发)。
	items := cloneUIComponents(designClipboardItems)
	minLeft, minTop := designClipboardMinPos()
	dx, dy := tx-minLeft, ty-minTop
	for i := range items {
		if isVisualType(items[i].Type) && i < len(designClipboardPos) {
			// 可视组件: Left/Top 重写为 画布绝对位置 + 组级平移。
			// 序列化 Left/Top 是相对原父容器的, 多选复制可能来自不同容器,
			// 粘贴统一挂到画布(FormRoot)下, 必须按绝对位置换算,
			// 否则面板内的组件粘贴到画布后位置偏移(面板原点差)。
			setPosProps(items[i].Properties, designClipboardPos[i].X+dx, designClipboardPos[i].Y+dy)
		} else {
			// 非可视组件: 包裹图标本就位于画布层, 仅随组平移
			shiftPosProps(items[i].Properties, dx, dy)
		}
	}
	var (
		created      []*TDesigningComponent
		createdItems []*uiBean.TUIComponent
	)
	for i := range items {
		if comp := m.createRecoverSubtree(items[i], m.FormRoot); comp != nil {
			created = append(created, comp)
			createdItems = append(createdItems, &items[i])
		}
	}
	if len(created) == 0 {
		logs.Warn("粘贴组件失败: 未创建任何组件")
		return false
	}
	// 应用属性(含几何/已修改项), 与 .ui 布局恢复同一管线。
	// 必须先恢复属性再开启设计模式(与加载路径一致):
	// 设计模式下每次真实 bounds 变更走组件消息通知路径(20ms+),
	// 先开启会让组件在旧位置停留数十毫秒后再跳到目标位置(可见闪跳)。
	m.RecoverComponentPropertyValue()
	EnableDeferredDesignMode(m)
	// Caption 同步: 默认状态(原 Caption == 原组件名)的组件
	// Caption 跟随新名字, 自定义 Caption 保持不变(递归含子组件)
	for i, comp := range created {
		syncPasteCaptions(comp, createdItems[i])
	}
	// 选中粘贴的顶层组件(主选组件经 SelectComponents→LoadPropertyToInspector
	// 加载检查器; 多选粘贴时其余组件没有属性页, 不能逐个加载检查器,
	// 否则 UpdatePropertyTreeWidth 解引用 nil propertyTree 崩溃)
	m.SelectComponents(created, created[0])
	triggerUIGeneration(m.FormRoot, nil, event.CodeGenUI)
	designClipboardPasteSeq++
	logs.Debug("粘贴组件:", len(created))
	return true
}

// cloneUIComponents 深拷贝序列化组件树(Properties/Child 底层数组一并复制)
func cloneUIComponents(items []uiBean.TUIComponent) []uiBean.TUIComponent {
	out := make([]uiBean.TUIComponent, len(items))
	for i, item := range items {
		item.Properties = append([]uiBean.TProperty(nil), item.Properties...)
		item.Child = cloneUIComponents(item.Child)
		out[i] = item
	}
	return out
}

// syncPasteCaptions 递归同步一组粘贴组件的 Caption(与剪贴板条目一一对应)
func syncPasteCaptions(comp *TDesigningComponent, item *uiBean.TUIComponent) {
	if comp == nil || item == nil {
		return
	}
	syncPasteCaption(comp)
	for i, child := range comp.Child {
		if i < len(item.Child) {
			syncPasteCaptions(child, &item.Child[i])
		}
	}
}

// syncPasteCaption 单个组件的 Caption 同步:
// 拥有 Caption 属性的组件(按钮/标签等显示 Caption 的控件),
// 粘贴后 Caption 一律更新为粘贴后的新名字, 与 Name 保持一致,
// 避免复制出来的组件与原组件显示相同文本;
// 没有 Caption 属性的组件(如 TEdit 使用 Text)不处理。
func syncPasteCaption(comp *TDesigningComponent) {
	if comp == nil || comp.PropertyList == nil {
		return
	}
	for _, prop := range comp.PropertyList {
		if !strings.EqualFold(prop.Name(), "caption") {
			continue
		}
		// 更新属性节点数据(.ui 序列化数据源)并调用 API 更新显示
		prop.SetEditValue(comp.Name())
		comp.recoverCallAPI("Caption", prop)
		return
	}
}

// createRecoverSubtree 从序列化组件树创建设计组件并挂接父子关系/组件树节点。
// 属性不在此处应用(暂存 RecoverProperty), 由 RecoverComponentPropertyValue 统一恢复;
// 但几何值(Left/Top/Width/Height)在此处先行提取 —— 直接在最终位置创建组件,
// 否则组件先以默认位置出现在画布左上角、属性恢复后再跳到目标位置,
// 期间差值会被画出来(复制粘贴/撤销恢复的"左上角闪现")。
func (m *FormTab) createRecoverSubtree(item uiBean.TUIComponent, parent *TDesigningComponent) *TDesigningComponent {
	left, top, width, height := recoverGeometry(item)
	// GetDesignerComponent 拒绝负坐标, 先夹到 0; 负值随后由属性恢复阶段应用
	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}
	newDesComp := GetDesignerComponent(m, left, top, item.ClassName)
	if newDesComp == nil {
		logs.Warn("创建设计组件失败, 未注册或未实现, 类名:", item.ClassName)
		return nil
	}
	if width > 0 && height > 0 {
		// 立即应用最终尺寸(设计模式未开启, SetBounds 走快速路径)
		newDesComp.SetBounds(left, top, width, height)
	}
	newDesComp.RecoverProperty = item.Properties
	newDesComp.SetParent(parent)
	parent.AddChild(newDesComp)
	for _, child := range item.Child {
		m.createRecoverSubtree(child, newDesComp)
	}
	return newDesComp
}

// recoverGeometry 从序列化属性中提取几何值(Left/Top/Width/Height)
func recoverGeometry(item uiBean.TUIComponent) (left, top, width, height int32) {
	for _, prop := range item.Properties {
		switch strings.ToLower(prop.Name) {
		case "left":
			left = anyToInt32(prop.Value)
		case "top":
			top = anyToInt32(prop.Value)
		case "width":
			width = anyToInt32(prop.Value)
		case "height":
			height = anyToInt32(prop.Value)
		}
	}
	return
}

// pasteTargetPos 粘贴目标位置(设计画布客户区坐标):
// 鼠标在画布内→鼠标位置; 否则从画布左上角级联偏移。
func (m *FormTab) pasteTargetPos() (int32, int32) {
	if pos, ok := overlayGuardMousePos(m); ok && posInDesignCanvas(m, pos) {
		return pos.X, pos.Y
	}
	off := designClipboardPasteSeq * 16
	return 24 + off, 24 + off
}

// designClipboardMinPos 计算剪贴板组件原包围盒左上角
// (取可视组件复制时画布绝对位置的最小值)
func designClipboardMinPos() (int32, int32) {
	minLeft, minTop := int32(0), int32(0)
	first := true
	for i := range designClipboardItems {
		if !isVisualType(designClipboardItems[i].Type) || i >= len(designClipboardPos) {
			continue
		}
		l, t := designClipboardPos[i].X, designClipboardPos[i].Y
		if first || l < minLeft {
			minLeft = l
		}
		if first || t < minTop {
			minTop = t
		}
		first = false
	}
	return minLeft, minTop
}

// isVisualType 是否按可视组件处理几何(常规可视 + 自定义可视)
func isVisualType(t consts.ComponentType) bool {
	return t == consts.CtVisual || t == consts.CtCustomVisual
}

// shiftPosProps 平移一组属性中的 Left/Top(粘贴位置调整)
func shiftPosProps(props []uiBean.TProperty, dx, dy int32) {
	for i := range props {
		switch strings.ToLower(props[i].Name) {
		case "left":
			props[i].Value = anyToInt32(props[i].Value) + dx
		case "top":
			props[i].Value = anyToInt32(props[i].Value) + dy
		}
	}
}

// setPosProps 重写一组属性中的 Left/Top(粘贴绝对位置换算)
func setPosProps(props []uiBean.TProperty, left, top int32) {
	for i := range props {
		switch strings.ToLower(props[i].Name) {
		case "left":
			props[i].Value = left
		case "top":
			props[i].Value = top
		}
	}
}

// removeDesignComponents 删除一组顶层设计组件(调用方负责撤销点捕获),
// 触发受影响父容器的 UI 重新生成, 返回受影响父容器集合。
func (m *FormTab) removeDesignComponents(top []*TDesigningComponent) map[*TDesigningComponent]struct{} {
	m.ClearSelection()
	parents := make(map[*TDesigningComponent]struct{})
	for _, comp := range top {
		if comp.parent == nil {
			continue
		}
		parents[comp.parent] = struct{}{}
		comp.Remove()
	}
	for parent := range parents {
		triggerUIGeneration(parent, nil, event.CodeGenUI)
	}
	return parents
}
