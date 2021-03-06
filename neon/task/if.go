package task

import (
	"github.com/c4s4/neon/neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "if",
		Func: ifFunc,
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
	If   bool        `neon:"expression"`
	Then build.Steps `neon:"steps"`
	Else build.Steps `neon:"optional,steps"`
}

func ifFunc(context *build.Context, args interface{}) error {
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
