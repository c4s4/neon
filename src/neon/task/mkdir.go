package task

import (
	"fmt"
	"neon/util"
	"os"
)

func init() {
	TasksMap["mkdir"] = Descriptor{
		Constructor: MkDir,
		Help:        "Make a directory",
	}
}

func MkDir(target *Target, args util.Object) (Task, error) {
	fields := []string{"mkdir"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	dir, ok := args["mkdir"].(string)
	if !ok {
		return nil, fmt.Errorf("argument to task mkdir must be a string")
	}
	return func() error {
		directory, err := target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("processing mkdir argument: %v", err)
		}
		fmt.Printf("Making directory '%s'\n", directory)
		err = os.MkdirAll(directory, DIR_FILE_MODE)
		if err != nil {
			return fmt.Errorf("making directory '%s': %s", directory, err)
		}
		return nil
	}, nil
}
