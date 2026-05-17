package main

import (
	"github.com/energye/designer/options"
	"github.com/energye/designer/options/bean"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"os"
)

func main() {
	bean.GProject = &bean.TProject{}
	api.SetDebug(true)
	form := options.TEnvForm{}
	lcl.Init()
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&form)
	form.WorkAreaCenter()
	lcl.Application.Run()
}

func init() {
	os.Setenv("--ws", "gtk3") // for linux
}
