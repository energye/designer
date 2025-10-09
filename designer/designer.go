package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
	"unsafe"
)

// 加载图片列表
func LoadImageList(owner lcl.IComponent, imageList []string, width, height int32) lcl.IImageList {
	images := lcl.NewImageList(owner)
	images.SetWidth(width)
	images.SetHeight(height)
	for _, image := range imageList {
		tool.ImageListAddPng(images, image)
	}
	return images
}

type TEngFormDesigner struct {
	designer lcl.IDesigner
}

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

func (m *TEngFormDesigner) Designer() lcl.IDesigner {
	return m.designer
}

func (m *TEngFormDesigner) onIsDesignMsg(sender lcl.IControl, message *types.TLMessage) bool {
	isDesign := sender.ComponentState().In(types.CsDesigning)
	result := true
	dispatchMsg := (*uintptr)(unsafe.Pointer(message))
	_ = dispatchMsg
	switch message.Msg {
	case messages.LM_PAINT:
		paint := (*types.TLMPaint)(unsafe.Pointer(dispatchMsg))
		fmt.Println("OnIsDesignMsg PAINT", message.Msg, isDesign, paint, sender.ToString())
	case messages.LM_LBUTTONDOWN, messages.LM_RBUTTONDOWN, messages.LM_LBUTTONDBLCLK:
		key := (*types.TLMKey)(unsafe.Pointer(dispatchMsg))
		fmt.Println("OnIsDesignMsg LBUTTONDOWN", message.Msg, isDesign, key, sender.ToString())
	case messages.LM_LBUTTONUP, messages.LM_RBUTTONUP:
		key := (*types.TLMKey)(unsafe.Pointer(dispatchMsg))
		fmt.Println("OnIsDesignMsg LBUTTONUP", message.Msg, isDesign, key, sender.ToString())
	case messages.LM_MOUSEMOVE:
		mouse := (*types.TLMMouse)(unsafe.Pointer(dispatchMsg))
		fmt.Println("OnIsDesignMsg MOUSEMOVE", message.Msg, isDesign, "mouse:", mouse, sender.ToString(), "X:", *mouse.XPos(), "Y:", *mouse.YPos())
	case messages.LM_SIZE:
		size := (*types.TLMSize)(unsafe.Pointer(dispatchMsg))
		fmt.Println("OnIsDesignMsg SIZE", message.Msg, isDesign, size, sender.ToString())
	case messages.LM_MOVE:
		move := (*types.TLMMove)(unsafe.Pointer(dispatchMsg))
		fmt.Println("OnIsDesignMsg MOVE", message.Msg, isDesign, move, sender.ToString())
	case messages.LM_ACTIVATE:
		fmt.Println("OnIsDesignMsg ACTIVATE", message.Msg, isDesign, sender.ToString())
	case messages.LM_CLOSEQUERY:
		fmt.Println("OnIsDesignMsg CLOSEQUERY", message.Msg, isDesign, sender.ToString())
	case messages.LM_SETCURSOR:
		fmt.Println("OnIsDesignMsg SETCURSOR", message.Msg, isDesign, sender.ToString())
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
	println("onNotification")
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
