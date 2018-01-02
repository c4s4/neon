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
		Func: Shell,
		Args: reflect.TypeOf(ShellArgs{}),
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run as a string or a list of strings.
- =: name of the variable to store trimmed output into (optional, output to
  console if not set).

Examples:

    # execute ls command and get result in 'files' variable
    - $: 'ls -al'
      =: 'files'
    # execute command as a list of strings
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

type ShellArgs struct {
	Shell []string `name:"$" wrap`
	To    string   `name:"=" optional`
}

func Shell(context *build.Context, args interface{}) error {
	params := args.(ShellArgs)
	output, err := Run(params.Shell, params.To == "", context)
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

// Run execute a command and return its output and an error (if command
// returned a value different from 0). Arguments:
// - pipe tells if we should print the output of the command on the console.
func Run(command []string, pipe bool, context *build.Context) (string, error) {
	if len(command) == 0 {
		return "", fmt.Errorf("empty command")
	} else if len(command) < 2 {
		return RunString(command[0], pipe, context)
	} else {
		return RunList(command, pipe, context)
	}
}

// Run execute the command and returns its output and an error. Arguments:
// - pipe tells if we should print output of the command on the console.
func RunList(cmd []string, pipe bool, context *build.Context) (string, error) {
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

// Run execute the command and returns its output and an error. Arguments:
// - pipe tells if we should print output of the command on the console.
func RunString(cmd string, pipe bool, context *build.Context) (string, error) {
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
