package builtin

import (
	"github.com/mattn/anko/vm"
)

type BuiltinDescriptor struct {
	Function interface{}
	Help     string
}

var Builtins map[string]BuiltinDescriptor = make(map[string]BuiltinDescriptor)

func AddBuiltins(vm *vm.Env) {
	for name, descriptor := range Builtins {
		vm.Define(name, descriptor.Function)
	}
}
