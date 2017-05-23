package builtin

import (
	"testing"
)

func TestFilter(t *testing.T) {
	filtered := Filter([]string{"foo", "bar"}, "foo")
	if len(filtered) != 1 {
		t.Errorf("Error builtin filter")
	}
	if filtered[0] != "bar" {
		t.Errorf("Error builtin filter")
	}
}
