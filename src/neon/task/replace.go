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
		eval, err := target.Build.Context.ReplaceProperties(pattern)
		if err != nil {
			return fmt.Errorf("evaluating pattern: %v", err)
		}
		pattern = eval
		eval, err = target.Build.Context.ReplaceProperties(with)
		if err != nil {
			return fmt.Errorf("evaluating with: %v", err)
		}
		with = eval
		eval, err = target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		dir = eval
		// find source files
		files, err := target.Build.Context.FindFiles(dir, includes, excludes)
		if err != nil {
			return fmt.Errorf("getting source files for copy task: %v", err)
		}
		if len(files) < 1 {
			return nil
		}
		target.Build.Info("Replacing text in %d file(s)", len(files))
		for _, file := range files {
			if dir != "" {
				file = filepath.Join(dir, file)
			}
			content, err := ioutil.ReadFile(file)
			if err != nil {
				return fmt.Errorf("reading file '%s': %v", file, err)
			}
			replaced := strings.Replace(string(content), pattern, with, -1)
			err = ioutil.WriteFile(file, []byte(replaced), FILE_MODE)
			if err != nil {
				return fmt.Errorf("writing file '%s': %v", file, err)
			}
		}
		return nil
	}, nil
}
