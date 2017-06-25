package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
)

func init() {
	build.TaskMap["mkdir"] = build.TaskDescriptor{
		Constructor: MkDir,
		Help: `Make a directory.

Arguments:

- mkdir: directory or list of directories to create.

Examples:

    # create a directory 'build'
    - mkdir: "build"`,
	}
}

func MkDir(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"mkdir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	dir, ok := args["mkdir"].(string)
	if !ok {
		return nil, fmt.Errorf("argument to task mkdir must be a string")
	}
	return func() error {
		_directory, _err := target.Build.Context.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("processing mkdir argument: %v", _err)
		}
		build.Message("Making directory '%s'", _directory)
		_err = os.MkdirAll(_directory, DIR_FILE_MODE)
		if _err != nil {
			return fmt.Errorf("making directory '%s': %s", _directory, _err)
		}
		return nil
	}, nil
}
