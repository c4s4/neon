package task

import (
	"neon/build"
	"neon/util"
	"time"
)

func init() {
	build.TaskMap["time"] = build.TaskDescriptor{
		Constructor: Time,
		Help: `Print duration to run a block of steps.

Arguments:

- time: the steps to measure execution duration.

Examples:

    # measure duration to say hello
    - time:
      - print: "Hello World!"`,
	}
}

func Time(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"time"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	steps, err := ParseSteps(target, args, "time")
	if err != nil {
		return nil, err
	}
	return func() error {
		_start := time.Now()
		build.Message("Starting timer...")
		_err := RunSteps(target.Build, steps)
		if _err != nil {
			return _err
		}
		_duration := time.Now().Sub(_start)
		build.Message("Duration: %s", _duration)
		return nil
	}, nil
}
