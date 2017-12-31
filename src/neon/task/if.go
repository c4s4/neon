package task

import (
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "if",
		Func: If,
		Args: reflect.TypeOf(IfArgs{}),
		Help: `If condition.

Arguments:

- if: the condition.
- then: the steps to execute if the condition is true.
- else: the steps to execute if the condition is false.

Examples:

    # print hello if x > 10 else print world
    - if: x > 10
      then:
      - print: "hello"
      else:
      - print: "world"`,
	})
}

type IfArgs struct {
	If   bool         `expression`
	Then []build.Step `steps`
	Else []build.Step `steps`
}

func If(context *build.Context, args interface{}) error {
	params := args.(IfArgs)
	if params.If {
		err := context.Run(params.Then)
		if err != nil {
			return err
		}
	} else {
		err := context.Run(params.Else)
		if err != nil {
			return err
		}
	}
	return nil
}
