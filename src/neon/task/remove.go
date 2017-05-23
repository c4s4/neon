package task

import (
	"fmt"
	zglob "github.com/mattn/go-zglob"
	"neon/build"
	"neon/util"
	"os"
	"sort"
)

func init() {
	build.TaskMap["remove"] = build.TaskDescriptor{
		Constructor: Remove,
		Help: `Remove file(s).

Arguments:

- remove: file or list of files to remove.

Examples:

    # remove all pyc files
    - remove: "**/*.pyc"`,
	}
}

func Remove(target *build.Target, args util.Object) (build.Task, error) {
	fields := []string{"remove"}
	if err := CheckFields(args, fields, fields); err != nil {
		return nil, err
	}
	patterns, err := args.GetListStringsOrString("remove")
	if err != nil {
		return nil, fmt.Errorf("remove argument must a string or list of strings")
	}
	return func() error {
		var _files []string
		for _, _p := range patterns {
			_pattern, _err := target.Build.Context.EvaluateString(_p)
			if _err != nil {
				return fmt.Errorf("evaluating argument in task remove: %v", _err)
			}
			_list, _ := zglob.Glob(_pattern)
			for _, _file := range _list {
				_files = append(_files, _file)
			}
		}
		sort.Strings(_files)
		if len(_files) > 0 {
			build.Info("Removing %d file(s)", len(_files))
			for _, _file := range _files {
				if _err := os.Remove(_file); _err != nil {
					return fmt.Errorf("removing file '%s': %v", _file, _err)
				}
			}
		}
		return nil
	}, nil
}
