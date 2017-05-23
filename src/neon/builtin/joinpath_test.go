package builtin

import (
	"testing"
)

func TestJoinpath(t *testing.T) {
	if Joinpath("foo", "bar") != "foo/bar" {
		t.Errorf("Error builtin joinpath")
	}
}
