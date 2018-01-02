package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "move",
		Func: Move,
		Args: reflect.TypeOf(MoveArgs{}),
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
	})
}

type MoveArgs struct {
	Move    []string `file wrap`
	Dir     string   `optional file`
	Exclude []string `optional file wrap`
	Tofile  string   `optional file`
	Todir   string   `optional file`
	Flat    bool     `optional`
}

func Move(context *build.Context, args interface{}) error {
	params := args.(MoveArgs)
	sources, err := util.FindFiles(params.Dir, params.Move, params.Exclude, true)
	if err != nil {
		return fmt.Errorf("getting source files for move task: %v", err)
	}
	if params.Tofile != "" && len(sources) > 1 {
		return fmt.Errorf("can't move more than one file to a given file, use todir instead")
	}
	if len(sources) < 1 {
		return nil
	}
	context.Message("Moving %d file(s)", len(sources))
	if params.Tofile != "" {
		file := filepath.Join(params.Dir, sources[0])
		if file != params.Tofile {
			err = os.Rename(file, params.Tofile)
			if err != nil {
				return fmt.Errorf("moving file: %v", err)
			}
		}
	}
	if params.Todir != "" {
		err = util.MoveFilesToDir(params.Dir, sources, params.Todir, params.Flat)
		if err != nil {
			return fmt.Errorf("moving file: %v", err)
		}
	}
	return nil
}
