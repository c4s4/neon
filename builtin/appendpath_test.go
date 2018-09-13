package builtin

import (
	"testing"
)

func TestAppendpath(t *testing.T) {
	p := appendPath("foo", []string{"spam", "eggs"})
	if p[0] != "foo/spam" || p[1] != "foo/eggs" {
		t.Errorf("Error builtin appendpath")
	}
}
