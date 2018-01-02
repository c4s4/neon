package task

import (
	"neon/build"
	"reflect"
	"time"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "sleep",
		Func: Sleep,
		Args: reflect.TypeOf(SleepArgs{}),
		Help: `Sleep a given number of seconds.
		
Arguments:

- sleep: the duration to sleep in seconds as a float or integer.

Examples:

    # sleep for 1.5 seconds
    - sleep: 1.5
    # sleep for 3 seconds
    - sleep: 3`,
	})
}

type SleepArgs struct {
	Sleep float64
}

func Sleep(context *build.Context, args interface{}) error {
	params := args.(SleepArgs)
	context.Message("Sleeping for %g seconds...", params.Sleep)
	time.Sleep(time.Duration(params.Sleep) * time.Second)
	return nil
}
