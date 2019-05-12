package builtin

import (
	"testing"
)

func TestLowercase(t *testing.T) {
	if lowercase("FOOBar") != "foobar" {
		t.Errorf("Error builtin lowercase")
	}
}
