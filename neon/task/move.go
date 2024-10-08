package task

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "move",
		Func: move,
		Args: reflect.TypeOf(moveArgs{}),
		Help: `Move file(s).

Arguments:

- move: globs of files to move (strings, file, wrap).
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
      flat:  true

Notes:

- Parameter 'tofile' is valid if only one file was selected by globs.
- One and only one of parameters 'tofile' and 'todir' might be set.`,
	})
}

type moveArgs struct {
	Move    []string `neon:"file,wrap"`
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file,wrap"`
	Tofile  string   `neon:"optional,file"`
	Todir   string   `neon:"optional,file"`
	Flat    bool     `neon:"optional"`
}

func move(context *build.Context, args interface{}) error {
	params := args.(moveArgs)
	if (params.Tofile != "" && params.Todir != "") ||
		(params.Tofile == "" && params.Todir == "") {
		return fmt.Errorf("one and only one of parameters 'tofile' an 'todir' may be set")
	}
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
	context.MessageArgs("Moving %d file(s)", len(sources))
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
