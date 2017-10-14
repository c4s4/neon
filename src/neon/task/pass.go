package task

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["pass"] = build.TaskDescriptor{
		Constructor: Pass,
		Help: `Does nothing.

Arguments:

- none

Examples:

    # do nothing
    - pass:

Notes:

- This implementation is super optimized for speed.`,
	}
}

func Pass(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"pass"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	return func(context *build.Context) error {
		return nil
	}, nil
}
