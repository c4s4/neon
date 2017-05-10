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
		cmd, err := target.Build.Context.ReplaceProperties(execute)
		if err != nil {
			return fmt.Errorf("processing execute argument: %v", err)
		}
		var command *exec.Cmd
		if runtime.GOOS == "windows" {
			command = exec.Command("cmd.exe", "/C", cmd)
		} else {
			command = exec.Command("sh", "-c", cmd)
		}
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current working directory: %v", err)
		}
		command.Dir = dir
		command.Env, err = target.Build.Context.GetEnvironment()
		if err != nil {
			return fmt.Errorf("building environment: %v", err)
		}
		bytes, err := command.CombinedOutput()
		target.Build.Context.SetProperty(output, strings.TrimSpace(string(bytes)))
		if err != nil {
			return fmt.Errorf("executing command: %v", err)
		}
		return nil
	}, nil
}
