package designer

import "sync/atomic"

// designerAltDown 设计窗体键盘事件维护的 Alt 状态；鼠标消息本身不带 Alt。
var designerAltDown atomic.Bool

func altHeld() bool {
	return designerAltDown.Load() || platformAltDown()
}
