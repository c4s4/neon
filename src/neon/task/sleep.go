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

- sleep: the duration to sleep in seconds as an integer.

Examples:

    # sleep for 10 seconds
    - sleep: 10`,
	}
}

func Sleep(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"sleep"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	duration, ok := args["sleep"].(int)
	if !ok {
		return nil, fmt.Errorf("argument of task sleep must be an integer")
	}
	return func() error {
		build.Message("Sleeping for %d seconds...", duration)
		time.Sleep(time.Duration(duration) * time.Second)
		return nil
	}, nil
}
