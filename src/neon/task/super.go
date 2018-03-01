package task

import (
	"fmt"
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
	ok, err := target.Build.RunParentTarget(context, target.Name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("no target '%s' found in parent build files", context.Stack.Last())
	}
	return nil
}
