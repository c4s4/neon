package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "read",
		Func: Read,
		Args: reflect.TypeOf(ReadArgs{}),
		Help: `Read given file as text and put its content in a variable.

Arguments:

- read: the file to read as a string.
- to: the name of the variable to set with the content.

Examples:

    # put content of LICENSE file on license variable
    - read: "LICENSE"
      to: license`,
	})
}

type ReadArgs struct {
	Read string `file`
	To   string
}

func Read(context *build.Context, args interface{}) error {
	params := args.(ReadArgs)
	content, err := ioutil.ReadFile(params.Read)
	if err != nil {
		return fmt.Errorf("reading content of file '%s': %v", params.Read, err)
	}
	context.SetProperty(params.To, string(content))
	return nil
}
