package task

import (
	"fmt"
	"neon/build"
	"os"
	"os/exec"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "java",
		Func: Java,
		Args: reflect.TypeOf(JavaArgs{}),
		Help: `Run Java virtual machine.

Arguments:

- javac: the main Java class name.
- cp: classpath for runtime.
- args: command line arguments (optional).

Examples:

    # run class foo.Bar with arguments foo and bar
    - javac: 'foo.Bar'
      cp:    'build/classes'
      args:  ['foo', 'bar']`,
	})
}

type JavaArgs struct {
	Java string
	Cp   string
	Args []string `optional`
}

func Java(context *build.Context, args interface{}) error {
	params := args.(JavaArgs)
	opt := []string{"-cp", params.Cp, params.Java}
	opt = append(opt, params.Args...)
	command := exec.Command("java", opt...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	// FIXME
	//command.Env, err = context.EvaluateEnvironment(context.target.Build)
	if err != nil {
		return fmt.Errorf("building environment: %v", err)
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		return fmt.Errorf("running Java virtual machine: %v", err)
	}
	return nil
}
