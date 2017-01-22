package task

import (
	"fmt"
	"neon/util"
)

func init() {
	TasksMap["for"] = Descriptor{
		Constructor: For,
		Help:        "For loop",
	}
}

func For(target *Target, args util.Object) (Task, error) {
	fields := []string{"for", "in", "do"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	variable, err := args.GetString("for")
	if err != nil {
		return nil, fmt.Errorf("'for' field of a 'for' loop must be a string")
	}
	list, err := args.GetList("in")
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
	return func() error {
		if expression != "" {
			result, err := target.Build.Context.Evaluate(expression)
			if err != nil {
				return fmt.Errorf("evaluating in field of for loop: %v", err)
			}
			list, err = util.ToList(result)
			if err != nil {
				return fmt.Errorf("'in' field of 'for' loop must be an expression that returns a list")
			}
		}
		for _, value := range list {
			target.Build.Context.SetProperty(variable, value)
			if err != nil {
				return err
			}
			err := RunSteps(steps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}
