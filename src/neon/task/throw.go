package task

import (
	"fmt"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "throw",
		Func: Throw,
		Args: reflect.TypeOf(ThrowArgs{}),
		Help: `Throws an error.

Arguments:

- throw: the message of the error (string).

Examples:

    # stop the build because tests failed
    - throw: "ERROR: tests failed"

Notes:

- You can catch raised errors with try/catch/finally task.
- Property _error is set with the error message.
- If the error was not catch, the error message will be printed on the console
  as the cause of the build failure.`,
	})
}

type ThrowArgs struct {
	Throw string
}

func Throw(context *build.Context, args interface{}) error {
	params := args.(ThrowArgs)
	return fmt.Errorf(params.Throw)
}
