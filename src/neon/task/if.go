package task

import (
	"fmt"
	"neon/util"
)

func init() {
	TasksMap["if"] = Descriptor{
		Constructor: If,
		Help:        "If condition",
	}
}

func If(target *Target, args util.Object) (Task, error) {
	fields := []string{"if", "then", "else"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	condition, err := args.GetString("if")
	if err != nil {
		return nil, fmt.Errorf("evaluating if construct: %v", err)
	}
	thenSteps, err := ParseSteps(target, args, "then")
	if err != nil {
		return nil, err
	}
	var elseSteps []Step
	if FieldPresent(args, "else") {
		elseSteps, err = ParseSteps(target, args, "else")
		if err != nil {
			return nil, err
		}
	}
	return func() error {
		result, err := target.Build.Context.Evaluate(condition)
		if err != nil {
			return fmt.Errorf("evaluating 'if' condition: %v", err)
		}
		boolean, ok := result.(bool)
		if !ok {
			return fmt.Errorf("evaluating if condition: must return a bool")
		}
		if boolean {
			err := RunSteps(thenSteps)
			if err != nil {
				return err
			}
		} else {
			err := RunSteps(elseSteps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}
