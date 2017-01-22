package task

import (
	"neon/util"
)

func init() {
	TasksMap["try"] = Descriptor{
		Constructor: Try,
		Help:        "Try/catch construct",
	}
}

func Try(target *Target, args util.Object) (Task, error) {
	fields := []string{"try", "catch"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	trySteps, err := ParseSteps(target, args, "try")
	if err != nil {
		return nil, err
	}
	catchSteps, err := ParseSteps(target, args, "catch")
	if err != nil {
		return nil, err
	}
	return func() error {
		err := RunSteps(trySteps)
		if err != nil {
			target.Build.Context.SetProperty("error", err.Error())
			err = RunSteps(catchSteps)
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}
