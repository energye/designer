package message

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

var (
	hint lcl.IHintWindow
)

func mustHint() {
	if hint == nil {
		hint = lcl.NewHintWindow(lcl.Application)
	}
}

func HintShow(content string) {
	mustHint()
	cursorPos := lcl.Mouse.CursorPos()
	displayRect := lcl.Screen.WorkAreaRect()
	x, y := cursorPos.X+15, cursorPos.Y+15
	if x+width > displayRect.Width() {
		x = x - (x + width - displayRect.Width())
	}
	if y+height > displayRect.Height() {
		y = y - (y + height - displayRect.Height())
	}
	mustMessage()
	rect := types.TRect{Left: x, Top: y}
	rect.SetWidth(int32(100))
	rect.SetHeight(int32(40))
	hint.ActivateHintWithRectString(rect, content)
}

func HintHide() {
	if hint != nil {
		hint.Hide()
	}
}
