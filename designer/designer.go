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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/api/misc"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/messages"
	"unsafe"
)

// TEngFormDesigner energy 窗体设计器
type TEngFormDesigner struct {
	designer           lcl.IDesigner
	componentList      map[uintptr]*TDesigningComponent // 设计中的组件列表, key: 组件实例ID, value: 设计组件
	canvas             lcl.ICanvas
	Form               lcl.IWinControl
	LookupRoot         *TDesigningComponent
	MouseMoveComponent *TDesigningComponent
	MouseDownComponent *TDesigningComponent
	MouseDownPos       types.TPoint
	MouseUpPos         types.TPoint
	LastMouseMovePos   types.TPoint
	band               designerBand
}

// 创建一个窗体设计器
func NewEngFormDesigner(formTab *FormTab) *TEngFormDesigner {
	m := new(TEngFormDesigner)
	newDesigner := lcl.NewDesigner()
	newDesigner.SetOnIsDesignMsg(m.designerOnIsDesignMsg)
	newDesigner.SetOnUTF8KeyPress(m.designerOnUTF8KeyPress)
	newDesigner.SetOnModified(m.designerOnModified)
	newDesigner.SetOnNotification(m.designerOnNotification)
	newDesigner.SetOnPaintGrid(m.designerOnPaintGrid)
	newDesigner.SetOnValidateRename(m.designerOnValidateRename)
	newDesigner.SetOnGetShiftState(m.designerOnGetShiftState)
	newDesigner.SetOnSelectOnlyThisComponent(m.designerOnSelectOnlyThisComponent)
	newDesigner.SetOnUniqueName(m.designerOnUniqueName)
	newDesigner.SetOnPrepareFreeDesigner(m.designerOnPrepareFreeDesigner)
	m.designer = newDesigner
	m.canvas = lcl.NewCanvas()
	return m
}

func (m *TEngFormDesigner) Free() {
	m.designer.SetOnIsDesignMsg(nil)
	m.designer.SetOnUTF8KeyPress(nil)
	m.designer.SetOnModified(nil)
	m.designer.SetOnNotification(nil)
	m.designer.SetOnPaintGrid(nil)
	m.designer.SetOnValidateRename(nil)
	m.designer.SetOnGetShiftState(nil)
	m.designer.SetOnSelectOnlyThisComponent(nil)
	m.designer.SetOnUniqueName(nil)
	m.designer.SetOnPrepareFreeDesigner(nil)
	m.designer.Free()
	m.componentList = nil
	m.canvas.Free()
}

// 返回设计器实例接口
func (m *TEngFormDesigner) Designer() lcl.IDesigner {
	return m.designer
}

// 添加设计组件到组件列表
func (m *TEngFormDesigner) AddComponentToList(component *TDesigningComponent) {
	if m.componentList == nil {
		m.componentList = make(map[uintptr]*TDesigningComponent)
	}
	//  以下在添加组件时根据不同类型使用实际设计可控制的组件对象

	if component.ComponentType == consts.CtForm {
		m.componentList[component.object.Instance()] = component
	} else if component.ComponentType == consts.CtNonVisual {
		m.componentList[component.objectNonWrap.Instance()] = component
	} else {
		m.componentList[component.Instance()] = component
	}
}

// 返回设计组件
func (m *TEngFormDesigner) GetComponentFormList(instance uintptr) *TDesigningComponent {
	if instance == 0 || m.componentList == nil {
		return nil
	}
	return m.componentList[instance]
}

// 删除一个设计组件
func (m *TEngFormDesigner) RemoveComponentFormList(instance uintptr) {
	delete(m.componentList, instance)
}

// message

