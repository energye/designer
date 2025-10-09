package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
)

type TFormDesigner struct {
	lcl.IDesigner
}

func NewCustomFormDesigner() *TFormDesigner {
	m := new(TFormDesigner)
	designer := lcl.NewDesigner()
	m.IDesigner = designer
	designer.SetOnIsDesignMsg(func(sender lcl.IControl, message *types.TLMessage) bool {
		isDesign := sender.ComponentState().In(types.CsDesigning)
		switch message.Msg {
		case messages.LM_PAINT:
			fmt.Println("OnIsDesignMsg PAINT", message.Msg, isDesign)
		case messages.LM_LBUTTONDOWN, messages.LM_RBUTTONDOWN, messages.LM_LBUTTONDBLCLK:
			fmt.Println("OnIsDesignMsg LBUTTONDOWN", message.Msg, isDesign)
		case messages.LM_LBUTTONUP, messages.LM_RBUTTONUP:
			fmt.Println("OnIsDesignMsg LBUTTONUP", message.Msg, isDesign)
		case messages.LM_MOUSEMOVE:
			fmt.Println("OnIsDesignMsg MOUSEMOVE", message.Msg, isDesign)
		case messages.LM_SIZE:
			fmt.Println("OnIsDesignMsg SIZE", message.Msg, isDesign)
		case messages.LM_MOVE:
			fmt.Println("OnIsDesignMsg MOVE", message.Msg, isDesign)
		case messages.LM_ACTIVATE:
			fmt.Println("OnIsDesignMsg ACTIVATE", message.Msg, isDesign)
		case messages.LM_CLOSEQUERY:
			fmt.Println("OnIsDesignMsg CLOSEQUERY", message.Msg, isDesign)
		case messages.LM_SETCURSOR:
			fmt.Println("OnIsDesignMsg SETCURSOR", message.Msg, isDesign)
		case messages.LM_CONTEXTMENU:
			fmt.Println("OnIsDesignMsg CONTEXTMENU", message.Msg, isDesign)
		case messages.CN_KEYDOWN, messages.CN_SYSKEYDOWN:
			fmt.Println("OnIsDesignMsg KEYDOWN", message.Msg, isDesign)
		case messages.CN_KEYUP, messages.CN_SYSKEYUP:
			fmt.Println("OnIsDesignMsg KEYUP", message.Msg, isDesign)
		default:
			return false
		}
		return true
	})
	return m
}

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
