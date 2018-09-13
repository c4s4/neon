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
		Func: java,
		Args: reflect.TypeOf(javaArgs{}),
		Help: `Run Java virtual machine.

Arguments:

- javac: main Java class name (string).
- cp: classpath to run main class (string).
- args: command line arguments (strings, optional, wrap).

Examples:

    # run class foo.Bar with arguments foo and bar
    - javac: 'foo.Bar'
      cp:    'build/classes'
      args:  ['foo', 'bar']`,
	})
}

type javaArgs struct {
	Java string
	Cp   string
	Args []string `neon:"optional,wrap"`
}

func java(context *build.Context, args interface{}) error {
	params := args.(javaArgs)
	opt := []string{"-cp", params.Cp, params.Java}
	opt = append(opt, params.Args...)
	command := exec.Command("java", opt...)
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment()
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
