package task

import (
	"fmt"
	"neon/build"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "$",
		Func: shell,
		Args: reflect.TypeOf(shellArgs{}),
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run (string or list of strings).
- =: name of the variable to set with command output, output to console if not
  set (string, optional).

Examples:

    # execute ls command and get result in 'files' variable
    - $: 'ls -al'
      =: 'files'
    # execute command as a list of strings and output on console
    - $: ['ls', '-al']

Notes:

- Commands defined as a string run in the shell defined by shell field at the
  root of the build file (or 'sh -c' on Unix and 'cmd /c' on Windows by
  default).
- Defining a command as a list of strings is useful on Windows. Default shell on
  Windows is 'cmd' which can't properly manage arguments with spaces.
- Argument of a command defined as a list won't be expanded by shell. Thus
  $USER won't be expanded for instance.`,
	})
}

type shellArgs struct {
	Shell []string `name:"$" wrap`
	To    string   `name:"=" optional`
}

func shell(context *build.Context, args interface{}) error {
	params := args.(shellArgs)
	output, err := run(params.Shell, params.To == "", context)
	if err != nil {
		if output != "" {
			context.Message(output)
		}
		return err
	}
	if params.To != "" {
		context.SetProperty(params.To, strings.TrimSpace(string(output)))
	}
	return nil
}

func run(command []string, pipe bool, context *build.Context) (string, error) {
	if len(command) == 0 {
		return "", fmt.Errorf("empty command")
	} else if len(command) < 2 {
		return runString(command[0], pipe, context)
	} else {
		return runList(command, pipe, context)
	}
}

func runList(cmd []string, pipe bool, context *build.Context) (string, error) {
	command := exec.Command(cmd[0], cmd[1:]...)
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment()
	if err != nil {
		return "", fmt.Errorf("building environment: %v", err)
	}
	if pipe {
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			return "", fmt.Errorf("executing command: %v", err)
		}
		return "", nil
	} else {
		bytes, err := command.CombinedOutput()
		if err != nil {
			return string(bytes), fmt.Errorf("executing command: %v", err)
		}
		return string(bytes), nil
	}
}

func runString(cmd string, pipe bool, context *build.Context) (string, error) {
	shell, err := context.Build.GetShell()
	if err != nil {
		return "", err
	}
	binary := shell[0]
	arguments := shell[1:]
	arguments = append(arguments, cmd)
	command := exec.Command(binary, arguments...)
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %v", err)
	}
	command.Dir = dir
	command.Env, err = context.EvaluateEnvironment()
	if err != nil {
		return "", fmt.Errorf("building environment: %v", err)
	}
	if pipe {
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			return "", fmt.Errorf("executing command: %v", err)
		}
		return "", nil
	} else {
		bytes, err := command.CombinedOutput()
		if err != nil {
			return string(bytes), fmt.Errorf("executing command: %v", err)
		}
		return string(bytes), nil
	}
}
