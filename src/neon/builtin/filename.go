package builtin

import (
	"neon/build"
	"path/filepath"
)

func init() {
	build.BuiltinMap["filename"] = build.BuiltinDescriptor{
		Function: Filename,
		Help: `Return filename of a given path.

Arguments:
- The path to get filename for as a string.
Returns:
- The filename of the path as a string.

Examples:
// get filename of path "/foo/bar/spam.txt"
file = filename("/foo/bar/spam.txt")`,
	}
}

func Filename(path string) string {
	return filepath.Base(path)
}
