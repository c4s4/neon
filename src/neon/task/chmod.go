package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
	"strconv"
)

func init() {
	build.TaskMap["chmod"] = build.TaskDescriptor{
		Constructor: Chmod,
		Help: `Changes mode of files.

Arguments:

- chmod: the list of globs of files to change mode (as a string or list of
  strings).
- mode: the mode in octal form (such as '0755') as a string
- dir: the root directory for glob (as a string, optional, defaults to '.').
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # make foo.sh executable for all users
    - chmod: "foo.sh"
      mod: "0755"
    # make all sh files in foo directory executable, except for bar.sh
    - chmod: "**/*.sh"
      mode: "0755"
      exclude: "**/bar.sh"`,
	}
}

func Chmod(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"chmod", "mode", "dir", "exclude"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("chmod")
	if err != nil {
		return nil, fmt.Errorf("argument chmod must be a string or list of strings")
	}
	mode, err := args.GetString("mode")
	if err != nil {
		return nil, fmt.Errorf("argument mode of task chmod must be a string")
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task chmod must be a string")
		}
	}
	var excludes []string
	if args.HasField("exclude") {
		excludes, err = args.GetListStringsOrString("exclude")
		if err != nil {
			return nil, fmt.Errorf("argument exclude of task chmod must be string or list of strings")
		}
	}
	return func() error {
		// evaluate arguments
		eval, err := target.Build.Context.EvaluateString(dir)
		if err != nil {
			return fmt.Errorf("evaluating directory: %v", err)
		}
		dir = eval
		eval, err = target.Build.Context.EvaluateString(mode)
		if err != nil {
			return fmt.Errorf("evaluating mode: %v", err)
		}
		mode = eval
		// find source files
		files, err := target.Build.Context.FindFiles(dir, includes, excludes)
		if err != nil {
			return fmt.Errorf("getting source files for chmod task: %v", err)
		}
		modeBase8, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return fmt.Errorf("converting mode '%s' in octal: %v", mode, err)
		}
		if len(files) < 1 {
			return nil
		}
		build.Info("Changing %d file(s) mode to %s", len(files), mode)
		for _, file := range files {
			if dir != "" {
				file = filepath.Join(dir, file)
			}
			err := os.Chmod(file, os.FileMode(modeBase8))
			if err != nil {
				return fmt.Errorf("changing mode of file '%s' to %s: %v", file, mode, err)
			}
		}
		return nil
	}, nil
}
