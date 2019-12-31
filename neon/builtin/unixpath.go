package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "unixpath",
		Func: unixPath,
		Help: `Convert a path to Unix format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path to unix
    uppercase("c:\foo\bar")
    # returns: "/c/foo/bar"`,
	})
}

func unixPath(path string) string {
	return util.PathToUnix(path)
}