// setCursor 设计模式 LM_SETCURSOR: 命中缩放控制点→设置缩放光标;
// 未命中→恢复默认箭头。必须显式管理: LCL 对设计模式窗口跳过默认
// 光标流程(WindowProcHelper.HandleSetCursor 的 csDesigning 分支),
// 且不为 WM_SETCURSOR 调 DefWindowProc(类光标), 不接管则光标停留。
// 光标经 Screen 全局立即生效(见 setOverlayCursor 注释),
// 残留由 overlayCursorGuard(OnUserInput)兜底复位。
// 始终返回 true(接管本消息, 阻止默认流程)。
func (m *TEngFormDesigner) setCursor(sender lcl.IControl, message *types.TLMessage) bool {
	formTab := m.formTab()
	if formTab == nil {
		return true
	}
	if comp := m.MouseDownComponent; comp != nil && comp.drag != nil && comp.drag.resizeDir >= 0 {
		// 缩放过程中固定缩放光标
		setOverlayCursor(comp.drag.resizeDir)
		message.Result = 1
		return true
	}
	if comp, dir, ok := formTab.dragHandleAt(m.GetMousePos()); ok && comp.drag != nil {
		setOverlayCursor(dir)
		message.Result = 1
		return true
	}
	resetOverlayCursor()
	message.Result = 1
	return true
}

// formTab 返回设计器所属的设计窗体
func (m *TEngFormDesigner) formTab() *FormTab {
	if m.LookupRoot == nil {
		return nil
	}
	return m.LookupRoot.FormTab
}

// messagePosToCanvas 将鼠标消息坐标(相对发送者控件客户区)换算为设计画布客户区坐标。
// 换算 = 发送者客户区原点 - 画布客户区原点 + 消息坐标, 两者同为 LCL 坐标体系,
// DPI 缩放/滚动位置下保持一致; 对任意发送者(组件/画布/滚动容器)均正确。
func (m *TEngFormDesigner) messagePosToCanvas(sender lcl.IControl, X, Y int32) types.TPoint {
	var so, co types.TPoint
	if sender != nil && sender.IsValid() {
		so = sender.ClientOrigin()
	}
	if root := m.LookupRoot; root != nil && root.Control() != nil && root.Control().IsValid() {
		co = root.Control().ClientOrigin()
	}
	return types.TPoint{X: so.X - co.X + X, Y: so.Y - co.Y + Y}
}

func (m *TEngFormDesigner) mouseDown(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("Designer message mouseDown", message.Msg, sender.ToString(), message)
	m.MouseDownPos = m.messagePosToCanvas(sender, int32(*message.XPos()), int32(*message.YPos()))
	// 上一次按下的 MouseUp 可能被缩放控制点截获, 先兜底结束残留的移动拖拽状态
	if m.MouseDownComponent != nil && m.MouseDownComponent.drag != nil {
		m.MouseDownComponent.drag.cancelStaleDrag()
	}
	// 缩放控制点命中(控制点绘制在组件外侧画布上, 命中优先于组件):
	// 只响应左键, 右键仍走组件菜单流程。
	if formTab := m.formTab(); formTab != nil && message.Msg == messages.LM_LBUTTONDOWN {
		comp, dir, hit := formTab.dragHandleAt(m.MouseDownPos)
		logs.Debug("PERF mouseDown handleHit:", hit, "dir:", dir, "pos:", m.MouseDownPos.X, m.MouseDownPos.Y,
			"sender:", sender.Name(), sender.ClassName())
		if hit && comp.drag != nil {
			m.MouseDownComponent = comp
			comp.drag.beginResize(dir, m.MouseDownPos)
			*message.Result() = 1
			return
		}
	}
	m.MouseDownComponent = m.GetComponentAtPos(m.MouseDownPos.X, m.MouseDownPos.Y)
	shift, button := m.GetMouseMsgShift(message)
	if button == types.MbLeft && SelectedComponent() == nil &&
		(m.MouseDownComponent == nil || m.MouseDownComponent.ComponentType == consts.CtForm) {
		// 空白画布按下：开始框选候选，单击在抬起时选中窗体
		m.bandBegin(m.MouseDownPos)
		*message.Result() = 1
		return
	}
	if m.MouseDownComponent != nil {
		// 统一传递设计画布坐标: 与 mouseMove/组件边界(BoundsRect)同一坐标系,
		// 避免按下(发送者相对坐标)与移动(画布坐标)混用导致组件跳位。
		pos := m.messagePosToCanvas(sender, int32(*message.XPos()), int32(*message.YPos()))
		m.MouseDownComponent.OnMouseDown(sender, button, shift, pos.X, pos.Y)
	}
	*message.Result() = 1
}

