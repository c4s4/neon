package task

import (
	"fmt"
	"neon/util"
)

func init() {
	TasksMap["while"] = Descriptor{
		Constructor: While,
		Help:        "While loop",
	}
}

func While(target *Target, args util.Object) (Task, error) {
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
			result, err := target.Build.Context.Evaluate(condition)
			if err != nil {
				return fmt.Errorf("evaluating 'while' field of 'while' loop: %v", err)
			}
			loop, ok := result.(bool)
			if !ok {
				return fmt.Errorf("evaluating 'while' condition: must return a bool")
			}
			if !loop {
				break
			}
			err = RunSteps(steps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}
