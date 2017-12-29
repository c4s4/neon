// +build ignore

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

- replace: the globs of files to work with (as a string or list of strings).
- with: map with replacements.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # replace foo with bar in file test.txt
    - replace: "test.txt"
      with:    {"foo": "bar"}`,
	}
}

func Replace(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"replace", "with", "dir", "exclude"}
	if err := CheckFields(args, fields, fields[:2]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("replace")
	if err != nil {
		return nil, fmt.Errorf("argument replace must be a string or list of strings")
	}
	with, err := args.GetMapStringString("with")
	if err != nil {
		return nil, fmt.Errorf("argument with of task replace must be a map")
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
	return func(context *build.Context) error {
		// evaluate arguments
		_eval, _err := context.EvaluateObject(with)
		if _err != nil {
			return fmt.Errorf("evaluating with: %v", _err)
		}
		_with, _err := util.ToMapStringString(_eval)
		if _err != nil {
			return fmt.Errorf("evaluating with: %v", _err)
		}
		_dir, _err := context.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_includes := make([]string, len(includes))
		for _index, _include := range includes {
			_includes[_index], _err = context.EvaluateString(_include)
			if _err != nil {
				return fmt.Errorf("evaluating includes: %v", _err)
			}
		}
		_excludes := make([]string, len(excludes))
		for _index, _exclude := range excludes {
			_excludes[_index], _err = context.EvaluateString(_exclude)
			if _err != nil {
				return fmt.Errorf("evaluating excludes: %v", _err)
			}
		}
		// find source files
		_files, _err := context.FindFiles(_dir, _includes, _excludes, false)
		if _err != nil {
			return fmt.Errorf("getting source files for copy task: %v", _err)
		}
		if len(_files) < 1 {
			return nil
		}
		for _, _file := range _files {
			context.Message("Replacing text in file '%s'", _file)
			if _dir != "" {
				_file = filepath.Join(_dir, _file)
			}
			_bytes, _err := ioutil.ReadFile(_file)
			if _err != nil {
				return fmt.Errorf("reading file '%s': %v", _file, _err)
			}
			_text := string(_bytes)
			for _old, _new := range _with {
				_text = strings.Replace(_text, _old, _new, -1)
			}
			_err = ioutil.WriteFile(_file, []byte(_text), FILE_MODE)
			if _err != nil {
				return fmt.Errorf("writing file '%s': %v", _file, _err)
			}
		}
		return nil
	}, nil
}
