package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"time"
)

func init() {
	build.TaskMap["sleep"] = build.TaskDescriptor{
		Constructor: Sleep,
		Help: `Sleep a given number of seconds.
		
Arguments:

- sleep: the duration to sleep in seconds as a float or integer.

Examples:

    # sleep for 1.5 seconds
    - sleep: 1.5
    # sleep for 3 seconds
    - sleep: 3`,
	}
}

func Sleep(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"sleep"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	var duration float64
	switch time := args["sleep"].(type) {
	case int:
		duration = float64(time)
	case float64:
		duration = time
	default:
		return nil, fmt.Errorf("argument of task sleep must be a float or an int")
	}
	return func(context *build.Context) error {
		context.Message("Sleeping for %g seconds...", duration)
		time.Sleep(time.Duration(duration) * time.Second)
		return nil
	}, nil
}
