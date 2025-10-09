package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
	"unsafe"
)

// TEngFormDesigner energy 窗体设计器
type TEngFormDesigner struct {
	designer lcl.IDesigner
}

// 创建一个窗体设计器
func NewEngFormDesigner() *TEngFormDesigner {
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

// message

func (m *TEngFormDesigner) setCursor(sender lcl.IControl, message *types.TLMessage) {
	//fmt.Println("OnIsDesignMsg setCursor", message.Msg, sender.ToString())
}

func (m *TEngFormDesigner) mouseDown(sender lcl.IControl, message *types.TLMKey) {
	fmt.Println("OnIsDesignMsg mouseDown", message.Msg, sender.ToString())
}

func (m *TEngFormDesigner) mouseUp(sender lcl.IControl, message *types.TLMKey) {
	fmt.Println("OnIsDesignMsg mouseUp", message.Msg, sender.ToString())
}

func (m *TEngFormDesigner) mouseMove(sender lcl.IControl, message *types.TLMMouse) {
	fmt.Println("OnIsDesignMsg mouseMove", message.Msg, message, sender.ToString(), "X:", *message.XPos(), "Y:", *message.YPos())
}

func (m *TEngFormDesigner) move(sender lcl.IControl, message *types.TLMMove) {
	fmt.Println("OnIsDesignMsg move", message.Msg, message, sender.ToString())
}

func (m *TEngFormDesigner) size(sender lcl.IControl, message *types.TLMSize) {
	fmt.Println("OnIsDesignMsg size", message.Msg, message, sender.ToString())
}

func (m *TEngFormDesigner) paint(sender lcl.IControl, message *types.TLMPaint) {
	fmt.Println("OnIsDesignMsg paint", message.Msg, message, sender.ToString())
}

// on event

func (m *TEngFormDesigner) onIsDesignMsg(sender lcl.IControl, message *types.TLMessage) bool {
	isDesign := sender.ComponentState().In(types.CsDesigning)
	result := true
	dispatchMsg := (*uintptr)(unsafe.Pointer(message))
	_ = dispatchMsg
	//sender.Dispatch(dispatchMsg)
	switch message.Msg {
	case messages.LM_PAINT:
		//paint := (*types.TLMPaint)(unsafe.Pointer(dispatchMsg))
		//m.paint(sender, paint)
	case messages.LM_LBUTTONDOWN, messages.LM_RBUTTONDOWN, messages.LM_LBUTTONDBLCLK:
		key := (*types.TLMKey)(unsafe.Pointer(dispatchMsg))
		m.mouseDown(sender, key)
	case messages.LM_LBUTTONUP, messages.LM_RBUTTONUP:
		key := (*types.TLMKey)(unsafe.Pointer(dispatchMsg))
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
		fmt.Println("OnIsDesignMsg ACTIVATE", message.Msg, isDesign, sender.ToString())
	case messages.LM_CLOSEQUERY:
		fmt.Println("OnIsDesignMsg CLOSEQUERY", message.Msg, isDesign, sender.ToString())
	case messages.LM_SETCURSOR:
		m.setCursor(sender, message)
	case messages.LM_CONTEXTMENU:
		fmt.Println("OnIsDesignMsg CONTEXTMENU", message.Msg, isDesign, sender.ToString())
	case messages.CN_KEYDOWN, messages.CN_SYSKEYDOWN:
		fmt.Println("OnIsDesignMsg KEYDOWN", message.Msg, isDesign, sender.ToString())
	case messages.CN_KEYUP, messages.CN_SYSKEYUP:
		fmt.Println("OnIsDesignMsg KEYUP", message.Msg, isDesign, sender.ToString())
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
