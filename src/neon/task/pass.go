package task

import (
	"neon/util"
)

func init() {
	TasksMap["pass"] = Descriptor{
		Constructor: Pass,
		Help:        "Does nothing",
	}
}

func Pass(target *Target, args util.Object) (Task, error) {
	fields := []string{"pass"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	return func() error {
		return nil
	}, nil
}
