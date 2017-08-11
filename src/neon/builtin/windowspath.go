package builtin

import (
	"neon/build"
	"neon/util"
)

func init() {
	build.BuiltinMap["windowspath"] = build.BuiltinDescriptor{
		Function: WindowsPath,
		Help: `Convert a path to Windows format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    // convert path to windows
    uppercase("/c/foo/bar")
    // returns: "c:\foo\bar"`,
	}
}

func WindowsPath(path string) string {
	return util.PathToWindows(path)
}
