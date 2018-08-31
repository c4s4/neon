package build

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()
	foo := &Target{Name: "foo"}
	err := stack.Push(foo)
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
	err = stack.Push(foo)
	if err == nil || err.Error() != "infinite loop: foo -> bar -> foo" {
		t.Errorf("Should have raised infinite loop error")
	}
	if stack.String() != "foo -> bar -> foo" {
		t.Errorf("Error ToString: %v", stack.String())
	}
	copy := stack.Copy()
	if copy.String() != "foo -> bar -> foo" {
		t.Errorf("Error ToString: %v", stack.String())
	}
	if stack.Last().Name != "foo" {
		t.Errorf("Error Last: %v", stack.Last())
	}
	err = stack.Pop()
	if err != nil {
		t.Errorf("Error poping")
	}
	if stack.Last().Name != "bar" {
		t.Errorf("Error Last: %v", stack.Last())
	}
	err = stack.Pop()
	if err != nil {
		t.Errorf("Error poping")
	}
	err = stack.Pop()
	if err != nil {
		t.Errorf("Error poping")
	}
	err = stack.Pop()
	if err == nil {
		t.Errorf("Error poping")
	}
	last := stack.Last()
	if last != nil {
		t.Errorf("Error on last")
	}
}
