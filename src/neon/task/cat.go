package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "cat",
		Func: Cat,
		Args: reflect.TypeOf(CatArgs{}),
		Help: `Print the content of a given file on the console.

Arguments:

- cat: the name of the file to print on console (string, file).

Examples:

    # print content of LICENSE file on the console
    - cat: "LICENSE"`,
	})
}

type CatArgs struct {
	Cat string `file`
}

func Cat(context *build.Context, args interface{}) error {
	params := args.(CatArgs)
	content, err := ioutil.ReadFile(params.Cat)
	if err != nil {
		return fmt.Errorf("printing content of file '%s': %v", params.Cat, err)
	}
	context.Message(string(content))
	return nil
}
