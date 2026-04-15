package main

import (
	"github.com/energye/designer/options"
	"github.com/energye/designer/options/bean"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
)

func main() {
	bean.GProject = &bean.TProject{}
	api.SetDebug(true)
	form := options.TConfigProjectForm{}
	lcl.Init(nil, nil)
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&form)
	form.WorkAreaCenter()
	lcl.Application.Run()
}
