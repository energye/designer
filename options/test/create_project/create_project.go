package main

import (
	"github.com/energye/designer/options"
	"github.com/energye/lcl/lcl"
)

func main() {
	form := options.TCreateProjectForm{}
	lcl.Init(nil, nil)
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&form)
	form.WorkAreaCenter()
	lcl.Application.Run()
}
