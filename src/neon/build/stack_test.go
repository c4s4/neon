package build

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()
	stack.Push("foo")
	if !stack.Contains("foo") {
		t.Errorf("Error contains")
	}
	if stack.Contains("bar") {
		t.Errorf("Error contains")
	}
	stack.Push("bar")
	if !stack.Contains("bar") {
		t.Errorf("Error contains")
	}
	if stack.String() != "foo -> bar" {
		t.Errorf("Error ToString: %v", stack.String())
	}
	if stack.Last() != "bar" {
		t.Errorf("Error Last: %v", stack.Last())
	}
}
