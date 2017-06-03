package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["super"] = build.TaskDescriptor{
		Constructor: Super,
		Help: `Call target with same name in parent build file.

Arguments:

- none

Examples:

    # call parent target
    - super:

Notes:

- This will raise en error if parent build files have no target with same name.`,
	}
}

func Super(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"super"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	return func() error {
		ok, err := target.Build.RunParentTarget(target.Name)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no target '%s' found in parent build files", target.Name)
		}
		return nil
	}, nil
}
