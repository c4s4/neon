package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"time"
)

func init() {
	build.TaskMap["time"] = build.TaskDescriptor{
		Constructor: Time,
		Help: `Record duration to run a block of steps.

Arguments:

- time: the steps to measure execution duration.
- to: the property to store duration in seconds as a float.

Examples:

    # measure duration to say hello
    - time:
      - print: "Hello World!"
      to: duration
    - print: 'duration: #{duration}s'`,
	}
}

func Time(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"time", "to"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	steps, err := ParseSteps(target, args, "time")
	if err != nil {
		return nil, err
	}
	to, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task time must be a string")
	}
	return func() error {
		_to, _err := target.Build.Context.EvaluateString(to)
		if _err != nil {
			return fmt.Errorf("evaluating property: %v", _err)
		}
		_start := time.Now()
		_err = RunSteps(target.Build, steps)
		if _err != nil {
			return _err
		}
		_duration := time.Now().Sub(_start).Seconds()
		target.Build.Context.SetProperty(_to, _duration)
		return nil
	}, nil
}
