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
	if !contains([]interface{}{"foo", "bar"}, "bar") {
		t.Errorf("Error builtin contains")
	}
	if contains([]interface{}{"foo", "bar"}, "spam") {
		t.Errorf("Error builtin contains")
	}
	var content interface{}
	content = []string{"foo", "bar"}
	if !contains(content, "bar") {
		t.Errorf("Error builtin contains")
	}
	if contains(content, "spam") {
		t.Errorf("Error builtin contains")
	}
}
