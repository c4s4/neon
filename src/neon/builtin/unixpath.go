package builtin

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "unixpath",
		Func: UnixPath,
		Help: `Convert a path to Unix format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    // convert path to unix
    uppercase("c:\foo\bar")
    // returns: "/c/foo/bar"`,
	})
}

func UnixPath(path string) string {
	return util.PathToUnix(path)
}
