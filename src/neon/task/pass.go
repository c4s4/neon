package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "pass",
		Func: Pass,
		Args: reflect.TypeOf(PassArgs{}),
		Help: `Does nothing.

Arguments:

- none

Examples:

    # do nothing
    - pass:

Notes:

- This implementation is super optimized for speed.`,
	})
}

type PassArgs struct {}

func Pass(context *build.Context, args interface{}) error {
	return nil
}
