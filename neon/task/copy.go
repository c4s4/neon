package task

import (
	"fmt"
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "copy",
		Func: copy,
		Args: reflect.TypeOf(copyArgs{}),
		Help: `Copy file(s).

Arguments:

- copy: globs of files to copy (strings, file, wrap).
- dir: root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: file to copy file to (string, optional, file).
- todir: directory to copy files to (string, optional, file).
- flat: tells if files should be flatten in destination directory, defaults to
  false (boolean, optional).

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
      flat: false

Notes:

- Parameter 'tofile' is valid if only one file was selected by globs.
- One and only one of parameters 'tofile' and 'todir' might be set.`,
	})
}

type copyArgs struct {
	Copy    []string `neon:"file,wrap"`
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file,wrap"`
	Tofile  string   `neon:"optional,file"`
	Todir   string   `neon:"optional,file"`
	Flat    bool     `neon:"optional"`
}

func copy(context *build.Context, args interface{}) error {
	params := args.(copyArgs)
	if (params.Tofile != "" && params.Todir != "") ||
		(params.Tofile == "" && params.Todir == "") {
		return fmt.Errorf("one and only one of parameters 'tofile' an 'todir' may be set")
	}
	sources, err := util.FindFiles(params.Dir, params.Copy, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for copy task: %v", err)
	}
	if params.Tofile != "" && len(sources) > 1 {
		return fmt.Errorf("can't copy more than one file to a given file, use todir instead")
	}
	if len(sources) < 1 {
		return nil
	}
	context.MessageArgs("Copying %d file(s)", len(sources))
	if params.Tofile != "" {
		file := filepath.Join(params.Dir, sources[0])
		err = util.CopyFile(file, params.Tofile)
		if err != nil {
			return fmt.Errorf("copying file: %v", err)
		}
	}
	if params.Todir != "" {
		err = util.CopyFilesToDir(params.Dir, sources, params.Todir, params.Flat)
		if err != nil {
			return fmt.Errorf("copying file: %v", err)
		}
	}
	return nil
}
