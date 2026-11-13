//go:build windows

package designer

import "syscall"

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

func platformAltDown() bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(0x12)) // VK_MENU
	return int16(r) < 0
}
