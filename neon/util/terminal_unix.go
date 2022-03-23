//go:build !windows
// +build !windows

package util

import (
	"syscall"
	"unsafe"
)

const (
	// DefaultTerminalWidth is the terminal width when undefined
	DefaultTerminalWidth = 80
)

// TerminalWidth returns the terminal width
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
		return DefaultTerminalWidth
	}
	return int(ws.Col)
}
