package clui

import (
	"os"
	"runtime"
	"strings"
)

func shouldUseSimpleUI() bool {
	if disabledByEnv() || isCI() || !isInteractiveTerminal() {
		return true
	}
	return runtime.GOOS == "windows" && isLegacyWindows()
}

func disabledByEnv() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("ENERGY_TUI")))
	return value == "0" || value == "false" || value == "off" || value == "no"
}

func isCI() bool {
	return os.Getenv("CI") != ""
}

func isInteractiveTerminal() bool {
	if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		return false
	}
	if stat, err := os.Stdout.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		return false
	}
	if strings.ToLower(os.Getenv("TERM")) == "dumb" {
		return false
	}
	return true
}
