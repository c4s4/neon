package task

import (
	"fmt"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "super",
		Func: Super,
		Args: reflect.TypeOf(SuperArgs{}),
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

type SuperArgs struct{}

func Super(context *build.Context, args interface{}) error {
	ok, err := context.Build.RunParentTarget(context.Stack.Last(), context)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("no target '%s' found in parent build files", context.Stack.Last())
	}
	return nil
}
