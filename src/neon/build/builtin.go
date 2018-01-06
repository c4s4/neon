package build

import (
	"github.com/c4s4/anko/vm"
)

// Descriptor for a builtin function
type BuiltinDesc struct {
	Name string
	Func interface{}
	Help string
}

// Map of builtin descriptors by name
var BuiltinMap = make(map[string]BuiltinDesc)

// AddBuitin adds given builtin to the map
// - desc: builtin description
func AddBuiltin(desc BuiltinDesc) {
	if _, ok := BuiltinMap[desc.Name]; ok {
		panic("Builtin function '" + desc.Name + "' already defined")
	}
	BuiltinMap[desc.Name] = desc
}

// LoadBuiltins loads defined builtins in the VM
// - vm: the VM to load builtins into
func LoadBuiltins(vm *vm.Env) {
	for name, descriptor := range BuiltinMap {
		vm.Define(name, descriptor.Func)
	}
}
