package task

import (
	"fmt"
	"neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "mkdir",
		Func: Mkdir,
		Args: reflect.TypeOf(MkdirArgs{}),
		Help: `Make a directory.

Arguments:

- mkdir: directories to create (strings, file, wrap).

Examples:

    # create a directory 'build'
    - mkdir: 'build'`,
	})
}

type MkdirArgs struct {
	Mkdir []string `file wrap`
}

func Mkdir(context *build.Context, args interface{}) error {
	params := args.(MkdirArgs)
	for _, dir := range params.Mkdir {
		context.Message("Making directory '%s'", dir)
		err := os.MkdirAll(dir, DIR_FILE_MODE)
		if err != nil {
			return fmt.Errorf("making directory '%s': %s", dir, err)
		}
	}
	return nil
}
