package builtin

import (
	"neon/build"
	"reflect"
	"strings"
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
    joinpath("foo", "bar", "spam.txt")
    // returns: "foo/bar/spam.txt" on a Linux box and "foo\bar\spam.txt" on
    // Windows`,
	}
}

func Joinpath(paths ...interface{}) string {
	s := make([]string, len(paths))
	for i, e := range paths {
		s[i] = reflect.ValueOf(e).String()
	}
	return strings.Join(s, "/")
}
