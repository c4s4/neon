package task

import (
	"github.com/c4s4/neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "pass",
		Func: pass,
		Args: reflect.TypeOf(passArgs{}),
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

type passArgs struct{}

func pass(context *build.Context, args interface{}) error {
	return nil
}
