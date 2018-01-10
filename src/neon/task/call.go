package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "call",
		Func: Call,
		Args: reflect.TypeOf(CallArgs{}),
		Help: `Call a build target.

Arguments:

- call: the name of the target(s) to call (strings, wrap).

Examples:

    # call target 'foo'
    - call: 'foo'`,
	})
}

type CallArgs struct {
	Call []string `wrap`
}

func Call(context *build.Context, args interface{}) error {
	params := args.(CallArgs)
	for _, target := range params.Call {
		context.Message("Calling target '%s'", target)
		context.Build.RunTarget(context, target)
	}
	return nil
}
