package task

import (
	"fmt"
	"neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "mkdir",
		Func: Mkdir,
		Args: reflect.TypeOf(MkdirArgs{}),
		Help: `Make a directory.

Arguments:

- mkdir: directory or list of directories to create.

Examples:

    # create a directory 'build'
    - mkdir: "build"`,
	})
}

type MkdirArgs struct {
	Mkdir string `file`
}

func Mkdir(context *build.Context, args interface{}) error {
	params := args.(MkdirArgs)
	context.Message("Making directory '%s'", params.Mkdir)
	err := os.MkdirAll(params.Mkdir, DIR_FILE_MODE)
	if err != nil {
		return fmt.Errorf("making directory '%s': %s", params.Mkdir, err)
	}
	return nil
}
