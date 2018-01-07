package builtin

import (
	"neon/build"
	"neon/util"
	"sort"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "find",
		Func: Find,
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
	})
}

func Find(dir string, patterns ...string) []string {
	files, err := util.FindFiles(dir, patterns, nil, true)
	if err != nil {
		return nil
	}
	sort.Strings(files)
	return files
}
