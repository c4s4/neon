package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "if",
		Func: if_,
		Args: reflect.TypeOf(ifArgs{}),
		Help: `If condition.

Arguments:

- if: the condition (boolean, expression).
- then: steps to execute if condition is true (steps).
- else: steps to execute if condition is false (optional, steps).

Examples:

    # print hello if x > 10 else print world
    - if: x > 10
      then:
      - print: "hello"
      else:
      - print: "world"`,
	})
}

type ifArgs struct {
	If   bool        `expression`
	Then build.Steps `steps`
	Else build.Steps `optional steps`
}

func if_(context *build.Context, args interface{}) error {
	params := args.(ifArgs)
	if params.If {
		err := params.Then.Run(context)
		if err != nil {
			return err
		}
	} else {
		err := params.Else.Run(context)
		if err != nil {
			return err
		}
	}
	return nil
}
