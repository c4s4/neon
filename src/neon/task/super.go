package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "super",
		Func: super,
		Args: reflect.TypeOf(superArgs{}),
		Help: `Call target with same name in parent build file.

Arguments:

- none

Examples:

    # call parent target
    - super:

Notes:

- This will raise en error if parent build files have no target with same name.`,
	})
}

type superArgs struct{}

func super(context *build.Context, args interface{}) error {
	target := context.Stack.Last()
	return target.Build.RunParentTarget(context, target.Name)
}
