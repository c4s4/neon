package builtin

import (
	"neon/build"
	"os"
)

func init() {
	build.BuiltinMap["exists"] = build.BuiltinDescriptor{
		Function: Exists,
		Help: `Tells if a given pat exists.

Arguments:
- The path to test as a string.
Returns:
- A boolean telling if path exists.

Examples:
// test if given path exists
if exists("/foo/bar") { ...`,
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
