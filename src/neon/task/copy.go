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
			return nil, fmt.Errorf("argument exclude mus be string or list of strings")
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
	return func() error {
		// evaluate arguments
		_eval, err := target.Build.Context.EvaluateString(dir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		_dir := _eval
		_eval, err = target.Build.Context.EvaluateString(tofile)
		if err != nil {
			return fmt.Errorf("evaluating destination file: %v", err)
		}
		_tofile := _eval
		_eval, err = target.Build.Context.EvaluateString(toDir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		_toDir := util.ExpandUserHome(_eval)
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
		_sources, err := target.Build.Context.FindFiles(_dir, _includes, _excludes)
		if err != nil {
			return fmt.Errorf("getting source files for copy task: %v", err)
		}
		if _tofile != "" && len(_sources) > 1 {
			return fmt.Errorf("can't copy more than one file to a given file, use todir instead")
		}
		if len(_sources) < 1 {
			return nil
		}
		build.Info("Copying %d file(s)", len(_sources))
		if _tofile != "" {
			file := filepath.Join(_dir, _sources[0])
			err = util.CopyFile(file, _tofile)
			if err != nil {
				return fmt.Errorf("copying file: %v", err)
			}
		}
		if _toDir != "" {
			err = util.CopyFilesToDir(_dir, _sources, _toDir, flat)
			if err != nil {
				return fmt.Errorf("copying file: %v", err)
			}
		}
		return nil
	}, nil
}
