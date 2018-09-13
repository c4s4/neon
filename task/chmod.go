package task

import (
	"fmt"
	"github.com/c4s4/neon/build"
	util "github.com/c4s4/neon/util"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "chmod",
		Func: chmod,
		Args: reflect.TypeOf(chmodArgs{}),
		Help: `Change mode of files.

Arguments:

- chmod: list of globs of files to change mode (strings, file, wrap).
- mode: mode to change to (integer).
- dir: the root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).

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

type chmodArgs struct {
	Chmod   []string `neon:"file,wrap"`
	Mode    int
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file,wrap"`
}

func chmod(context *build.Context, args interface{}) error {
	params := args.(chmodArgs)
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
			return fmt.Errorf("changing mode of file '%s' to %#o: %v", file, params.Mode, err)
		}
	}
	return nil
}
