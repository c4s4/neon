package builtin

import (
	"errors"

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

- Nothing, but the error message is stored in '_error' property if the error
  is caught with try task.

Examples:

    # stop the script with an error message
    throw("Some tests failed")
    # returns: nothing, the script is interrupted on error`,
	})
}

func throw(message string) error {
	return errors.New(message)
}
