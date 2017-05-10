package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["print"] = build.TaskDescriptor{
		Constructor: Print,
		Help: `Print a message on the console.

Arguments:

- print: the text to print as a string.

Examples:

    # say hello
    - print: "Hello World!"`,
	}
}

func Print(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"print"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	message, ok := args["print"].(string)
	if !ok {
		return nil, fmt.Errorf("argument of task print must be a string")
	}
	return func() error {
		evaluated, err := target.Build.Context.ReplaceProperties(message)
		if err != nil {
			return fmt.Errorf("processing print argument: %v", err)
		}
		build.Info(evaluated)
		return nil
	}, nil
}
