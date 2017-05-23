package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	build.TaskMap["execute"] = build.TaskDescriptor{
		Constructor: Execute,
		Help: `Execute a command and return output and value.

Arguments:

- execute: command to run.
- output: name of the variable to store trimed output into.

Examples:

    # execute ls command and get result in 'files' variable
    - execute: 'ls'
      output:  'files'`,
	}
}

func Execute(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"execute", "output"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	execute, ok := args["execute"].(string)
	if !ok {
		return nil, fmt.Errorf("argument of task execute must be a string")
	}
	output, ok := args["output"].(string)
	if !ok {
		return nil, fmt.Errorf("argument output of task execute must be a string")
	}
	return func() error {
		_cmd, _err := target.Build.Context.EvaluateString(execute)
		if _err != nil {
			return fmt.Errorf("processing execute argument: %v", _err)
		}
		var _command *exec.Cmd
		if runtime.GOOS == "windows" {
			_command = exec.Command("cmd.exe", "/C", _cmd)
		} else {
			_command = exec.Command("sh", "-c", _cmd)
		}
		_dir, _err := os.Getwd()
		if _err != nil {
			return fmt.Errorf("getting current working directory: %v", _err)
		}
		_command.Dir = _dir
		_command.Env, _err = target.Build.Context.EvaluateEnvironment()
		if _err != nil {
			return fmt.Errorf("building environment: %v", _err)
		}
		_bytes, _err := _command.CombinedOutput()
		target.Build.Context.SetProperty(output, strings.TrimSpace(string(_bytes)))
		if _err != nil {
			return fmt.Errorf("executing command: %v", _err)
		}
		return nil
	}, nil
}
