package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
)

func init() {
	build.TaskMap["move"] = build.TaskDescriptor{
		Constructor: Move,
		Help: `Move file(s).

Arguments:

- move: the list of globs of files to move (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the file to move to (as a string, optional, only if glob selects a
  single file).
- todir: directory to move file(s) to (as a string, optional).
- flat: tells if files should be flatten in destination directory (as a boolean,
  optional, defaults to true).

Examples:

    # move file foo to bar
    - move:   "foo"
      tofile: "bar"
    # move text files in directory 'book' (except 'foo.txt') to directory 'text'
    - move: "**/*.txt"
      dir: "book"
      exclude: "**/foo.txt"
      todir: "text"
    # move all go sources to directory 'src', preserving directory structure
    - move: "**/*.go"
      todir: "src"
      flat: false`,
	}
}

func Move(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"move", "dir", "exclude", "tofile", "todir", "flat"}
	if err := CheckFields(args, fields, fields[:1]); err != nil {
		return nil, err
	}
	includes, err := args.GetListStringsOrString("move")
	if err != nil {
		return nil, fmt.Errorf("argument move must be a string or list of strings")
	}
	var dir string
	if args.HasField("dir") {
		dir, err = args.GetString("dir")
		if err != nil {
			return nil, fmt.Errorf("argument dir of task move must be a string")
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
			return nil, fmt.Errorf("argument tofile of task move must be a string")
		}
	}
	var toDir string
	if args.HasField("todir") {
		toDir, err = args.GetString("todir")
		if err != nil {
			return nil, fmt.Errorf("argument todir of task move must be a string")
		}
	}
	flat := true
	if args.HasField("flat") {
		flat, err = args.GetBoolean("flat")
		if err != nil {
			return nil, fmt.Errorf("argument flat of task move must be a boolean")
		}
	}
	if (tofile == "" && toDir == "") || (tofile != "" && toDir != "") {
		return nil, fmt.Errorf("move task must have one of 'to' or 'toDir' argument")
	}
	return func() error {
		// evaluate arguments
		_dir, _err := target.Build.Context.EvaluateString(dir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_tofile, _err := target.Build.Context.EvaluateString(tofile)
		if _err != nil {
			return fmt.Errorf("evaluating destination file: %v", _err)
		}
		_toDir, _err := target.Build.Context.EvaluateString(toDir)
		if _err != nil {
			return fmt.Errorf("evaluating destination directory: %v", _err)
		}
		_toDir = util.ExpandUserHome(_toDir)
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
		_sources, _err := target.Build.Context.FindFiles(_dir, _includes, _excludes, true)
		if _err != nil {
			return fmt.Errorf("getting source files for move task: %v", _err)
		}
		if _tofile != "" && len(_sources) > 1 {
			return fmt.Errorf("can't move more than one file to a given file, use todir instead")
		}
		if len(_sources) < 1 {
			return nil
		}
		build.Message("Moving %d file(s)", len(_sources))
		if _tofile != "" {
			_file := filepath.Join(_dir, _sources[0])
			if _file != _tofile {
				_err = os.Rename(_file, _tofile)
				if _err != nil {
					return fmt.Errorf("moving file: %v", _err)
				}
			}
		}
		if _toDir != "" {
			_err = util.MoveFilesToDir(_dir, _sources, _toDir, flat)
			if _err != nil {
				return fmt.Errorf("moving file: %v", _err)
			}
		}
		return nil
	}, nil
}
