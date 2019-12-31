package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"path/filepath"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "filename",
		Func: filename,
		Help: `Return filename of a given path.

Arguments:

- The path to get filename for as a string.

Returns:

- The filename of the path as a string.

Examples:

    # get filename of path "/foo/bar/spam.txt"
    filename("/foo/bar/spam.txt")
    # returns: "spam.txt"`,
	})
}

func filename(path string) string {
	return filepath.Base(path)
}
