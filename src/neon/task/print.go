package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "print",
		Func: Print,
		Args: reflect.TypeOf(PrintArgs{}),
		Help: `Print a message on the console.

Arguments:

- print: the text to print as a string.

Examples:

    # say hello
    - print: "Hello World!"`,
	})
}

type PrintArgs struct {
	Print string
}

func Print(context *build.Context, args interface{}) error {
	params := args.(PrintArgs)
	context.Message(params.Print)
	return nil
}
