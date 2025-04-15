package task

import (
	"os"
	"reflect"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "setenv",
		Func: setenv,
		Args: reflect.TypeOf(setenvArgs{}),
		Help: `Set environment variable with given value.

Arguments:

- setenv: environment variable name (string).
- value: value of this environment variable (string).

Examples:

    # set environment variable VERSION to value "1.2.3"
    - setenv: 'VERSION'
      value:  '1.2.3'`,
	})
}

type setenvArgs struct {
	Setenv string
	Value  string
}

func setenv(context *build.Context, args interface{}) error {
	params := args.(setenvArgs)
	return os.Setenv(params.Setenv, params.Value)
}
