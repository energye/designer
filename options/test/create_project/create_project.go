package main

import (
	"github.com/energye/designer/options"
	"github.com/energye/lcl/lcl"
	"os"
)

func main() {
	form := options.TCreateProjectForm{}
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
