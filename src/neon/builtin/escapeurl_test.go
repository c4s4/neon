package builtin

import (
	"testing"
)

func TestEscapeUrl(t *testing.T) {
	if EscapeUrl("/foo bar") != "/foo%20bar" {
		t.Errorf("Error builtin escapeulr")
	}
}
