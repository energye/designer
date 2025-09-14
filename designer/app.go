package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
)

// 设计器应用窗口
type TAppWindow struct {
	lcl.TEngForm
}

var (
	mainWindow TAppWindow
)

func Run() {
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&mainWindow)
	lcl.Application.Run()
}

func (m *TAppWindow) FormCreate(sender lcl.IObject) {
	log.Println("FormCreate")
	cfg := config.Config
	m.SetDoubleBuffered(true)
	m.SetWidth(int32(cfg.Window.Width))
	m.SetHeight(int32(cfg.Window.Height))
	m.WorkAreaCenter()
	m.SetWindowIcon()
	m.SetOnShow(m.OnShow)
}

func (m *TAppWindow) OnShow(sender lcl.IObject) {
	log.Println("OnShow")
	m.ShowInMonitor()
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
