package task

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"io/ioutil"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "cat",
		Func: cat,
		Args: reflect.TypeOf(catArgs{}),
		Help: `Print the content of a given file on the console.

Arguments:

- cat: the name of the file to print on console (string, file).

Examples:

    # print content of LICENSE file on the console
    - cat: "LICENSE"`,
	})
}

type catArgs struct {
	Cat string `neon:"file"`
}

func cat(context *build.Context, args interface{}) error {
	params := args.(catArgs)
	content, err := ioutil.ReadFile(params.Cat)
	if err != nil {
		return fmt.Errorf("printing content of file '%s': %v", params.Cat, err)
	}
	context.Message(string(content))
	return nil
}
