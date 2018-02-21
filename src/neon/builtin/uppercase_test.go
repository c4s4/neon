package builtin

import (
	"testing"
)

func TestUppercase(t *testing.T) {
	if uppercase("fooBar") != "FOOBAR" {
		t.Errorf("Error builtin uppercase")
	}
}
