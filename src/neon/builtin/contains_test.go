package builtin

import (
	"testing"
)

func TestContains(t *testing.T) {
	if !contains([]string{"foo", "bar"}, "bar") {
		t.Errorf("Error builtin contains")
	}
	if contains([]string{"foo", "bar"}, "spam") {
		t.Errorf("Error builtin contains")
	}
}
