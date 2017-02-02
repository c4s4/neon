package builtin

import (
	"neon/build"
	"path/filepath"
)

func init() {
	build.BuiltinMap["joinpath"] = build.BuiltinDescriptor{
		Function: Joinpath,
		Help: `Join file paths.

Arguments:
- The paths to join as a list of strings.
Returns:
- Joined path as a string.

Examples:
// join paths "/foo", "bar" and "spam.txt"
joinpath("foo", "bar", "spam.txt")`,
	}
}

func Joinpath(paths ...string) string {
	return filepath.Join(paths...)
}
