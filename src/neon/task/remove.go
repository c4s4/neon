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
		var files []string
		for _, patt := range patterns {
			pattern, err := target.Build.Context.ReplaceProperties(patt)
			if err != nil {
				return fmt.Errorf("evaluating argument in task remove: %v", err)
			}
			list, _ := zglob.Glob(pattern)
			for _, file := range list {
				files = append(files, file)
			}
		}
		sort.Strings(files)
		if len(files) > 0 {
			target.Build.Info("Removing %d file(s)", len(files))
			for _, file := range files {
				if err = os.Remove(file); err != nil {
					return fmt.Errorf("removing file '%s': %v", file, err)
				}
			}
		}
		return nil
	}, nil
}
