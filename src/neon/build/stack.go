package build

import (
	"fmt"
	"strings"
)

// Stack is structure for a stack. A stack lists all targets that run during a build.
type Stack struct {
	Targets []*Target
}

// NewStack makes a new stack
// Returns: a pointer to the stack
func NewStack() *Stack {
	stack := Stack{
		Targets: make([]*Target, 0),
	}
	return &stack
}

// Contains tells if the stack contains given target
// - target: target to test
// Returns: a boolean telling if target is in the stack
func (stack *Stack) Contains(name string) bool {
	for _, target := range stack.Targets {
		if name == target.Name {
			return true
		}
	}
	return false
}

// Push a target on the stack
// - target: target to push on the stack
// Return: an error if we are in an infinite loop
func (stack *Stack) Push(target *Target) error {
	for _, t := range stack.Targets {
		if t == target {
			stack.Targets = append(stack.Targets, target)
			return fmt.Errorf("infinite loop: %v", stack.String())
		}
	}
	stack.Targets = append(stack.Targets, target)
	return nil
}

// Pop target on the stack
// Return: error if something went wrong
func (stack *Stack) Pop() error {
	length := len(stack.Targets)
	if length == 0 {
		return fmt.Errorf("no target on stack")
	}
	stack.Targets = stack.Targets[:length-1]
	return nil
}

// Last gets the last target on stack
// Return: last target on stack
func (stack *Stack) Last() *Target {
	if len(stack.Targets) == 0 {
		return nil
	}
	return stack.Targets[len(stack.Targets)-1]
}

// ToString returns string representation of the stack, such as:
// "foo -> bar -> spam"
// Return: the stack as a string
func (stack *Stack) String() string {
	names := make([]string, len(stack.Targets))
	for i, target := range stack.Targets {
		names[i] = target.Name
	}
	return strings.Join(names, " -> ")
}

// Copy returns a copy of the stack
// Return: pointer to a copy of the stack
func (stack *Stack) Copy() *Stack {
	another := make([]*Target, len(stack.Targets))
	for i := 0; i < len(stack.Targets); i++ {
		another[i] = stack.Targets[i]
	}
	return &Stack{another}
}
