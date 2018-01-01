package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"path/filepath"
	"strings"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "path",
		Func: Path,
		Args: reflect.TypeOf(PathArgs{}),
		Help: `Build a path from files and put it in a variable.

Arguments:

- path: the list of globs of files to build the path (as a string or list of
  strings).
- to: the variable to put path into.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # build classpath with jar files in lib directory
    - path: "lib/*.jar"
      to: "classpath"`,
	})
}

type PathArgs struct {
	Path []string    `file wrap`
	To   string
	Dir  string      `file optional`
	Exclude []string `optional file wrap`
}

func Path(context *build.Context, args interface{}) error {
	params := args.(PathArgs)
	files, err := util.FindFiles(params.Dir, params.Path, params.Exclude, true)
	if err != nil {
		return fmt.Errorf("getting source files for path task: %v", err)
	}
	if len(files) < 1 {
		return nil
	}
	context.Message("Building path with %d file(s)", len(files))
	path := strings.Join(files, string(filepath.ListSeparator))
	context.SetProperty(params.To, path)
	return nil
}
