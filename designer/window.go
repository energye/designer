package designer

import (
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
)

func (m *TApp) SetWindowIcon() {
	if iconData := resources.Assets("window-icon_256x256.png"); iconData != nil {
		stream := lcl.NewMemoryStream()
		lcl.StreamHelper.Write(stream, iconData)
		stream.SetPosition(0)
		png := lcl.NewPortableNetworkGraphic()
		png.LoadFromStreamWithStream(stream)
		lcl.Application.Icon().Assign(png)
		png.Free()
		stream.Free()
	}
}

func (m *TApp) ShowInMonitor() {
	// 控制窗口显示鼠标所在显示器
	centerOnMonitor := func(monitor lcl.IMonitor) {
		m.SetLeft(monitor.Left() + (monitor.Width()-m.Width())/2)
		top := monitor.Top() + (monitor.Height()-m.Height())/2
		m.SetTop(top)
	}
	mousePos := lcl.Mouse.CursorPos()
	var (
		i         int32 = 0
		defaultOK       = true
	)
	for ; i < lcl.Screen.MonitorCount(); i++ {
		if tempMonitor := lcl.Screen.Monitors(i); tempMonitor.WorkareaRect().PtInRect(mousePos) {
			defaultOK = false
			centerOnMonitor(tempMonitor)
			break
		}
	}
	if defaultOK {
		centerOnMonitor(lcl.Screen.PrimaryMonitor())
	}
}
