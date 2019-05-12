package builtin

import (
	"testing"
)

func TestJoinpath(t *testing.T) {
	if joinPath("foo", "bar") != "foo/bar" {
		t.Errorf("Error builtin joinpath")
	}
}
