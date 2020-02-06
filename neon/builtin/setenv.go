package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"os"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "setenv",
		Func: setenv,
		Help: `Set environment variable.

Arguments:

- The variable name.
- The variable value.

Examples:

    # set foo to value bar
    setenv("foo", "bar")`})
}

func setenv(name, value string) {
	os.Setenv(name, value)
}
