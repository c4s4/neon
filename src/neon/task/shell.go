package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"strings"
	"runtime"
)

func init() {
	build.TaskMap["$"] = build.TaskDescriptor{
		Constructor: Shell,
		Help: `Execute a command and return output and value.

Arguments:

- $: command to run or a map of commands per operating system.
- =: name of the variable to store trimed output into (optional, output to
  console if not set).

Examples:

    # execute ls command and get result in 'files' variable
    - $: 'ls'
      =: 'files'
    # execute dir command on windows and ls on other OS
    - $:
    	windows: 'dir'
    	default: 'ls'`,
	}
}

func Shell(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"$", "="}
	var err error
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	var cmds map[string]string
	if util.IsString(args["$"]) {
		cmds = map[string]string {
			"default": args["$"].(string),
		}
	} else if util.IsMap(args["$"]) {
		cmds, err = util.ToMapStringString(args["$"])
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Errorf("shell command must be a string or a map of strings")
	}
	var output string
	var ok bool
	if args.HasField("=") {
		output, ok = args["="].(string)
		if !ok {
			return nil, fmt.Errorf("argument = of task $ must be a string")
		}
	}
	return func() error {
		_cmd, _ok := cmds[runtime.GOOS]
		if !_ok {
			_cmd, _ok = cmds["default"]
			if !_ok {
				return fmt.Errorf("no command found for '%s'", runtime.GOOS)
			}
		}
		_cmd, _err := target.Build.Context.EvaluateString(_cmd)
		if _err != nil {
			return fmt.Errorf("processing $ argument: %v", _err)
		}
		_output, _err := target.Build.Context.EvaluateString(output)
		if _err != nil {
			return fmt.Errorf("processing output argument: %v", _err)
		}
		var _command *exec.Cmd
		_shell, _err := target.Build.GetShell()
		if _err != nil {
			return _err
		}
		_binary := _shell[0]
		_arguments := _shell[1:]
		_arguments = append(_arguments, _cmd)
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
