package task

import (
	"fmt"
	"io/ioutil"
	"neon/build"
	"neon/util"
	"path/filepath"
	"strings"
)

func init() {
	build.TaskMap["replace"] = build.TaskDescriptor{
		Constructor: Replace,
		Help: `Replace pattern in text files.

Arguments:

- replace: the list of globs of files to work with (as a string or list of strings).
- pattern: the text to replace.
- with: the replacement text.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # replace foo with bar in file test.txt
    - replace: "test.txt"
      pattern: "foo"
      with: "bar"`,
	}
}

func Replace(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"replace", "pattern", "with", "dir", "exclude"}
	if err := CheckFields(args, fields, fields[:3]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("replace")
	if err != nil {
		return nil, fmt.Errorf("argument replace must be a string or list of strings")
	}
	pattern, err := args.GetString("pattern")
	if err != nil {
		return nil, fmt.Errorf("argument pattern of task replace must be a string")
	}
	with, err := args.GetString("with")
	if err != nil {
		return nil, fmt.Errorf("argument with of task replace must be a string")
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task copy must be a string")
		}
	}
	var excludes []string
	if args.HasField("exclude") {
		excludes, err = args.GetListStringsOrString("exclude")
		if err != nil {
			return nil, fmt.Errorf("argument exclude mus be string or list of strings")
		}
	}
	return func() error {
		// evaluate arguments
		_pattern, _err := target.Build.Context.EvaluateString(pattern)
		if _err != nil {
			return fmt.Errorf("evaluating pattern: %v", _err)
		}
		_with, _err := target.Build.Context.EvaluateString(with)
		if _err != nil {
			return fmt.Errorf("evaluating with: %v", _err)
		}
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
		_files, _err := target.Build.Context.FindFiles(_dir, _includes, _excludes)
		if _err != nil {
			return fmt.Errorf("getting source files for copy task: %v", _err)
		}
		if len(_files) < 1 {
			return nil
		}
		build.Info("Replacing text in %d file(s)", len(_files))
		for _, _file := range _files {
			if _dir != "" {
				_file = filepath.Join(_dir, _file)
			}
			_content, _err := ioutil.ReadFile(_file)
			if _err != nil {
				return fmt.Errorf("reading file '%s': %v", _file, _err)
			}
			_replaced := strings.Replace(string(_content), _pattern, _with, -1)
			_err = ioutil.WriteFile(_file, []byte(_replaced), FILE_MODE)
			if _err != nil {
				return fmt.Errorf("writing file '%s': %v", _file, _err)
			}
		}
		return nil
	}, nil
}
