// +build ignore

package task

import (
	"fmt"
	"neon/build"
	"neon/util"
)

func init() {
	build.TaskMap["throw"] = build.TaskDescriptor{
		Constructor: Throw,
		Help: `Throws an error.

Arguments:

- throw: the message of the error.

Examples:

    # stop the build because tests don't run
    - throw: "ERROR: tests don't run"

Notes:

- The error message will be printed on the console as the source of the build
  failure.`,
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
	return func(context *build.Context) error {
		_message, _err := context.EvaluateString(message)
		if _err != nil {
			return fmt.Errorf("processing thow argument: %v", _err)
		}
		return fmt.Errorf(_message)
	}, nil
}