func (m *TEngFormDesigner) mouseUp(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("Designer message mouseUp", message.Msg, sender.ToString())
	shift, button := m.GetMouseMsgShift(message)
	if button == types.MbLeft && m.band.down {
		m.bandEnd(m.messagePosToCanvas(sender, int32(*message.XPos()), int32(*message.YPos())))
		*message.Result() = 1
		m.MouseDownComponent = nil
		return
	}
	control := m.MouseDownComponent
	if control == nil {
		m.MouseUpPos = m.messagePosToCanvas(sender, int32(*message.XPos()), int32(*message.YPos()))
		control = m.GetComponentAtPos(m.MouseUpPos.X, m.MouseUpPos.Y)
	}
	if control != nil {
		x, y := int32(*message.XPos()), int32(*message.YPos())
		control.OnMouseUp(sender, button, shift, x, y)
	}
	*message.Result() = 1
	m.MouseDownComponent = nil
}

func (m *TEngFormDesigner) mouseMove(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("Designer message mouseMove", message.Msg, message, sender.ToString())
	m.MouseMoveComponent = m.MouseDownComponent
	// 左键已不在按下状态: 上一次 MouseUp 可能被缩放控制点截获,
	// 兜底结束残留的移动拖拽, 避免后续 mousemove 驱动陈旧拖拽移动整组组件。
	if m.MouseDownComponent != nil && message.Keys&consts.MK_LBUTTON == 0 {
		if m.MouseDownComponent.drag != nil {
			m.MouseDownComponent.drag.cancelStaleDrag()
		}
		m.MouseDownComponent = nil
		m.MouseMoveComponent = nil
		*message.Result() = 1
		return
	}
	if m.band.down {
		bandPos := m.messagePosToCanvas(sender, int32(*message.XPos()), int32(*message.YPos()))
		if message.Keys&consts.MK_LBUTTON == 0 {
			// 鼠标按键已在窗口外释放，按当前位置结束框选
			m.bandEnd(bandPos)
		} else {
			m.bandUpdate(bandPos)
		}
		*message.Result() = 1
		return
	}
	if m.MouseMoveComponent != nil {
		shift, _ := m.GetMouseMsgShift(message)
		pos := m.messagePosToCanvas(sender, int32(*message.XPos()), int32(*message.YPos()))
		m.MouseMoveComponent.OnMouseMove(sender, shift, pos.X, pos.Y)
	}
	*message.Result() = 1
}

func (m *TEngFormDesigner) move(sender lcl.IControl, message *types.TLMMove) {
	//logs.Debug("Designer message move", message.Msg, message, sender.ToString(), message.MoveType, *message.XPos(), *message.YPos())
}

func (m *TEngFormDesigner) size(sender lcl.IControl, message *types.TLMSize) {
	//logs.Debug("Designer message size", message.Msg, message, sender.ToString())
}

func (m *TEngFormDesigner) paint(sender lcl.IControl, message *types.TLMPaint) {
	isAcceptsControl := sender.ControlStyle().In(types.CsAcceptsControls)
	//logs.Debug("Designer message paint", message.Msg, "DC:", message.DC, sender.ToString(), "isAcceptsControl:", isAcceptsControl)
	//comp := m.GetComponentFormList(sender.Instance())
	if message.DC != 0 && isAcceptsControl {
		m.canvas.SetHandle(message.DC)
		pen := m.canvas.PenToPen()
		pen.SetColor(bgDarkColor)
		pen.SetWidth(1)
		pen.SetStyle(types.PsSolid)
		r := sender.ClientRect()
		misc.DrawGrid(m.canvas.Handle(), r, 8, 8)
		m.canvas.SetHandle(0)
	}
}
func (m *TEngFormDesigner) popupMenu(sender lcl.IControl, message *types.TLMContextMenu) {
	//instance := sender.Instance()
	//comp := m.GetComponentFormList(instance)
	//if comp != nil {
	//	x, y := int32(*message.XPos()), int32(*message.YPos())
	//	if x == -1 && y == -1 {
	//		point := sender.ClientToScreenWithPoint(types.TPoint{}) // todo types.TPoint{} 参数需要传递一个正确的值
	//		x, y = point.X, point.Y
	//	}
	//	comp.FormTab.componentMenu.treePopupMenu.PopUpWithIntX2(x, y)
	//}
	*message.Result() = 1
}

