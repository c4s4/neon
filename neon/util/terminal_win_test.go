//go:build windows
// +build windows

package util

import (
	"testing"
)

func TestTerminalWidth(t *testing.T) {
	width := TerminalWidth()
	if width != 80 {
		t.Errorf("Bad terminal width: %d", width)
	}
}
