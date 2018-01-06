package task

import (
	"fmt"
	"neon/build"
	"os"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "link",
		Func: Link,
		Args: reflect.TypeOf(LinkArgs{}),
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

type LinkArgs struct {
	Link string `file`
	To   string `file`
}

func Link(context *build.Context, args interface{}) error {
	params := args.(LinkArgs)
	context.Message("Linking file '%s' to '%s'", params.Link, params.To)
	err := os.Symlink(params.Link, params.To)
	if err != nil {
		return fmt.Errorf("linking files: %v", err)
	}
	return nil
}
