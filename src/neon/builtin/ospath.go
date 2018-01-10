package builtin

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "ospath",
		Func: OsPath,
		Help: `Convert path to running OS.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path foo/bar to OS format
    path = ospath("foo/bar")
    # will return foo/bar on Unix and foo\bar on Windows`,
	})
}

func OsPath(path string) string {
	if util.Windows() {
		return util.PathToWindows(path)
	} else {
		return util.PathToUnix(path)
	}
}
