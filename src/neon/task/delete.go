package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
)

func init() {
	build.TaskMap["delete"] = build.TaskDescriptor{
		Constructor: Delete,
		Help: `Delete files or directories (recursively).

Arguments:

- delete: glob to select files or directory to delete.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # delete build directory
    - delete: "#{BUILD_DIR}"
    # delete all XML files except 'foo.xml'
    - delete:  "**/*.xml"
      exclude: "**/foo.xml"

Notes:

- Handle with care, directories are deleted recursively!`,
	}
}

func Delete(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"delete", "dir", "exclude"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("delete")
	if err != nil {
		return nil, fmt.Errorf("delete argument must be string or list of strings")
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task delete must be a string")
		}
	}
	var excludes []string
	if args.HasField("exclude") {
		excludes, err = args.GetListStringsOrString("exclude")
		if err != nil {
			return nil, fmt.Errorf("argument exclude must be string or list of strings")
		}
	}
	return func(context *build.Context) error {
		// evaluate arguments
		_dir, _err := context.VM.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_includes := make([]string, len(includes))
		for _index, _include := range includes {
			_includes[_index], _err = context.VM.EvaluateString(_include)
			if _err != nil {
				return fmt.Errorf("evaluating includes: %v", _err)
			}
		}
		_excludes := make([]string, len(excludes))
		for _index, _exclude := range excludes {
			_excludes[_index], _err = context.VM.EvaluateString(_exclude)
			if _err != nil {
				return fmt.Errorf("evaluating excludes: %v", _err)
			}
		}
		// find files to delete
		_files, _err := context.VM.FindFiles(_dir, _includes, _excludes, true)
		if _err != nil {
			return fmt.Errorf("getting source files for delete task: %v", _err)
		}
		if len(_files) < 1 {
			return nil
		}
		build.Message("Deleting %d file(s) or directory(ies)", len(_files))
		for _, _file := range _files {
			_path := filepath.Join(_dir, _file)
			if util.DirExists(_path) {
				_err = os.RemoveAll(_path)
				if _err != nil {
					return fmt.Errorf("deleting directory '%s': %v", _path, _err)
				}
			} else {
				if _err := os.Remove(_path); _err != nil {
					return fmt.Errorf("deleting file '%s': %v", _path, _err)
				}
			}
		}
		return nil
	}, nil
}
