package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"path/filepath"
	"strings"
)

func init() {
	build.TaskMap["path"] = build.TaskDescriptor{
		Constructor: Path,
		Help: `Build a path from files and put it in a variable.

Arguments:

- path: the list of globs of files to build the path (as a string or list of
  strings).
- to: the variable to put path into.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # build classpath with jar files in lib directory
    - path: "lib/*.jar"
      to: "classpath"`,
	}
}

func Path(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"path", "to", "dir", "exclude"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("path")
	if err != nil {
		return nil, fmt.Errorf("argument path must be a string or list of strings")
	}
	to, err := args.GetString("to")
	if err != nil {
		return nil, fmt.Errorf("argument to of task replace must be a string")
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task path must be a string")
		}
	}
	var excludes []string
	if args.HasField("exclude") {
		excludes, err = args.GetListStringsOrString("exclude")
		if err != nil {
			return nil, fmt.Errorf("argument exclude ot task path must be string or list of strings")
		}
	}
	return func() error {
		// evaluate arguments
		_dir, _err := target.Build.Context.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_includes := make([]string, len(includes))
		for _index, _include := range includes {
			_includes[_index], _err = target.Build.Context.EvaluateString(_include)
			if _err != nil {
				return fmt.Errorf("evaluating includes: %v", _err)
			}
		}
		_excludes := make([]string, len(excludes))
		for _index, _exclude := range excludes {
			_excludes[_index], _err = target.Build.Context.EvaluateString(_exclude)
			if _err != nil {
				return fmt.Errorf("evaluating excludes: %v", _err)
			}
		}
		// find source files
		_files, _err := target.Build.Context.FindFiles(_dir, _includes, _excludes, true)
		if _err != nil {
			return fmt.Errorf("getting source files for path task: %v", _err)
		}
		if len(_files) < 1 {
			return nil
		}
		build.Info("Building path with %d file(s)", len(_files))
		path := strings.Join(_files, string(filepath.ListSeparator))
		target.Build.Context.SetProperty(to, path)
		return nil
	}, nil
}
