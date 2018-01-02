package task

import (
	"neon/build"
	"reflect"
	"time"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "time",
		Func: Time,
		Args: reflect.TypeOf(TimeArgs{}),
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
	})
}

type TimeArgs struct {
	Time []build.Step `steps`
	To   string
}

func Time(context *build.Context, args interface{}) error {
	params := args.(TimeArgs)
	start := time.Now()
	err := context.Run(params.Time)
	if err != nil {
		return err
	}
	duration := time.Now().Sub(start).Seconds()
	if params.To != "" {
		context.SetProperty(params.To, duration)
	} else {
		context.Message("Duration: %gs", duration)
	}
	return nil
}
