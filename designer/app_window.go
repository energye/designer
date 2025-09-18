package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"log"
)

var (
	mainWindow       TAppWindow
	toolbarHeight    int32 = 66
	bgColor                = colors.RGBToColor(56, 57, 60)
	windowShowEvents []func()
)

// 设计器应用窗口
type TAppWindow struct {
	lcl.TEngForm
	mainMenu            lcl.IMainMenu              // 主菜单
	componentProperties lcl.IApplicationProperties //
	toolbar             *TopToolbar                // 顶部工具栏
	box                 *BottomBox                 // 底部布局盒子
	bar                 *StatusBar
}

func (m *TAppWindow) FormCreate(sender lcl.IObject) {
	log.Println("FormCreate")
	cfg := config.Config
	// 属性
	m.SetCaption(cfg.Title + " " + cfg.Version)
	m.SetDoubleBuffered(true)
	m.SetWidth(int32(cfg.Window.Width))
	m.SetHeight(int32(cfg.Window.Height))
	constra := m.Constraints()
	constra.SetMinWidth(400)
	constra.SetMinHeight(200)
	// 设置窗口图标
	m.setWindowIcon()
	// 窗口显示事件
	m.SetOnShow(m.OnShow)
	// 创建设计器布局
	m.createDesignerLayout()
	// status bar
	m.createStatusBar()
}

func (m *TAppWindow) OnShow(sender lcl.IObject) {
	log.Println("OnShow")
	// 窗口显示在鼠标所在的窗口
	m.showInMonitor()
	for _, fn := range windowShowEvents {
		fn()
	}
}

func (m *TAppWindow) FormAfterCreate(sender lcl.IObject) {
	log.Println("FormAfterCreate")
}

func (m *TAppWindow) CreateParams(params *types.TCreateParams) {
	log.Println("CreateParams")
}

func (m *TAppWindow) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	log.Println("OnCloseQuery")
}

func (m *TAppWindow) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	log.Println("OnClose")
}

func AddOnShow(fn func()) {
	windowShowEvents = append(windowShowEvents, fn)
}
