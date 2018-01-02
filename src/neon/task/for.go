package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "for",
		Func: For,
		Args: reflect.TypeOf(ForArgs{}),
		Help: `For loop.

Arguments:

- for: the name of the variable to set at each loop iteration.
- in: the list of values or expression that generates this list.
- do: the block of steps to execute at each loop iteration.

Examples:

    # create empty files
    - for: file
      in:  ["foo", "bar"]
      do:
    - touch: "#{file}"
    # print first 10 integers
    - for: i
      in: range(10)
      do:
      - print: "#{i}"`,
	})
}

type ForArgs struct {
	For string
	In  []interface{} `expression`
	Do  []build.Step  `steps`
}

func For(context *build.Context, args interface{}) error {
	params := args.(ForArgs)
	for _, value := range params.In {
		context.SetProperty(params.For, value)
		err := context.Run(params.Do)
		if err != nil {
			return err
		}
	}
	return nil
}
