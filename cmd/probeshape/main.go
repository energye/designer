// LCL SetShapeWithRegion 跨平台能力探针:
// 验证 liblcl TWinControl.SetShapeWithRegion(TRegion) 在 Windows 上
// 是否等价于 SetWindowRgn (GetWindowRgn 返回 1/2 即生效, 0=无区域)。
// 仅用于本地诊断, 不参与设计器构建。
package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

var (
	user32Mod         = syscall.NewLazyDLL("user32.dll")
	procGetWindowRgn  = user32Mod.NewProc("GetWindowRgn")
	gdi32Mod          = syscall.NewLazyDLL("gdi32.dll")
	procCreateRectRgn = gdi32Mod.NewProc("CreateRectRgn")
	procDeleteObject  = gdi32Mod.NewProc("DeleteObject")
)

type ShapeProbe struct {
	lcl.TEngForm
	pane lcl.IPanel
}

// regionState 读取窗口当前区域: 0=NULLREGION/错误, 1=SIMPLEREGION, 2=COMPLEXREGION
func (m *ShapeProbe) regionState() int {
	h, _, _ := procCreateRectRgn.Call(0, 0, 0, 0)
	if h == 0 {
		return -1
	}
	got, _, _ := procGetWindowRgn.Call(m.pane.Handle(), h)
	procDeleteObject.Call(h)
	return int(got)
}

func (m *ShapeProbe) applyRegion(rects [][4]int32) {
	rgn := lcl.NewRegion()
	for _, r := range rects {
		rgn.AddRectangle(r[0], r[1], r[2], r[3])
	}
	m.pane.SetShapeWithRegion(rgn)
}

func (m *ShapeProbe) FormCreate(sender lcl.IObject) {
	m.SetWidth(800)
	m.SetHeight(600)
	scroll := lcl.NewScrollBox(m)
	scroll.SetParent(m)
	scroll.SetAlign(types.AlClient)
	m.pane = lcl.NewPanel(scroll)
	m.pane.SetParent(scroll)
	m.pane.SetAlign(types.AlClient)
	m.pane.BringToFront()

	fmt.Println("before regionState:", m.regionState())
	m.applyRegion([][4]int32{{10, 10, 30, 30}, {100, 50, 140, 90}})
	fmt.Println("PROBE-1 after 2 rects regionState:", m.regionState())
	m.applyRegion([][4]int32{{40, 40, 60, 60}})
	fmt.Println("PROBE-2 after 1 rect regionState:", m.regionState())
	m.applyRegion(nil)
	fmt.Println("PROBE-3 after empty regionState:", m.regionState())
	fmt.Println("PROBE DONE")
}

func main() {
	libname.LibName = lib.ExtractLibrary(os.TempDir())
	lcl.Init()
	lcl.Application.Initialize()
	lcl.Application.SetTitle("probeshape")
	var probe ShapeProbe
	lcl.Application.NewForms(&probe)
	lcl.Application.Run()
}
