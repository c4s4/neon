package task

import (
	"fmt"
	"neon/util"
)

func init() {
	TasksMap["script"] = Descriptor{
		Constructor: Script,
		Help:        "Run an Anko script",
	}
}

func Script(target *Target, args util.Object) (Task, error) {
	fields := []string{"script"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	source, err := args.GetString("script")
	if err != nil {
		return nil, fmt.Errorf("parsing script task: %v", err)
	}
	return func() error {
		_, err := target.Build.Context.Evaluate(source)
		if err != nil {
			return fmt.Errorf("evaluating script: %v", err)
		}
		return nil
	}, nil
}
