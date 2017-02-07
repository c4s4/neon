package builtin

import (
	"neon/build"
	"os/user"
	"path/filepath"
	"strings"
)

func init() {
	build.BuiltinMap["expand"] = build.BuiltinDescriptor{
		Function: Expand,
		Help: `Exapand file name by replace ~/ with home directory.

Arguments:
- The path to expand as a string.
Returns:
- The expanded path as a string.

Examples:
// expand path ~/.profile
profile = expand("~/.profile")`,
	}
}

func Expand(path string) string {
	if strings.HasPrefix(path, "~/") {
		user, _ := user.Current()
		home := user.HomeDir
		path = filepath.Join(home, path[2:])
	}
	return path
}