func (m *TEngFormDesigner) GetMouseMsgShift(message *types.TLMMouse) (shift types.TShiftState, button types.TMouseButton) {
	if message.Keys&consts.MK_SHIFT == consts.MK_SHIFT {
		shift = shift.Include(types.SsShift)
	}
	if message.Keys&consts.MK_CONTROL == consts.MK_CONTROL {
		shift = shift.Include(types.SsCtrl)
	}
	if message.Keys&consts.MK_ALT == consts.MK_ALT {
		shift = shift.Include(types.SsAlt)
	}
	switch message.Msg {
	case messages.LM_LBUTTONUP, messages.LM_LBUTTONDBLCLK, messages.LM_LBUTTONTRIPLECLK, messages.LM_LBUTTONQUADCLK:
		// 左键状态由 button 返回，SsShift 仅表示真实 Shift 键，避免普通点击被当成 Shift+点击。
		button = types.MbLeft
	case messages.LM_MBUTTONUP, messages.LM_MBUTTONDBLCLK, messages.LM_MBUTTONTRIPLECLK, messages.LM_MBUTTONQUADCLK:
		shift = shift.Include(types.SsMiddle)
		button = types.MbMiddle
	case messages.LM_RBUTTONUP, messages.LM_RBUTTONDBLCLK, messages.LM_RBUTTONTRIPLECLK, messages.LM_RBUTTONQUADCLK:
		shift = shift.Include(types.SsRight)
		button = types.MbRight
	default:
		if message.Keys&consts.MK_MBUTTON != 0 {
			shift = shift.Include(types.SsMiddle)
			button = types.MbMiddle
		} else if message.Keys&consts.MK_RBUTTON != 0 {
			shift = shift.Include(types.SsRight)
			button = types.MbRight
		} else if message.Keys&consts.MK_LBUTTON != 0 {
			button = types.MbLeft
		} else if message.Keys&consts.MK_XBUTTON1 != 0 {
			shift = shift.Include(types.SsExtra1)
			button = types.MbExtra1
		} else if message.Keys&consts.MK_XBUTTON2 != 0 {
			shift = shift.Include(types.SsExtra2)
			button = types.MbExtra2
		}
	}
	switch message.Msg {
	case messages.LM_LBUTTONDBLCLK, messages.LM_MBUTTONDBLCLK, messages.LM_RBUTTONDBLCLK, messages.LM_XBUTTONDBLCLK:
		shift = shift.Include(types.SsDouble)
	case messages.LM_LBUTTONTRIPLECLK, messages.LM_MBUTTONTRIPLECLK, messages.LM_RBUTTONTRIPLECLK, messages.LM_XBUTTONTRIPLECLK:
		shift = shift.Include(types.SsTriple)
	case messages.LM_LBUTTONQUADCLK, messages.LM_MBUTTONQUADCLK, messages.LM_RBUTTONQUADCLK, messages.LM_XBUTTONQUADCLK:
		shift = shift.Include(types.SsQuad)
	}
	return
}

