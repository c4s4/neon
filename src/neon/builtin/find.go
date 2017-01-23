package builtin

import (
	zglob "github.com/mattn/go-zglob"
	"neon/build"
	"os"
	"sort"
)

func init() {
	build.BuiltinMap["find"] = build.BuiltinDescriptor{
		Function: Find,
		Help:     "Find files",
	}
}

func Find(dir string, patterns ...string) []string {
	oldDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	defer os.Chdir(oldDir)
	err = os.Chdir(dir)
	if err != nil {
		return nil
	}
	var files []string
	for _, pattern := range patterns {
		f, _ := zglob.Glob(pattern)
		for _, e := range f {
			files = append(files, e)
		}
	}
	sort.Strings(files)
	return files
}
