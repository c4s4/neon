package builtin

import (
	"github.com/c4s4/neon/build"
	"path/filepath"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "followlink",
		Func: followLink,
		Help: `Follow symbolic link.

Arguments:

- The symbolic link to follow.

Returns:

- The path with symbolic links followed.

Examples:

    # follow symbolic link 'foo'
    followlink("foo")
    # returns: 'bar'`,
	})
}

func followLink(file string) string {
	path, err := filepath.EvalSymlinks(file)
	if err != nil {
		return file
	}
	return path
}
