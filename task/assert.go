package task

import (
	"fmt"
	"github.com/c4s4/neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "assert",
		Func: assert,
		Args: reflect.TypeOf(assertArgs{}),
		Help: `Make an assertion and fail if assertion is false.

Arguments:

- assert: the assertion to perform (boolean, expression).

Examples:

    # assert that foo == "bar", and fail otherwise
    - assert: 'foo == "bar"'`,
	})
}

type assertArgs struct {
	Assert bool `neon:"expression"`
}

func assert(context *build.Context, args interface{}) error {
	params := args.(assertArgs)
	if !params.Assert {
		return fmt.Errorf("assertion failed")
	}
	return nil
}
