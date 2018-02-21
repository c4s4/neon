package builtin

import (
	"testing"
)

func TestFilename(t *testing.T) {
	if filename("/foo/bar/spam.txt") != "spam.txt" {
		t.Errorf("Error builtin filename")
	}
}
