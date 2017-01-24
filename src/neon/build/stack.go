package build

import (
	"strings"
)

type Stack struct {
	Targets []string
}

func NewStack() *Stack {
	stack := Stack{
		Targets: make([]string, 0),
	}
	return &stack
}

func (stack *Stack) Contains(target string) bool {
	for _, name := range stack.Targets {
		if target == name {
			return true
		}
	}
	return false
}

func (stack *Stack) Push(target string) {
	stack.Targets = append(stack.Targets, target)
}

func (stack *Stack) ToString() string {
	return strings.Join(stack.Targets, " -> ")
}
