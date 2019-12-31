package task

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "link",
		Func: link,
		Args: reflect.TypeOf(linkArgs{}),
		Help: `Create a symbolic link.

Arguments:

- link: source file (string, file).
- to: destination of the link (string, file).

Examples:

    # create a link from file 'foo' to 'bar'
    - link: 'foo''
      to:   'bar''`,
	})
}

type linkArgs struct {
	Link string `neon:"file"`
	To   string `neon:"file"`
}

func link(context *build.Context, args interface{}) error {
	params := args.(linkArgs)
	context.Message("Linking file '%s' to '%s'", params.Link, params.To)
	err := os.Symlink(params.Link, params.To)
	if err != nil {
		return fmt.Errorf("linking files: %v", err)
	}
	return nil
}
