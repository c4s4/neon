package builtin

import (
	"testing"
)

func TestEscapeUrl(t *testing.T) {
	if escapeUrl("/foo bar") != "/foo%20bar" {
		t.Errorf("Error builtin escapeulr")
	}
}
