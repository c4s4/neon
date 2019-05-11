package task

import (
	"fmt"
	"github.com/c4s4/neon/build"
	"io/ioutil"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "read",
		Func: read,
		Args: reflect.TypeOf(readArgs{}),
		Help: `Read given file as text and put its content in a variable.

Arguments:

- read: file to read (string, file).
- to: name of the variable to set with its content (string).

Examples:

    # put content of LICENSE file in license variable
    - read: 'LICENSE'
      to:   'license'`,
	})
}

type readArgs struct {
	Read string `neon:"file"`
	To   string
}

func read(context *build.Context, args interface{}) error {
	params := args.(readArgs)
	content, err := ioutil.ReadFile(params.Read)
	if err != nil {
		return fmt.Errorf("reading content of file '%s': %v", params.Read, err)
	}
	context.SetProperty(params.To, string(content))
	return nil
}
