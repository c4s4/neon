package builtin

import (
	"github.com/c4s4/neon/build"

	zglob "github.com/mattn/go-zglob"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "filter",
		Func: filter,
		Help: `Filter a list of files with excludes.

Arguments:

- includes: the list of files to filter.
- excludes: a list of patterns for files to exclude.

Returns:

- The list if filtered files as a list of strings.

Examples:

    # filter text files removing those in build directory
    filter(find(".", "**/*.txt"), "build/**/*")
    # returns: files with extension "txt" in current directory and
    # subdirectories, except those in "build" directory

Notes:

- Works great with find() builtin.`,
	})
}

func filter(includes []string, excludes ...string) []string {
	var filtered []string
	for _, include := range includes {
		excluded := false
		for _, pattern := range excludes {
			exclude, err := zglob.Match(pattern, include)
			if exclude && err == nil {
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, include)
		}
	}
	return filtered
}
