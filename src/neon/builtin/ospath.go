package builtin

import (
	"github.com/c4s4/neon/build"
	"github.com/c4s4/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "ospath",
		Func: osPath,
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

func osPath(path string) string {
	if util.Windows() {
		return util.PathToWindows(path)
	}
	return util.PathToUnix(path)
}
