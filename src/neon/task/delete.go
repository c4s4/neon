package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
)

func init() {
	build.TaskMap["delete"] = build.TaskDescriptor{
		Constructor: Delete,
		Help: `Delete a directory recursively.

Arguments:

- delete: directory or list of directories to delete.

Examples:

    # delete build directory
    - delete: "#{BUILD_DIR}"

Notes:

- Handle with care, this is recursive!`,
	}
}

func Delete(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"delete"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	directories, err := args.GetListStringsOrString("delete")
	if err != nil {
		return nil, fmt.Errorf("delete argument must be string or list of strings")
	}
	return func() error {
		for _, dir := range directories {
			directory, err := target.Build.Context.ReplaceProperties(dir)
			if err != nil {
				return fmt.Errorf("evaluating directory in task delete: %v", err)
			}
			if _, err := os.Stat(directory); err == nil {
				target.Build.Info("Deleting directory '%s'", directory)
				err = os.RemoveAll(directory)
				if err != nil {
					return fmt.Errorf("deleting directory '%s': %v", directory, err)
				}
			}
		}
		return nil
	}, nil
}
