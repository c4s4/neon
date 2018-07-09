package task

import (
	"neon/build"
	"net"
	"reflect"
	t "time"
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
- wait: tells if we wait until resource if released (if true) or stop on error
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

func singleton(context *build.Context, args interface{}) error {
	params := args.(singletonArgs)
	var listener net.Listener
	var err error
	if params.Wait {
		for listener == nil {
			listener, _ = build.ListenPort(params.Singleton)
		}
		if listener == nil {
			t.Sleep(1.0)
		}
	} else {
		listener, err = build.ListenPort(params.Singleton)
		if err != nil {
			return err
		}
	}
	if listener != nil {
		defer listener.Close()
	}
	err = params.Steps.Run(context)
	if err != nil {
		return err
	}
	return nil
}
