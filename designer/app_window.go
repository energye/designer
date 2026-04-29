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
	"fmt"
	"github.com/energye/designer/designer/dependmod"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"sync"
	"time"
)

var (
	MainWindow       TAppWindow
	bgDarkColor      = colors.RGBToColor(56, 57, 60)
	bgLightColor     = colors.ClWhite // colors.TColor(0xF3F4F6)
	windowShowEvents []func()
	imageActions     *tool.ImageList
	imageComponents  *tool.ImageList
	imageItem        *tool.ImageList
	imageMenu        *tool.ImageList
	imageTabComp     *tool.ImageList
	themeControls    tool.HashMap[string, lcl.IWinControl]
	splitterWidth    = int32(5)
	leftToolsWidth   = int32(110)
	gOnShow          = sync.Once{}
	gAppEGPPath      string
)

// 设计器应用窗口
type TAppWindow struct {
	lcl.TEngForm
	mainMenu              *TMainMenu                 // 主菜单
	toolBtnLayout         *TToolBtnLayout            // 工具具栏
	contentLayout         *ContentLayout             // 底部布局盒子
	openDialog            lcl.IOpenDialog            // 打开对话框
	saveDialog            lcl.ISaveDialog            // 保存对话框
	selectDirectoryDialog lcl.ISelectDirectoryDialog // 选择文件夹对话框
	closing               bool                       // 关闭中
}

func SetComponentDefaultColor(control lcl.IWinControl) {
	control.SetColor(bgLightColor)
}

// 添加组件到主题控件集合
func AddComponentTheme(control lcl.IWinControl) {
	themeControls.Add(tool.IntToString(control.Instance()), control)
}

// 设置应用配置文件路径
func SetAppEGPPath(path string) {
	gAppEGPPath = path
}

// 更新设计器窗口标题
// 打开项目后, 新建项目后
func UpdateDesignerTitle(title string) {
	var (
		windowTitle string
	)
	if title == "" {
		windowTitle = fmt.Sprintf("%v %v", config.DesignerConfig.Title, config.DesignerConfig.Version)
	} else {
		windowTitle = fmt.Sprintf("%v %v - %v", title, config.DesignerConfig.Title, config.DesignerConfig.Version)
	}
	lcl.RunOnMainThreadAsync(func(id uint32) {
		logs.Println("UpdateDesignerTitle:", windowTitle)
		MainWindow.SetCaption(windowTitle)
		ProjectTreeSetProjectName(title)
	})
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
	logs.Window = m // 用于调试, 窗口指针问题
	logs.Info("Designer FormCreate")
	cfg := config.DesignerConfig
	if cfg.WindowLayout.WindowBoundsRect.Width() <= 400 {
		cfg.WindowLayout.WindowBoundsRect.SetWidth(400)
	}
	if cfg.WindowLayout.WindowBoundsRect.Height() <= 200 {
		cfg.WindowLayout.WindowBoundsRect.SetHeight(200)
	}
	// 属性
	m.SetCaption(cfg.Title + " " + cfg.Version)
	m.SetDoubleBuffered(true)
	m.SetBoundsRect(cfg.WindowLayout.WindowBoundsRect)
	m.SetColor(bgLightColor)
	constra := m.Constraints()
	constra.SetMinWidth(400)
	constra.SetMinHeight(200)
	if cfg.WindowLayout.WindowBoundsRect.Left == 0 && cfg.WindowLayout.WindowBoundsRect.Top == 0 {
		// 窗口显示在鼠标所在的窗口
		m.showInMonitor()
	}
	if cfg.WindowLayout.WindowState != types.WsNormal {
		m.SetWindowState(cfg.WindowLayout.WindowState)
	}

	m.initAllImageList()
	// 设置窗口图标
	m.setWindowIcon()
	// 窗口大小改变事件
	m.SetOnResize(m.WindowOnResize)
	m.SetOnWindowStateChange(m.WindowOnWindowStateChange)
	// 创建设计器布局
	m.initDesignerLayoutV2()
}

func (m *TAppWindow) initAllImageList() {
	imageActions = tool.NewImageList(m, "actions", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}, Image150: types.TSize{Cx: 24, Cy: 24}, Image200: types.TSize{Cx: 32, Cy: 32}})
	imageComponents = tool.NewImageList(m, "components", tool.ImageRect{Image50: types.TSize{Cx: 16, Cy: 16}, Image100: types.TSize{Cx: 24, Cy: 24}, Image150: types.TSize{Cx: 36, Cy: 36}, Image200: types.TSize{Cx: 48, Cy: 48}})
	imageItem = tool.NewImageList(m, "item", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}, Image150: types.TSize{Cx: 24, Cy: 24}, Image200: types.TSize{Cx: 32, Cy: 32}})
	imageMenu = tool.NewImageList(m, "menu", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}, Image150: types.TSize{Cx: 24, Cy: 24}, Image200: types.TSize{Cx: 32, Cy: 32}})
	imageTabComp = tool.NewImageList(m, "tab-comp", tool.ImageRect{Image100: types.TSize{Cx: 16, Cy: 16}})
	tool.AppendSVGToImageList(imageComponents, "icons/svg_files", tool.ImageRect{Image50: types.TSize{Cx: 16, Cy: 16}, Image100: types.TSize{Cx: 24, Cy: 24}, Image150: types.TSize{Cx: 36, Cy: 36}, Image200: types.TSize{Cx: 48, Cy: 48}})
	tool.AppendPNGToImageList(imageComponents, "icons/png_files", tool.ImageRect{Image50: types.TSize{Cx: 16, Cy: 16}, Image100: types.TSize{Cx: 24, Cy: 24}, Image150: types.TSize{Cx: 36, Cy: 36}, Image200: types.TSize{Cx: 48, Cy: 48}})
}

