package builtin

import (
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
	"path/filepath"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "directory",
		Func: directory,
		Help: `Return directory of a given path.

Arguments:

- The path to get directory for as a string.

Returns:

- The directory of the path as a string.

Examples:

    # get directory of path "/foo/bar/spam.txt"
    dir = directory("/foo/bar/spam.txt")
    # returns: "/foo/bar"`,
	})
}

func directory(path string) string {
	return util.PathToUnix(filepath.Dir(path))
}
