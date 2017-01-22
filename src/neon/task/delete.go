package task

import (
	"fmt"
	"neon/util"
	"os"
)

func init() {
	TasksMap["delete"] = Descriptor{
		Constructor: Delete,
		Help:        "Delete a directory recursively",
	}
}

func Delete(target *Target, args util.Object) (Task, error) {
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
				fmt.Printf("Deleting directory '%s'\n", directory)
				err = os.RemoveAll(directory)
				if err != nil {
					return fmt.Errorf("deleting directory '%s': %v", directory, err)
				}
			}
		}
		return nil
	}, nil
}
