package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"path/filepath"
)

func init() {
	build.TaskMap["copy"] = build.TaskDescriptor{
		Constructor: Copy,
		Help: `Copy file(s).

Arguments:

- copy: the list of globs of files to copy (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the file to copy to (as a string, optional, only if glob selects a
  single file).
- todir: directory to copy file(s) to (as a string, optional).
- flat: tells if files should be flatten in destination directory (as a boolean,
  optional, defaults to true).

Examples:

    # copy file foo to bar
    - copy:   "foo"
      tofile: "bar"
    # copy text files in directory 'book' (except 'foo.txt') to directory 'text'
    - copy: "**/*.txt"
      dir: "book"
      exclude: "**/foo.txt"
      todir: "text"
    # copy all go sources to directory 'src', preserving directory structure
    - copy: "**/*.go"
      todir: "src"
      flat: false`,
	}
}

func Copy(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"copy", "dir", "exclude", "tofile", "todir", "flat"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("copy")
	if err != nil {
		return nil, fmt.Errorf("argument copy must be a string or list of strings")
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
			return nil, fmt.Errorf("argument exclude must be string or list of strings")
		}
	}
	var tofile string
	if args.HasField("tofile") {
		tofile, err = args.GetString("tofile")
		if err != nil {
			return nil, fmt.Errorf("argument tofile of task copy must be a string")
		}
	}
	var toDir string
	if args.HasField("todir") {
		toDir, err = args.GetString("todir")
		if err != nil {
			return nil, fmt.Errorf("argument todir of task copy must be a string")
		}
	}
	flat := true
	if args.HasField("flat") {
		flat, err = args.GetBoolean("flat")
		if err != nil {
			return nil, fmt.Errorf("argument flat of task copy must be a boolean")
		}
	}
	if (tofile == "" && toDir == "") || (tofile != "" && toDir != "") {
		return nil, fmt.Errorf("copy task must have one of 'to' or 'toDir' argument")
	}
	return func(context *build.Context) error {
		// evaluate arguments
		_eval, _err := context.VM.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_dir := _eval
		_eval, _err = context.VM.EvaluateString(tofile)
		if _err != nil {
			return fmt.Errorf("evaluating destination file: %v", _err)
		}
		_tofile := _eval
		_eval, _err = context.VM.EvaluateString(toDir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_toDir := util.ExpandUserHome(_eval)
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
		// find source files
		_sources, _err := context.VM.FindFiles(_dir, _includes, _excludes, false)
		if _err != nil {
			return fmt.Errorf("getting source files for copy task: %v", _err)
		}
		if _tofile != "" && len(_sources) > 1 {
			return fmt.Errorf("can't copy more than one file to a given file, use todir instead")
		}
		if len(_sources) < 1 {
			return nil
		}
		build.Message("Copying %d file(s)", len(_sources))
		if _tofile != "" {
			file := filepath.Join(_dir, _sources[0])
			_err = util.CopyFile(file, _tofile)
			if _err != nil {
				return fmt.Errorf("copying file: %v", _err)
			}
		}
		if _toDir != "" {
			_err = util.CopyFilesToDir(_dir, _sources, _toDir, flat)
			if _err != nil {
				return fmt.Errorf("copying file: %v", _err)
			}
		}
		return nil
	}, nil
}
