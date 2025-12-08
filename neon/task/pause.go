package task

import (
	"fmt"
	"reflect"
	t "time"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "pause",
		Func: pause,
		Args: reflect.TypeOf(pauseArgs{}),
		Help: `Wait for given duration.

Arguments:

- pause: duration to pause (string).
- mute: if set to true, do not print a message (bool, optional).

Examples:

    # pause for 3 seconds
    - pause: 3s
    # pause for 1 minutes without message
    - pause: 1m
      mute:  true`,
	})
}

type pauseArgs struct {
	Pause string
	Mute bool `neon:"optional"`
}

func pause(context *build.Context, args interface{}) error {
	params := args.(pauseArgs)
	if !params.Mute {
		context.MessageArgs("Pausing for %s seconds...", params.Pause)
	}
	duration, err := t.ParseDuration(params.Pause)
	if err != nil {
		return fmt.Errorf("invalid pause duration '%s': %v", params.Pause, err)
	}
	t.Sleep(duration)
	return nil
}
