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
		eval, err := target.Build.Context.EvaluateString(dir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		dir = eval
		eval, err = target.Build.Context.EvaluateString(tofile)
		if err != nil {
			return fmt.Errorf("evaluating destination file: %v", err)
		}
		tofile = eval
		eval, err = target.Build.Context.EvaluateString(toDir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		toDir = util.ExpandUserHome(eval)
		// find source files
		sources, err := target.Build.Context.FindFiles(dir, includes, excludes)
		if err != nil {
			return fmt.Errorf("getting source files for move task: %v", err)
		}
		if tofile != "" && len(sources) > 1 {
			return fmt.Errorf("can't move more than one file to a given file, use todir instead")
		}
		if len(sources) < 1 {
			return nil
		}
		build.Info("Moving %d file(s)", len(sources))
		if tofile != "" {
			file := filepath.Join(dir, sources[0])
			err = os.Rename(file, tofile)
			if err != nil {
				return fmt.Errorf("moving file: %v", err)
			}
		}
		if toDir != "" {
			err = util.MoveFilesToDir(dir, sources, toDir, flat)
			if err != nil {
				return fmt.Errorf("moving file: %v", err)
			}
		}
		return nil
	}, nil
}
