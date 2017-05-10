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
		Help: `Find files.

Arguments:

- The directory of files to find.
- The list of pattern for files to find.

Returns:

- Files as a list of strings.

Examples:

    // find all text files in book directory
    find("book", "**/*.txt")
    // returns: list of files with extension "txt"
    // find all xml and yml files in src directory
    find("src", "**/*.xml", "**/*.yml")
    // returns: list of "xml" and "yml" files

Notes:

- Files may be filtered with filter() builtin.`,
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
