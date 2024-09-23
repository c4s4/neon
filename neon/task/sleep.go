package task

import (
	"github.com/c4s4/neon/neon/build"
	"reflect"
	t "time"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "sleep",
		Func: sleep,
		Args: reflect.TypeOf(sleepArgs{}),
		Help: `Sleep given number of seconds.

Arguments:

- sleep: duration to sleep in seconds (float).
- mute: if set to true, do not print a message (bool, optional).

Examples:

    # sleep for 1.5 seconds
    - sleep: 1.5
    # sleep for 3 seconds without message
    - sleep: 3.0
      mute: true`,
	})
}

type sleepArgs struct {
	Sleep float64
	Mute  bool `neon:"optional"`
}

func sleep(context *build.Context, args interface{}) error {
	params := args.(sleepArgs)
	if !params.Mute {
		context.Message("Sleeping for %g seconds...", params.Sleep)
	}
	t.Sleep(t.Duration(params.Sleep) * t.Second)
	return nil
}
