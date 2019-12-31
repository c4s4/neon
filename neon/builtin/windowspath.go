package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "windowspath",
		Func: windowsPath,
		Help: `Convert a path to Windows format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path to windows
    uppercase("/c/foo/bar")
    # returns: "c:\foo\bar"`,
	})
}

func windowsPath(path string) string {
	return util.PathToWindows(path)
}
