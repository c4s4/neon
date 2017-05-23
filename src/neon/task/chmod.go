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
		_dir, _err := target.Build.Context.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("evaluating directory: %v", _err)
		}
		_mode, _err := target.Build.Context.EvaluateString(mode)
		if _err != nil {
			return fmt.Errorf("evaluating _mode: %v", _err)
		}
		_includes := make([]string, len(includes))
		for index, _include := range includes {
			_includes[index], err = target.Build.Context.EvaluateString(_include)
			if err != nil {
				return fmt.Errorf("evaluating includes: %v", err)
			}
		}
		_excludes := make([]string, len(excludes))
		for index, _exclude := range excludes {
			_excludes[index], err = target.Build.Context.EvaluateString(_exclude)
			if err != nil {
				return fmt.Errorf("evaluating excludes: %v", err)
			}
		}
		// find source files
		_files, _err := target.Build.Context.FindFiles(_dir, _includes, _excludes)
		if _err != nil {
			return fmt.Errorf("getting source files for chmod task: %v", _err)
		}
		_modeBase8, _err := strconv.ParseUint(_mode, 8, 32)
		if _err != nil {
			return fmt.Errorf("converting mode '%s' in octal: %v", _mode, _err)
		}
		if len(_files) < 1 {
			return nil
		}
		build.Info("Changing %d file(s) mode to %s", len(_files), _mode)
		for _, _file := range _files {
			if _dir != "" {
				_file = filepath.Join(_dir, _file)
			}
			_err := os.Chmod(_file, os.FileMode(_modeBase8))
			if _err != nil {
				return fmt.Errorf("changing mode of file '%s' to %s: %v", _file, _mode, _err)
			}
		}
		return nil
	}, nil
}
