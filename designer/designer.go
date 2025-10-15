package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
	"unsafe"
)

// TEngFormDesigner energy 窗体设计器
type TEngFormDesigner struct {
	designer      lcl.IDesigner
	componentList map[uintptr]*DesigningComponent // 设计中的组件列表, key: 组件实例ID, value: 设计组件
}

// 创建一个窗体设计器
func NewEngFormDesigner(form *FormTab) *TEngFormDesigner {
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
	return m
}

// 返回设计器实例接口
func (m *TEngFormDesigner) Designer() lcl.IDesigner {
	return m.designer
}

// 添加设计组件到组件列表
func (m *TEngFormDesigner) AddComponentToList(component *DesigningComponent) {
	if m.componentList == nil {
		m.componentList = make(map[uintptr]*DesigningComponent)
	}
	m.componentList[component.Instance()] = component
}

// 返回设计组件
func (m *TEngFormDesigner) GetComponentFormList(instance uintptr) *DesigningComponent {
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
	//logs.Debug("OnIsDesignMsg setCursor", message.Msg, sender.ToString())
	//lcl.Screen.SetCursor(types.CrDefault)
	//message.Result = 1
}

func (m *TEngFormDesigner) mouseDown(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("OnIsDesignMsg mouseDown", message.Msg, sender.ToString(), message)
	instance := sender.Instance()
	comp := m.GetComponentFormList(instance)
	if comp != nil {
		shift, button := m.GetMouseMsgShift(message)
		x, y := int32(*message.XPos()), int32(*message.YPos())
		comp.OnMouseDown(sender, button, shift, x, y)
	}
}

func (m *TEngFormDesigner) mouseUp(sender lcl.IControl, message *types.TLMMouse) {
	//logs.Debug("OnIsDesignMsg mouseUp", message.Msg, sender.ToString())
	instance := sender.Instance()
	comp := m.GetComponentFormList(instance)
	if comp != nil {
		shift, button := m.GetMouseMsgShift(message)
		x, y := int32(*message.XPos()), int32(*message.YPos())
		comp.OnMouseUp(sender, button, shift, x, y)
	}
}

func (m *TEngFormDesigner) mouseMove(sender lcl.IControl, message *types.TLMMouse) {
	logs.Debug("OnIsDesignMsg mouseMove", message.Msg, message, sender.ToString())
	instance := sender.Instance()
	comp := m.GetComponentFormList(instance)
	if comp != nil {
		shift, _ := m.GetMouseMsgShift(message)
		x, y := int32(*message.XPos()), int32(*message.YPos())
		comp.OnMouseMove(sender, shift, x, y)
	}
	println("isNil:", comp == nil)
}

func (m *TEngFormDesigner) move(sender lcl.IControl, message *types.TLMMove) {
	//logs.Debug("OnIsDesignMsg move", message.Msg, message, sender.ToString(), message.MoveType, *message.XPos(), *message.YPos())
}

func (m *TEngFormDesigner) size(sender lcl.IControl, message *types.TLMSize) {
	//logs.Debug("OnIsDesignMsg size", message.Msg, message, sender.ToString())
}

func (m *TEngFormDesigner) paint(sender lcl.IControl, message *types.TLMPaint) {
	//`logs.Debug("OnIsDesignMsg paint", message.Msg, message, sender.ToString())
}

const (
	// Mouse message key states
	MK_LBUTTON  = 1
	MK_RBUTTON  = 2
	MK_SHIFT    = 4
	MK_CONTROL  = 8
	MK_MBUTTON  = 0x10
	MK_XBUTTON1 = 0x20
	MK_XBUTTON2 = 0x40
	// following are "virtual" key states
	MK_DOUBLECLICK = 0x80
	MK_TRIPLECLICK = 0x100
	MK_QUADCLICK   = 0x200
	MK_ALT         = 0x20000000
)

func (m *TEngFormDesigner) GetMouseMsgShift(message *types.TLMMouse) (shift types.TShiftState, button types.TMouseButton) {
	if message.Keys&MK_SHIFT == MK_SHIFT {
		shift = shift.Include(types.SsShift)
	}
	if message.Keys&MK_CONTROL == MK_CONTROL {
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
		if message.Keys&MK_MBUTTON != 0 {
			shift = shift.Include(types.SsMiddle)
			button = types.MbMiddle
		} else if message.Keys&MK_RBUTTON != 0 {
			shift = shift.Include(types.SsRight)
			button = types.MbRight
		} else if message.Keys&MK_LBUTTON != 0 {
			shift = shift.Include(types.SsShift)
			button = types.MbLeft
		} else if message.Keys&MK_XBUTTON1 != 0 {
			shift = shift.Include(types.SsExtra1)
			button = types.MbExtra1
		} else if message.Keys&MK_XBUTTON2 != 0 {
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
	logs.Debug("IsDesignMsg", message.Msg)
	//isDesign := sender.ComponentState().In(types.CsDesigning)
	result := true
	dispatchMsg := (*uintptr)(unsafe.Pointer(message))
	//_ = dispatchMsg
	//sender.Dispatch(dispatchMsg)
	switch message.Msg {
	case messages.LM_PAINT:
		paint := (*types.TLMPaint)(unsafe.Pointer(dispatchMsg))
		m.paint(sender, paint)
	case messages.LM_LBUTTONDOWN, messages.LM_RBUTTONDOWN, messages.LM_LBUTTONDBLCLK:
		key := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
		m.mouseDown(sender, key)
	case messages.LM_LBUTTONUP, messages.LM_RBUTTONUP:
		key := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
		m.mouseUp(sender, key)
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
		//logs.Debug("OnIsDesignMsg ACTIVATE", message.Msg, isDesign, sender.ToString())
	case messages.LM_CLOSEQUERY:
		//logs.Debug("OnIsDesignMsg CLOSEQUERY", message.Msg, isDesign, sender.ToString())
	case messages.LM_SETCURSOR:
		m.setCursor(sender, message)
	case messages.LM_CONTEXTMENU:
		//logs.Debug("OnIsDesignMsg CONTEXTMENU", message.Msg, isDesign, sender.ToString())
	case messages.CN_KEYDOWN, messages.CN_SYSKEYDOWN:
		//logs.Debug("OnIsDesignMsg KEYDOWN", message.Msg, isDesign, sender.ToString())
	case messages.CN_KEYUP, messages.CN_SYSKEYUP:
		//logs.Debug("OnIsDesignMsg KEYUP", message.Msg, isDesign, sender.ToString())
	default:
		result = false
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
	println("onPaintGrid")
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
