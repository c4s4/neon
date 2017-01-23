package build

import (
	"github.com/mattn/anko/vm"
)

type BuiltinDescriptor struct {
	Function interface{}
	Help     string
}

var BuiltinMap map[string]BuiltinDescriptor = make(map[string]BuiltinDescriptor)

func LoadBuiltins(vm *vm.Env) {
	for name, descriptor := range BuiltinMap {
		vm.Define(name, descriptor.Function)
	}
}
