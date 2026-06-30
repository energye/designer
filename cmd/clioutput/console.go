package clioutput

import (
	"strings"

	"github.com/energye/designer/cmd/clui"
	"github.com/energye/designer/event"
)

func Bind(ui clui.UI) {
	if ui == nil {
		return
	}
	event.SetConsoleWriter(func(level event.ConsoleLevel, message string) {
		message = strings.TrimRight(message, "\r\n")
		if message == "" {
			return
		}
		switch level {
		case event.ConsoleLevelWarn:
			ui.Warn(message)
		case event.ConsoleLevelError:
			ui.Error(message)
		case event.ConsoleLevelDebug:
			ui.Info("[DEBUG]", message)
		default:
			ui.Info(message)
		}
	})
}
