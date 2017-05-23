package builtin

import (
	"testing"
)

func TestLowercase(t *testing.T) {
	if Lowercase("FOOBar") != "foobar" {
		t.Errorf("Error builtin lowercase")
	}
}
