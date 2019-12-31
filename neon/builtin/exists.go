package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"github.com/c4s4/neon/neon/util"
	"os"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "exists",
		Func: exists,
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

func exists(path string) bool {
	path = util.ExpandUserHome(path)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
