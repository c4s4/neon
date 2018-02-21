package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "for",
		Func: for_,
		Args: reflect.TypeOf(forArgs{}),
		Help: `For loop.

Arguments:

- for: variable name to set at each loop iteration (string).
- in: values or expression to generate values to iterate on (list or
  expression).
- do: steps to execute at each loop iteration (steps).

Examples:

    # create empty files
    - for: file
      in:  ["foo", "bar"]
      do:
    - touch: =file
    # print first 10 integers
    - for: i
      in: range(10)
      do:
      - print: '={i}'`,
	})
}

type forArgs struct {
	For string
	In  []interface{} `expression`
	Do  build.Steps   `steps`
}

func for_(context *build.Context, args interface{}) error {
	params := args.(forArgs)
	for _, value := range params.In {
		context.SetProperty(params.For, value)
		err := params.Do.Run(context)
		if err != nil {
			return err
		}
	}
	return nil
}
