package task

import (
	"fmt"
	"neon/build"
	"neon/util"
	"path/filepath"
	"reflect"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "copy",
		Func: Copy,
		Args: reflect.TypeOf(CopyArgs{}),
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
  optional, defaults to false).

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
	})
}

type CopyArgs struct {
	Copy    []string `file wrap`
	Dir     string   `optional file`
	Exclude []string `optional file wrap`
	Tofile  string   `optional file`
	Todir   string   `optional file`
	Flat    bool     `optional`
}

func Copy(context *build.Context, args interface{}) error {
	params := args.(CopyArgs)
	// find source files
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
	context.Message("Copying %d file(s)", len(sources))
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
