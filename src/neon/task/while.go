package task

import (
	"fmt"
	"neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "while",
		Func: While,
		Args: reflect.TypeOf(WhileArgs{}),
		Help: `While loop.

Arguments:

- while: condition evaluated at each iteration (string).
- do: steps that run while condition is true (steps).

Examples:

    # loop until i >= 10
    - while: 'i < 10'
      do:
      - script: 'println(i++)'`,
	})
}

type WhileArgs struct {
	While string
	Do    []build.Step `steps`
}

func While(context *build.Context, args interface{}) error {
	params := args.(WhileArgs)
	for {
		result, err := context.EvaluateExpression(params.While)
		if err != nil {
			return fmt.Errorf("evaluating 'while' expression: %v", err)
		}
		loop, ok := result.(bool)
		if !ok {
			return fmt.Errorf("evaluating 'while' expression: must return a bool")
		}
		if !loop {
			break
		}
		err = context.Run(params.Do)
		if err != nil {
			return err
		}
	}
	return nil
}
