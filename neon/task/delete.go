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
		Name: "delete",
		Func: delete,
		Args: reflect.TypeOf(deleteArgs{}),
		Help: `Delete files or directories (recursively).

Arguments:

- delete: glob of files or directories to delete (strings, file, wrap).
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).

Examples:

    # delete build directory
    - delete: =BUILD_DIR
    # delete all XML files except 'foo.xml'
    - delete:  "**/*.xml"
      exclude: "**/foo.xml"

Notes:

- Handle with care: if globs select directories, they are deleted recursively!`,
	})
}

type deleteArgs struct {
	Delete  []string `neon:"file,wrap"`
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file,wrap"`
}

func delete(context *build.Context, args interface{}) error {
	params := args.(deleteArgs)
	files, err := util.FindFiles(params.Dir, params.Delete, params.Exclude, true)
	if err != nil {
		return fmt.Errorf("getting source files for delete task: %v", err)
	}
	if len(files) < 1 {
		return nil
	}
	context.MessageArgs("Deleting %d file(s) or directory(ies)", len(files))
	for _, file := range files {
		path := filepath.Join(params.Dir, file)
		if util.DirExists(path) {
			err = os.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("deleting directory '%s': %v", path, err)
			}
		} else {
			if err = os.Remove(path); err != nil {
				return fmt.Errorf("deleting file '%s': %v", path, err)
			}
		}
	}
	return nil
}
