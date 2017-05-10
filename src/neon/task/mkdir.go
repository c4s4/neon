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
		directory, err := target.Build.Context.EvaluateString(dir)
		if err != nil {
			return fmt.Errorf("processing mkdir argument: %v", err)
		}
		build.Info("Making directory '%s'", directory)
		err = os.MkdirAll(directory, DIR_FILE_MODE)
		if err != nil {
			return fmt.Errorf("making directory '%s': %s", directory, err)
		}
		return nil
	}, nil
}
