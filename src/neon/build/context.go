package build

import (
	"fmt"
)

type Context struct {
	VM    *VM
	Index *Index
	Stack *Stack
}

func NewContext(build *Build) (*Context, error) {
	vm, err := NewVM(build)
	if err != nil {
		return nil, fmt.Errorf("evaluating context: %v", err)
	}
	context := Context{
		VM:    vm,
		Index: NewIndex(),
		Stack: NewStack(),
	}
	return &context, nil
}

func (context *Context) Message(text string, args ...interface{}) {
	printGrey(text, args...)
}