// paintContainerOverlays 在每个设计模式组件自身绘制完成后叠加覆盖层:
//   - 次选组件的选中边框
//   - 主选组件的缩放控制点(仅画布控制点; 窗体根节点为窗体外侧面板, handleRects 为空)
//
// 绘制面 = 组件所属容器的画布 + 与覆盖层矩形重叠的其它组件表面(补丁):
// 解决 Z 序遮挡 —— 被前面组件覆盖的后面的组件, 其锚点/边框绘制在前面的组件表面,
// 依旧可见可命中(命中测试见 dragHandleAt, 优先于选中前面的组件)。
// Z 序由绘制顺序自然解析: 最上层的组件最后绘制, 补丁最终落在最上层表面。
// 每个表面重绘时自动重画补丁; 几何变化后由 invalidateOverlaySurfaces 失效相关表面。
func (m *TEngFormDesigner) paintContainerOverlays(sender lcl.IControl, paint *types.TLMPaint) {
	if paint == nil || paint.DC == 0 {
		return
	}
	formTab := m.formTab()
	if formTab == nil || formTab.selectionDragActive {
		return
	}
	comp := formTab.GetComponentFormList(sender.Instance())
	if comp == nil {
		return
	}
	primary := formTab.PrimarySelected()
	secondary := formTab.overlaySecondary()
	if primary == nil && len(secondary) == 0 {
		return
	}
	ctrl := comp.Control()
	if ctrl == nil || !ctrl.IsValid() {
		return
	}
	canvas := m.canvas
	if canvas == nil || !canvas.IsValid() {
		logs.Debug("PERF paintOverlays canvas invalid")
		return
	}
	// 本组件客户区矩形(过滤不相交的覆盖层)
	cw, ch := ctrl.ClientWidth(), ctrl.ClientHeight()
	canvas.SetHandle(paint.DC)
	brush := canvas.BrushToBrush()
	brush.SetColor(designerSelectionColor)
	for _, x := range secondary {
		xctrl := x.Control()
		if xctrl == nil || !xctrl.IsValid() {
			continue
		}
		br := xctrl.BoundsRect()
		off := formTab.overlayOffset(x, comp)
		l, t := br.Left+off.X, br.Top+off.Y
		w, h := br.Width(), br.Height()
		if l >= cw || t >= ch || l+w <= 0 || t+h <= 0 {
			continue
		}
		canvas.FillRectWithRect(types.Rect(l-1, t-1, l+w+1, t))     // top
		canvas.FillRectWithRect(types.Rect(l-1, t+h, l+w+1, t+h+1)) // bottom
		canvas.FillRectWithRect(types.Rect(l-1, t, l, t+h))         // left
		canvas.FillRectWithRect(types.Rect(l+w, t, l+w+1, t+h))     // right
	}
	if primary != nil && primary.drag != nil {
		rects := primary.drag.handleRects()
		if len(rects) > 0 {
			off := formTab.overlayOffset(primary, comp)
			brush.SetColor(colors.ClBlack)
			for _, item := range rects {
				l, t := item.rect.Left+off.X, item.rect.Top+off.Y
				w, h := item.rect.Width(), item.rect.Height()
				if l >= cw || t >= ch || l+w <= 0 || t+h <= 0 {
					continue
				}
				canvas.FillRectWithRect(types.Rect(l, t, l+w, t+h))
			}
		}
	}
	canvas.SetHandle(0)
}

// on event

func (m *TEngFormDesigner) designerOnIsDesignMsg(sender lcl.IControl, message *types.TLMessage) bool {
	result := false // 标记行为 false: 未处理, true: 已处理
	isDesign := sender.ComponentState().In(types.CsDesigning)
	//println("IsDesignMsg", message.Msg, "isDesign:", isDesign)
	if isDesign {
		result = true
		dispatchMsg := (*uintptr)(unsafe.Pointer(message))
		switch message.Msg {
		case messages.LM_PAINT:
			paint := *(*types.TLMPaint)(unsafe.Pointer(dispatchMsg))
			pintPtr := uintptr(unsafe.Pointer(&paint))
			sender.Dispatch(&pintPtr)
			// 覆盖层平台(Win32/GTK2/Qt5): 锚点/边框绘制在置顶覆盖层上, 无需画布补丁。
			// 其余平台(Cocoa 无窗口区域实现): 容器画布绘制回退。
			if !overlayPaneSupported() {
				m.paintContainerOverlays(sender, &paint)
			}
		case messages.LM_LBUTTONDOWN, messages.LM_RBUTTONDOWN, messages.LM_LBUTTONDBLCLK:
			key := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
			m.mouseDown(sender, key)
		case messages.LM_LBUTTONUP, messages.LM_RBUTTONUP:
			key := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
			m.mouseUp(sender, key)
			//m.Form.SetFocus()
		case messages.LM_MOUSEMOVE:
			mouse := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
			m.mouseMove(sender, mouse)
		case messages.LM_SIZE:
			size := (*types.TLMSize)(unsafe.Pointer(dispatchMsg))
			m.size(sender, size)
		case messages.LM_MOVE:
			move := (*types.TLMMove)(unsafe.Pointer(dispatchMsg))
			m.move(sender, move)
		case messages.LM_ACTIVATE:
			//logs.Debug("IsDesignMsg ACTIVATE", message.Msg, isDesign, sender.ToString())
		case messages.LM_CLOSEQUERY:
			//logs.Debug("IsDesignMsg CLOSEQUERY", message.Msg, isDesign, sender.ToString())
		case messages.LM_SETCURSOR:
			// 未命中控制点: 交给 LCL 默认光标流程(控件 Cursor 属性恢复箭头)
			if !m.setCursor(sender, message) {
				result = false
			}
		case messages.LM_CONTEXTMENU:
			//println("IsDesignMsg CONTEXTMENU", message.Msg, sender.ToString(), sender.Name())
			contextMenu := (*types.TLMContextMenu)(unsafe.Pointer(dispatchMsg))
			m.popupMenu(sender, contextMenu)
		case messages.CN_KEYDOWN, messages.CN_SYSKEYDOWN:
			// 不再在这里无条件吞掉按键：设计快捷键由 DesignFormOnKeyDown 处理，
			// 其余按键交给菜单快捷键，避免“有时好用有时不好用”。
			result = false
		case messages.CN_KEYUP, messages.CN_SYSKEYUP:
			//println("IsDesignMsg KEYUP", message.Msg, sender.ToString(), sender.Name())
			//case messages.LM_HSCROLL, messages.LM_VSCROLL:
		case messages.LM_SETFOCUS:
			//println("messages.LM_SETFOCUS")
			m.Form.SetFocus() // TODO 防止控件获得焦点, 还没找到其它处理方试
			//if m.MouseDownComponent != nil && sender.Instance() == m.MouseDownComponent.Instance() {
			//	println("setfocus", message.Msg, sender.ToString(), sender.Name(), m.Form.CanFocus(), m.Form.CanSetFocus())
			//	if api.IsObjectInstanceOf(sender.Instance(), lcl.TCustomComboBoxClass()) ||
			//		api.IsObjectInstanceOf(sender.Instance(), lcl.TCustomEditClass()) ||
			//		api.IsObjectInstanceOf(sender.Instance(), lcl.TCustomListBoxClass()) {
			//	}
			//}
			message.Result = 1
		//case messages.CM_DESIGNHITTEST:
		//	logs.Debug("designhittest")
		//	message.Result = 1
		default:
			// 其它让组件自动处理消息
			result = false
		}
	}
	return result
}

