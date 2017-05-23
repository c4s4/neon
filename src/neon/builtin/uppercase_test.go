package builtin

import (
	"testing"
)

func TestUppercase(t *testing.T) {
	if Uppercase("fooBar") != "FOOBAR" {
		t.Errorf("Error builtin uppercase")
	}
}
