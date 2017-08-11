package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"strings"
)

func init() {
	build.TaskMap["$"] = build.TaskDescriptor{
		Constructor: Shell,
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run.
- =: name of the variable to store trimed output into (optional, output to
  console if not set).

Examples:

    # execute ls command and get result in 'files' variable
    - $: 'ls'
      =: 'files'`,
	}
}

func Shell(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"$", "="}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	shell, ok := args["$"].(string)
	if !ok {
		return nil, fmt.Errorf("argument of task $ must be a string")
	}
	var output string
	if args.HasField("=") {
		output, ok = args["="].(string)
		if !ok {
			return nil, fmt.Errorf("argument = of task $ must be a string")
		}
	}
	return func() error {
		_shell, _err := target.Build.Context.EvaluateString(shell)
		if _err != nil {
			return fmt.Errorf("processing $ argument: %v", _err)
		}
		_output, _err := target.Build.Context.EvaluateString(output)
		if _err != nil {
			return fmt.Errorf("processing output argument: %v", _err)
		}
		var _command *exec.Cmd
		_binary := target.Build.Shell[0]
		_arguments := target.Build.Shell[1:]
		_arguments = append(_arguments, _shell)
		_command = exec.Command(_binary, _arguments...)
		_dir, _err := os.Getwd()
		if _err != nil {
			return fmt.Errorf("getting current working directory: %v", _err)
		}
		_command.Dir = _dir
		_command.Env, _err = target.Build.Context.EvaluateEnvironment()
		if _err != nil {
			return fmt.Errorf("building environment: %v", _err)
		}
		if _output == "" {
			_command.Stdin = os.Stdin
			_command.Stdout = os.Stdout
			_command.Stderr = os.Stderr
			_err = _command.Run()
			if _err != nil {
				return fmt.Errorf("executing command: %v", _err)
			}
		} else {
			_bytes, _err := _command.CombinedOutput()
			target.Build.Context.SetProperty(output, strings.TrimSpace(string(_bytes)))
			if _err != nil {
				return fmt.Errorf("executing command: %v", _err)
			}
		}
		return nil
	}, nil
}
