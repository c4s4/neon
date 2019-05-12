// +build windows

package util

const (
	DEFAULT_TERMINAL_WIDTH = 80
)

// Get terminal width
func TerminalWidth() int {
	return DEFAULT_TERMINAL_WIDTH
}
