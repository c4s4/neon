package builtin

import (
	"neon/build"
	"path/filepath"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc {
		Name: "absolute",
		Func: Absolute,
		Help: `Return absolute value of a given path.

Arguments:

- The path to get absolute value.

Returns:

- The absolute value of the path.

Examples:

    // get absolute value of path "foo/../bar/spam.txt"
    path = absolute("foo/../bar/spam.txt")
    // returns: "/home/user/build/bar/spam.txt"`,
	})
}

func Absolute(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return abs
}
