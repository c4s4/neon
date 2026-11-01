package task

import (
	"fmt"
	"reflect"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "chdir",
		Func: chdir,
		Args: reflect.TypeOf(chdirArgs{}),
		Help: `Change current working directory.

Arguments:

- chdir: the directory to change to (string, file).

Examples:

    # go to build directory
    - chdir: "build"

Notes:

- Working directory is set to the build file directory before each target.
- Changing to another directory while threads are running raises an error, as
  current directory is shared by all threads.`,
	})
}

type chdirArgs struct {
	Chdir string `neon:"file"`
}

func chdir(context *build.Context, args interface{}) error {
	params := args.(chdirArgs)
	context.MessageArgs("Changing working directory to '%s'", params.Chdir)
	err := build.Chdir(params.Chdir)
	if err != nil {
		return fmt.Errorf("changing working directory to '%s': %s", params.Chdir, err)
	}
	return nil
}
