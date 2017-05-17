package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
)

func init() {
	build.TaskMap["chdir"] = build.TaskDescriptor{
		Constructor: Chdir,
		Help: `Change current working directory.

Arguments:

- chdir: the directory to change to (as a string).

Examples:

    # go to build directory
    - chdir: "build"

Notes:

- Working directory is set to the build file directory before each target.`,
	}
}

func Chdir(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"chdir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	dir, ok := args["chdir"].(string)
	if !ok {
		return nil, fmt.Errorf("argument to task chdir must be a string")
	}
	return func() error {
		_directory, _err := target.Build.Context.EvaluateString(dir)
		build.Info("Changing to _directory '%s'", _directory)
		if _err != nil {
			return fmt.Errorf("processing chdir argument: %v", _err)
		}
		_err = os.Chdir(_directory)
		if _err != nil {
			return fmt.Errorf("changing to _directory '%s': %s", _directory, _err)
		}
		return nil
	}, nil
}
