package task

import (
	"fmt"
	"neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "chdir",
		Func: Chdir,
		Args: reflect.TypeOf(ChdirArgs{}),
		Help: `Change current working directory.

Arguments:

- chdir: the directory to change to (as a string).

Examples:

    # go to build directory
    - chdir: "build"

Notes:

- Working directory is set to the build file directory before each target.`,
	})
}

type ChdirArgs struct {
	Chdir string `file`
}

func Chdir(context *build.Context, args interface{}) error {
	params := args.(ChdirArgs)
	context.Message("Changing to directory '%s'", params.Chdir)
	err := os.Chdir(params.Chdir)
	if err != nil {
		return fmt.Errorf("changing to directory '%s': %s", params.Chdir, err)
	}
	return nil
}
