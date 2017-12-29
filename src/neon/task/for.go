// +build ignore

package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["for"] = build.TaskDescriptor{
		Constructor: For,
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
	}
}

func For(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"for", "in", "do"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	variable, err := args.GetString("for")
	if err != nil {
		return nil, fmt.Errorf("'for' field of a 'for' loop must be a string")
	}
	_list, err := args.GetList("in")
	expression := ""
	if err != nil {
		expression, err = args.GetString("in")
		if err != nil {
			return nil, fmt.Errorf("'in' field of 'for' loop must be a list or string")
		}
	}
	steps, err := ParseSteps(target, args, "do")
	if err != nil {
		return nil, err
	}
	return func(context *build.Context) error {
		if expression != "" {
			_result, _err := context.EvaluateExpression(expression)
			if _err != nil {
				return fmt.Errorf("evaluating in field of for loop: %v", _err)
			}
			_list, _err = util.ToList(_result)
			if _err != nil {
				return fmt.Errorf("'in' field of 'for' loop must be an expression that returns a list")
			}
		}
		for _, _value := range _list {
			context.SetProperty(variable, _value)
			_err := context.Run(steps)
			if _err != nil {
				return _err
			}
		}
		return nil
	}, nil
}
