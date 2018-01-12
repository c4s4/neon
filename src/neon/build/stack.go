package build

import (
	"strings"
	"fmt"
)

// Structure for a stack. A stack lists all targets that run during a build.
type Stack struct {
	Targets []string
}

// NewStack makes a new stack
// Returns: a pointer to the stack
func NewStack() *Stack {
	stack := Stack{
		Targets: make([]string, 0),
	}
	return &stack
}

// Contains tells if the stack contains given target
// - target: the name of the target
// Returns: a boolean telling if target is in the stack
func (stack *Stack) Contains(target string) bool {
	for _, name := range stack.Targets {
		if target == name {
			return true
		}
	}
	return false
}

// Push a target on the stack
// - target: the name of the target to push on the stack
// Return: an error if we are in an infinite loop
func (stack *Stack) Push(target string) error {
	for _, t := range stack.Targets {
		if t == target {
			stack.Targets = append(stack.Targets, target)
			loop := strings.Join(stack.Targets, " -> ")
			return fmt.Errorf("infinite loop: %v", loop)
		}
	}
	stack.Targets = append(stack.Targets, target)
	return nil
}

// Last gets the last target on stack
// Return: name of the last target on stack
func (stack *Stack) Last() string {
	if len(stack.Targets) == 0 {
		return ""
	} else {
		return stack.Targets[len(stack.Targets)-1]
	}
}

// ToString returns string representation of the stack, such as:
// "foo -> bar -> spam"
// Return: the stack as a string
func (stack *Stack) String() string {
	return strings.Join(stack.Targets, " -> ")
}

// Copy returns a copy of the stack
// Return: pointer to a copy of the stack
func (stack *Stack) Copy() *Stack {
	another := make([]string, len(stack.Targets))
	for i := 0; i < len(stack.Targets); i++ {
		another[i] = stack.Targets[i]
	}
	return &Stack{another}
}
