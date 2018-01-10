package builtin

import (
	"neon/build"
	"neon/util"
	"os"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "exists",
		Func: Exists,
		Help: `Tells if a given path exists.

Arguments:

- The path to test as a string.

Returns:

- A boolean telling if path exists.

Examples:

    # test if given path exists
    exists("/foo/bar")
    # returns: true if file "/foo/bar" exists`,
	})
}

func Exists(path string) bool {
	path = util.ExpandUserHome(path)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
