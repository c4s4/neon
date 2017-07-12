// +build !windows

package util

import (
	"syscall"
	"unsafe"
)

const (
	DEFAULT_TERMINAL_WIDTH = 80
)

// Get terminal width
func TerminalWidth() int {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 {
		return DEFAULT_TERMINAL_WIDTH
	} else {
		return int(ws.Col)
	}
}
