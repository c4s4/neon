package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["while"] = build.TaskDescriptor{
		Constructor: While,
		Help: `While loop.

Arguments:

- while: the condition that is evaluated at each loop.
- do: steps that run while condition is true.

Examples:

    # loop until i >= 10
    - while: 'i < 10'
      do:
      - script: 'println(i++)'`,
	}
}

func While(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"while", "do"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	condition, err := args.GetString("while")
	if err != nil {
		return nil, fmt.Errorf("'while' field of a 'while' loop must be a string")
	}
	steps, err := ParseSteps(target, args, "do")
	if err != nil {
		return nil, err
	}
	return func() error {
		for {
			_result, _err := target.Build.Context.EvaluateExpression(condition)
			if _err != nil {
				return fmt.Errorf("evaluating 'while' field of 'while' loop: %v", _err)
			}
			_loop, _ok := _result.(bool)
			if !_ok {
				return fmt.Errorf("evaluating 'while' condition: must return a bool")
			}
			if !_loop {
				break
			}
			_err = RunSteps(target.Build, steps)
			if _err != nil {
				return _err
			}
		}
		return nil
	}, nil
}
