package task

import (
	"fmt"
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "replace",
		Func: replace,
		Args: reflect.TypeOf(replaceArgs{}),
		Help: `Replace text matching patterns in files.

Arguments:

- replace: globs of files to process (strings, file, wrap).
- with: map with replacements (map with string keys and values).
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, files).

Examples:

    # replace foo with bar in file test.txt
    - replace: 'test.txt'
      with:    {'foo': 'bar'}`,
	})
}

type replaceArgs struct {
	Replace []string `neon:"file,wrap"`
	With    map[string]string
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file"`
}

func replace(context *build.Context, args interface{}) error {
	params := args.(replaceArgs)
	files, err := util.FindFiles(params.Dir, params.Replace, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for copy task: %v", err)
	}
	if len(files) < 1 {
		return nil
	}
	for _, file := range files {
		context.Message("Replacing text in file '%s'", file)
		if params.Dir != "" {
			file = filepath.Join(params.Dir, file)
		}
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file '%s': %v", file, err)
		}
		text := string(bytes)
		for old, new := range params.With {
			text = strings.Replace(text, old, new, -1)
		}
		err = ioutil.WriteFile(file, []byte(text), FileMode)
		if err != nil {
			return fmt.Errorf("writing file '%s': %v", file, err)
		}
	}
	return nil
}
