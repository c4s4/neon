package task

import (
	"fmt"
	"neon/build"
	"os"
	"path/filepath"
	"reflect"
	util "neon/util"
)

func init() {
	build.AddTask(build.TaskDesc {
		Name: "chmod",
		Func: Chmod,
		Args: reflect.TypeOf(ChmodArgs{}),
		Help: `Change mode of files.

Arguments:

- chmod: the list of globs of files to change mode (as a string or list of
  strings).
- mode: the mode in octal form (such as '0755') as an integer.
- dir: the root directory for glob (as a string, optional, defaults to '.').
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # make foo.sh executable for all users
    - chmod: "foo.sh"
      mode:  0755
    # make all sh files in foo directory executable, except for bar.sh
    - chmod:   "**/*.sh"
      mode:    0755
      exclude: "**/bar.sh"

Notes:
- The mode is an integer, thus must not be surrounded with quotes, or it would
  be a string and parsing of the task would fail.
- We usually set mode with octal integers, starting with '0'. If you don't put
  starting '0', this is decimal integer and you won't probably have expected
  result.`,
	})
}

type ChmodArgs struct {
	Chmod []string   `file wrap`
	Mode  int
	Dir   string     `file optional`
	Exclude []string `optional`
}

func Chmod(context *build.Context, args interface{}) error {
	params := args.(ChmodArgs)
	files, err := util.FindFiles(params.Dir, params.Chmod, params.Exclude, true)
	if err != nil {
		return fmt.Errorf("getting source files for chmod task: %v", err)
	}
	if len(files) < 1 {
		return nil
	}
	context.Message("Changing %d file(s) mode to %#o", len(files), params.Mode)
	for _, file := range files {
		if params.Dir != "" {
			file = filepath.Join(params.Dir, file)
		}
		err := os.Chmod(file, os.FileMode(params.Mode))
		if err != nil {
			return fmt.Errorf("changing mode of file '%s' to %s: %v", file, params.Mode, err)
		}
	}
	return nil
}
