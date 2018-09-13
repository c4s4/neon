package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "print",
		Func: print,
		Args: reflect.TypeOf(printArgs{}),
		Help: `Print a message on the console.

Arguments:

- print: text to print (string).

Examples:

    # say hello
    - print: 'Hello World!'`,
	})
}

type printArgs struct {
	Print string
}

func print(context *build.Context, args interface{}) error {
	params := args.(printArgs)
	context.Message(params.Print)
	return nil
}
