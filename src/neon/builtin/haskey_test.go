package builtin

import (
	"testing"
)

func TestHaskey(t *testing.T) {
	m := map[interface{}]interface{}{"key": 1}
	if !haskey(m, "key") {
		t.Errorf("Error builtin haskey")
	}
	if haskey(m, "not") {
		t.Errorf("Error builtin haskey")
	}
}
