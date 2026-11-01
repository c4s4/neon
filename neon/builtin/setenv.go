package builtin

import (
	"os"

	"github.com/c4s4/neon/neon/build"
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

func setenv(name, value string) error {
	return os.Setenv(name, value)
}
