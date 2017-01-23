package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["script"] = build.TaskDescriptor{
		Constructor: Script,
		Help:        "Run an Anko script",
	}
}

func Script(target *build.Target, args util.Object) (build.Task, error) {
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