func (m *TEngFormDesigner) designerOnUTF8KeyPress(uTF8Key *string) {
	//println("designerOnUTF8KeyPress")
	*uTF8Key = ""
}

func (m *TEngFormDesigner) designerOnModified() {
	//println("designerOnModified")
}

func (m *TEngFormDesigner) designerOnNotification(component lcl.IComponent, operation types.TOperation) {
	//println("designerOnNotification operation:", operation)
}

func (m *TEngFormDesigner) designerOnPaintGrid() {
	//println("designerOnPaintGrid")
}

func (m *TEngFormDesigner) designerOnValidateRename(component lcl.IComponent, curName string, newName string) {
	//println("designerOnValidateRename")
}

func (m *TEngFormDesigner) designerOnGetShiftState() types.TShiftState {
	//println("designerOnGetShiftState")
	return types.NewSet()
}

func (m *TEngFormDesigner) designerOnSelectOnlyThisComponent(component lcl.IComponent) {
	//println("designerOnSelectOnlyThisComponent")
}

func (m *TEngFormDesigner) designerOnUniqueName(baseName string) string {
	//println("designerOnUniqueName")
	return ""
}

func (m *TEngFormDesigner) designerOnPrepareFreeDesigner(freeComponent bool) {
	//println("designerOnPrepareFreeDesigner")
}

// GetDesignControl 获取设计时控件的可选择父控件
// 当控件处于设计模式下时，某些控件可能被标记为不可选择，
// 此函数会向上查找直到找到可选择的父控件为止
//
//	control: 需要检查的设计控件
func (m *TEngFormDesigner) GetDesignControl(control lcl.IControl) lcl.IControl {
	if control == nil || control.Instance() == m.Form.Instance() /*|| control.Owner().Instance() == m.Form.Instance() */ {
		return control
	}
	if control.Instance() == m.Form.Instance() {
		return control
	}
	owner := control.Owner()
	isControl := api.IsObjectInstanceOf(owner.Instance(), lcl.TControlClass())
	if isControl {
		ownerControl := lcl.AsControl(owner)
		if !ownerControl.ControlStyle().In(types.CsOwnedChildrenNotSelectable) {
			return control
		}
		return m.GetDesignControl(ownerControl)
	} else {
		return control
	}
}

