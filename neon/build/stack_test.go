package build

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()
	foo := &Target{Name: "foo"}
	bar := &Target{Name: "bar"}
	err := stack.Push(foo)
	if err != nil || !stack.Contains("foo") {
		t.Errorf("Error contains")
	}
	if stack.Contains("bar") {
		t.Errorf("Error contains")
	}
	err = stack.Push(bar)
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
}

func TestStackPop(t *testing.T) {
	stack := NewStack()
	foo := &Target{Name: "foo"}
	bar := &Target{Name: "bar"}
	if err := stack.Push(foo); err != nil {
		t.Fatalf("pushing: %v", err)
	}
	if err := stack.Push(bar); err != nil {
		t.Fatalf("pushing: %v", err)
	}
	err := stack.Pop()
	if err != nil {
		t.Errorf("Error poping")
	}
	if stack.Last().Name != "foo" {
		t.Errorf("Error Last: %v", stack.Last())
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
