package builtin

import (
	"testing"
	"strings"
)

func TestAbsolute(t *testing.T) {
	if !strings.HasSuffix(Absolute("foo/../bar/spam.txt"), "bar/spam.txt") {
		t.Errorf("TestAbsolute failed")
	}
}
