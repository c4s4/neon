package build

import (
	"github.com/mattn/anko/env"
)

// BuiltinDesc is a descriptor for a builtin function
type BuiltinDesc struct {
	Name string
	Func interface{}
	Help string
}

// BuiltinMap is a map of builtin descriptors by name
var BuiltinMap = make(map[string]BuiltinDesc)

// AddBuiltin adds given builtin to the map
// - desc: builtin description
func AddBuiltin(desc BuiltinDesc) {
	if _, ok := BuiltinMap[desc.Name]; ok {
		panic("Builtin function '" + desc.Name + "' already defined")
	}
	BuiltinMap[desc.Name] = desc
}

// LoadBuiltins loads defined builtins in the VM
// - vm: the VM to load builtins into
func LoadBuiltins(vm *env.Env) {
	for name, descriptor := range BuiltinMap {
		if err := vm.Define(name, descriptor.Func); err != nil {
			panic("loading builtin " + name)
		}
	}
}
