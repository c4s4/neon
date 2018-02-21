package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "call",
		Func: call,
		Args: reflect.TypeOf(callArgs{}),
		Help: `Call a build target.

Arguments:

- call: the name of the target(s) to call (strings, wrap).

Examples:

    # call target 'foo'
    - call: 'foo'`,
	})
}

type callArgs struct {
	Call []string `wrap`
}

func call(context *build.Context, args interface{}) error {
	params := args.(callArgs)
	for _, target := range params.Call {
		context.Message("Calling target '%s'", target)
		err := context.Build.RunTarget(context, target)
		if err != nil {
			return err
		}
	}
	return nil
}
