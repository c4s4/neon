package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["throw"] = build.TaskDescriptor{
		Constructor: Throw,
		Help:        "Throws an error",
	}
}

func Throw(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"throw"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	message, ok := args["throw"].(string)
	if !ok {
		return nil, fmt.Errorf("argument of throw print must be a string")
	}
	return func() error {
		return fmt.Errorf(message)
	}, nil
}
