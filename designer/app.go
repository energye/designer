package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
)

type TApp struct {
	lcl.TEngForm
}

var app TApp

func Run() {
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&app)
	lcl.Application.Run()
}

func (m *TApp) FormCreate(sender lcl.IObject) {
	log.Println("FormCreate")
	cfg := config.Config
	m.SetWidth(int32(cfg.Window.Width))
	m.SetHeight(int32(cfg.Window.Height))
	m.WorkAreaCenter()
}

func (m *TApp) FormAfterCreate(sender lcl.IObject) {
	log.Println("FormAfterCreate")
}

func (m *TApp) CreateParams(params *types.TCreateParams) {
	log.Println("CreateParams")
}

func (m *TApp) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	log.Println("OnCloseQuery")
}

func (m *TApp) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	log.Println("OnClose")
}
