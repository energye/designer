// ----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// # Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// ----------------------------------------

package designer

import (
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/application"
)

var (
	trayExit           *application.TTrayMenuItem
	trayShowMainWindow *application.TTrayMenuItem
)

func (m *TAppWindow) initTray() {
	tray := application.NewTrayIcon()
	trayMenu := tray.Menu()
	trayShowMainWindow = trayMenu.AddMenuItem("显示主界面").SetOnClick(func() {
		m.Show()
	})
	trayMenu.AddSeparator()
	trayExit = trayMenu.AddMenuItem("退出 Designer").SetOnClick(func() {
		m.Close()
	})
	trayIconData := resources.Images("icons/window-icon_64x64.png")
	tray.SetIconBytes(trayIconData)
	tray.Show()
}
