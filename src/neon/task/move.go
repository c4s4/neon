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

- move: globs of files to move (strings, file, wrap)
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: file to move file to (string, optional, file).
- todir: directory to move file(s) to (string, optional, file).
- flat: tells if files should be flatten in destination directory, defaults to
  false (boolean, optional).

Examples:

    # move file foo to bar
    - move:   'foo'
      tofile: 'bar'
    # move text files in directory 'book' (except 'foo.txt') to directory 'text'
    - move:    '**/*.txt'
      dir:     'book'
      exclude: '**/foo.txt'
      todir:   'text'
    # move all go sources to directory 'src', flattening structure
    - move:  '**/*.go'
      todir: 'src'
      flat:  true`,
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
