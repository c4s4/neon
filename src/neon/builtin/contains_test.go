package builtin

import (
	"testing"
)

func TestContains(t *testing.T) {
	if !Contains([]string{"foo", "bar"}, "bar") {
		t.Errorf("Error builtin contains")
	}
	if Contains([]string{"foo", "bar"}, "spam") {
		t.Errorf("Error builtin contains")
	}
}
