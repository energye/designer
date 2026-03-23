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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mainWindow       TAppWindow
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
	componentProperties   lcl.IApplicationProperties //
	box                   *BottomBox                 // 底部布局盒子
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
	if title == "" {
		title = fmt.Sprintf("%v %v", config.FormConfig.Title, config.FormConfig.Version)
	} else {
		title = fmt.Sprintf("%v %v - %v", title, config.FormConfig.Title, config.FormConfig.Version)
	}
	lcl.RunOnMainThreadSync(func() {
		logs.Println("UpdateDesignerTitle:", title)
		mainWindow.SetCaption(title)
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
	cfg := config.FormConfig
	// 属性
	m.SetCaption(cfg.Title + " " + cfg.Version)
	m.SetDoubleBuffered(true)
	m.SetLeft(cfg.Window.X)
	m.SetTop(cfg.Window.Y)
	m.SetWidth(cfg.Window.Width)
	m.SetHeight(cfg.Window.Height)
	m.SetColor(bgLightColor)
	constra := m.Constraints()
	constra.SetMinWidth(400)
	constra.SetMinHeight(200)
	if cfg.Window.X == 0 && cfg.Window.Y == 0 {
		// 窗口显示在鼠标所在的窗口
		m.showInMonitor()
	}
	if cfg.Window.WindowState != 0 {
		m.SetWindowState(cfg.Window.WindowState)
	}

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
	gOnShow.Do(func() {
		// 默认禁用组件功能
		SetEnableFuncComponent(false)
		// 窗口显示在鼠标所在的窗口
		//m.showInMonitor()
		for _, fn := range windowShowEvents {
			fn()
		}
		// 向消息输出基本信息
		cfg := config.FormConfig
		_, _, _, _, _, v := api.LCLVersion()
		consoleText := tool.Buffer{}
		consoleText.WriteString(cfg.Title, ":", cfg.Version, " LCL:v", v)
		WriteConsole(consoleText.String())

		// 检查框架 frameworks 框架是否设置安装目录，并正确安装
		m.checkInstallFrameworks()
		if true { // 一个开关, 动态配置
			isEgp := strings.HasSuffix(os.Args[len(os.Args)-1], consts.EGPExt)
			if isEgp {
				// 自动打开 energy 项目
				filePath := os.Args[len(os.Args)-1]
				event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: filePath}})
			} else if config.Config.LastProject != "" && tool.IsExist(config.Config.LastProject) {
				// 自动打开 最后一次打开的项目
				event.Emit(event.TTrigger{Name: event.Project, Payload: event.TPayload{Type: event.ProjectLoad, Data: config.Config.LastProject}})
			}
		}
	})
}

// 检查框架是否设置 frameworks 框架安装目录，并正确安装
func (m *TAppWindow) checkInstallFrameworks() {
	frameworksDir := config.Config.FrameworkDir
	if config.Config.FrameworkDir == "" || !tool.IsExist(frameworksDir) {
		event.ConsoleWriteError("未设置 framework 框架安装目录")
		lcl.RunOnMainThreadAsync(func(id uint32) {
			// 禁用所有功能组件
			SetEnableFuncComponent(false)
			// 弹出一个引导窗口
			form := NewInstallFrameworkForm()
			form.ShowModal()
		})
	}
}

//func (m *TAppWindow) FormAfterCreate(sender lcl.IObject) {
//	logs.Info("FormAfterCreate")
//}

//func (m *TAppWindow) CreateParams(params *types.TCreateParams) {
//	logs.Info("CreateParams")
//}

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
	br := m.BoundsRect()
	// 更新设计器窗口
	config.UpdateWindow(br.Left, br.Top, br.Width(), br.Height(), m.WindowState())
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
