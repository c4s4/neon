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
		eval, err := target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		dir = eval
		eval, err = target.Build.Context.ReplaceProperties(tofile)
		if err != nil {
			return fmt.Errorf("evaluating destination file: %v", err)
		}
		tofile = eval
		eval, err = target.Build.Context.ReplaceProperties(toDir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		toDir = util.ExpandUserHome(eval)
		// find source files
		sources, err := target.Build.Context.FindFiles(dir, includes, excludes)
		if err != nil {
			return fmt.Errorf("getting source files for copy task: %v", err)
		}
		if tofile != "" && len(sources) > 1 {
			return fmt.Errorf("can't copy more than one file to a given file, use todir instead")
		}
		if len(sources) < 1 {
			return nil
		}
		build.Info("Copying %d file(s)", len(sources))
		if tofile != "" {
			file := filepath.Join(dir, sources[0])
			err = util.CopyFile(file, tofile)
			if err != nil {
				return fmt.Errorf("copying file: %v", err)
			}
		}
		if toDir != "" {
			err = util.CopyFilesToDir(dir, sources, toDir, flat)
			if err != nil {
				return fmt.Errorf("copying file: %v", err)
			}
		}
		return nil
	}, nil
}
