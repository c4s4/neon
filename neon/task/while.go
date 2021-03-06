package task

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "while",
		Func: while,
		Args: reflect.TypeOf(whileArgs{}),
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

type whileArgs struct {
	While string
	Do    build.Steps `neon:"steps"`
}

func while(context *build.Context, args interface{}) error {
	params := args.(whileArgs)
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
		err = params.Do.Run(context)
		if err != nil {
			return err
		}
	}
	return nil
}
