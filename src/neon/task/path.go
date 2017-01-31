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
- path: the list of globs of files to build the path (as a string or list of strings).
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
		eval, err := target.Build.Context.ReplaceProperties(dir)
		if err != nil {
			return fmt.Errorf("evaluating destination directory: %v", err)
		}
		dir = eval
		// find source files
		files, err := target.Build.Context.FindFiles(dir, includes, excludes)
		if err != nil {
			return fmt.Errorf("getting source files for path task: %v", err)
		}
		if len(files) < 1 {
			return nil
		}
		target.Build.Info("Building path with %d file(s)", len(files))
		path := strings.Join(files, string(filepath.ListSeparator))
		target.Build.Context.SetProperty(to, path)
		return nil
	}, nil
}
