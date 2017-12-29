// +build ignore

package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
)

func init() {
	build.TaskMap["java"] = build.TaskDescriptor{
		Constructor: Java,
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
	}
}

func Java(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"java", "cp", "args"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	main, err := args.GetString("java")
	if err != nil {
		return nil, fmt.Errorf("argument java must be a string")
	}
	cp, err := args.GetString("cp")
	if err != nil {
		return nil, fmt.Errorf("argument cp must be a string")
	}
	var params []string
	if args.HasField("args") {
		params, err = args.GetListStringsOrString("args")
		if err != nil {
			return nil, fmt.Errorf("argument args must be a string or list of strings")
		}
	}
	return func(context *build.Context) error {
		// find java source files
		_main, _err := context.EvaluateString(main)
		if _err != nil {
			return fmt.Errorf("evaluating main java class: %v", _err)
		}
		_cp, _err := context.EvaluateString(cp)
		if _err != nil {
			return fmt.Errorf("evaluating cp: %v", _err)
		}
		var _params []string
		for _, _param := range params {
			_p, _err := context.EvaluateString(_param)
			if _err != nil {
				return fmt.Errorf("evaluating argument: %v", _err)
			}
			_params = append(_params, _p)
		}
		// run java command
		_args := []string{"-cp", _cp, _main}
		for _, _p := range _params {
			_args = append(_args, _p)
		}
		command := exec.Command("java", _args...)
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current working directory: %v", err)
		}
		command.Dir = dir
		command.Env, err = context.EvaluateEnvironment(target.Build)
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
	}, nil
}