func (m *TAppWindow) WindowOnResize(sender lcl.IObject) {
	//println("WindowOnResize")
	DebounceContentLayoutBoxReAlign()
}

func (m *TAppWindow) WindowOnWindowStateChange(sender lcl.IObject) {
	//println("WindowOnWindowStateChange")
	ImmediatelyContentLayoutBoxReAlign()
}

func (m *TAppWindow) OnShow(sender lcl.IObject) {
	logs.Info("OnShow")
	gOnShow.Do(func() {
		// 默认禁用组件功能
		SetEnableFuncComponent(false)
		// 窗口显示在鼠标所在的窗口
		//m.showInMonitor()
		for _, fn := range windowShowEvents {
			fn()
		}
		// 向消息输出基本信息
		cfg := config.DesignerConfig
		_, _, _, _, _, v := api.LCLVersion()
		consoleText := tool.Buffer{}
		consoleText.WriteString(cfg.Title, ":", cfg.Version, " LCL:v", v)
		WriteConsole(consoleText.String())

		// 初始化依赖模块信息 ast
		dependmod.InitDependencyModule(m.ReadBounds(), func(ok bool) {
			if ok { // 一个开关, 动态配置
				autoAssociateProjectLoad()
			}
		})
	})
}

func (m *TAppWindow) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	logs.Info("OnCloseQuery closing:", m.closing)
	*canClose = m.closing
	if !m.closing {
		go m.handleClose()
	}
}

// 处理设计器窗口关闭
func (m *TAppWindow) handleClose() {
	logs.Info("closeHandle")
	m.closing = true

	// 收集窗口布局信息
	// 收集窗口信息
	br := m.BoundsRect()
	// 收集菜单-视图
	viewWidgetsChecked := m.mainMenu.viewWidgets.Checked()
	viewProjectChecked := m.mainMenu.viewProject.Checked()
	viewInspectorChecked := m.mainMenu.viewInspector.Checked()
	viewConsoleChecked := m.mainMenu.viewConsole.Checked()
	viewStatusbarChecked := m.mainMenu.viewStatusbar.Checked()
	// 收集布局 width height
	widgetPanelWidth := m.contentLayout.widgetPanel.Width()
	projectPanelWidth := m.contentLayout.projectPanel.Width()
	inspectorPanelWidth := m.contentLayout.inspectorPanel.Width()
	consoleLogPanelHeight := m.contentLayout.consoleLogPanel.Height()

	// 持久化设计器窗口布局关键数据
	windowLayout := &config.Config.WindowLayout
	if windowLayout.MenuView == nil {
		windowLayout.MenuView = &config.StorageMenuView{}
	}
	if windowLayout.ContentLayout == nil {
		windowLayout.ContentLayout = &config.StorageContentLayout{}
	}
	if m.WindowState() == types.WsNormal {
		windowLayout.WindowBoundsRect = br
	}
	windowLayout.WindowState = m.WindowState()
	windowLayout.ContentLayout.WidgetPanelWidth = widgetPanelWidth
	windowLayout.ContentLayout.ProjectPanelWidth = projectPanelWidth
	windowLayout.ContentLayout.InspectorPanelWidth = inspectorPanelWidth
	windowLayout.ContentLayout.ConsoleLogPanelHeight = consoleLogPanelHeight
	windowLayout.MenuView.WidgetsChecked = viewWidgetsChecked
	windowLayout.MenuView.ProjectChecked = viewProjectChecked
	windowLayout.MenuView.InspectorChecked = viewInspectorChecked
	windowLayout.MenuView.ConsoleChecked = viewConsoleChecked
	windowLayout.MenuView.StatusbarChecked = viewStatusbarChecked
	windowLayout.ContentLayout.InspectorLayout.PropertyTreeWidth = gDefaultPropertyNameColumnTreeWidth
	windowLayout.ContentLayout.InspectorLayout.EventTreeWidth = gDefaultEventNameColumnTreeWidth

	// 更新最后打开的项目
	config.UpdateLastProject(gAppEGPPath)
	// 更新配置文件
	config.UpdateConfig()
	// 取消所有生成事件
	event.CancelAll()
	// 延迟关闭
	time.AfterFunc(time.Second/10, func() {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			// 最后在UI线程关闭
			m.Close()
		})
	})
}

func (m *TAppWindow) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	logs.Info("OnClose closing:", m.closing)
}

func AddOnShow(fn func()) {
	windowShowEvents = append(windowShowEvents, fn)
}
