// +build !windows

package util

import (
	"testing"
)

func TestTerminalWidth(t *testing.T) {
	width := TerminalWidth()
	if width < 0 || width > 1000 {
		t.Errorf("Bad terminal width: %d", width)
	}
}
