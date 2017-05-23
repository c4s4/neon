package builtin

import (
	"testing"
)

func TestDirectory(t *testing.T) {
	if Directory("/foo/bar/spam.txt") != "/foo/bar" {
		t.Errorf("Error builtin directory")
	}
}
