package build

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()
	err := stack.Push(&Target{Name: "foo"})
	if err != nil || !stack.Contains("foo") {
		t.Errorf("Error contains")
	}
	if stack.Contains("bar") {
		t.Errorf("Error contains")
	}
	err = stack.Push(&Target{Name: "bar"})
	if err != nil || !stack.Contains("bar") {
		t.Errorf("Error contains")
	}
	if stack.String() != "foo -> bar" {
		t.Errorf("Error ToString: %v", stack.String())
	}
	if stack.Last().Name != "bar" {
		t.Errorf("Error Last: %v", stack.Last())
	}
}
