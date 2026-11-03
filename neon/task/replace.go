package task

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddTask(build.TaskDesc{
		Name: "replace",
		Func: replace,
		Args: reflect.TypeOf(replaceArgs{}),
		Help: `Replace strings in files.

Arguments:

- replace: globs of files to process (strings, file, wrap).
- with: map with strings to replace as keys and replacements as values (map
  with string keys and values).
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file).

Examples:

    # replace foo with bar in file test.txt
    - replace: 'test.txt'
      with:    {'foo': 'bar'}`,
	})
}

type replaceArgs struct {
	Replace []string `neon:"file,wrap"`
	With    map[interface{}]interface{}
	Dir     string   `neon:"optional,file"`
	Exclude []string `neon:"optional,file"`
}

func replace(context *build.Context, args interface{}) error {
	params := args.(replaceArgs)
	with, err := util.ToMapStringString(params.With)
	if err != nil {
		return fmt.Errorf("bad with argument: %v", err)
	}
	files, err := util.FindFiles(params.Dir, params.Replace, params.Exclude, false)
	if err != nil {
		return fmt.Errorf("getting source files for replace task: %v", err)
	}
	if len(files) < 1 {
		return nil
	}
	// longest strings are replaced first, so that result doesn't depend on
	// map iteration order when strings overlap
	olds := make([]string, 0, len(with))
	for old := range with {
		olds = append(olds, old)
	}
	sort.Slice(olds, func(i, j int) bool {
		if len(olds[i]) != len(olds[j]) {
			return len(olds[i]) > len(olds[j])
		}
		return olds[i] < olds[j]
	})
	for _, file := range files {
		context.MessageArgs("Replacing text in file '%s'", file)
		if params.Dir != "" {
			file = filepath.Join(params.Dir, file)
		}
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("getting mode of file '%s': %v", file, err)
		}
		bytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file '%s': %v", file, err)
		}
		text := string(bytes)
		for _, old := range olds {
			text = strings.ReplaceAll(text, old, with[old])
		}
		// WriteFile keeps mode of existing file
		err = os.WriteFile(file, []byte(text), info.Mode().Perm())
		if err != nil {
			return fmt.Errorf("writing file '%s': %v", file, err)
		}
	}
	return nil
}