// GetComponentAtPos 根据给定的坐标位置查找对应的组件
//
//	x - 相对于根组件的X坐标
//	y - 相对于根组件的Y坐标
//
//	返回在指定位置的最上层组件，如果未找到则返回nil
func (m *TEngFormDesigner) GetComponentAtPos(x, y int32) *TDesigningComponent {
	var (
		root  = m.LookupRoot
		found *TDesigningComponent
	)
	var iteratorChild func(*TDesigningComponent, int32)
	iteratorChild = func(component *TDesigningComponent, level int32) {
		br := m.GetRelativeBounds(component)
		in := br.PtInRect(types.Point(x, y))
		if in {
			found = component
		}
		for _, child := range component.Child {
			iteratorChild(child, level+1)
		}
	}
	iteratorChild(root, 0)
	return found
}

// GetRelativeBounds 获取指定设计组件相对于设计器的边界矩形
//
//	component: 要获取边界的设计器组件指针
//	result: 包含左、上、右、下坐标的矩形结构体
func (m *TEngFormDesigner) GetRelativeBounds(component *TDesigningComponent) (result types.TRect) {
	topLeft := m.GetRelativeTopLeft(component)
	result.Left = topLeft.X
	result.Top = topLeft.Y
	result.Right = result.Left + component.Control().Width()
	result.Bottom = result.Top + component.Control().Height()
	return
}

// GetRelativeTopLeft 计算组件相对于表单设计根节点的左上角坐标
// component: 要计算位置的设计组件
func (m *TEngFormDesigner) GetRelativeTopLeft(component *TDesigningComponent) (result types.TPoint) {
	form := m.LookupRoot
	parent := component.Parent()
	if form == nil || parent == nil {
		return types.NewPoint(0, 0)
	}
	result = parent.Control().ClientOrigin()
	formOrigin := form.Control().ClientOrigin()
	result.X = result.X - formOrigin.X + component.Control().Left()
	result.Y = result.Y - formOrigin.Y + component.Control().Top()
	return
}

// GetMousePos 获取鼠标在窗体设计器中的相对位置
//
//	result - 鼠标相对于窗体客户区的坐标点
func (m *TEngFormDesigner) GetMousePos() (result types.TPoint) {
	form := m.LookupRoot
	mousePos := lcl.Mouse.CursorPos()
	formClientOrigin := form.Control().ClientOrigin()
	result.X = mousePos.X - formClientOrigin.X
	result.Y = mousePos.Y - formClientOrigin.Y
	return
}

// 启用或禁用功能组件
func SetEnableFuncComponent(enable bool) {
	MainWindow.mainMenu.SetEnableMenuItems(enable)
	if MainWindow.toolBtnLayout != nil {
		MainWindow.toolBtnLayout.toolbarBtn.SetEnableToolButtons(enable)
	}
}

// 启用或禁用保存功能
//func SetEnableSaveFunc(enable bool) {
//	lcl.RunOnMainThreadAsync(func(id uint32) {
//		//mainWindow.mainMenu.save.SetEnabled(enable)
//		//toolbar.toolbarBtn.saveBtn.SetEnabled(enable)
//	})
//}

// 全局初始化, 非线程安全
func Init() {
	// 注册控制台消息处理器事件
	initConsoleEvent()
	// 注册项目源码扫描处理器事件
	initProjectSrcEvent()
	// 注册组件
	initRegisterComponent()

	// 初始化依赖模块信息 ast
	//dependmod.InitDependencyModule()
	// 提取所有启用的框架
	//frameworks.ExtractFrameworks()
}

// initGlobalCursorGuard 注册锚点缩放光标残留兜底(需在 lcl.Application
// 初始化后调用)。双通道复位: OnUserInput 每次用户输入即时复位(先于
// 控件事件); 定时器 120ms 轮询覆盖不派发用户输入事件的窗口路径
// (部分应用面板的鼠标移动不触发 OnUserInput)。
func initGlobalCursorGuard() {
	lcl.Application.SetOnUserInput(func(_ lcl.IObject, _ uint32) {
		overlayCursorRestore()
	})
	cursorTimer := lcl.NewTimer(nil)
	cursorTimer.SetInterval(120)
	cursorTimer.SetOnTimer(func(lcl.IObject) {
		overlayCursorRestore()
	})
}
