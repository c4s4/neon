package builtin

import (
	"strings"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "split",
		Func: split,
		Help: `Split a string.

Arguments:

- The string to split.
- The separator for splitting.

Returns:

- A list of strings.

Examples:

    # split "foo bar" with space
    split("foo bar", " ")
    # returns: ["foo", "bar"]`,
	})
}

func split(str, sep string) []string {
	return strings.Split(str, sep)
}
