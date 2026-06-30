package clui

import "golang.org/x/sys/windows"

func isLegacyWindows() bool {
	version := windows.RtlGetVersion()
	return version.MajorVersion < 10
}
