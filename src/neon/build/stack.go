package build

import (
	"strings"
)

// Structure for a stack. A stack lists all targets that run during a build.
type Stack struct {
	Targets []string
}

// Make a stack
func NewStack() *Stack {
	stack := Stack{
		Targets: make([]string, 0),
	}
	return &stack
}

// Tells if the stack contains given target
func (stack *Stack) Contains(target string) bool {
	for _, name := range stack.Targets {
		if target == name {
			return true
		}
	}
	return false
}

// Push a target on the stack
func (stack *Stack) Push(target string) {
	stack.Targets = append(stack.Targets, target)
}

// Return string representation of the stack ("foo -> bar -> spam" for instance)
func (stack *Stack) ToString() string {
	return strings.Join(stack.Targets, " -> ")
}
