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

- time: steps we want to measure execution duration (steps).
- to: property to store duration in seconds as a float, if not set, duration is
  printed on the console (string, optional).

Examples:

    # print duration to say hello
    - time:
      - print: 'Hello World!'
      to: duration
    - print: 'duration: ={duration}s'`,
	})
}

type TimeArgs struct {
	Time build.Steps `steps`
	To   string      `optional`
}

func Time(context *build.Context, args interface{}) error {
	params := args.(TimeArgs)
	start := time.Now()
	err := params.Time.Run(context)
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
