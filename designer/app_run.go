package designer

import "github.com/energye/lcl/lcl"

func Run() {
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&mainWindow)
	lcl.Application.Run()
}
