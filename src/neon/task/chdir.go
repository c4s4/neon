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
		Help:        "Change current working directory",
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
		directory, err := target.Build.Context.ReplaceProperties(dir)
		fmt.Printf("Changing to directory '%s'\n", directory)
		if err != nil {
			return fmt.Errorf("processing chdir argument: %v", err)
		}
		err = os.Chdir(directory)
		if err != nil {
			return fmt.Errorf("changing to directory '%s': %s", directory, err)
		}
		return nil
	}, nil
}
