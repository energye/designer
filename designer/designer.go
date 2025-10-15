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
	designer lcl.IDesigner
	form     *FormTab
}

// 创建一个窗体设计器
func NewEngFormDesigner(form *FormTab) *TEngFormDesigner {
	m := new(TEngFormDesigner)
	m.form = form
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
	logs.Debug("OnIsDesignMsg setCursor", message.Msg, sender.ToString())
	//lcl.Screen.SetCursor(types.CrDefault)
	//message.Result = 1
}

func (m *TEngFormDesigner) mouseDown(sender lcl.IControl, message *types.TLMKey) {
	instance := sender.Instance()
	comp := m.form.GetComponentFormList(instance)
	logs.Debug("OnIsDesignMsg mouseDown", message.Msg, sender.ToString(), "选中控件:", comp)

}

func (m *TEngFormDesigner) mouseUp(sender lcl.IControl, message *types.TLMKey) {
	logs.Debug("OnIsDesignMsg mouseUp", message.Msg, sender.ToString())
}

func (m *TEngFormDesigner) mouseMove(sender lcl.IControl, message *types.TLMMouse) {
	logs.Debug("OnIsDesignMsg mouseMove", message.Msg, message, sender.ToString(), "X:", *message.XPos(), "Y:", *message.YPos())
}

func (m *TEngFormDesigner) move(sender lcl.IControl, message *types.TLMMove) {
	logs.Debug("OnIsDesignMsg move", message.Msg, message, sender.ToString(), message.MoveType, *message.XPos(), *message.YPos())
}

func (m *TEngFormDesigner) size(sender lcl.IControl, message *types.TLMSize) {
	logs.Debug("OnIsDesignMsg size", message.Msg, message, sender.ToString())
}

func (m *TEngFormDesigner) paint(sender lcl.IControl, message *types.TLMPaint) {
	logs.Debug("OnIsDesignMsg paint", message.Msg, message, sender.ToString())
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
		paint := (*types.TLMPaint)(unsafe.Pointer(dispatchMsg))
		m.paint(sender, paint)
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
		logs.Debug("OnIsDesignMsg ACTIVATE", message.Msg, isDesign, sender.ToString())
	case messages.LM_CLOSEQUERY:
		logs.Debug("OnIsDesignMsg CLOSEQUERY", message.Msg, isDesign, sender.ToString())
	case messages.LM_SETCURSOR:
		m.setCursor(sender, message)
	case messages.LM_CONTEXTMENU:
		logs.Debug("OnIsDesignMsg CONTEXTMENU", message.Msg, isDesign, sender.ToString())
	case messages.CN_KEYDOWN, messages.CN_SYSKEYDOWN:
		logs.Debug("OnIsDesignMsg KEYDOWN", message.Msg, isDesign, sender.ToString())
	case messages.CN_KEYUP, messages.CN_SYSKEYUP:
		logs.Debug("OnIsDesignMsg KEYUP", message.Msg, isDesign, sender.ToString())
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
