// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

var (
	mainWindow       TAppWindow
	bgDarkColor      = colors.RGBToColor(56, 57, 60)
	bgLightColor     = colors.TColor(0xF3F4F6)
	windowShowEvents []func()
	imageActions     *tool.ImageList
	imageComponents  *tool.ImageList
	imageItem        *tool.ImageList
	imageMenu        *tool.ImageList
	imageTabComp     *tool.ImageList
	themeControls    tool.HashMap[lcl.IWinControl]
	splitterWidth    = int32(4)
)

// 设计器应用窗口
type TAppWindow struct {
	lcl.TEngForm
	mainMenu            *TMainMenu                 // 主菜单
	componentProperties lcl.IApplicationProperties //
	box                 *BottomBox                 // 底部布局盒子
	bar                 *StatusBar
}

func SetComponentDefaultColor(control lcl.IWinControl) {
	control.SetColor(bgLightColor)
}

// 添加组件到主题控件集合
func AddComponentTheme(control lcl.IWinControl) {
	themeControls.Add(tool.IntToString(control.Instance()), control)
}

// 切换组件主题
func SwitchAllTheme(dark bool) {
	themeControls.Iterate(func(key string, control lcl.IWinControl) bool {
		if dark {
			control.SetColor(bgDarkColor)
		} else {
			control.SetColor(bgLightColor)
		}
		return false
	})
}

func (m *TAppWindow) FormCreate(sender lcl.IObject) {
	vtedit.MainForm = m
	logs.Info("FormCreate")
	cfg := config.Config
	// 属性
	m.SetCaption(cfg.Title)
	m.SetDoubleBuffered(true)
	m.SetWidth(int32(cfg.Window.Width))
	m.SetHeight(int32(cfg.Window.Height))
	m.SetColor(bgLightColor)
	constra := m.Constraints()
	constra.SetMinWidth(400)
	constra.SetMinHeight(200)
	// 窗口显示在鼠标所在的窗口
	m.showInMonitor()
	m.initAllImageList()
	// 设置窗口图标
	m.setWindowIcon()
	// 窗口显示事件
	m.SetOnShow(m.OnShow)
	// 创建设计器布局
	m.createDesignerLayout()
	// status bar
	//newStatusBar(m)
}

func (m *TAppWindow) initAllImageList() {
	imageActions = tool.NewImageList(m, "actions", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}, Image150: types.TSize{Cx: 24, Cy: 24}, Image200: types.TSize{Cx: 32, Cy: 32}})
	imageComponents = tool.NewImageList(m, "components", tool.ImageRect{Image100: types.TSize{Cx: 24, Cy: 24}, Image150: types.TSize{Cx: 36, Cy: 36}, Image200: types.TSize{Cx: 48, Cy: 48}})
	imageItem = tool.NewImageList(m, "item", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}, Image150: types.TSize{Cx: 24, Cy: 24}, Image200: types.TSize{Cx: 32, Cy: 32}})
	imageMenu = tool.NewImageList(m, "menu", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}, Image150: types.TSize{Cx: 24, Cy: 24}, Image200: types.TSize{Cx: 32, Cy: 32}})
	imageTabComp = tool.NewImageList(m, "tab-comp", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}})
}

func (m *TAppWindow) OnShow(sender lcl.IObject) {
	logs.Info("OnShow")
	// 窗口显示在鼠标所在的窗口
	m.showInMonitor()
	for _, fn := range windowShowEvents {
		fn()
	}
}

func (m *TAppWindow) FormAfterCreate(sender lcl.IObject) {
	logs.Info("FormAfterCreate")
}

func (m *TAppWindow) CreateParams(params *types.TCreateParams) {
	logs.Info("CreateParams")
}

func (m *TAppWindow) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	logs.Info("OnCloseQuery")
}

func (m *TAppWindow) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	logs.Info("OnClose")
	// 取消所有生成事件
	event.CancelAll()
}

func AddOnShow(fn func()) {
	windowShowEvents = append(windowShowEvents, fn)
}
