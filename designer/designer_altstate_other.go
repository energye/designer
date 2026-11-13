//go:build !windows

package designer

func platformAltDown() bool {
	return false
}
