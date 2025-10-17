package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

var (
	mainWindow       TAppWindow
	toolbarHeight    int32 = 66
	bgDrakColor            = colors.RGBToColor(56, 57, 60)
	bgLightColor           = colors.ClWhite
	windowShowEvents []func()
	imageActions     *tool.ImageList
	imageComponents  *tool.ImageList
	imageItem        *tool.ImageList
	imageMenu        *tool.ImageList
)

// 设计器应用窗口
type TAppWindow struct {
	lcl.TEngForm
	mainMenu            lcl.IMainMenu              // 主菜单
	componentProperties lcl.IApplicationProperties //
	box                 *BottomBox                 // 底部布局盒子
	bar                 *StatusBar
}

func (m *TAppWindow) FormCreate(sender lcl.IObject) {
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
	imageActions = tool.NewImageList(m, "actions")
	imageComponents = tool.NewImageList(m, "components")
	imageItem = tool.NewImageList(m, "item")
	imageMenu = tool.NewImageList(m, "menu")
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
}

func AddOnShow(fn func()) {
	windowShowEvents = append(windowShowEvents, fn)
}
