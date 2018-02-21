package builtin

import (
	"neon/build"
	"reflect"
	"strings"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "joinpath",
		Func: joinPath,
		Help: `Join file paths.

Arguments:

- The paths to join as a list of strings.

Returns:

- Joined path as a string.

Examples:

    # join paths "/foo", "bar" and "spam.txt"
    joinpath("foo", "bar", "spam.txt")
    # returns: "foo/bar/spam.txt" on a Linux box and "foo\bar\spam.txt" on
    # Windows`,
	})
}

func joinPath(paths ...interface{}) string {
	s := make([]string, len(paths))
	for i, e := range paths {
		s[i] = reflect.ValueOf(e).String()
	}
	return strings.Join(s, "/")
}
