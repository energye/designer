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
	"github.com/energye/lcl/types/messages"
	"unsafe"
)

// TEngFormDesigner energy 窗体设计器
type TEngFormDesigner struct {
	designer           lcl.IDesigner
	componentList      map[uintptr]*TDesigningComponent // 设计中的组件列表, key: 组件实例ID, value: 设计组件
	canvas             lcl.ICanvas
	Form               lcl.IWinControl
	LookupRoot         lcl.IComponent
	MouseMoveComponent lcl.IComponent
	MouseDownComponent lcl.IComponent
}

// 创建一个窗体设计器
func NewEngFormDesigner(formTab *FormTab) *TEngFormDesigner {
	m := new(TEngFormDesigner)
	newDesigner := lcl.NewDesigner()
	newDesigner.SetOnIsDesignMsg(m.onIsDesignMsg)
	newDesigner.SetOnUTF8KeyPress(m.onUTF8KeyPress)
	newDesigner.SetOnModified(m.onModified)
	newDesigner.SetOnNotification(m.onNotification)
	newDesigner.SetOnPaintGrid(m.onPaintGrid)
	newDesigner.SetOnValidateRename(m.onValidateRename)
	newDesigner.SetOnGetShiftState(m.onGetShiftState)
	newDesigner.SetOnSelectOnlyThisComponent(m.onSelectOnlyThisComponent)
	newDesigner.SetOnUniqueName(m.onUniqueName)
	newDesigner.SetOnPrepareFreeDesigner(m.onPrepareFreeDesigner)
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

func (m *TEngFormDesigner) setCursor(sender lcl.IControl, message *types.TLMessage) {
	//logs.Debug("Design Msg setCursor", message.Msg, sender.ToString())
	lcl.Screen.SetCursor(types.CrDefault)
	message.Result = 1
}

func (m *TEngFormDesigner) mouseDown(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("Designer message mouseDown", message.Msg, sender.ToString(), message)
	m.MouseDownComponent = m.GetDesignControl(sender)
	instance := m.MouseDownComponent.Instance()
	comp := m.GetComponentFormList(instance)
	if comp != nil {
		shift, button := m.GetMouseMsgShift(message)
		x, y := int32(*message.XPos()), int32(*message.YPos())
		comp.OnMouseDown(sender, button, shift, x, y)
	}
	*message.Result() = 1
}

func (m *TEngFormDesigner) mouseUp(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("Designer message mouseUp", message.Msg, sender.ToString())
	control := m.MouseDownComponent
	if control == nil {
		control = m.GetDesignControl(sender)
	}
	instance := control.Instance()
	comp := m.GetComponentFormList(instance)
	if comp != nil {
		shift, button := m.GetMouseMsgShift(message)
		x, y := int32(*message.XPos()), int32(*message.YPos())
		comp.OnMouseUp(sender, button, shift, x, y)
	}
	*message.Result() = 1
	m.MouseDownComponent = nil
}

func (m *TEngFormDesigner) mouseMove(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("Designer message mouseMove", message.Msg, message, sender.ToString())
	m.MouseMoveComponent = m.MouseDownComponent
	if m.MouseMoveComponent == nil {
		m.MouseMoveComponent = m.GetDesignControl(sender)
	}

	control := m.MouseMoveComponent
	instance := control.Instance()
	comp := m.GetComponentFormList(instance)
	if comp != nil {
		shift, _ := m.GetMouseMsgShift(message)
		x, y := int32(*message.XPos()), int32(*message.YPos())
		comp.OnMouseMove(sender, shift, x, y)
	}
	//*message.Result() = 1
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
	switch message.Msg {
	case messages.LM_LBUTTONUP, messages.LM_LBUTTONDBLCLK, messages.LM_LBUTTONTRIPLECLK, messages.LM_LBUTTONQUADCLK:
		shift = shift.Include(types.SsShift)
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
			shift = shift.Include(types.SsShift)
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

// on event

func (m *TEngFormDesigner) onIsDesignMsg(sender lcl.IControl, message *types.TLMessage) bool {
	//logs.Debug("IsDesignMsg", message.Msg)
	result := false // 标记行为 false: 未处理, true: 已处理
	isDesign := sender.ComponentState().In(types.CsDesigning)
	if isDesign {
		result = true
		dispatchMsg := (*uintptr)(unsafe.Pointer(message))
		switch message.Msg {
		case messages.LM_PAINT:
			//println("paint")
			paint := *(*types.TLMPaint)(unsafe.Pointer(dispatchMsg))
			pintPtr := uintptr(unsafe.Pointer(&paint))
			sender.Dispatch(&pintPtr)
			//m.paint(sender, &paint)
		case messages.LM_LBUTTONDOWN, messages.LM_RBUTTONDOWN, messages.LM_LBUTTONDBLCLK:
			//println("down", sender.ToString())
			key := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
			m.mouseDown(sender, key)
		case messages.LM_LBUTTONUP, messages.LM_RBUTTONUP:
			//println("up")
			key := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
			m.mouseUp(sender, key)
		case messages.LM_MOUSEMOVE:
			//println("move", sender.ToString())
			mouse := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
			m.mouseMove(sender, mouse)
		case messages.LM_SIZE:
			//println("size")
			size := (*types.TLMSize)(unsafe.Pointer(dispatchMsg))
			m.size(sender, size)
		case messages.LM_MOVE:
			move := (*types.TLMMove)(unsafe.Pointer(dispatchMsg))
			m.move(sender, move)
		case messages.LM_ACTIVATE:
			//logs.Debug("Designer message ACTIVATE", message.Msg, isDesign, sender.ToString())
		case messages.LM_CLOSEQUERY:
			//logs.Debug("Designer message CLOSEQUERY", message.Msg, isDesign, sender.ToString())
		case messages.LM_SETCURSOR:
			//m.setCursor(sender, message)
		case messages.LM_CONTEXTMENU:
			logs.Debug("Designer message CONTEXTMENU", message.Msg, sender.ToString())
			contextMenu := (*types.TLMContextMenu)(unsafe.Pointer(dispatchMsg))
			m.popupMenu(sender, contextMenu)
		case messages.CN_KEYDOWN, messages.CN_SYSKEYDOWN:
			logs.Debug("Designer message KEYDOWN", message.Msg, sender.ToString())
		case messages.CN_KEYUP, messages.CN_SYSKEYUP:
			logs.Debug("Designer message KEYUP", message.Msg, sender.ToString())
		//case messages.LM_HSCROLL, messages.LM_VSCROLL:
		case messages.LM_ENTER, messages.LM_SETFOCUS:
			//println("enter", sender.ToString(), m.Form.CanFocus(), m.Form.CanSetFocus())
			if api.IsObjectInstanceOf(sender.Instance(), lcl.TCustomComboBoxClass()) ||
				api.IsObjectInstanceOf(sender.Instance(), lcl.TCustomEditClass()) ||
				api.IsObjectInstanceOf(sender.Instance(), lcl.TCustomListBoxClass()) {
				m.Form.SetFocus()
			}
		default:
			// 其它让组件自动处理消息
			result = false
		}
	}
	return result
}

func (m *TEngFormDesigner) onUTF8KeyPress(uTF8Key *string) {
	println("onUTF8KeyPress")
}

func (m *TEngFormDesigner) onModified() {
	println("onModified")
}

func (m *TEngFormDesigner) onNotification(component lcl.IComponent, operation types.TOperation) {
	//println("onNotification")
}

func (m *TEngFormDesigner) onPaintGrid() {
	//println("onPaintGrid")
}

func (m *TEngFormDesigner) onValidateRename(component lcl.IComponent, curName string, newName string) {
	println("onValidateRename")
}

func (m *TEngFormDesigner) onGetShiftState() types.TShiftState {
	println("onGetShiftState")
	return types.NewSet()
}

func (m *TEngFormDesigner) onSelectOnlyThisComponent(component lcl.IComponent) {
	println("onSelectOnlyThisComponent")
}

func (m *TEngFormDesigner) onUniqueName(baseName string) string {
	println("onUniqueName")
	return ""
}

func (m *TEngFormDesigner) onPrepareFreeDesigner(freeComponent bool) {
	println("onPrepareFreeDesigner")
}

func (m *TEngFormDesigner) GetDesignControl(control lcl.IControl) lcl.IControl {
	if control == nil || control.Instance() == m.LookupRoot.Instance() /*|| control.Owner().Instance() == m.LookupRoot.Instance() */ {
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

func (m *TEngFormDesigner) GetDesignedComponent(component lcl.IComponent) lcl.IComponent {
	return component
}

// 启用或禁用功能组件
func SetEnableFuncComponent(enable bool) {
	mainWindow.mainMenu.SetEnableMenuItems(enable)
	if toolbar != nil {
		toolbar.toolbarBtn.SetEnableToolButtons(enable)
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
func init() {
	// 注册控制台消息处理器事件
	initConsoleEvent()
	// 初始化工具栏相关配置
	initConfigToolbar()
	// 注册组件
	initRegisterComponent()
	initRegisterComponent2()
}
