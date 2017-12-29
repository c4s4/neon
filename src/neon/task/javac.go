// +build ignore

package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"os/exec"
	"path"
)

func init() {
	build.TaskMap["javac"] = build.TaskDescriptor{
		Constructor: Javac,
		Help: `Compile Java source files.

Arguments:

- javac: the glob for Java source files.
- source: directory for source files.
- exclude: glob for source files to exclude (optional).
- dest: destination directory for generated classes.
- cp: classpath for compilation.

Examples:

    # compile Java source files in src directory
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
    # compile Java source files in src directory with given classpath
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
      cp:     '#{classpath}'`,
	}
}

func Javac(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"javac", "source", "dest", "exclude", "cp"}
	if err := CheckFields(args, fields, fields[:3]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("javac")
	if err != nil {
		return nil, fmt.Errorf("argument javac must be a string or list of strings")
	}
	source, err := args.GetString("source")
	if err != nil {
		return nil, fmt.Errorf("argument source must be a string")
	}
	dest, err := args.GetString("dest")
	if err != nil {
		return nil, fmt.Errorf("argument dest must be a string")
	}
	var excludes []string
	if args.HasField("exclude") {
		includes, err = args.GetListStringsOrString("exclude")
		if err != nil {
			return nil, fmt.Errorf("argument exclude must be a string or list of strings")
		}
	}
	var cp string
	if args.HasField("cp") {
		cp, err = args.GetString("cp")
		if err != nil {
			return nil, fmt.Errorf("argument cp must be a string")
		}
	}
	return func(context *build.Context) error {
		// find java source files
		_source, _err := context.EvaluateString(source)
		if _err != nil {
			return fmt.Errorf("evaluating source diectory: %v", _err)
		}
		_dest, _err := context.EvaluateString(dest)
		if _err != nil {
			return fmt.Errorf("evaluating destination diectory: %v", _err)
		}
		if !util.DirExists(_dest) {
			_err = os.MkdirAll(_dest, util.DIR_FILE_MODE)
			if _err != nil {
				return fmt.Errorf("making destination diectory: %v", _err)
			}
		}
		_sources, _err := context.FindFiles(source, includes, excludes, false)
		if _err != nil {
			return fmt.Errorf("finding java source files: %v", _err)
		}
		_cp, _err := context.EvaluateString(cp)
		if _err != nil {
			return fmt.Errorf("evaluating classpath: %v", _err)
		}
		// run javac command
		context.Message("Compiling %d Java source file(s)", len(_sources))
		_args := []string{"-d", _dest}
		if _cp != "" {
			_args = append(_args, []string{"-cp", _cp}...)
		}
		for _, _s := range _sources {
			_args = append(_args, path.Join(_source, _s))
		}
		command := exec.Command("javac", _args...)
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
			return fmt.Errorf("compiling java source files: %v", err)
		}
		return nil
	}, nil
}
