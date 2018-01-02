package task

import (
	"fmt"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "assert",
		Func: Assert,
		Args: reflect.TypeOf(AssertArgs{}),
		Help: `Make an assertion and fail if assertion condition is false.

Arguments:

- assert: the assertion to perform as an expression.

Examples:

    # assert that foo == "bar", and fail otherwise
    - assert: 'foo == "bar"'`,
	})
}

type AssertArgs struct {
	Assert bool `expression`
}

func Assert(context *build.Context, args interface{}) error {
	params := args.(AssertArgs)
	if !params.Assert {
		return fmt.Errorf("assertion failed")
	}
	return nil
}
