package task

import (
	"fmt"
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "mkdir",
		Func: mkdir,
		Args: reflect.TypeOf(mkdirArgs{}),
		Help: `Make a directory.

Arguments:

- mkdir: directories to create (strings, file, wrap).

Examples:

    # create a directory 'build'
    - mkdir: 'build'`,
	})
}

type mkdirArgs struct {
	Mkdir []string `neon:"file,wrap"`
}

func mkdir(context *build.Context, args interface{}) error {
	params := args.(mkdirArgs)
	for _, dir := range params.Mkdir {
		if !util.DirExists(dir) {
			context.Message("Making directory '%s'", dir)
			err := os.MkdirAll(dir, DirFileMode)
			if err != nil {
				return fmt.Errorf("making directory '%s': %s", dir, err)
			}
		}
	}
	return nil
}
