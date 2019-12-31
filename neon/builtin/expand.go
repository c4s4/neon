package builtin

import (
	"github.com/c4s4/neon/neon/build"
	"os/user"
	"path/filepath"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "expand",
		Func: expand,
		Help: `Expand file name replacing ~/ with home directory.

Arguments:

- The path to expand as a string.

Returns:

- The expanded path as a string.

Examples:

    # expand path ~/.profile
    profile = expand("~/.profile")
    # returns: "/home/casa/.profile" on my machine`,
	})
}

func expand(path string) string {
	if strings.HasPrefix(path, "~/") {
		user, _ := user.Current()
		home := user.HomeDir
		path = filepath.Join(home, path[2:])
	}
	return path
}
