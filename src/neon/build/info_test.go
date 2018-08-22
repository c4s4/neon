package build

import (
	"testing"
)

func TestInfoThemes(t *testing.T) {
	themes := InfoThemes()
	if themes != "bee blue bold cyan fire green magenta marine nature red reverse rgb yellow" {
		t.Errorf("Bad themes")
	}
}

func TestInfoBuiltins(t *testing.T) {
	builtins := InfoBuiltins()
	if builtins != "test" {
		t.Errorf("Bad builtins")
	}
}
