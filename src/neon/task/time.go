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
- to: the property to store duration in seconds as a float (optional,
  print duration on console if not set).

Examples:

    # print duration to say hello
    - time:
      - print: "Hello World!"
      to: duration
    - print: 'duration: #{duration}s'`,
	}
}

func Time(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"time", "to"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	steps, err := ParseSteps(target, args, "time")
	if err != nil {
		return nil, err
	}
	var to string
	if args.HasField("to") {
		to, err = args.GetString("to")
		if err != nil {
			return nil, fmt.Errorf("argument to of task time must be a string")
		}
	}
	return func(context *build.Context) error {
		_to, _err := context.VM.EvaluateString(to)
		if _err != nil {
			return fmt.Errorf("evaluating property: %v", _err)
		}
		_start := time.Now()
		_err = RunSteps(steps, context)
		if _err != nil {
			return _err
		}
		_duration := time.Now().Sub(_start).Seconds()
		if to != "" {
			context.VM.SetProperty(_to, _duration)
		} else {
			context.Message("Duration: %gs", _duration)
		}
		return nil
	}, nil
}
