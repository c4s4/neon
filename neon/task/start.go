package task

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "start",
		Func: start,
		Args: reflect.TypeOf(startArgs{}),
		Help: `Start a command in background and put its PID in a variable.

Arguments:

- start: command to run (string).
- pid: name of the variable to put PID into.

Examples:

    # start a command in background and put its PID in 'pid' variable
    - start: 'ls -al'
	  pid:  pid

Notes:

- Commands run in the shell defined by shell field at the root of the build file
  (or 'sh -c' on Unix and 'cmd /c' on Windows by default).`,
	})
}

type startArgs struct {
	Start string `neon:"name=start,string"`
	Pid   string `neon:"name=pid,string"`
}

func start(context *build.Context, args interface{}) error {
	params := args.(startArgs)
	cmd := exec.Command(params.Start)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var err error
	cmd.Env, err = context.EvaluateEnvironment()
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	cmd.Environ()
	if err := cmd.Start(); err != nil {
		return err
	}
	context.SetProperty(params.Pid, cmd.Process.Pid)
	return nil
}
