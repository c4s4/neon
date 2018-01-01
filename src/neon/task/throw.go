package task

import (
	"fmt"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "throw",
		Func: Throw,
		Args: reflect.TypeOf(ThrowArgs{}),
		Help: `Throws an error.

Arguments:

- throw: the message of the error.

Examples:

    # stop the build because tests don't run
    - throw: "ERROR: tests don't run"

Notes:

- The error message will be printed on the console as the source of the build
  failure.`,
	})
}

type ThrowArgs struct {
	Throw string
}

func Throw(context *build.Context, args interface{}) error {
	params := args.(ThrowArgs)
	return fmt.Errorf(params.Throw)
}
