package task

import (
	"neon/build"
	"reflect"
	"time"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "sleep",
		Func: sleep,
		Args: reflect.TypeOf(sleepArgs{}),
		Help: `Sleep given number of seconds.
		
Arguments:

- sleep: duration to sleep in seconds (float).

Examples:

    # sleep for 1.5 seconds
    - sleep: 1.5
    # sleep for 3 seconds (3.0 as a float)
    - sleep: 3.0`,
	})
}

type sleepArgs struct {
	Sleep float64
}

func sleep(context *build.Context, args interface{}) error {
	params := args.(sleepArgs)
	context.Message("Sleeping for %g seconds...", params.Sleep)
	time.Sleep(time.Duration(params.Sleep) * time.Second)
	return nil
}
