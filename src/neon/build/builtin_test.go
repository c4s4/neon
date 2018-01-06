package build

import (
	"testing"
)

func TestAddBuiltin(t *testing.T) {
	BuiltinMap = make(map[string]BuiltinDesc)
	AddBuiltin(BuiltinDesc{
		Name: "test",
		Func: TestAddBuiltin,
		Help: "Test builtin",
	})
	Assert(len(BuiltinMap), 1, t)
	Assert(BuiltinMap["test"].Name, "test", t)
	Assert(BuiltinMap["test"].Help, "Test builtin", t)
}
