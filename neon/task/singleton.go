package task

import (
	"fmt"
	"net"
	"reflect"
	t "time"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "singleton",
		Func: singleton,
		Args: reflect.TypeOf(singletonArgs{}),
		Help: `Ensure that only one instance of a block of steps is running.

Arguments:

- singleton: port we are listening to, should be between 1024 and 65535 (integer).
- steps: steps we want to run (steps).
- wait: tells if we wait until resource is released (if true) or stop on error
  (if false, which is default) (bool, optional).

Examples:

    # ensure one single instance is waiting
    - singleton: 12345
      steps:
      - sleep: 10.0`,
	})
}

type singletonArgs struct {
	Singleton int
	Steps     build.Steps `neon:"steps"`
	Wait      bool        `neon:"optional"`
}

// singletonRetryDelay is the delay between attempts to listen singleton port
const singletonRetryDelay = 20 * t.Millisecond

func singleton(context *build.Context, args interface{}) error {
	params := args.(singletonArgs)
	var listener net.Listener
	var err error
	if params.Wait {
		if params.Singleton < 1 || params.Singleton > 65535 {
			return fmt.Errorf("singleton port must be between 1 and 65535")
		}
		for {
			listener, err = build.ListenPort(params.Singleton)
			if err == nil {
				break
			}
			// short pause, to avoid busy loop while reacting quickly
			t.Sleep(singletonRetryDelay)
		}
	} else {
		listener, err = build.ListenPort(params.Singleton)
		if err != nil {
			return err
		}
	}
	defer func() {
		_ = listener.Close()
	}()
	return params.Steps.Run(context)
}
