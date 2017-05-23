package builtin

import (
	"testing"
)

func TestSplit(t *testing.T) {
	splitted := Split("foo bar", " ")
	if len(splitted) != 2 {
		t.Errorf("Error builtin split")
	}
	if splitted[0] != "foo" {
		t.Errorf("Error builtin split")
	}
	if splitted[1] != "bar" {
		t.Errorf("Error builtin split")
	}
}
