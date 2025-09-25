package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
)

func Run() {
	logs.Info("ENERGY Designer RUN")
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.SetScaled(true)
	lcl.Application.NewForms(&mainWindow)
	lcl.Application.Run()
	logs.Info("ENERGY Designer RUN END.")
}
