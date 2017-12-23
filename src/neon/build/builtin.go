package build

import (
	"github.com/c4s4/anko/vm"
)

// Descriptor for a builtin function
type BuiltinDescriptor struct {
	Function interface{}
	Help     string
}

// Map of builtin descriptors by name
var BuiltinMap = make(map[string]BuiltinDescriptor)

// Load defined builtins
func LoadBuiltins(vm *vm.Env) {
	for name, descriptor := range BuiltinMap {
		vm.Define(name, descriptor.Function)
	}
}
