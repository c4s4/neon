package builtin

import (
	"neon/build"
	"path/filepath"
)

func init() {
	build.BuiltinMap["directory"] = build.BuiltinDescriptor{
		Function: Directory,
		Help: `Return directory of a given path.

Arguments:

- The path to get directory for as a string.

Returns:

- The directory of the path as a string.

Examples:

    // get directory of path "/foo/bar/spam.txt"
    dir = directory("/foo/bar/spam.txt")
    // returns: "/foo/bar"`,
	}
}

func Directory(path string) string {
	return filepath.Dir(path)
}
