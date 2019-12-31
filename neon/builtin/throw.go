package builtin

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "throw",
		Func: throw,
		Help: `Throw an error that will cause script failure.

Arguments:

- The error message of the failure.

Returns:

- Nothing, but sets the variable 'error' with the error message.

Examples:

    # stop the script with an error message
    throw("Some tests failed")
    # returns: nothing, the script is interrupted on error`,
	})
}

func throw(message string) error {
	return fmt.Errorf(message)
}
